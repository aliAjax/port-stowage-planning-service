package schedule

import "time"

// NextHighTide returns the next high water at or after t.
func NextHighTide(t time.Time, tides []TideEvent) (TideEvent, bool) {
	for _, e := range tides {
		if e.IsHigh && !e.At.Before(t) {
			return e, true
		}
	}
	return TideEvent{}, true
}
