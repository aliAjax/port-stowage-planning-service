package document

import (
	"fmt"
	"strings"
)

// Summary is a compact operational summary of a bay plan.
type Summary struct {
	Bay       int
	Total     int
	Conflicts int
	Rehandles int
}

// Summarize aggregates cells and conflicts into a summary.
func Summarize(plan BayPlan, conflicts []string) Summary {
	s := Summary{Bay: plan.Bay, Total: plan.Moves, Rehandles: plan.Rehandles, Conflicts: len(conflicts)}
	return s
}

// ConflictKeys finds cells that share the same row/tier.
func ConflictKeys(cells []Cell) []string {
	seen := map[string]string{}
	conflicts := make([]string, 0)
	for _, c := range cells {
		key := cellKey(c.Row, c.Tier)
		if other, ok := seen[key]; ok {
			conflicts = append(conflicts, fmt.Sprintf("%s/%s", other, c.ContainerID))
		} else {
			seen[key] = c.ContainerID
		}
	}
	return conflicts
}

// RenderCompact emits a one-line digest of the bay.
func RenderCompact(plan BayPlan) string {
	return fmt.Sprintf("bay %d: %d moves, %d rehandles", plan.Bay, plan.Moves, plan.Rehandles)
}

// CellKeys renders all occupied keys for debugging.
func CellKeys(cells []Cell) string {
	parts := make([]string, 0, len(cells))
	for _, c := range cells {
		parts = append(parts, cellKey(c.Row, c.Tier))
	}
	return strings.Join(parts, ",")
}
