package reefer

import "fmt"

// AlarmPolicy defines when a reading becomes an alarm.
type AlarmPolicy struct {
	DeviationC float64
	RecoveryOK int
}

// DefaultPolicy tolerates 3 degrees and recovers after 3 good readings.
func DefaultPolicy() AlarmPolicy {
	return AlarmPolicy{DeviationC: 3.0, RecoveryOK: 3}
}

// Evaluate classifies a reading relative to the setpoint.
func (p AlarmPolicy) Evaluate(tempC, setpointC float64) bool {
	d := tempC - setpointC
	if d < 0 {
		d = -d
	}
	return d > p.DeviationC
}

// RecoveryThreshold returns the number of good readings needed to clear.
func (p AlarmPolicy) RecoveryThreshold() int {
	if p.RecoveryOK <= 0 {
		return 1
	}
	return p.RecoveryOK
}

// ValidateSetpoint rejects out-of-range reefer setpoints.
func ValidateSetpoint(setpointC float64) error {
	if setpointC < -40 || setpointC > 30 {
		return fmt.Errorf("setpoint %.1fC out of supported range", setpointC)
	}
	return nil
}
