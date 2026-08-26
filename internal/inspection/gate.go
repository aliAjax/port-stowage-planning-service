package inspection

import (
	"fmt"
	"time"
)

const (
	DirectionIn  = "in"
	DirectionOut = "out"
)

// GateRecord is one container pass through the terminal gate.
type GateRecord struct {
	ContainerID      string
	SealNo           string
	Direction        string
	DeclaredWeightKg float64
	VGMWeightKg      float64
	HazardClass      string
	At               time.Time
}

// GateResult is the inspector verdict for a gate pass.
type GateResult struct {
	Accepted     bool
	Holds        []string
	TolerancePct float64
}

// Inspector validates gate passes and tracks holds per container.
type Inspector struct {
	holds map[string][]string
}

func NewInspector() *Inspector {
	return &Inspector{holds: map[string][]string{}}
}

// WeightTolerance is the maximum relative difference between declared
// weight and VGM before the pass is held.
const WeightTolerance = 0.05

// Inspect evaluates one gate pass.
func (g *Inspector) Inspect(r GateRecord) GateResult {
	res := GateResult{TolerancePct: WeightTolerance * 100}
	holds := make([]string, 0)
	if r.Direction != DirectionIn && r.Direction != DirectionOut {
		holds = append(holds, "invalid direction")
	}
	if r.Direction == DirectionOut && r.SealNo == "" {
		holds = append(holds, "missing seal")
	}
	if r.VGMWeightKg <= 0 {
		holds = append(holds, "missing VGM weight")
	} else if r.DeclaredWeightKg > 0 {
		diff := r.VGMWeightKg - r.DeclaredWeightKg
		if diff < 0 {
			diff = -diff
		}
		if diff/r.DeclaredWeightKg > WeightTolerance {
			holds = append(holds, "weight discrepancy")
		}
	}
	if r.Direction == DirectionIn && r.HazardClass == "" {
		holds = append(holds, "hazard class not declared")
	}
	if len(holds) > 0 {
		g.holds[r.ContainerID] = append(g.holds[r.ContainerID], holds...)
	}
	res.Holds = append([]string(nil), g.holds[r.ContainerID]...)
	res.Accepted = len(res.Holds) == 0
	return res
}

// Release clears holds once the operator resolves them.
func (g *Inspector) Release(containerID, sealNo string) error {
	holds, ok := g.holds[containerID]
	if !ok || len(holds) == 0 {
		return fmt.Errorf("container %s has no holds", containerID)
	}
	for _, h := range holds {
		if h == "missing seal" && sealNo == "" {
			return fmt.Errorf("seal still missing for %s", containerID)
		}
	}
	delete(g.holds, containerID)
	return nil
}

// HoldsFor returns the container's current holds.
func (g *Inspector) HoldsFor(containerID string) []string {
	return append([]string(nil), g.holds[containerID]...)
}
