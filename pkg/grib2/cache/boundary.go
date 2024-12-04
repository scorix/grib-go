package cache

func NewBoundary(minLat, maxLat, minLon, maxLon float32, datasource GridDataSource, cacheStore Store) GridCache {
	inCache := func(lat, lon float32) bool {
		return lat >= minLat && lat <= maxLat && lon >= minLon && lon <= maxLon
	}

	return NewCustom(inCache, datasource, cacheStore)
}
