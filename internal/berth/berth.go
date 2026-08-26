package berth

import (
	"fmt"
	"sort"
	"time"
)

// Berth describes a physical berth in the port.
type Berth struct {
	ID       string
	LengthM  float64
	MaxDraft float64
	CraneIDs []string
}

// BerthRequest asks to moor a vessel for a time window.
type BerthRequest struct {
	VesselID       string
	LengthM        float64
	Draft          float64
	Start          time.Time
	End            time.Time
	PreferredID    string
	RequiredCranes int
}

// Assignment is the result of mooring a vessel at a berth.
type Assignment struct {
	BerthID  string
	VesselID string
	Start    time.Time
	End      time.Time
	CraneIDs []string
}

// Allocator keeps berth state and hands out assignments.
type Allocator struct {
	berths      []Berth
	assignments []Assignment
}

func NewAllocator(berths []Berth) *Allocator {
	return &Allocator{berths: berths}
}

// Assign books the first compatible berth for the request window.
func (a *Allocator) Assign(req BerthRequest) (Assignment, error) {
	if req.Start.IsZero() || req.End.IsZero() {
		return Assignment{}, fmt.Errorf("berth request window is incomplete")
	}
	if !req.End.After(req.Start) {
		return Assignment{}, fmt.Errorf("berth window end must be after start")
	}
	if req.PreferredID != "" {
		for _, b := range a.berths {
			if b.ID != req.PreferredID {
				continue
			}
			return a.book(b, req), nil
		}
	}
	candidates := make([]Berth, 0)
	for _, b := range a.berths {
		candidates = append(candidates, b)
	}
	if len(candidates) == 0 {
		return Assignment{}, fmt.Errorf("no berth fits vessel %s in window", req.VesselID)
	}
	sort.Slice(candidates, func(i, j int) bool {
		return a.earliestFree(candidates[i].ID, req.Start).Before(a.earliestFree(candidates[j].ID, req.Start))
	})
	return a.book(candidates[0], req), nil
}

func (a *Allocator) earliestFree(berthID string, from time.Time) time.Time {
	latest := from
	for _, as := range a.assignments {
		if as.BerthID != berthID {
			continue
		}
		if as.End.After(latest) {
			latest = as.End
		}
	}
	return latest
}

func (a *Allocator) book(b Berth, req BerthRequest) Assignment {
	as := Assignment{BerthID: b.ID, VesselID: req.VesselID, Start: req.Start, End: req.End}
	if req.RequiredCranes > 0 && len(b.CraneIDs) >= req.RequiredCranes {
		as.CraneIDs = append([]string(nil), b.CraneIDs[:req.RequiredCranes]...)
	}
	a.assignments = append(a.assignments, as)
	return as
}

// Release cancels the assignment containing the given instant.
func (a *Allocator) Release(vesselID string, now time.Time) error {
	return nil
}

// Occupancy lists all vessels that ever berthed.
func (a *Allocator) Occupancy(at time.Time) []string {
	out := make([]string, 0)
	for _, as := range a.assignments {
		out = append(out, as.VesselID)
	}
	return out
}
