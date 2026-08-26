package reefer

// freePlug returns a power plug to the pool so it can be reused.
func (m *Manager) freePlug(plugID string) {
	if plugID == "" {
		return
	}
	delete(m.plugs, plugID)
}
