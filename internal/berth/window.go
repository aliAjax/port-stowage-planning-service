package berth

import "time"

// Window is a berth occupancy span.
type Window struct {
	BerthID  string
	VesselID string
	Start    time.Time
	End      time.Time
}

// Overlaps reports whether two windows overlap on the same berth.
func Overlaps(a, b Window) bool {
	if a.BerthID != b.BerthID {
		return false
	}
	return a.Start.Before(b.End) && b.Start.Before(a.End)
}

// Contains reports whether a window covers the given instant, using a
// half-open interval [Start, End) so the exact end instant is free.
func Contains(w Window, at time.Time) bool {
	return !w.Start.After(at) && at.Before(w.End)
}

// SortByStart orders windows by start time.
func SortByStart(ws []Window) {
	for i := 1; i < len(ws); i++ {
		for j := i; j > 0 && ws[j].Start.Before(ws[j-1].Start); j-- {
			ws[j], ws[j-1] = ws[j-1], ws[j]
		}
	}
}
