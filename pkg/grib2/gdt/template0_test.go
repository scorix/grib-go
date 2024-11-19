package gdt_test

import (
	"testing"

	"github.com/scorix/grib-go/pkg/grib2/gdt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplate0_GetGridIndex(t *testing.T) {
	tpldef := gdt.Template0FixedPart{
		LatitudeOfFirstGridPoint:  90000000,
		LongitudeOfFirstGridPoint: 0,
		LatitudeOfLastGridPoint:   -90000000,
		LongitudeOfLastGridPoint:  359000000,
		IDirectionIncrement:       1000000,
		JDirectionIncrement:       1000000,
		ScanningMode:              0,
	}
	tpl := tpldef.AsTemplate()

	assert.Equal(t, 0, tpl.GetGridIndex(90, 0))
	assert.Equal(t, 0, tpl.GetGridIndex(89.999999, 0))
	assert.Equal(t, 0, tpl.GetGridIndex(90, 0.000001))

	i := 0
	for lat := 90; lat > -90; lat-- {
		for lon := 0; lon < 360; lon++ {
			require.Equal(t, i, tpl.GetGridIndex(float32(lat), float32(lon)), "lat: %d, lon: %d", lat, lon)
			i++
		}
	}
}
