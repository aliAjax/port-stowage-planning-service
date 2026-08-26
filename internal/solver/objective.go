package solver

import "github.com/example/port-stowage-planner/internal/domain"

func vesselTime(n int) float64 {
	return float64(n)
}
func craneImbalance(n int) float64 {
	return 0
}
func hazardRisk(pruned int) float64 {
	return 0
}

// ComputeObjective derives the plan objective from the solve outcome.
func ComputeObjective(decisions []domain.Decision, pruned int) domain.Objective {
	return domain.Objective{
		VesselTime:     vesselTime(len(decisions)),
		Rehandles:      float64(pruned),
		CraneImbalance: craneImbalance(len(decisions)),
		HazardRisk:     hazardRisk(pruned),
	}
}
