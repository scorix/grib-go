package cache

type GridDataSource interface {
	ReadGridAt(grid int) (float32, error)
}

type GridCache interface {
	ReadGridAt(grid int, lat, lon float32) (float32, error)
}

type noCache struct {
	datasource GridDataSource
}

func NewNoCache(datasource GridDataSource) GridCache {
	return &noCache{datasource: datasource}
}

func (n *noCache) ReadGridAt(grid int, lat, lon float32) (float32, error) {
	return n.datasource.ReadGridAt(grid)
}
