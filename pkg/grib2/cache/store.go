package cache

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Store interface {
	Get(ctx context.Context, key int) (float32, bool)
	Set(ctx context.Context, key int, value float32)
}

type mapStore struct {
	mu    sync.RWMutex
	cache map[int]float32
}

func NewMapStore() Store {
	return &mapStore{cache: make(map[int]float32)}
}

func (m *mapStore) Get(ctx context.Context, grid int) (float32, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.cache[grid]

	sp := trace.SpanFromContext(ctx)
	sp.SetAttributes(attribute.Int("cache.grid", grid), attribute.Bool("cache.hit", ok))

	return v, ok
}

func (m *mapStore) Set(ctx context.Context, grid int, value float32) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache[grid] = value
}
