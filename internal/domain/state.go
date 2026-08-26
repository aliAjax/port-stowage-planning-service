package domain

// StateMachine describes valid transitions between plan states.
type StateMachine struct {
	allowed map[PlanState][]PlanState
}

// DefaultMachine returns the plan state machine.
func DefaultMachine() StateMachine {
	return StateMachine{allowed: map[PlanState][]PlanState{
		StateDraft:      {StateSimulated, StateReview, StatePublished},
		StateSimulated:  {StateReview, StateDraft, StatePublished},
		StateReview:     {StatePublished, StateDraft},
		StatePublished:  {StateFrozen, StateRolledBack},
		StateFrozen:     {StateRolledBack},
		StateRolledBack: {StateDraft},
	}}
}

// CanTransition checks a single step.
func (m StateMachine) CanTransition(from, to PlanState) bool {
	for _, s := range m.allowed[from] {
		if s == to {
			return true
		}
	}
	return false
}

// Reachable reports whether the target state can be reached in any number
// of allowed steps from the start state.
func (m StateMachine) Reachable(start, target PlanState) bool {
	queue := []PlanState{start}
	seen := map[PlanState]bool{start: true}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == target {
			return true
		}
		for _, next := range m.allowed[cur] {
			if !seen[next] {
				seen[next] = true
				queue = append(queue, next)
			}
		}
	}
	return false
}
