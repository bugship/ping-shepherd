package store

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

// Memory is an in-process target store for tests and local runs without Postgres.
type Memory struct {
	mu      sync.Mutex
	targets map[string]Target
}

func NewMemory() *Memory {
	return &Memory{targets: make(map[string]Target)}
}

func (m *Memory) Ping(context.Context) error { return nil }

func (m *Memory) CreateTarget(_ context.Context, in NewTarget) (Target, error) {
	norm, err := NormalizeTarget(in)
	if err != nil {
		return Target{}, err
	}
	t := Target{
		ID:        newID(),
		Name:      norm.Name,
		URL:       norm.URL,
		Enabled:   true,
		CreatedAt: time.Now().UTC(),
	}
	m.mu.Lock()
	m.targets[t.ID] = t
	m.mu.Unlock()
	return t, nil
}

func (m *Memory) ListTargets(_ context.Context) ([]Target, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Target, 0, len(m.targets))
	for _, t := range m.targets {
		out = append(out, t)
	}
	return out, nil
}

func (m *Memory) GetTarget(_ context.Context, id string) (Target, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.targets[id]
	if !ok {
		return Target{}, ErrNotFound
	}
	return t, nil
}

func (m *Memory) DeleteTarget(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.targets[id]; !ok {
		return ErrNotFound
	}
	delete(m.targets, id)
	return nil
}

func newID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return fmt.Sprintf("%x", b[:])
}
