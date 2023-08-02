package favorites

import (
	"context"
	"fmt"
	"time"

	"github.com/gerladeno/favorites-mechanics/pkg/list"
)

type logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

//go:generate options-gen -out-filename=manager_options.gen.go -from-struct=Options
type Options struct {
	inMemory         bool          `option:"mandatory" validate:"required"`
	configPath       string        `option:"mandatory" validate:"required"`
	syncConfigPeriod time.Duration `option:"mandatory" validate:"required"`
	maxDisplayLen    int           `option:"mandatory" validate:"required"`
}

type Manager struct {
	log              logger
	opts             Options
	root             *list.DeLinkedList[entry]
	syncNotification chan struct{}
	entryIDs         map[int]*list.Node[entry]
	maxID            int
}

type entry struct {
	ID        int                       `yaml:"id"`
	Name      string                    `yaml:"name"`
	Exec      string                    `yaml:"exec"`
	Parent    *list.DeLinkedList[entry] `yaml:"parent"`
	Entries   *list.DeLinkedList[entry] `yaml:"entries"`
	IsDir     bool                      `yaml:"isDir"`
	CreatedAt time.Time                 `yaml:"createdAt"`
	UpdatedAt time.Time                 `yaml:"updatedAt"`
}

func NewManager(ctx context.Context, log logger, opts Options) (*Manager, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("opts.Validate(): %w", err)
	}

	manager := Manager{
		log:              log,
		opts:             opts,
		root:             new(list.DeLinkedList[entry]),
		syncNotification: make(chan struct{}),
		entryIDs:         make(map[int]*list.Node[entry]),
		maxID:            0,
	}

	if !opts.inMemory {
		manager.syncNotification = make(chan struct{}, 1)

		if err := manager.readConfig(); err != nil {
			return nil, fmt.Errorf("readConfig(): %w", err)
		}

		manager.runSyncer(ctx)
	}

	return &manager, nil
}

func (m *Manager) readConfig() error {
	return nil
}

func (m *Manager) runSyncer(ctx context.Context) {
	ticker := time.NewTicker(m.opts.syncConfigPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.syncIn()
		case <-m.syncNotification:
			m.syncOut()
			ticker.Reset(m.opts.syncConfigPeriod)
		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) syncIn() {
}

func (m *Manager) syncOut() {
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

	e, ok := m.entryIDs[id]
	if !ok {
		return nil
	}

	return e
}

func (m *Manager) DisplayEntry(e entry) string {
	if e.Name != "" {
		return e.Name
	}

	if len(e.Exec) <= m.opts.maxDisplayLen {
		return e.Exec
	}

	return e.Exec[:m.opts.maxDisplayLen-3] + "..."
}

func (m *Manager) newEntry(name, exec string, isDir bool, parent *list.DeLinkedList[entry]) entry {
	m.maxID++

	return entry{
		ID:        m.maxID,
		Name:      name,
		Exec:      exec,
		Parent:    parent,
		Entries:   nil,
		IsDir:     isDir,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (m *Manager) registerEntry(e *list.Node[entry]) {
	m.entryIDs[e.Value.ID] = e
}

func (m *Manager) unregisterEntry(id int) {
	delete(m.entryIDs, id)
}
