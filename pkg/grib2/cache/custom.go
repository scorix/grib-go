package cache

import (
	"context"
	"strconv"

	"golang.org/x/sync/singleflight"
)

type custom struct {
	inCache func(lat, lon float32) bool

	cache      Store
	datasource GridDataSource
	sfg        singleflight.Group
}

func NewCustom(inCache func(lat, lon float32) bool, datasource GridDataSource, cacheStore Store) GridCache {
	return &custom{inCache: inCache, datasource: datasource, cache: cacheStore}
}

func (c *custom) InCache(lat, lon float32) bool {
	return c.inCache(lat, lon)
}

func (c *custom) ReadGridAt(ctx context.Context, grid int, lat, lon float32) (float32, error) {
	if !c.InCache(lat, lon) {
		return c.datasource.ReadGridAt(ctx, grid)
	}

	v, err, _ := c.sfg.Do(strconv.Itoa(grid), func() (interface{}, error) {
		vFromCache, ok := c.cache.Get(ctx, grid)
		if !ok {
			vFromSource, err := c.datasource.ReadGridAt(ctx, grid)
			if err != nil {
				return 0, err
			}

			c.cache.Set(ctx, grid, vFromSource)
			vFromCache = vFromSource
		}

		return vFromCache, nil
	})

	return v.(float32), err
}
