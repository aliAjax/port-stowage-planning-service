package solver

import (
	"testing"

	"github.com/example/port-stowage-planner/internal/domain"
)

func dec(id string) domain.Decision {
	return domain.Decision{ContainerID: id, Slot: domain.Slot{Bay: 1, Row: 1, Tier: 1}, Reasons: []string{"weight within slot limit"}}
}

func TestExplainReasonsCopy(t *testing.T) {
	d := dec("c1")
	got := Explain(d)
	got["constraints"].([]string)[0] = "tampered"
	again := Explain(d)
	if again["constraints"].([]string)[0] == "tampered" {
		t.Fatalf("Explain must return a copy of the reasons slice")
	}
}

func TestObjectiveCraneImbalance(t *testing.T) {
	obj := ComputeObjective([]domain.Decision{dec("c1"), dec("c2"), dec("c3"), dec("c4"), dec("c5")}, 0)
	if obj.CraneImbalance == 0 {
		t.Fatalf("crane imbalance must be non-zero for an uneven load")
	}
}

func TestObjectiveHazardRisk(t *testing.T) {
	obj := ComputeObjective([]domain.Decision{dec("c1")}, 4)
	if obj.HazardRisk == 0 {
		t.Fatalf("hazard risk must reflect pruned containers")
	}
}

func TestObjectiveVesselTime(t *testing.T) {
	obj := ComputeObjective([]domain.Decision{dec("c1"), dec("c2"), dec("c3")}, 0)
	if obj.VesselTime != 6 {
		t.Fatalf("vessel time must scale with the number of decisions, got %v", obj.VesselTime)
	}
}
