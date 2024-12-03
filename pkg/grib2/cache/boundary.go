package cache

import (
	"fmt"

	"golang.org/x/sync/singleflight"
)

type boundary struct {
	minLat float32
	maxLat float32
	minLon float32
	maxLon float32

	cache      map[int]float32
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
		cache:      make(map[int]float32),
	}
}

func (b *boundary) ReadGridAt(grid int, lat, lon float32) (float32, error) {
	if lat < b.minLat || lat > b.maxLat || lon < b.minLon || lon > b.maxLon {
		return b.datasource.ReadGridAt(grid)
	}

	v, err, _ := b.sfg.Do(fmt.Sprintf("%d", grid), func() (interface{}, error) {
		vFromCache, ok := b.cache[grid]
		if !ok {
			vFromSource, err := b.datasource.ReadGridAt(grid)
			if err != nil {
				return 0, err
			}

			b.cache[grid] = vFromSource
			vFromCache = vFromSource
		}

		return vFromCache, nil
	})

	return v.(float32), err
}
