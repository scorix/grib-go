package gdt_test

import (
	"fmt"
	"testing"

	"github.com/scorix/grib-go/pkg/grib2/gdt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tpl_scanmode0 = gdt.Template0FixedPart{
	Ni:                        1440,
	Nj:                        721,
	LatitudeOfFirstGridPoint:  90000000,
	LongitudeOfFirstGridPoint: 0,
	LatitudeOfLastGridPoint:   -90000000,
	LongitudeOfLastGridPoint:  359750000,
	IDirectionIncrement:       250000,
	JDirectionIncrement:       250000,
}

var tpl_scanmode64 = gdt.Template0FixedPart{
	Ni:                        363,
	Nj:                        373,
	LatitudeOfFirstGridPoint:  33046875,
	LongitudeOfFirstGridPoint: 346007813,
	LatitudeOfLastGridPoint:   67921875,
	LongitudeOfLastGridPoint:  36914063,
	IDirectionIncrement:       140625,
	JDirectionIncrement:       93750,
	ScanningMode:              64,
}

func TestGetGridPointLL(t *testing.T) {
	tests := []struct {
		tpl      gdt.Template0FixedPart
		lat, lon float32
		n        int
	}{
		{tpl: tpl_scanmode0, lat: 90, lon: 0, n: 0},
		{tpl: tpl_scanmode0, lat: 90, lon: 0.25, n: 1},
		{tpl: tpl_scanmode0, lat: 90, lon: 359.75, n: 1439},
		{tpl: tpl_scanmode0, lat: 89.75, lon: 0, n: 1440},
		{tpl: tpl_scanmode0, lat: -90, lon: 359.75, n: 1038239},
		{tpl: tpl_scanmode0, lat: 90, lon: 269.25, n: 1077},

		// {tpl: tpl_135400, lat: 33.046875, lon: 346.007813, n: 0},
		// {tpl: tpl_135400, lat: 33.046875, lon: 345.91406, n: 1},
		// {tpl: tpl_135400, lat: 33.046875, lon: 36.914063, n: 362},
		// {tpl: tpl_135400, lat: 67.921875, lon: 36.914063, n: 135399},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%.2f,%.2f", tt.lat, tt.lon), func(t *testing.T) {
			t.Parallel()

			template := gdt.Template0{
				Template0FixedPart: tt.tpl,
			}

			sm, err := template.GetScanningMode()
			require.NoError(t, err)

			lat, lon := sm.GetGridPointLL(tt.n)
			assert.Equal(t, tt.lat, lat)
			assert.Equal(t, tt.lon, lon)

			n := sm.GetGridPointFromLL(tt.lat, tt.lon)
			assert.Equal(t, tt.n, n)
		})
	}
}
