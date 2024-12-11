package cache_test

import (
	"context"
	"testing"

	"github.com/scorix/grib-go/pkg/grib2/cache"
)

// 实现一个简单的 GridDataSource
type testDataSource struct{}

func (t *testDataSource) ReadGridAt(ctx context.Context, index int) (float32, error) {
	return 1.0, nil
}

func BenchmarkCustom_ReadGridAt_Parallel(b *testing.B) {
	store := cache.NewMapStore()

	c := cache.NewCustom(
		func(lat, lon float32) bool {
			return true
		},
		&testDataSource{},
		store,
	)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := c.ReadGridAt(context.Background(), 1, 1, 1)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
