package berth

import "fmt"

// AssignmentState is the lifecycle state of a berth assignment.
type AssignmentState string

const (
	StateBooked    AssignmentState = "booked"
	StateActive    AssignmentState = "active"
	StateCompleted AssignmentState = "completed"
	StateReleased  AssignmentState = "released"
)

// CanTransition validates a single lifecycle step.
func CanTransition(from, to AssignmentState) bool {
	switch from {
	case StateBooked:
		return to == StateActive || to == StateReleased
	case StateActive:
		return to == StateCompleted
	case StateCompleted, StateReleased:
		return false
	}
	return false
}

// ValidateTransition returns a descriptive error for illegal moves.
func ValidateTransition(from, to AssignmentState) error {
	if CanTransition(from, to) {
		return nil
	}
	return fmt.Errorf("illegal berth assignment transition %s -> %s", from, to)
}
