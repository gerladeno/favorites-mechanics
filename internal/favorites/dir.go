package favorites

func (m *Manager) AddDir(name string, parentID int, nextID int) {
	if name == "" {
		return
	}

	defer m.notifySinker()

	dir := m.getDirByID(parentID)
	next := m.getEntryByID(nextID)
	node := dir.AddElement(m.newEntry(name, "", true, dir), nil, next)
	m.registerEntry(node)
}

func (m *Manager) DeleteDir(id int) {
	node := m.getEntryByID(id)
	if node == nil {
		return
	}

	for _, elem := range node.Value.Entries.List() {
		m.unregisterEntry(elem.ID)
	}

	node.Value.Parent.DeleteElement(node)
	m.unregisterEntry(node.Value.ID)
}
