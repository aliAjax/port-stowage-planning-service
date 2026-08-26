package reefer

import "testing"

func TestReeferAckThenAlarm(t *testing.T) {
	m := NewManager()
	if _, err := m.Attach("c1", 5, "p1"); err != nil {
		t.Fatal(err)
	}
	if err := m.Report("c1", 20); err != nil {
		t.Fatal(err)
	}
	if err := m.Acknowledge("c1"); err != nil {
		t.Fatal(err)
	}
	if err := m.Report("c1", 25); err != nil {
		t.Fatal(err)
	}
	u, _ := m.UnitOf("c1")
	if u.State != StateAlarm {
		t.Fatalf("over-temperature after acknowledge must re-raise the alarm, got %s", u.State)
	}
}

func TestReeferDuplicatePlugRejected(t *testing.T) {
	m := NewManager()
	if _, err := m.Attach("c1", 5, "p1"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Attach("c2", 5, "p1"); err == nil {
		t.Fatalf("same plug for two units must be rejected")
	}
}

func TestReeferPowerOffFreesPlug(t *testing.T) {
	m := NewManager()
	if _, err := m.Attach("c1", 5, "p1"); err != nil {
		t.Fatal(err)
	}
	if m.PlugFree("p1") {
		t.Fatalf("plug must be in use after attach")
	}
	if err := m.PowerOff("c1"); err != nil {
		t.Fatal(err)
	}
	if !m.PlugFree("p1") {
		t.Fatalf("power-off must free the plug")
	}
}

func TestReeferUnknownErrors(t *testing.T) {
	m := NewManager()
	if err := m.Report("nope", 10); err == nil {
		t.Fatalf("reporting an unknown unit must error")
	}
}
