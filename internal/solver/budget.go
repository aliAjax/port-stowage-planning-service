package solver

import (
	"context"
	"time"
)

// Budget is the solve time budget tracker.
type Budget struct {
	start time.Time
	limit time.Duration
}

// NewBudget starts a budget for a solve run.
func NewBudget(limit time.Duration) Budget {
	return Budget{start: time.Now(), limit: limit}
}

// Exhausted reports whether the budget has run out.
func (b Budget) Exhausted() bool { return time.Since(b.start) > b.limit }

// Remaining returns the remaining budget.
func (b Budget) Remaining() time.Duration {
	left := b.limit - time.Since(b.start)
	if left < 0 {
		return 0
	}
	return left
}

// Cancelled reports whether the context is done.
func Cancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
