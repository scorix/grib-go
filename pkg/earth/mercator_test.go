package earth_test

import (
	"testing"

	"github.com/scorix/grib-go/pkg/earth"
	"github.com/stretchr/testify/assert"
)

func TestLatLonToMercator(t *testing.T) {
	x, y := earth.LatLonToMercator(38.8, 113.6)
	assert.InDelta(t, 12645894.154, x, 0.01)
	assert.InDelta(t, 4693063.644, y, 0.01)

}

func TestMercatorToLatLon(t *testing.T) {
	lon, lat := earth.MercatorToLatLon(12645894.154, 4693063.644)
	assert.InDelta(t, 113.6, lon, 0.01)
	assert.InDelta(t, 38.8, lat, 0.01)
}
