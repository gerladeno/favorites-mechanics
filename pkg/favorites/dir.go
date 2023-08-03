package favorites

func (m *Manager) AddDir(name string, parentID int, nextID int) {
	if name == "" {
		return
	}

	defer m.notifySinker()

	dir := m.getDirByID(parentID)
	next := m.getEntryByID(nextID)
	node := dir.AddElement(m.newEntry(name, "", true, parentID), nil, next)
	m.registerEntry(node)
}

func (m *Manager) DeleteDir(id int) {
	node := m.getEntryByID(id)
	if node == nil {
		return
	}

	defer m.notifySinker()

	for _, elem := range node.Value.Entries.List() {
		if elem.IsDir {
			m.DeleteDir(elem.ID)
		}

		m.unregisterEntry(elem.ID)
	}

	dir := m.getDirByID(node.Value.ParentID)
	dir.DeleteElement(node)
	m.unregisterEntry(node.Value.ID)
}
