package cache

import (
	"context"
	"fmt"

	"golang.org/x/sync/singleflight"
)

type boundary struct {
	minLat float32
	maxLat float32
	minLon float32
	maxLon float32

	cache      Store
	datasource GridDataSource
	sfg        singleflight.Group
}

func NewBoundary(minLat, maxLat, minLon, maxLon float32, datasource GridDataSource) GridCache {
	return &boundary{
		minLat:     minLat,
		maxLat:     maxLat,
		minLon:     minLon,
		maxLon:     maxLon,
		datasource: datasource,
		cache:      NewMapStore(),
	}
}

func (b *boundary) ReadGridAt(ctx context.Context, grid int, lat, lon float32) (float32, error) {
	if lat < b.minLat || lat > b.maxLat || lon < b.minLon || lon > b.maxLon {
		return b.datasource.ReadGridAt(ctx, grid)
	}

	v, err, _ := b.sfg.Do(fmt.Sprintf("%d", grid), func() (interface{}, error) {
		vFromCache, ok := b.cache.Get(ctx, fmt.Sprintf("%d", grid))
		if !ok {
			vFromSource, err := b.datasource.ReadGridAt(ctx, grid)
			if err != nil {
				return 0, err
			}

			b.cache.Set(ctx, fmt.Sprintf("%d", grid), vFromSource)
			vFromCache = vFromSource
		}

		return vFromCache, nil
	})

	return v.(float32), err
}
