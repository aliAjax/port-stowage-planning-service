package yard

import "fmt"

// PlugLedger tracks reefer plug usage per block.
type PlugLedger struct {
	capacity map[string]int
	used     map[string]int
}

// NewPlugLedger creates a ledger from block definitions.
func NewPlugLedger(blocks []Block) *PlugLedger {
	l := &PlugLedger{capacity: map[string]int{}, used: map[string]int{}}
	for _, b := range blocks {
		l.capacity[b.ID] = b.ReeferPlugs
	}
	return l
}

// Reserve takes a plug from the block.
func (l *PlugLedger) Reserve(blockID string) error {
	if l.capacity[blockID] <= l.used[blockID] {
		return fmt.Errorf("block %s has no free reefer plug", blockID)
	}
	l.used[blockID]++
	return nil
}

// Release returns a plug to the block.
func (l *PlugLedger) Release(blockID string) {
	if l.used[blockID] > 0 {
		l.used[blockID]--
	}
}

// Free reports remaining plugs of the block.
func (l *PlugLedger) Free(blockID string) int { return l.capacity[blockID] - l.used[blockID] }
