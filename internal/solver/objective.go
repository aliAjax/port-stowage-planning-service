package solver

import (
	"math"

	"github.com/example/port-stowage-planner/internal/domain"
)

// moveMinutesPerContainer is the deterministic per-move handling time used
// when no crane-specific rate is available. It matches the historical inline
// objective so existing regressions stay stable.
const moveMinutesPerContainer = 2.0

// vesselTime estimates total berth time as the number of placed containers
// times the per-move handling time. Pruned (unplaced) containers do not add
// moves but represent lost throughput, so they are not counted here; the
// rehandle and hazard terms capture their cost separately.
func vesselTime(n int) float64 {
	return float64(n) * moveMinutesPerContainer
}

// craneImbalance measures how evenly the bay load is spread across the
// assigned quay cranes. With no explicit crane list the solver assumes every
// vessel is worked by quayCranes cranes, each covering a contiguous bay range.
// The imbalance is the difference between the busiest and idlest crane's
// move count; a perfectly balanced spread yields 0.
const quayCranes = 3

func craneImbalance(decisions []domain.Decision) float64 {
	if len(decisions) == 0 {
		return 0
	}
	loads := make([]int, quayCranes)
	maxBay := 0
	for _, d := range decisions {
		if d.Slot.Bay > maxBay {
			maxBay = d.Slot.Bay
		}
	}
	if maxBay == 0 {
		return 0
	}
	for _, d := range decisions {
		// Assign each bay to a crane by dividing the bay range evenly.
		c := (d.Slot.Bay - 1) * quayCranes / maxBay
		if c >= quayCranes {
			c = quayCranes - 1
		}
		if c < 0 {
			c = 0
		}
		loads[c]++
	}
	min, max := loads[0], loads[0]
	for _, l := range loads {
		if l < min {
			min = l
		}
		if l > max {
			max = l
		}
	}
	return math.Abs(float64(max - min))
}

// hazardRisk grows with the number of pruned containers, since pruning a
// hazardous container for an isolation violation is itself a risk signal: it
// means the plan could not safely stow that cargo. Each pruned unit adds a
// fixed penalty so the term varies with the prune count instead of staying
// flat at zero.
const hazardPenaltyPerPrune = 1.5

func hazardRisk(pruned int) float64 {
	return float64(pruned) * hazardPenaltyPerPrune
}

// ComputeObjective derives the plan objective from the solve outcome. It is
// the single source of truth for objective scoring; Solve must call it rather
// than recomputing an inline copy that can drift out of sync.
func ComputeObjective(decisions []domain.Decision, pruned int) domain.Objective {
	return domain.Objective{
		VesselTime:     vesselTime(len(decisions)),
		Rehandles:      float64(pruned),
		CraneImbalance: craneImbalance(decisions),
		HazardRisk:     hazardRisk(pruned),
	}
}
