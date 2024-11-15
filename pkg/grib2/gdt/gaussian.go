package gdt

// GetRegularGGGridIndex calculates the index position of a given latitude and longitude in a Gaussian grid.
// This function uses an approximate algorithm to determine the indices.
// Parameters:
// - lat, lon: The target position's latitude and longitude (in degrees)
// - lat0, lon0: The coordinates of the first grid point (scaled by 1e6)
// - lat1, lon1: The coordinates of the last grid point (scaled by 1e6)
// - n, ni: The number of grid points in the latitude and longitude directions
// Returns:
// - i, j: The grid indices
// - nIdx: The linear index
func GetRegularGGGridIndex(lat, lon float32, lat0, lon0, lat1, lon1 int32, n, ni int32) (i, j, nIdx int) {
	// Handle longitude index
	intLat, intLon := adjustCoordinates(lat, lon, lon0)
	var jLon int32
	if lon1 > lon0 {
		jLon = (intLon - lon0) / ((lon1 - lon0) / ni)
	} else {
		jLon = (lon0 - intLon) / ((lon0 - lon1) / ni)
	}

	// Calculate latitude index using the actual grid boundaries
	totalLats := 2 * n
	latRange := float64(lat0 - lat1)
	latDiff := float64(lat0 - intLat)

	// Linear interpolation using the actual grid boundaries
	iLat := int32((latDiff * float64(totalLats)) / latRange)

	// Ensure indices are within valid ranges
	if iLat < 0 {
		iLat = 0
	}
	if iLat >= totalLats {
		iLat = totalLats - 1
	}
	if jLon < 0 {
		jLon = 0
	}
	if jLon >= ni {
		jLon = ni - 1
	}

	i = int(iLat)
	j = int(jLon)
	nIdx = i*int(ni) + j

	return i, j, nIdx
}

// GetRegularGGGridPointByIndex calculates the actual latitude and longitude of a grid point by its index.
// This function uses an approximate algorithm to determine the coordinates.
// Parameters:
// - i, j: The grid indices
// - lat0, lon0: The coordinates of the first grid point (scaled by 1e6)
// - lat1, lon1: The coordinates of the last grid point (scaled by 1e6)
// - n, ni: The number of grid points in the latitude and longitude directions
// Returns:
// - gLat, gLon: The actual latitude and longitude of the grid point (in degrees)
func GetRegularGGGridPointByIndex(i, j int, lat0, lon0, lat1, lon1 int32, n, ni int32) (gLat, gLon float32) {
	totalLats := 2 * n

	// Modify the calculation of the actual coordinates of the grid point
	// Calculate the actual latitude and longitude of the grid point
	gLat = float32(float64(lat0)/1e6 - float64(i)*float64(lat0-lat1)/float64(totalLats)/1e6)

	// Calculate longitude
	if lon1 > lon0 {
		gLon = float32(float64(lon0)+float64(j)*float64(lon1-lon0)/float64(ni)) / 1e6
	} else {
		gLon = float32(float64(lon0)-float64(j)*float64(lon0-lon1)/float64(ni)) / 1e6
	}

	return gLat, gLon
}
