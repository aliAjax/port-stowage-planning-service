package schedule

import (
	"testing"
	"time"
)

func th(h, m int, height float64, high bool) TideEvent {
	return TideEvent{At: time.Date(2026, 8, 27, h, m, 0, 0, time.UTC), HeightM: height, IsHigh: high}
}

func TestTideEmptyErrors(t *testing.T) {
	_, err := AvailableDraft(time.Now(), nil)
	if err == nil {
		t.Fatalf("empty tide table must report an error")
	}
}

func TestTideOutOfOrderErrors(t *testing.T) {
	tides := []TideEvent{th(10, 0, 5.0, true), th(9, 0, 1.0, false)}
	_, err := AvailableDraft(tides[0].At, tides)
	if err == nil {
		t.Fatalf("out-of-order tide table must report an error")
	}
}

func TestFindWindowEmptyErrors(t *testing.T) {
	_, err := FindWindow(3.0, time.Hour, nil, time.Now())
	if err == nil {
		t.Fatalf("FindWindow on an empty table must report an error")
	}
}

func TestNextHighTideNotFound(t *testing.T) {
	tides := []TideEvent{th(10, 0, 5.0, false), th(12, 0, 1.0, false)}
	_, ok := NextHighTide(time.Now(), tides)
	if ok {
		t.Fatalf("table without high tide must report not-found")
	}
}
