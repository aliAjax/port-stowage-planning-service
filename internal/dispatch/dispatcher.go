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
	if a.LeaseToken != i.LeaseToken || time.Now().After(i.LeaseUntil) {
		return fmt.Errorf("invalid or expired lease")
	}
	i.Status = a.Result
	s.Seen[i.ID] = i
	return nil
}
