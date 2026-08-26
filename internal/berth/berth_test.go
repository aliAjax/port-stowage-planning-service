package berth

import (
	"testing"
	"time"
)

func at(h, m int) time.Time {
	return time.Date(2026, 8, 27, h, m, 0, 0, time.UTC)
}

func berths() []Berth {
	return []Berth{
		{ID: "b1", LengthM: 400, MaxDraft: 16, CraneIDs: []string{"qc1", "qc2"}},
		{ID: "b2", LengthM: 350, MaxDraft: 14, CraneIDs: []string{"qc3"}},
	}
}

func TestBerthOverlapRejected(t *testing.T) {
	a := NewAllocator([]Berth{{ID: "b1", LengthM: 400, MaxDraft: 16, CraneIDs: []string{"qc1"}}})
	if _, err := a.Assign(BerthRequest{VesselID: "V1", LengthM: 300, Draft: 12, Start: at(10, 0), End: at(12, 0)}); err != nil {
		t.Fatal(err)
	}
	if _, err := a.Assign(BerthRequest{VesselID: "V2", LengthM: 300, Draft: 12, Start: at(11, 0), End: at(13, 0)}); err == nil {
		t.Fatalf("overlapping assignments on the same berth must be rejected")
	}
}

func TestBerthDraftCompatibility(t *testing.T) {
	a := NewAllocator([]Berth{{ID: "b1", LengthM: 400, MaxDraft: 10, CraneIDs: []string{"qc1"}}})
	if _, err := a.Assign(BerthRequest{VesselID: "V1", LengthM: 300, Draft: 15, Start: at(10, 0), End: at(12, 0)}); err == nil {
		t.Fatalf("vessel draft deeper than berth must be rejected")
	}
}

func TestBerthReleaseFrees(t *testing.T) {
	a := NewAllocator(berths())
	if _, err := a.Assign(BerthRequest{VesselID: "V1", LengthM: 300, Draft: 12, Start: at(10, 0), End: at(12, 0)}); err != nil {
		t.Fatal(err)
	}
	if err := a.Release("V1", at(11, 0)); err != nil {
		t.Fatal(err)
	}
	if got := a.Occupancy(at(11, 30)); len(got) != 0 {
		t.Fatalf("released vessel must not occupy the berth, occupancy=%v", got)
	}
}

func TestBerthOccupancyCurrent(t *testing.T) {
	a := NewAllocator(berths())
	if _, err := a.Assign(BerthRequest{VesselID: "V1", LengthM: 300, Draft: 12, Start: at(10, 0), End: at(12, 0)}); err != nil {
		t.Fatal(err)
	}
	if got := a.Occupancy(at(14, 0)); len(got) != 0 {
		t.Fatalf("vessel whose window ended must not occupy the berth, occupancy=%v", got)
	}
}

func TestWindowContainsEndExclusive(t *testing.T) {
	w := Window{BerthID: "b1", Start: at(10, 0), End: at(12, 0)}
	if Contains(w, at(12, 0)) {
		t.Fatalf("window end instant must not be occupied")
	}
	if !Contains(w, at(11, 0)) {
		t.Fatalf("instant inside the window must be occupied")
	}
}
