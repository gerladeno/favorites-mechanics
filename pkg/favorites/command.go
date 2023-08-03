package favorites

func (m *Manager) AddCommand(name, exec string, parentID int, nextID int) {
	if name == "" && exec == "" {
		return
	}

	defer m.notifySinker()

	dir := m.getDirByID(parentID)
	next := m.getEntryByID(nextID)
	node := dir.AddElement(m.newEntry(name, exec, false, parentID), nil, next)
	m.registerEntry(node)
}

func (m *Manager) DeleteCommand(id int) {
	node := m.getEntryByID(id)
	if node == nil {
		return
	}

	defer m.notifySinker()

	dir := m.getDirByID(node.Value.ParentID)
	dir.DeleteElement(node)
	m.unregisterEntry(node.Value.ID)
}

func (m *Manager) ModifyExec(id int, exec string) {
	node := m.getEntryByID(id)
	if node == nil {
		return
	}

	defer m.notifySinker()

	node.Value.Exec = exec
}
