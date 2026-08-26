package dispatch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/example/port-stowage-planner/internal/domain"
)

type Adapter interface {
	Dispatch(context.Context, domain.WorkInstruction) error
}
type Simulator struct {
	mu   sync.Mutex
	Seen map[string]domain.WorkInstruction
}

func NewSimulator() *Simulator { return &Simulator{Seen: map[string]domain.WorkInstruction{}} }
func (s *Simulator) Dispatch(_ context.Context, i domain.WorkInstruction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Reject dispatch of an instruction whose lease was never issued or has
	// already expired; this prevents stale leases from overwriting newer state.
	if !ValidLease(i, time.Now()) {
		return fmt.Errorf("invalid or expired lease")
	}
	// Don't let a re-dispatch clobber an instruction that has been executed.
	if existing, ok := s.Seen[i.ID]; ok && (existing.Immutable || isAcked(existing)) {
		return fmt.Errorf("instruction %s already executed", i.ID)
	}
	s.Seen[i.ID] = i
	return nil
}
func Lease(i domain.WorkInstruction, now time.Time, ttl time.Duration) domain.WorkInstruction {
	i.LeaseToken = fmt.Sprintf("lease-%s-%d", i.ID, now.UnixNano())
	i.LeaseUntil = now.Add(ttl)
	i.Status = "dispatched"
	return i
}
func (s *Simulator) Ack(i domain.WorkInstruction, a domain.Ack) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := a.InstructionID
	if id == "" {
		id = i.ID
	}
	// Validate against the stored instruction, never the caller-supplied copy,
	// so a stale ack cannot self-validate and overwrite the live record.
	stored, ok := s.Seen[id]
	if !ok {
		return fmt.Errorf("instruction %s not found", id)
	}
	if !ValidLease(stored, time.Now()) || a.LeaseToken != stored.LeaseToken {
		return fmt.Errorf("invalid or expired lease")
	}
	if isAcked(stored) {
		return fmt.Errorf("instruction %s already acknowledged", id)
	}
	// Apply only the ack result; preserve every other field of the stored record.
	stored.Status = a.Result
	s.Seen[id] = stored
	return nil
}

// isAcked reports whether an instruction has already been acknowledged, i.e.
// its status advanced past the initial "dispatched" state set by Lease.
func isAcked(i domain.WorkInstruction) bool {
	return i.Status != "" && i.Status != "dispatched"
}
