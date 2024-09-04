package gdt_test

import (
	"fmt"
	"testing"

	"github.com/scorix/grib-go/pkg/grib2/gdt"
	"github.com/stretchr/testify/assert"
)

func TestGetGridPoint(t *testing.T) {
	template := gdt.Template0{
		Template0FixedPart: gdt.Template0FixedPart{
			Ni:                        1440,
			Nj:                        721,
			LatitudeOfFirstGridPoint:  90000000,
			LongitudeOfFirstGridPoint: 0,
			LatitudeOfLastGridPoint:   -90000000,
			LongitudeOfLastGridPoint:  359750000,
			IDirectionIncrement:       250000,
			JDirectionIncrement:       250000,
		},
	}

	tests := []struct {
		lat, lon float32
		n        int
	}{
		{lat: 90, lon: 0, n: 0},
		{lat: 90, lon: 0.25, n: 1},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%.2f,%.2f", tt.lat, tt.lon), func(t *testing.T) {
			lat, lon := template.GetGridPoint(tt.n)
			assert.Equal(t, tt.lat, lat)
			assert.Equal(t, tt.lon, lon)
		})
	}
}
