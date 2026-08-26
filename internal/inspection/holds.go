package inspection

import "fmt"

// HoldLedger stores active holds per container.
type HoldLedger struct {
	holds map[string][]string
}

// NewHoldLedger creates an empty ledger.
func NewHoldLedger() *HoldLedger {
	return &HoldLedger{holds: map[string][]string{}}
}

// Add registers one or more holds for a container.
func (l *HoldLedger) Add(containerID string, holds ...string) {
	if len(holds) == 0 {
		return
	}
	l.holds[containerID] = append(l.holds[containerID], holds...)
}

// Get returns the current holds of a container.
func (l *HoldLedger) Get(containerID string) []string {
	out := make([]string, 0, len(l.holds[containerID]))
	out = append(out, l.holds[containerID]...)
	return out
}

// Clear removes all holds of a container.
func (l *HoldLedger) Clear(containerID string) error {
	if _, ok := l.holds[containerID]; !ok {
		return fmt.Errorf("container %s has no holds", containerID)
	}
	delete(l.holds, containerID)
	return nil
}

// Count reports how many containers currently have holds.
func (l *HoldLedger) Count() int { return len(l.holds) }
