package gdt

import "cmp"

// GetRegularLLGridIndex calculates the index position of a given latitude and longitude in a regular latitude-longitude grid.
// Parameters:
// - lat, lon: The target position's latitude and longitude (in degrees)
// - lat0, lon0: The coordinates of the first grid point (scaled by 1e6)
// - lat1, lon1: The coordinates of the last grid point (scaled by 1e6)
// - dlat, dlon: The increments of latitude and longitude (scaled by 1e6)
// Returns:
// - i, j, n: The grid indices and the linear index
func GetRegularLLGridIndex(lat, lon float32, lat0, lon0, lat1, lon1, dlat, dlon int32) (i, j, n int) {
	// Process input coordinates
	intLat, intLon := adjustCoordinates(lat, lon, lon0)

	// Calculate grid indices
	iLat := calculateLatIndex(intLat, lat0, lat1, dlat)
	jLon, ni := calculateLonIndex(intLon, lon0, lon1, dlon)

	// Convert to integer indices
	i, j = int(iLat), int(jLon)
	n = i*int(ni) + j

	return i, j, n
}

// GetRegularLLGridPointByIndex calculates the actual latitude and longitude of a grid point by its index.
// Parameters:
// - i, j: The grid indices
// - lat0, lon0: The coordinates of the first grid point (scaled by 1e6)
// - lat1, lon1: The coordinates of the last grid point (scaled by 1e6)
// - dlat, dlon: The increments of latitude and longitude (scaled by 1e6)
// Returns:
// - gLat, gLon: The actual latitude and longitude of the grid point (in degrees)
func GetRegularLLGridPointByIndex(i, j int, lat0, lon0, lat1, lon1, dlat, dlon int32) (gLat, gLon float32) {
	// Determine the scan direction
	di, dj := cmp.Compare(lat1, lat0), cmp.Compare(lon1, lon0)

	// Calculate the actual coordinates of the grid point
	gLat, gLon = calculateGridCoordinates(int32(i), int32(j), lat0, lon0, dlat, dlon, di, dj)

	return gLat, gLon
}

// adjustCoordinates adjusts the input coordinates
func adjustCoordinates(lat, lon float32, lon0 int32) (intLat, intLon int32) {
	intLat, intLon = int32(lat*1e6), int32(lon*1e6)
	if intLon < 0 && lon0 >= 0 {
		intLon += 360 * 1e6
	}
	return intLat, intLon
}

// calculateLatIndex calculates the latitude index
func calculateLatIndex(intLat, lat0, lat1, dlat int32) int32 {
	if lat1 > lat0 {
		// Scanning from south to north
		return ((intLat - lat0) + (dlat / 2)) / dlat
	}
	// Scanning from north to south
	return ((lat0 - intLat) + (dlat / 2)) / dlat
}

// calculateLonIndex calculates the longitude index and the number of grid points per row
func calculateLonIndex(intLon, lon0, lon1, dlon int32) (index, ni int32) {
	if lon1 > lon0 {
		// Scanning from west to east
		index = ((intLon - lon0) + (dlon / 2)) / dlon
		ni = (lon1-lon0)/dlon + 1
	} else {
		// Scanning from east to west
		index = ((lon0 - intLon) + (dlon / 2)) / dlon
		ni = (lon0-lon1)/dlon + 1
	}

	return index, ni
}

// calculateGridCoordinates calculates the actual coordinates of the grid point
func calculateGridCoordinates(iLat, jLon, lat0, lon0, dlat, dlon int32, di, dj int) (gLat, gLon float32) {
	gLat = float32(float64(lat0+iLat*dlat*int32(di)) / 1e6)
	gLon = float32(float64(lon0+jLon*dlon*int32(dj)) / 1e6)
	return gLat, gLon
}
