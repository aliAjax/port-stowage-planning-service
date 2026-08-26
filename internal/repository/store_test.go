package repository

import (
	"context"
	"testing"

	"github.com/example/port-stowage-planner/internal/domain"
)

func TestSavePlanStaleRejected(t *testing.T) {
	st := New()
	p1 := domain.Plan{ID: "p1", Revision: 0}
	p2 := domain.Plan{ID: "p1", Revision: 1}
	if err := st.SavePlan(context.Background(), p1); err != nil {
		t.Fatal(err)
	}
	if err := st.SavePlan(context.Background(), p2); err != nil {
		t.Fatal(err)
	}
	stale := domain.Plan{ID: "p1", Revision: 0}
	if err := st.SavePlan(context.Background(), stale); err == nil {
		t.Fatalf("saving a stale revision must be rejected")
	}
}

func TestPlanDeepCopy(t *testing.T) {
	st := New()
	p := domain.Plan{ID: "p1", Revision: 0}
	p.Instructions = []domain.WorkInstruction{{ID: "i1", Status: "dispatched"}}
	if err := st.SavePlan(context.Background(), p); err != nil {
		t.Fatal(err)
	}
	got, _ := st.Plan(context.Background(), "p1")
	got.Instructions[0].Status = "mutated"
	again, _ := st.Plan(context.Background(), "p1")
	if again.Instructions[0].Status == "mutated" {
		t.Fatalf("Plan must return a deep copy, store instruction was mutated")
	}
}

func TestListPlansDeepCopy(t *testing.T) {
	st := New()
	p := domain.Plan{ID: "p2", Revision: 0}
	p.Explanations = map[string][]string{"c1": {"reason"}}
	if err := st.SavePlan(context.Background(), p); err != nil {
		t.Fatal(err)
	}
	got := st.ListPlans(context.Background())
	got[0].Explanations["c1"] = append(got[0].Explanations["c1"], "tampered")
	again := st.ListPlans(context.Background())
	if len(again[0].Explanations["c1"]) != 1 {
		t.Fatalf("ListPlans must deep copy, store explanations were polluted")
	}
}

func TestSnapshotDeepCopy(t *testing.T) {
	st := New()
	p := domain.Plan{ID: "p3", Revision: 0}
	p.Instructions = []domain.WorkInstruction{{ID: "i1", Status: "done"}}
	if err := st.SavePlan(context.Background(), p); err != nil {
		t.Fatal(err)
	}
	got := st.Snapshot()
	got[0].Instructions[0].Status = "mutated"
	again := st.Snapshot()
	if again[0].Instructions[0].Status == "mutated" {
		t.Fatalf("Snapshot must deep copy, store instruction was mutated")
	}
}
