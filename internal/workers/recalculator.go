package workers

import (
	"context"
	"sync"
	"time"
)

type Recalculator struct {
	mu      sync.Mutex
	pending map[string]time.Time
}

func New() *Recalculator { return &Recalculator{pending: map[string]time.Time{}} }
func (r *Recalculator) Trigger(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pending[id] = time.Now()
}
func (r *Recalculator) Run(ctx context.Context, fn func(context.Context, string) error) error {
	r.mu.Lock()
	ids := make([]string, 0, len(r.pending))
	for id := range r.pending {
		ids = append(ids, id)
	}
	r.pending = map[string]time.Time{}
	r.mu.Unlock()
	for _, id := range ids {
		if err := fn(ctx, id); err != nil {
			return err
		}
	}
	return nil
}
