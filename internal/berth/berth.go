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

// Assign books a compatible, conflict-free berth for the request window.
func (a *Allocator) Assign(req BerthRequest) (Assignment, error) {
	if req.Start.IsZero() || req.End.IsZero() {
		return Assignment{}, fmt.Errorf("berth request window is incomplete")
	}
	if !req.End.After(req.Start) {
		return Assignment{}, fmt.Errorf("berth window end must be after start")
	}
	if req.PreferredID != "" {
		var preferred *Berth
		for i := range a.berths {
			if a.berths[i].ID == req.PreferredID {
				preferred = &a.berths[i]
				break
			}
		}
		if preferred == nil {
			return Assignment{}, fmt.Errorf("berth %s not found", req.PreferredID)
		}
		return a.book(*preferred, req)
	}
	candidates := make([]Berth, 0, len(a.berths))
	for _, b := range a.berths {
		candidates = append(candidates, b)
	}
	sort.Slice(candidates, func(i, j int) bool {
		return a.earliestFree(candidates[i].ID, req.Start).Before(a.earliestFree(candidates[j].ID, req.Start))
	})
	for i := range candidates {
		as, err := a.book(candidates[i], req)
		if err == nil {
			return as, nil
		}
	}
	return Assignment{}, fmt.Errorf("no berth fits vessel %s in window", req.VesselID)
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

// book reserves a berth after validating compatibility and absence of
// overlapping bookings. It returns an error instead of mutating state when
// the vessel exceeds the berth's dimensions or the window collides with an
// existing assignment on the same berth.
func (a *Allocator) book(b Berth, req BerthRequest) (Assignment, error) {
	if req.LengthM > b.LengthM {
		return Assignment{}, fmt.Errorf("vessel %s length %.1f exceeds berth %s length %.1f", req.VesselID, req.LengthM, b.ID, b.LengthM)
	}
	if req.Draft > b.MaxDraft {
		return Assignment{}, fmt.Errorf("vessel %s draft %.1f exceeds berth %s max draft %.1f", req.VesselID, req.Draft, b.ID, b.MaxDraft)
	}
	cand := Window{BerthID: b.ID, VesselID: req.VesselID, Start: req.Start, End: req.End}
	for _, as := range a.assignments {
		if Overlaps(cand, Window{BerthID: as.BerthID, VesselID: as.VesselID, Start: as.Start, End: as.End}) {
			return Assignment{}, fmt.Errorf("berth %s occupied by vessel %s during window", b.ID, as.VesselID)
		}
	}
	out := Assignment{BerthID: b.ID, VesselID: req.VesselID, Start: req.Start, End: req.End}
	if req.RequiredCranes > 0 && len(b.CraneIDs) >= req.RequiredCranes {
		out.CraneIDs = append([]string(nil), b.CraneIDs[:req.RequiredCranes]...)
	}
	a.assignments = append(a.assignments, out)
	return out, nil
}

// Release removes the booking for vesselID whose window covers now. It
// returns an error when no matching booking exists.
func (a *Allocator) Release(vesselID string, now time.Time) error {
	for i, as := range a.assignments {
		if as.VesselID != vesselID {
			continue
		}
		if Contains(Window{BerthID: as.BerthID, VesselID: as.VesselID, Start: as.Start, End: as.End}, now) {
			a.assignments = append(a.assignments[:i], a.assignments[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("no active booking for vessel %s at %s", vesselID, now.Format(time.RFC3339))
}

// Occupancy lists vessels whose booking window covers at.
func (a *Allocator) Occupancy(at time.Time) []string {
	out := make([]string, 0)
	for _, as := range a.assignments {
		if Contains(Window{BerthID: as.BerthID, VesselID: as.VesselID, Start: as.Start, End: as.End}, at) {
			out = append(out, as.VesselID)
		}
	}
	return out
}
