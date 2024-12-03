package cache

import (
	"context"
	"sync"
)

type Store interface {
	Get(ctx context.Context, key string) (float32, bool)
	Set(ctx context.Context, key string, value float32)
}

type mapStore struct {
	mu    sync.RWMutex
	cache map[string]float32
}

func NewMapStore() Store {
	return &mapStore{cache: make(map[string]float32)}
}

func (m *mapStore) Get(ctx context.Context, key string) (float32, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.cache[key]
	return v, ok
}

func (m *mapStore) Set(ctx context.Context, key string, value float32) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache[key] = value
}
