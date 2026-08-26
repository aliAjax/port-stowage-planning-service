package solver

import "github.com/example/port-stowage-planner/internal/domain"

// Explain renders a decision for the explain endpoint. The returned map is
// independent of the underlying Decision: copying Reasons into a fresh slice
// prevents a caller from mutating the decision's constraints when it edits
// one explanation entry.
func Explain(d domain.Decision) map[string]any {
	reasons := make([]string, len(d.Reasons))
	copy(reasons, d.Reasons)
	return map[string]any{"container_id": d.ContainerID, "slot": d.Slot, "score": d.Score, "constraints": reasons}
}
