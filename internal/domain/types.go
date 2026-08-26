package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

type PlanState string

const (
	StateDraft      PlanState = "draft"
	StateSimulated  PlanState = "simulated"
	StateReview     PlanState = "review"
	StatePublished  PlanState = "published"
	StateFrozen     PlanState = "frozen"
	StateRolledBack PlanState = "rolled_back"
)

type Port struct{ ID, Name, Timezone string }
type Vessel struct {
	ID, Name string
	IMO      string
	Bays     int
	MaxDraft float64
}
type Voyage struct {
	ID        string    `json:"id"`
	VesselID  string    `json:"vessel_id"`
	PortID    string    `json:"port_id"`
	ETA       time.Time `json:"eta"`
	ETD       time.Time `json:"etd"`
	TideStart time.Time `json:"tide_start"`
	TideEnd   time.Time `json:"tide_end"`
}
type YardBlock struct {
	ID                string
	Rows, Bays, Tiers int
	ReeferPlugs       int
	HazardClasses     []string
}
type CraneCapability struct {
	ContainerSize   []int
	MaxWeight       float64
	TwinLift        bool
	MinSafetyMeters float64
}
type QuayCrane struct {
	ID, PortID    string
	Capability    CraneCapability
	AvailableFrom time.Time
	AvailableTo   time.Time
}
type Container struct {
	ID            string  `json:"id"`
	VoyageID      string  `json:"voyage_id"`
	ISOSize       string  `json:"iso_size"`
	Destination   string  `json:"destination"`
	WeightKg      float64 `json:"weight_kg"`
	HazardClass   string  `json:"hazard_class"`
	Reefer        bool    `json:"reefer"`
	OnDeck        bool    `json:"on_deck"`
	RequiresPower bool    `json:"requires_power"`
	Priority      int     `json:"priority"`
	Sequence      int     `json:"sequence"`
}
type Slot struct {
	Bay, Row, Tier int
	OnDeck         bool
	CoveredByHatch bool
	HazardClass    string
	MaxWeight      float64
}
type Decision struct {
	ContainerID string
	Slot        Slot
	Score       float64
	Reasons     []string
}
type WorkInstruction struct {
	ID, PlanID, CraneID, ContainerID string
	Sequence                         int
	LeaseToken                       string
	LeaseUntil                       time.Time
	Status                           string
	CreatedAt                        time.Time
	Immutable                        bool
}
type Plan struct {
	ID, VoyageID, Version, ContentHash string
	State                              PlanState
	Decisions                          []Decision
	Instructions                       []WorkInstruction
	Explanations                       map[string][]string
	Objective                          Objective
	CreatedAt, UpdatedAt               time.Time
	Revision                           int
}
type Objective struct{ VesselTime, Rehandles, CraneImbalance, HazardRisk float64 }
type Ack struct {
	InstructionID, LeaseToken, Result, Message string
	At                                         time.Time
}

func (p Plan) StableHash() string {
	ids := make([]string, 0, len(p.Decisions))
	for _, d := range p.Decisions {
		ids = append(ids, fmt.Sprintf("%s:%d:%d:%d", d.ContainerID, d.Slot.Bay, d.Slot.Row, d.Slot.Tier))
	}
	sort.Strings(ids)
	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%v", p.VoyageID, p.Version, ids)))
	return hex.EncodeToString(h[:])
}

func ValidTransition(from, to PlanState) bool {
	allowed := map[PlanState][]PlanState{StateDraft: {StateSimulated, StateReview, StatePublished}, StateSimulated: {StateReview, StateDraft, StatePublished}, StateReview: {StatePublished, StateDraft}, StatePublished: {StateFrozen, StateRolledBack}, StateFrozen: {StateRolledBack}, StateRolledBack: {StateDraft}}
	for _, s := range allowed[from] {
		if s == to {
			return true
		}
	}
	return false
}
