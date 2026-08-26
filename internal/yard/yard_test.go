package yard

import "testing"

func yardBlocks() []Block {
	return []Block{
		{ID: "yb1", Capacity: 2, ReeferPlugs: 1, HazardClasses: []string{"1", "2"}, Destinations: []string{"TYO"}},
		{ID: "yb2", Capacity: 5, ReeferPlugs: 0, Destinations: []string{"OSA"}},
	}
}

func TestYardReeferPlugAccounting(t *testing.T) {
	a := NewAllocator(yardBlocks())
	if _, err := a.Allocate(Request{ContainerID: "r1", Reefer: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := a.Allocate(Request{ContainerID: "r2", Reefer: true}); err == nil {
		t.Fatalf("second reefer must be rejected when the only plug is taken")
	}
}

func TestYardCapacityLimit(t *testing.T) {
	a := NewAllocator([]Block{{ID: "yb1", Capacity: 2, ReeferPlugs: 1, Destinations: []string{"TYO"}}})
	for i := 0; i < 2; i++ {
		if _, err := a.Allocate(Request{ContainerID: string(rune('a' + i)), Reefer: false, Destination: "TYO"}); err != nil {
			t.Fatalf("container %d should fit: %v", i, err)
		}
	}
	if _, err := a.Allocate(Request{ContainerID: "overflow", Reefer: false, Destination: "TYO"}); err == nil {
		t.Fatalf("block capacity must be strictly enforced")
	}
}

func TestYardReleaseReturnsPlug(t *testing.T) {
	a := NewAllocator([]Block{{ID: "yb1", Capacity: 5, ReeferPlugs: 1, Destinations: []string{"TYO"}}})
	if _, err := a.Allocate(Request{ContainerID: "r1", Reefer: true, Destination: "TYO"}); err != nil {
		t.Fatal(err)
	}
	if err := a.Release("r1"); err != nil {
		t.Fatal(err)
	}
	if got := a.FreePlugs("yb1"); got != 1 {
		t.Fatalf("released plug must be returned to the block, free=%d", got)
	}
}

func TestPlugLedgerReserveBoundary(t *testing.T) {
	l := NewPlugLedger(yardBlocks())
	if err := l.Reserve("yb1"); err != nil {
		t.Fatal(err)
	}
	if err := l.Reserve("yb1"); err == nil {
		t.Fatalf("the last plug must not be over-reserved")
	}
}
