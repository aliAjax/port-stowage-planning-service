package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/port-stowage-planner/internal/dispatch"
	"github.com/example/port-stowage-planner/internal/repository"
)

func TestHTTPSolveNilMap(t *testing.T) {
	srv := New(repository.New(), dispatch.NewSimulator())
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	voy := `{"id":"voy-1","vessel_id":"vessel-demo","port_id":"port-1","eta":"2026-08-27T00:00:00Z","etd":"2026-08-28T00:00:00Z"}`
	if resp, err := http.Post(ts.URL+"/api/v1/voyages", "application/json", bytes.NewBufferString(voy)); err != nil {
		t.Fatal(err)
	} else {
		resp.Body.Close()
	}
	box := `{"id":"box-1","voyage_id":"voy-1","iso_size":"40","weight_kg":12000,"destination":"TYO","on_deck":true,"priority":1}`
	if resp, err := http.Post(ts.URL+"/api/v1/containers", "application/json", bytes.NewBufferString(box)); err != nil {
		t.Fatal(err)
	} else {
		resp.Body.Close()
	}
	body := `{"id":"plan-1","voyage_id":"voy-1"}`
	resp, err := http.Post(ts.URL+"/api/v1/plans/solve", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("solve request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		var e map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&e)
		t.Fatalf("solve must return 201, got %d (%v)", resp.StatusCode, e)
	}
}

func TestHTTPExplainNoNull(t *testing.T) {
	srv := New(repository.New(), dispatch.NewSimulator())
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()
	if err := srv.Store.SavePlan(nil, emptyPlan()); err != nil {
		t.Fatal(err)
	}
	resp, err := http.Get(ts.URL + "/api/v1/plans/p1/explain")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("explain must return a JSON object, got error %v", err)
	}
}
