package inspection

import (
	"testing"
	"time"
)

func outbound(seal string, declared, vgm float64) GateRecord {
	return GateRecord{ContainerID: "c1", SealNo: seal, Direction: DirectionOut, DeclaredWeightKg: declared, VGMWeightKg: vgm, At: time.Now()}
}

func TestGateHoldsNoSharedBackingArray(t *testing.T) {
	g := NewInspector()
	res := g.Inspect(outbound("", 1000, 2000))
	if len(res.Holds) == 0 {
		t.Fatalf("expected holds")
	}
	res.Holds[0] = "tampered"
	after := g.HoldsFor("c1")
	for _, h := range after {
		if h == "tampered" {
			t.Fatalf("mutating the returned holds polluted internal state")
		}
	}
}

func TestGateReleaseClearsHolds(t *testing.T) {
	g := NewInspector()
	g.Inspect(outbound("", 1000, 2000))
	if err := g.Release("c1", "seal-9"); err != nil {
		t.Fatalf("release failed: %v", err)
	}
	if got := g.HoldsFor("c1"); len(got) != 0 {
		t.Fatalf("holds must be cleared after release, got %v", got)
	}
}

func TestHoldsForNoShared(t *testing.T) {
	g := NewInspector()
	g.Inspect(outbound("", 1000, 2000))
	got := g.HoldsFor("c1")
	got[0] = "mutated"
	if again := g.HoldsFor("c1"); again[0] == "mutated" {
		t.Fatalf("HoldsFor returned a shared slice")
	}
}

func TestHoldLedgerNoPollution(t *testing.T) {
	l := NewHoldLedger()
	l.Add("c1", "missing seal", "weight discrepancy")
	got := l.Get("c1")
	got[0] = "polluted"
	again := l.Get("c1")
	if again[0] == "polluted" {
		t.Fatalf("HoldLedger.Get returned a shared backing array")
	}
}
