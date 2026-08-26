package lashing

import (
	"fmt"
	"math"
)

// Motion captures the environmental conditions a lashing plan must resist.
type Motion struct {
	RollDeg   float64
	PitchDeg  float64
	WindKnots float64
}

// LashingPoint is a physical securing point on a cell.
type LashingPoint struct {
	Row       int
	Tier      int
	Available bool
}

// Requirement is the computed lashing need for one cell.
type Requirement struct {
	Row          int
	Tier         int
	PointsNeeded int
	ForceKN      float64
	OK           bool
	Reason       string
}

// pointCapacityKN is how much force a single lashing point can take.
const pointCapacityKN = 50.0

// ForceKN computes the horizontal securing force for a container.
func ForceKN(weightKg float64, m Motion) float64 {
	base := weightKg * 0.04 / 1000.0
	roll := m.RollDeg * 18.0
	pitch := m.PitchDeg * 8.0
	wind := m.WindKnots * m.WindKnots / 200.0
	return base + roll + pitch + wind
}

// Compute derives the securing requirement for a cell.
func Compute(row, tier int, weightKg float64, m Motion, points []LashingPoint) Requirement {
	force := ForceKN(weightKg, m)
	needed := int(math.Ceil(force / pointCapacityKN))
	if needed < 2 {
		needed = 2
	}
	if needed > 8 {
		needed = 8
	}
	req := Requirement{Row: row, Tier: tier, PointsNeeded: needed, ForceKN: force}
	avail := 0
	for _, p := range points {
		if p.Row == row && p.Tier == tier && p.Available {
			avail++
		}
	}
	if avail < needed {
		req.OK = false
		req.Reason = fmt.Sprintf("only %d lashing points available, need %d", avail, needed)
		return req
	}
	req.OK = true
	req.Reason = fmt.Sprintf("%d lashing points allocated", needed)
	return req
}

// ValidatePlan checks every cell requirement and collects warnings.
func ValidatePlan(reqs []Requirement) []string {
	warnings := make([]string, 0)
	for _, r := range reqs {
		if !r.OK {
			warnings = append(warnings, fmt.Sprintf("row %d tier %d: %s", r.Row, r.Tier, r.Reason))
		}
	}
	return warnings
}

// LashingPointsFor generates the physical points of a cell.
func LashingPointsFor(row, tier, count int) []LashingPoint {
	out := make([]LashingPoint, 0, count)
	for i := 0; i < count; i++ {
		out = append(out, LashingPoint{Row: row, Tier: tier, Available: true})
	}
	return out
}
