package repository

import "github.com/example/port-stowage-planner/internal/domain"

// Snapshot returns a point-in-time view of all plans. The returned plans are
// deep copies and may be mutated freely without affecting the store.
func (s *Store) Snapshot() []domain.Plan {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Plan, 0, len(s.plans))
	for _, p := range s.plans {
		out = append(out, clonePlan(p))
	}
	return out
}
