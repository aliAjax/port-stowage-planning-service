package dispatch

import (
	"fmt"
	"time"

	"github.com/example/port-stowage-planner/internal/domain"
)

// IssueLease mints a fresh lease token for an instruction.
func IssueLease(i domain.WorkInstruction, now time.Time, ttl time.Duration) domain.WorkInstruction {
	i.LeaseToken = fmt.Sprintf("lease-%s-%d", i.ID, now.UnixNano())
	i.LeaseUntil = now.Add(ttl)
	i.Status = "dispatched"
	return i
}

// ValidLease checks token and expiry of an instruction.
func ValidLease(i domain.WorkInstruction, now time.Time) bool {
	if i.LeaseToken == "" {
		return false
	}
	return now.Before(i.LeaseUntil)
}
