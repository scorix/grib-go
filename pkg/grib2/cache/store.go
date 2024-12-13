package cache

import (
	"container/list"
	"context"
	"sync"
)

type Store interface {
	Get(ctx context.Context, key int) (float32, bool)
	Set(ctx context.Context, key int, value float32)
}

// LRU Cache 实现
type lruStore struct {
	mu        sync.RWMutex
	capacity  int
	cache     map[int]*list.Element
	lru       *list.List
	entryPool sync.Pool
}

type entry struct {
	key   int
	value float32
}

func NewLRUStore(capacity int) Store {
	return &lruStore{
		capacity: capacity,
		cache:    make(map[int]*list.Element, capacity),
		lru:      list.New(),
		entryPool: sync.Pool{
			New: func() any {
				return &entry{}
			},
		},
	}
}

func (l *lruStore) Get(ctx context.Context, key int) (float32, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if elem, ok := l.cache[key]; ok {
		l.lru.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}

	return 0, false
}

func (l *lruStore) Set(ctx context.Context, key int, value float32) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if elem, ok := l.cache[key]; ok {
		l.lru.MoveToFront(elem)
		elem.Value.(*entry).value = value
		return
	}

	// 如果缓存已满，删除最久未使用的条目
	if l.lru.Len() >= l.capacity {
		oldest := l.lru.Back()
		if oldest != nil {
			delete(l.cache, oldest.Value.(*entry).key)
			l.lru.Remove(oldest)
			l.entryPool.Put(oldest.Value)
		}
	}

	// 添加新条目
	newEntry := l.entryPool.Get().(*entry)
	newEntry.key = key
	newEntry.value = value
	elem := l.lru.PushFront(newEntry)
	l.cache[key] = elem
}

// 简单的 map 实现
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

	return v, ok
}

func (m *mapStore) Set(ctx context.Context, grid int, value float32) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache[grid] = value
}
