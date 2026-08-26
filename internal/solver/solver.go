package solver

import (
	"context"
	"fmt"
	"github.com/example/port-stowage-planner/internal/domain"
	"sort"
	"time"
)

type Candidate struct {
	Slots  []domain.Slot
	Budget time.Duration
}
type Result struct {
	Decisions []domain.Decision
	Objective domain.Objective
	Feasible  bool
	Pruned    int
	BestKnown bool
}

func GenerateSlots(vessel domain.Vessel) []domain.Slot {
	slots := make([]domain.Slot, 0, vessel.Bays*4*3)
	for b := 1; b <= vessel.Bays; b++ {
		for r := 1; r <= 4; r++ {
			for t := 1; t <= 3; t++ {
				slots = append(slots, domain.Slot{Bay: b, Row: r, Tier: t, OnDeck: t == 3, CoveredByHatch: t < 3, MaxWeight: 50000})
			}
		}
	}
	return slots
}

func Solve(ctx context.Context, vessel domain.Vessel, voyage domain.Voyage, cs []domain.Container, cranes []domain.QuayCrane, budget time.Duration) Result {
	_ = voyage
	_ = cranes
	start := time.Now()
	sort.Slice(cs, func(i, j int) bool {
		if cs[i].Priority != cs[j].Priority {
			return cs[i].Priority > cs[j].Priority
		}
		return cs[i].ID < cs[j].ID
	})
	slots := GenerateSlots(vessel)
	used := map[string]bool{}
	decisions := []domain.Decision{}
	pruned := 0
	for _, c := range cs {
		if err := c.Validate(); err != nil {
			pruned++
			continue
		}
		placed := false
		for _, sl := range slots {
			if time.Since(start) > budget {
				break
			}
			select {
			case <-ctx.Done():
				return Result{Decisions: decisions, Pruned: pruned, BestKnown: true}
			default:
			}
			if used[fmt.Sprintf("%d/%d/%d", sl.Bay, sl.Row, sl.Tier)] || c.WeightKg > sl.MaxWeight || (!c.OnDeck && sl.OnDeck) || (c.RequiresPower && sl.Tier < 2) {
				pruned++
				continue
			}
			if c.HazardClass != "" && sl.Row == 1 {
				pruned++
				continue
			}
			used[fmt.Sprintf("%d/%d/%d", sl.Bay, sl.Row, sl.Tier)] = true
			reason := []string{"weight within slot limit", "deterministic first-fit slot"}
			if c.HazardClass != "" {
				reason = append(reason, "hazard isolation row rule")
			}
			decisions = append(decisions, domain.Decision{ContainerID: c.ID, Slot: sl, Score: float64(sl.Bay*10 + sl.Row), Reasons: reason})
			placed = true
			break
		}
		if !placed {
			pruned++
		}
	}
	obj := ComputeObjective(decisions, pruned)
	return Result{Decisions: decisions, Objective: obj, Feasible: len(decisions) == len(cs), Pruned: pruned, BestKnown: time.Since(start) > budget}
}
