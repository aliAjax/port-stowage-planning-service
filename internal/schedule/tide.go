package schedule

import (
	"fmt"
	"time"
)

// TideEvent is a high or low water mark in chronological order.
type TideEvent struct {
	At      time.Time
	HeightM float64
	IsHigh  bool
}

// Window is a contiguous span where the available draft never drops below
// the vessel requirement.
type Window struct {
	Start     time.Time
	End       time.Time
	MaxDraftM float64
}

// AvailableDraft returns the interpolated tide height at the given instant.
func AvailableDraft(t time.Time, tides []TideEvent) (float64, error) {
	if len(tides) == 0 {
		return 0, fmt.Errorf("tide table is empty")
	}
	for i := 0; i < len(tides)-1; i++ {
		if tides[i+1].At.Before(tides[i].At) {
			return 0, fmt.Errorf("tide table out of order at index %d", i+1)
		}
	}
	if t.Before(tides[0].At) {
		return tides[0].HeightM, nil
	}
	for i := 0; i < len(tides)-1; i++ {
		a, b := tides[i], tides[i+1]
		if t.Equal(a.At) || (t.After(a.At) && t.Before(b.At)) {
			if b.At.Equal(a.At) {
				return a.HeightM, nil
			}
			f := float64(t.Sub(a.At)) / float64(b.At.Sub(a.At))
			return a.HeightM + f*(b.HeightM-a.HeightM), nil
		}
	}
	return tides[len(tides)-1].HeightM, nil
}

// FindWindow locates the earliest window of at least minDuration during
// which the available draft stays at or above minDraftM.
func FindWindow(minDraftM float64, minDuration time.Duration, tides []TideEvent, from time.Time) (Window, error) {
	if minDuration <= 0 {
		return Window{}, fmt.Errorf("minimum duration must be positive")
	}
	if len(tides) == 0 {
		return Window{}, fmt.Errorf("tide table is empty")
	}
	start := time.Time{}
	var peak float64
	for i := 0; i < len(tides); i++ {
		t := tides[i]
		if t.At.Before(from) {
			continue
		}
		h, err := AvailableDraft(t.At, tides)
		if err != nil {
			return Window{}, err
		}
		above := h >= minDraftM
		if above && start.IsZero() {
			start = t.At
			peak = h
		}
		if above && h > peak {
			peak = h
		}
		if !above && !start.IsZero() {
			if t.At.Sub(start) >= minDuration {
				return Window{Start: start, End: t.At, MaxDraftM: peak}, nil
			}
			start = time.Time{}
			peak = 0
		}
	}
	if !start.IsZero() {
		last := tides[len(tides)-1].At
		if last.Sub(start) >= minDuration {
			return Window{Start: start, End: last, MaxDraftM: peak}, nil
		}
	}
	return Window{}, nil
}
