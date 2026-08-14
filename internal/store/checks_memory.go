package store

import (
	"context"
	"sort"
	"time"
)

func (m *Memory) RecordCheck(_ context.Context, targetID string, c Check) (Check, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.targets[targetID]; !ok {
		return Check{}, ErrNotFound
	}
	if c.ID == "" {
		c.ID = newID()
	}
	c.TargetID = targetID
	if c.CheckedAt.IsZero() {
		c.CheckedAt = time.Now().UTC()
	}
	m.checks[targetID] = append(m.checks[targetID], c)
	return c, nil
}

func (m *Memory) LatestCheck(_ context.Context, targetID string) (Check, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.targets[targetID]; !ok {
		return Check{}, ErrNotFound
	}
	list := m.checks[targetID]
	if len(list) == 0 {
		return Check{}, ErrNotFound
	}
	return list[len(list)-1], nil
}

func (m *Memory) ListChecks(_ context.Context, targetID string, limit int) ([]Check, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.targets[targetID]; !ok {
		return nil, ErrNotFound
	}
	list := append([]Check(nil), m.checks[targetID]...)
	sort.Slice(list, func(i, j int) bool {
		return list[i].CheckedAt.After(list[j].CheckedAt)
	})
	if limit > 0 && len(list) > limit {
		list = list[:limit]
	}
	if list == nil {
		list = []Check{}
	}
	return list, nil
}
