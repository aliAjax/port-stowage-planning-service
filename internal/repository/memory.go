package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/example/port-stowage-planner/internal/domain"
)

// ErrStalePlan is returned when a plan being saved is older than the plan
// already stored under the same ID. Callers must re-read the current plan and
// reapply their changes on top of the latest revision before saving again.
var ErrStalePlan = errors.New("plan is stale; a newer revision already exists")

type Store struct {
	mu         sync.RWMutex
	ports      map[string]domain.Port
	vessels    map[string]domain.Vessel
	voyages    map[string]domain.Voyage
	containers map[string]domain.Container
	cranes     map[string]domain.QuayCrane
	plans      map[string]domain.Plan
	acks       []domain.Ack
}

func New() *Store {
	return &Store{ports: map[string]domain.Port{}, vessels: map[string]domain.Vessel{}, voyages: map[string]domain.Voyage{}, containers: map[string]domain.Container{}, cranes: map[string]domain.QuayCrane{}, plans: map[string]domain.Plan{}}
}
func (s *Store) PutPort(_ context.Context, v domain.Port) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ports[v.ID] = v
}
func (s *Store) PutVessel(_ context.Context, v domain.Vessel) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.vessels[v.ID] = v
}
func (s *Store) PutVoyage(_ context.Context, v domain.Voyage) error {
	if err := v.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.voyages[v.ID] = v
	return nil
}
func (s *Store) PutContainer(_ context.Context, c domain.Container) error {
	if err := c.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.containers[c.ID] = c
	return nil
}
func (s *Store) PutCrane(_ context.Context, c domain.QuayCrane) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cranes[c.ID] = c
}
func (s *Store) Voyage(_ context.Context, id string) (domain.Voyage, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.voyages[id]
	return v, ok
}
func (s *Store) Containers(_ context.Context, voyage string) []domain.Container {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Container{}
	for _, c := range s.containers {
		if c.VoyageID == voyage {
			out = append(out, c)
		}
	}
	return out
}
func (s *Store) Cranes(_ context.Context, port string) []domain.QuayCrane {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.QuayCrane{}
	for _, c := range s.cranes {
		if c.PortID == port {
			out = append(out, c)
		}
	}
	return out
}
func (s *Store) SavePlan(_ context.Context, p domain.Plan) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cur, ok := s.plans[p.ID]; ok {
		// Reject writes that would regress the stored revision: a caller
		// operating on a stale copy must not clobber a newer revision.
		if p.Revision < cur.Revision {
			return ErrStalePlan
		}
	}
	s.plans[p.ID] = clonePlan(p)
	return nil
}
func (s *Store) Plan(_ context.Context, id string) (domain.Plan, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.plans[id]
	if !ok {
		return domain.Plan{}, false
	}
	return clonePlan(p), true
}
func (s *Store) ListPlans(_ context.Context) []domain.Plan {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Plan, 0, len(s.plans))
	for _, p := range s.plans {
		out = append(out, clonePlan(p))
	}
	return out
}
func (s *Store) AddAck(_ context.Context, a domain.Ack) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.acks = append(s.acks, a)
}

// clonePlan returns a deep copy of p so that callers cannot mutate the store's
// internal slices or maps by editing a returned plan.
func clonePlan(p domain.Plan) domain.Plan {
	if p.Decisions != nil {
		d := make([]domain.Decision, len(p.Decisions))
		for i, dec := range p.Decisions {
			d[i] = dec
			if dec.Reasons != nil {
				d[i].Reasons = append([]string(nil), dec.Reasons...)
			}
		}
		p.Decisions = d
	}
	if p.Instructions != nil {
		p.Instructions = append([]domain.WorkInstruction(nil), p.Instructions...)
	}
	if p.Explanations != nil {
		e := make(map[string][]string, len(p.Explanations))
		for k, v := range p.Explanations {
			e[k] = append([]string(nil), v...)
		}
		p.Explanations = e
	}
	return p
}
