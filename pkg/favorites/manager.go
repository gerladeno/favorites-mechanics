package favorites

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/gerladeno/favorites-mechanics/pkg/list"
)

type logger interface {
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
}

//go:generate options-gen -out-filename=manager_options.gen.go -from-struct=Options
type Options struct {
	inMemory         bool          `option:"mandatory"`
	configPath       string        `option:"mandatory" validate:"required"`
	syncConfigPeriod time.Duration `option:"mandatory" validate:"required"`
	maxDisplayLen    int           `option:"mandatory" validate:"required"`
}

type Manager struct {
	log              logger
	opts             Options
	root             *list.DeLinkedList[entry]
	mu               sync.RWMutex
	syncNotification chan struct{}
	EntryIDs         map[int]*list.Node[entry]
	maxID            int
}

// entry is an internal type for management.
type entry struct {
	ID        int
	Name      string
	Exec      string
	ParentID  int
	Entries   *list.DeLinkedList[entry]
	IsDir     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewManager(ctx context.Context, log logger, opts Options) (*Manager, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("opts.Validate(): %w", err)
	}

	manager := Manager{ //nolint:exhaustruct
		log:              log,
		opts:             opts,
		root:             new(list.DeLinkedList[entry]),
		syncNotification: make(chan struct{}),
		EntryIDs:         make(map[int]*list.Node[entry]),
		maxID:            0,
	}

	if !opts.inMemory {
		manager.syncNotification = make(chan struct{}, 1)

		if err := manager.readConfig(); err != nil {
			return nil, fmt.Errorf("readConfig(): %w", err)
		}

		go manager.runSyncer(ctx)
	}

	return &manager, nil
}

func (m *Manager) readConfig() error {
	file, err := os.Open(m.opts.configPath)
	if err != nil {
		pathErr := &os.PathError{} //nolint:exhaustruct
		if errors.As(err, &pathErr) && errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("os.Open(m.opts.configPath): %w", err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			m.log.Warn("file.Close():", err)
		}
	}()

	var entries []Entry

	if err = yaml.NewDecoder(file).Decode(&entries); err != nil {
		return fmt.Errorf("yaml.NewDecoder(file).Decode(&bytes): %w", err)
	}

	m.setRoot(entries)

	return nil
}

func (m *Manager) runSyncer(ctx context.Context) {
	ticker := time.NewTicker(m.opts.syncConfigPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.SyncIn()
		case <-m.syncNotification:
			m.SyncOut()
			ticker.Reset(m.opts.syncConfigPeriod)
		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) SyncIn() {
	if err := m.readConfig(); err != nil {
		m.log.Warn("m.readConfig():", err)
	}
}

func (m *Manager) SyncOut() {
	entries := make([]Entry, 0, m.root.Len())

	m.mu.RLock()
	for _, elem := range m.root.List() {
		entries = append(entries, m.entry2ExternalEntry(elem, true))
	}
	m.mu.RUnlock()

	bytes, err := yaml.Marshal(entries)
	if err != nil {
		m.log.Warn("yaml.Marshal(m.root):", err)

		return
	}

	file, err := os.Create(m.opts.configPath)
	if err != nil {
		m.log.Warn("os.Create(m.opts.configPath):", err)

		return
	}

	defer func() {
		if err = file.Close(); err != nil {
			m.log.Warn("file.Close():", err)
		}
	}()

	if _, err = file.Write(bytes); err != nil {
		m.log.Warn("file.Write(bytes):", err)

		return
	}
}

func (m *Manager) notifySinker() {
	if m.opts.inMemory {
		return
	}

	select {
	case m.syncNotification <- struct{}{}:
	default:
	}
}

func (m *Manager) getDirByID(id int) *list.DeLinkedList[entry] {
	e := m.getEntryByID(id)
	if e == nil {
		return m.root
	}

	return e.Value.Entries
}

func (m *Manager) getEntryByID(id int) *list.Node[entry] {
	if id == 0 {
		return nil
	}

	e, ok := m.EntryIDs[id]
	if !ok {
		return nil
	}

	return e
}

func (m *Manager) DisplayEntry(entry *Entry) string {
	if entry == nil {
		return ""
	}

	if entry.Name != "" {
		return entry.Name
	}

	if len(entry.Exec) <= m.opts.maxDisplayLen {
		return entry.Exec
	}

	return entry.Exec[:m.opts.maxDisplayLen-3] + "..."
}

func (m *Manager) newEntry(name, exec string, isDir bool, parentID int) entry {
	m.maxID++

	dir := list.DeLinkedList[entry]{}

	return entry{
		ID:        m.maxID,
		Name:      name,
		Exec:      exec,
		ParentID:  parentID,
		Entries:   &dir,
		IsDir:     isDir,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (m *Manager) registerEntry(e *list.Node[entry]) {
	m.EntryIDs[e.Value.ID] = e
}

func (m *Manager) unregisterEntry(id int) {
	delete(m.EntryIDs, id)
}

func (m *Manager) MoveEntry(targetID, parentID, nextID int) {
	node := m.getEntryByID(targetID)
	if node == nil {
		return
	}

	defer m.notifySinker()
	defer func() {
		node.Value.UpdatedAt = time.Now()
	}()

	dir := m.getDirByID(parentID)
	if node.Value.ParentID == parentID {
		dir.MoveItem(node, nil, m.getEntryByID(nextID))

		return
	}

	currentDir := m.getDirByID(node.Value.ParentID)
	currentDir.DeleteElement(node)
	m.unregisterEntry(targetID)

	node = dir.AddElement(node.Value, nil, m.getEntryByID(nextID))
	m.registerEntry(node)
	node.Value.ParentID = parentID
}

func (m *Manager) RenameEntry(targetID int, name string) {
	node := m.getEntryByID(targetID)
	if node == nil {
		return
	}

	defer m.notifySinker()

	node.Value.Name = name
}

func (m *Manager) ListDirectory(id int) []Entry {
	m.mu.RLock()
	l := m.getDirByID(id).List()
	m.mu.RUnlock()

	result := make([]Entry, 0, len(l))

	for _, elem := range l {
		result = append(result, m.entry2ExternalEntry(elem, false))
	}

	return result
}

func (m *Manager) entry2ExternalEntry(entry entry, withSubDirs bool) Entry {
	var entries []Entry

	if withSubDirs {
		for _, elem := range entry.Entries.List() {
			entries = append(entries, m.entry2ExternalEntry(elem, true))
		}
	}

	return Entry{
		ID:        entry.ID,
		Name:      entry.Name,
		Exec:      entry.Exec,
		ParentID:  entry.ParentID,
		Entries:   entries,
		IsDir:     entry.IsDir,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
	}
}

func (m *Manager) entry2InternalEntry(exEntry Entry, entryIDs map[int]*list.Node[entry]) entry {
	entries := &list.DeLinkedList[entry]{}

	for _, elem := range exEntry.Entries {
		e := entries.AddElement(m.entry2InternalEntry(elem, entryIDs), nil, nil)
		entryIDs[elem.ID] = e
	}

	return entry{
		ID:        exEntry.ID,
		Name:      exEntry.Name,
		Exec:      exEntry.Exec,
		ParentID:  exEntry.ParentID,
		Entries:   entries,
		IsDir:     exEntry.IsDir,
		CreatedAt: exEntry.CreatedAt,
		UpdatedAt: exEntry.UpdatedAt,
	}
}

// Entry is entry representation for external use.
type Entry struct {
	ID        int       `yaml:"id"`
	Name      string    `yaml:"name"`
	Exec      string    `yaml:"exec"`
	ParentID  int       `yaml:"parentId"`
	Entries   []Entry   `yaml:"entries"`
	IsDir     bool      `yaml:"isDir"`
	CreatedAt time.Time `yaml:"createdAt"`
	UpdatedAt time.Time `yaml:"updatedAt"`
}

func (m *Manager) setRoot(entries []Entry) {
	entryIDs := make(map[int]*list.Node[entry])
	root := &list.DeLinkedList[entry]{}

	for _, elem := range entries {
		node := root.AddElement(m.entry2InternalEntry(elem, entryIDs), nil, nil)
		entryIDs[node.Value.ID] = node
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.root = root
	m.EntryIDs = entryIDs
}
