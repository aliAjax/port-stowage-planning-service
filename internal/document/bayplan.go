package document

import (
	"fmt"
	"sort"
	"strings"

	"github.com/example/port-stowage-planner/internal/domain"
)

// Cell identifies one occupied slot in a bay.
type Cell struct {
	Row         int
	Tier        int
	ContainerID string
}

// BayPlan is a compiled view of one bay for reporting.
type BayPlan struct {
	Bay       int
	Cells     []Cell
	Moves     int
	Rehandles int
}

// Build compiles decisions for one bay into a bay plan and counts the
// rehandles a discharge of the bay would cause.
func Build(decisions []domain.Decision, bay int) BayPlan {
	plan := BayPlan{Bay: bay}
	byKey := map[string]domain.Decision{}
	for _, d := range decisions {
		if d.Slot.Bay != bay {
			continue
		}
		key := cellKey(d.Slot.Row, d.Slot.Tier)
		byKey[key] = d
	}
	for key := range byKey {
		var row, tier int
		fmt.Sscanf(key, "%d/%d", &row, &tier)
		plan.Cells = append(plan.Cells, Cell{Row: row, Tier: tier, ContainerID: byKey[key].ContainerID})
	}
	sort.Slice(plan.Cells, func(i, j int) bool {
		if plan.Cells[i].Tier != plan.Cells[j].Tier {
			return plan.Cells[i].Tier < plan.Cells[j].Tier
		}
		return plan.Cells[i].Row < plan.Cells[j].Row
	})
	plan.Moves = len(plan.Cells)
	for _, c := range plan.Cells {
		plan.Rehandles += RehandlesAbove(byKey, c.Row, c.Tier)
	}
	return plan
}

// RehandlesAbove counts containers stacked in higher tiers of the same row
// that must be moved before this cell can be discharged.
func RehandlesAbove(byKey map[string]domain.Decision, row, tier int) int {
	n := 0
	for key := range byKey {
		var r, t int
		fmt.Sscanf(key, "%d/%d", &r, &t)
		if r == row && t > tier {
			n++
		}
	}
	return n
}

// Render produces a compact text representation of the bay.
func Render(plan BayPlan) string {
	var b strings.Builder
	fmt.Fprintf(&b, "bay %d moves=%d rehandles=%d\n", plan.Bay, plan.Moves, plan.Rehandles)
	for _, c := range plan.Cells {
		fmt.Fprintf(&b, "  r%d t%d %s\n", c.Row, c.Tier, c.ContainerID)
	}
	return b.String()
}

func cellKey(row, tier int) string { return fmt.Sprintf("%d/%d", row, tier) }
