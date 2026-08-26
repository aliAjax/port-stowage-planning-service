package solver

import "github.com/example/port-stowage-planner/internal/domain"

func Explain(d domain.Decision) map[string]any {
	return map[string]any{"container_id": d.ContainerID, "slot": d.Slot, "score": d.Score, "constraints": d.Reasons}
}
