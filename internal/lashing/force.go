package lashing

import "fmt"

// ForceModel computes securing force from weight and motion.
type ForceModel struct {
	RollFactor  float64
	PitchFactor float64
	WindFactor  float64
}

// DefaultModel is the standard force model.
func DefaultModel() ForceModel {
	return ForceModel{RollFactor: 18.0, PitchFactor: 8.0, WindFactor: 0.005}
}

// HorizontalForce returns the horizontal securing force in kN.
func (m ForceModel) HorizontalForce(weightKg float64, windKnots float64, rollDeg, pitchDeg float64) (float64, error) {
	if weightKg <= 0 {
		return 0, fmt.Errorf("container weight must be positive")
	}
	if windKnots < 0 {
		return 0, fmt.Errorf("wind speed cannot be negative")
	}
	base := weightKg * 0.04 / 1000.0
	roll := rollDeg * m.RollFactor
	pitch := pitchDeg * m.PitchFactor
	wind := windKnots * windKnots * m.WindFactor
	return base + roll + pitch + wind, nil
}
