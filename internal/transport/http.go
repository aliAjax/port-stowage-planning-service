package transport

import (
	"encoding/json"
	"fmt"
	"github.com/example/port-stowage-planner/internal/dispatch"
	"github.com/example/port-stowage-planner/internal/domain"
	"github.com/example/port-stowage-planner/webui"
	"github.com/example/port-stowage-planner/internal/repository"
	"github.com/example/port-stowage-planner/internal/simulation"
	"github.com/example/port-stowage-planner/internal/solver"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	Store   *repository.Store
	Adapter *dispatch.Simulator
	mux     *http.ServeMux
	vessel  domain.Vessel
}

func New(st *repository.Store, a *dispatch.Simulator) *Server {
	s := &Server{Store: st, Adapter: a, vessel: domain.Vessel{ID: "vessel-demo", Name: "Demo Vessel", Bays: 12, MaxDraft: 14.2}, mux: http.NewServeMux()}
	s.routes()
	return s
}
func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"status":"ok"}`)) })
	s.mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"status":"ready"}`)) })
	s.mux.HandleFunc("/api/v1/voyages", s.voyages)
	s.mux.HandleFunc("/api/v1/containers", s.containers)
	s.mux.HandleFunc("/api/v1/plans", s.plans)
	s.mux.HandleFunc("/api/v1/plans/solve", s.solve)
	s.mux.HandleFunc("/api/v1/plans/", s.planAction)
	s.mux.Handle("/", webui.Handler())
}
func (s *Server) Handler() http.Handler { return logging(s.mux) }
func write(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
func decode(r *http.Request, v any) error { return json.NewDecoder(r.Body).Decode(v) }
func (s *Server) voyages(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		write(w, s.Store.ListPlans(r.Context()), 200)
		return
	}
	var v domain.Voyage
	if err := decode(r, &v); err != nil {
		write(w, map[string]string{"error": err.Error()}, 400)
		return
	}
	if err := s.Store.PutVoyage(r.Context(), v); err != nil {
		write(w, map[string]string{"error": err.Error()}, 400)
		return
	}
	write(w, v, 201)
}
func (s *Server) containers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		write(w, map[string]string{"error": "method not allowed"}, 405)
		return
	}
	var c domain.Container
	if err := decode(r, &c); err != nil {
		write(w, map[string]string{"error": err.Error()}, 400)
		return
	}
	if err := s.Store.PutContainer(r.Context(), c); err != nil {
		write(w, map[string]string{"error": err.Error()}, 400)
		return
	}
	write(w, c, 201)
}

type solveRequest struct {
	ID       string `json:"id"`
	VoyageID string `json:"voyage_id"`
}

func (s *Server) solve(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		write(w, map[string]string{"error": "method not allowed"}, 405)
		return
	}
	var req solveRequest
	if err := decode(r, &req); err != nil {
		write(w, map[string]string{"error": err.Error()}, 400)
		return
	}
	v, ok := s.Store.Voyage(r.Context(), req.VoyageID)
	if !ok {
		write(w, map[string]string{"error": "voyage not found"}, 404)
		return
	}
	cs := s.Store.Containers(r.Context(), req.VoyageID)
	res := solver.Solve(r.Context(), s.vessel, v, cs, nil, 500*time.Millisecond)
	p := domain.Plan{ID: req.ID, VoyageID: req.VoyageID, Version: "v1", State: domain.StateDraft, Decisions: res.Decisions, Objective: res.Objective, Explanations: map[string][]string{}, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Revision: 0}
	for _, d := range res.Decisions {
		p.Explanations[d.ContainerID] = d.Reasons
	}
	p.ContentHash = p.StableHash()
	if err := s.Store.SavePlan(r.Context(), p); err != nil {
		write(w, map[string]string{"error": err.Error()}, 409)
		return
	}
	write(w, map[string]any{"plan": p, "feasible": res.Feasible, "pruned": res.Pruned, "best_known": res.BestKnown}, 201)
}
func (s *Server) plans(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		write(w, s.Store.ListPlans(r.Context()), 200)
		return
	}
	write(w, map[string]string{"error": "method not allowed"}, 405)
}
func (s *Server) planAction(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 5 {
		write(w, map[string]string{"error": "invalid path"}, 400)
		return
	}
	id, action := parts[3], parts[4]
	p, ok := s.Store.Plan(r.Context(), id)
	if !ok {
		write(w, map[string]string{"error": "plan not found"}, 404)
		return
	}
	switch action {
	case "publish":
		if !domain.ValidTransition(p.State, domain.StatePublished) {
			write(w, map[string]string{"error": "invalid transition"}, 409)
			return
		}
		p.State = domain.StatePublished
		p.Revision++
		p.UpdatedAt = time.Now().UTC()
		for i, d := range p.Decisions {
			p.Instructions = append(p.Instructions, dispatch.Lease(domain.WorkInstruction{ID: fmt.Sprintf("%s-i-%d", p.ID, i+1), PlanID: p.ID, ContainerID: d.ContainerID, Sequence: i + 1, CreatedAt: time.Now().UTC(), Immutable: true}, time.Now().UTC(), 10*time.Minute))
		}
		_ = s.Store.SavePlan(r.Context(), p)
		write(w, p, 200)
	case "freeze":
		if !domain.ValidTransition(p.State, domain.StateFrozen) {
			write(w, map[string]string{"error": "invalid transition"}, 409)
			return
		}
		p.State = domain.StateFrozen
		p.Revision++
		_ = s.Store.SavePlan(r.Context(), p)
		write(w, p, 200)
	case "rollback":
		if !domain.ValidTransition(p.State, domain.StateRolledBack) {
			write(w, map[string]string{"error": "invalid transition"}, 409)
			return
		}
		p.State = domain.StateRolledBack
		p.Revision++
		_ = s.Store.SavePlan(r.Context(), p)
		write(w, p, 200)
	case "simulate":
		write(w, simulation.Run(r.Context(), p), 200)
	case "explain":
		write(w, p.Explanations, 200)
	case "instructions":
		write(w, p.Instructions, 200)
	default:
		write(w, map[string]string{"error": "unknown action"}, 404)
	}
}
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf(`{"method":%q,"path":%q,"duration_ms":%d}`+"\n", r.Method, r.URL.Path, time.Since(start).Milliseconds())
	})
}
