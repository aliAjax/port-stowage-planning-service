package reefer

import (
	"fmt"
	"math"
)

const (
	StateMonitoring   = "monitoring"
	StateAlarm        = "alarm"
	StateAcknowledged = "acknowledged"
	StatePoweredOff   = "powered_off"
)

type Unit struct {
	ContainerID string
	SetpointC   float64
	TempC       float64
	State       string
	AlarmCount  int
	OKReports   int
	Plugged     bool
	plugID      string
}

type Manager struct {
	units map[string]*Unit
	plugs map[string]bool
}

func NewManager() *Manager {
	return &Manager{units: map[string]*Unit{}, plugs: map[string]bool{}}
}

// Attach registers a reefer unit on a power plug. A plug already in use by
// another unit is rejected so two reefers can never share one outlet.
func (m *Manager) Attach(containerID string, setpointC float64, plugID string) (*Unit, error) {
	if _, ok := m.units[containerID]; ok {
		return nil, fmt.Errorf("reefer %s already attached", containerID)
	}
	if plugID == "" {
		return nil, fmt.Errorf("reefer %s requires a non-empty plug id", containerID)
	}
	if !m.PlugFree(plugID) {
		return nil, fmt.Errorf("plug %s already in use", plugID)
	}
	u := &Unit{ContainerID: containerID, SetpointC: setpointC, TempC: setpointC, State: StateMonitoring, Plugged: true, plugID: plugID}
	m.units[containerID] = u
	m.plugs[plugID] = true
	return u, nil
}

// Report ingests a temperature reading. A reading for a unit that is not
// registered returns an error instead of being silently dropped.
func (m *Manager) Report(containerID string, tempC float64) error {
	u, ok := m.units[containerID]
	if !ok {
		return fmt.Errorf("reefer %s unknown", containerID)
	}
	if u.State == StatePoweredOff {
		return fmt.Errorf("reefer %s is powered off", containerID)
	}
	u.TempC = tempC
	delta := math.Abs(tempC - u.SetpointC)
	if delta > 3.0 {
		u.AlarmCount++
		u.OKReports = 0
		// A new over-temp condition always re-raises the alarm, even when the
		// previous alarm was only acknowledged. Without this the platform stays
		// silent after a confirmed alarm re-heats.
		u.State = StateAlarm
		return nil
	}
	if u.State == StateAlarm || u.State == StateAcknowledged {
		u.OKReports++
		if u.OKReports >= 3 {
			u.State = StateMonitoring
			u.OKReports = 0
		}
	}
	return nil
}

// Acknowledge silences an active alarm.
func (m *Manager) Acknowledge(containerID string) error {
	u, ok := m.units[containerID]
	if !ok {
		return fmt.Errorf("reefer %s unknown", containerID)
	}
	if u.State != StateAlarm {
		return fmt.Errorf("reefer %s has no active alarm", containerID)
	}
	u.State = StateAcknowledged
	u.OKReports = 0
	return nil
}

// PowerOff detaches the unit from power and returns its plug to the pool.
func (m *Manager) PowerOff(containerID string) error {
	u, ok := m.units[containerID]
	if !ok {
		return fmt.Errorf("reefer %s unknown", containerID)
	}
	u.State = StatePoweredOff
	u.Plugged = false
	m.freePlug(u.plugID)
	return nil
}

// UnitOf returns the unit by container id.
func (m *Manager) UnitOf(containerID string) (*Unit, bool) {
	u, ok := m.units[containerID]
	return u, ok
}

// InAlarm lists units currently in an alarm or acknowledged state.
func (m *Manager) InAlarm() []string {
	out := make([]string, 0)
	for _, u := range m.units {
		if u.State == StateAlarm || u.State == StateAcknowledged {
			out = append(out, u.ContainerID)
		}
	}
	return out
}

// PlugFree reports whether a power plug is available.
func (m *Manager) PlugFree(plugID string) bool {
	return !m.plugs[plugID]
}
