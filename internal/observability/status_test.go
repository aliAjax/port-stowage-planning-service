package observability

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestOverallAllOk(t *testing.T) {
	r := NewRegistry()
	r.Set("db", "ok")
	r.Set("queue", "degraded")
	if got := r.Overall(); got != "degraded" {
		t.Fatalf("degraded component must be reflected, got %q", got)
	}
}

func TestSetOverwrites(t *testing.T) {
	r := NewRegistry()
	r.Set("db", "ok")
	r.Set("db", "degraded")
	if got := r.Overall(); got != "degraded" {
		t.Fatalf("Set must overwrite the previous value, got %q", got)
	}
}

func TestOverallEmptyUnknown(t *testing.T) {
	r := NewRegistry()
	if got := r.Overall(); got != "unknown" {
		t.Fatalf("empty registry must report unknown, got %q", got)
	}
}

func TestHealthReflectsRegistry(t *testing.T) {
	r := NewRegistry()
	r.Set("db", "degraded")
	h := Health(r)
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest("GET", "/healthz", nil))
	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "degraded" {
		t.Fatalf("health must reflect registry status, got %v", body["status"])
	}
}

func TestRegistryComponentsCopy(t *testing.T) {
	r := NewRegistry()
	r.Set("db", "ok")
	got := r.Components()
	got["db"] = "degraded"
	if r.Components()["db"] != "ok" {
		t.Fatalf("Components must return a copy, mutation leaked into the registry")
	}
}
