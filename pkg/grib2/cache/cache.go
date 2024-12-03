package cache

import "context"

type GridDataSource interface {
	ReadGridAt(ctx context.Context, grid int) (float32, error)
}

type GridCache interface {
	ReadGridAt(ctx context.Context, grid int, lat, lon float32) (float32, error)
}

type noCache struct {
	datasource GridDataSource
}

func NewNoCache(datasource GridDataSource) GridCache {
	return &noCache{datasource: datasource}
}

func (n *noCache) ReadGridAt(ctx context.Context, grid int, lat, lon float32) (float32, error) {
	return n.datasource.ReadGridAt(ctx, grid)
}
