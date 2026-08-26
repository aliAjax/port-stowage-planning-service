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

// ValidLease reports whether an instruction has a non-empty lease token that
// has not yet expired. An instruction issued without a lease token is never
// valid, regardless of the expiry field.
func ValidLease(i domain.WorkInstruction, now time.Time) bool {
	return i.LeaseToken != "" && now.Before(i.LeaseUntil)
}
