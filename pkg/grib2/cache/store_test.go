package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkStore_Parallel(b *testing.B) {
	b.Run("LRU", func(b *testing.B) {
		store := NewLRUStore(1000)
		ctx := context.Background()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				store.Set(ctx, 1, 1.0)
				_, _ = store.Get(ctx, 1)
			}
		})
	})

	b.Run("Map", func(b *testing.B) {
		store := NewMapStore()
		ctx := context.Background()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				store.Set(ctx, 1, 1.0)
				_, _ = store.Get(ctx, 1)
			}
		})
	})
}

func TestLRUStore(t *testing.T) {
	store := NewLRUStore(1)
	ctx := context.Background()

	store.Set(ctx, 1, 1.0)
	v, ok := store.Get(ctx, 1)
	assert.True(t, ok)
	assert.Equal(t, float32(1.0), v)

	store.Set(ctx, 2, 2.0)
	v, ok = store.Get(ctx, 2)
	assert.True(t, ok)
	assert.Equal(t, float32(2.0), v)

	v, ok = store.Get(ctx, 3)
	assert.False(t, ok)
	assert.Equal(t, float32(0.0), v)
}
