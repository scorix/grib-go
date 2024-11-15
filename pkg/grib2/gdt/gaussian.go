package gdt

import (
	"math"
)

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
	// Process input coordinates
	_, intLon := adjustCoordinates(lat, lon, lon0)

	// Calculate longitude index
	var jLon int32
	if lon1 > lon0 {
		jLon = (intLon - lon0) / ((lon1 - lon0) / ni)
	} else {
		jLon = (lon0 - intLon) / ((lon0 - lon1) / ni)
	}

	// Convert latitude to radians
	latRad := float64(lat) * math.Pi / 180.0

	// Use an improved nonlinear mapping
	normalizedLat := math.Sin(latRad)
	totalLats := 2 * n

	// Adjust coefficients to better match the distribution of the Gaussian grid
	// Use the asin(sin) combination to adjust the nonlinear distribution
	adjustedLat := math.Asin(normalizedLat) / (math.Pi / 2.0)
	iLat := int32((1.0 - adjustedLat) * float64(totalLats) / 2.0)

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
