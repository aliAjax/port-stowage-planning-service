package schedule

import (
	"fmt"
	"time"
)

// BuildWindows derives all usable tidal windows from the table.
func BuildWindows(minDraftM float64, tides []TideEvent) ([]Window, error) {
	if len(tides) < 2 {
		return nil, fmt.Errorf("tide table too short to derive windows")
	}
	ws := make([]Window, 0)
	start := time.Time{}
	var peak float64
	for _, t := range tides {
		h, err := AvailableDraft(t.At, tides)
		if err != nil {
			return nil, err
		}
		if h >= minDraftM {
			if start.IsZero() {
				start = t.At
			}
			if h > peak {
				peak = h
			}
			continue
		}
		if !start.IsZero() {
			ws = append(ws, Window{Start: start, End: t.At, MaxDraftM: peak})
			start = time.Time{}
			peak = 0
		}
	}
	if !start.IsZero() {
		last := tides[len(tides)-1].At
		ws = append(ws, Window{Start: start, End: last, MaxDraftM: peak})
	}
	return ws, nil
}

// ClosestHighTide finds the high tide nearest to t.
func ClosestHighTide(t time.Time, tides []TideEvent) (TideEvent, error) {
	var best TideEvent
	found := false
	for _, e := range tides {
		if !e.IsHigh {
			continue
		}
		if !found {
			best = e
			found = true
			continue
		}
		if absDur(e.At.Sub(t)) < absDur(best.At.Sub(t)) {
			best = e
		}
	}
	if !found {
		return TideEvent{}, fmt.Errorf("no high tide in table")
	}
	return best, nil
}

func absDur(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
