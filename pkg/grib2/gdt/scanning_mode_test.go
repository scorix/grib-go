package gdt_test

import (
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

var tpl_scanmode40 = gdt.Template40FixedPart{
	Ni:                        3072,
	Nj:                        1536,
	LatitudeOfFirstGridPoint:  89910324,
	LongitudeOfFirstGridPoint: 0,
	LatitudeOfLastGridPoint:   -89910324,
	LongitudeOfLastGridPoint:  359882813,
	IDirectionIncrement:       117188,
	N:                         768,
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

func TestScanningMode0000(t *testing.T) {
	tests := []struct {
		name      string
		tpl       gdt.Template
		i, j      int
		wantLat   float32
		wantLon   float32
		wantGrid  int
		approxLat float32
		approxLon float32
		delta     float64
	}{
		{
			name:     "First point",
			tpl:      &tpl_scanmode0,
			i:        0,
			j:        0,
			wantLat:  90,
			wantLon:  0,
			wantGrid: 0,
		},
		{
			name:     "Last point",
			tpl:      &tpl_scanmode0,
			i:        720,
			j:        1439,
			wantLat:  -90,
			wantLon:  359.75,
			wantGrid: 1038239,
		},
		{
			name:     "Middle point",
			tpl:      &tpl_scanmode0,
			i:        360,
			j:        720,
			wantLat:  0,
			wantLon:  180,
			wantGrid: 519120,
		},
		{
			name:      "Approximate point",
			tpl:       &tpl_scanmode0,
			i:         236,
			j:         488,
			wantLat:   31,
			wantLon:   122,
			wantGrid:  340328,
			approxLat: 31.1,
			approxLon: 122.01,
		},
		{
			name:     "regular_gg 1",
			tpl:      &tpl_scanmode40,
			i:        427,
			j:        993,
			wantLat:  39.9042,
			wantLon:  116.4074,
			wantGrid: 1312737,
			delta:    0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := tt.tpl

			sm, err := template.GetScanningMode()
			require.NoError(t, err)

			i, j, gotGrid := sm.GetGridPointFromLL(tt.wantLat, tt.wantLon)
			assert.Equal(t, tt.i, i, "i mismatch")
			assert.Equal(t, tt.j, j, "j mismatch")
			assert.Equal(t, tt.wantGrid, gotGrid, "grid point mismatch")

			gotLat, gotLon := sm.GetGridPointLL(tt.i, tt.j)
			assert.InDelta(t, tt.wantLat, gotLat, tt.delta, "latitude mismatch")
			assert.InDelta(t, tt.wantLon, gotLon, tt.delta, "longitude mismatch")

			if tt.approxLat != 0 && tt.approxLon != 0 {
				i, j, gotGrid = sm.GetGridPointFromLL(tt.approxLat, tt.approxLon)
				assert.Equal(t, tt.i, i, "i mismatch")
				assert.Equal(t, tt.j, j, "j mismatch")
				assert.Equal(t, tt.wantGrid, gotGrid, "approximate point mismatch")
			}
		})
	}
}
