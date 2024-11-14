package gdt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRegularLLGridIndex(t *testing.T) {
	tests := []struct {
		name         string
		lat, lon     float32
		lat0, lon0   int32
		lat1, lon1   int32
		dlat, dlon   int32
		wantI, wantJ int
		wantN        int
		wantGLat     float32
		wantGLon     float32
	}{
		{
			name:     "simple grid",
			lat:      10.0,
			lon:      20.0,
			lat0:     0,
			lon0:     0,
			lat1:     90000000,
			lon1:     90000000,
			dlat:     10000000,
			dlon:     10000000,
			wantI:    1,
			wantJ:    2,
			wantN:    12,
			wantGLat: 10.0,
			wantGLon: 20.0,
		},
		{
			name:     "negative coordinates",
			lat:      -20.0,
			lon:      -30.0,
			lat0:     -90000000,
			lon0:     -180000000,
			lat1:     90000000,
			lon1:     180000000,
			dlat:     10000000,
			dlon:     10000000,
			wantI:    7,
			wantJ:    15,
			wantN:    274,
			wantGLat: -20.0,
			wantGLon: -30.0,
		},
		{
			name:     "global 0.25deg grid",
			lat:      90.0, // First point
			lon:      0.0,
			lat0:     90000000,  // latitudeOfFirstGridPoint (90 degrees)
			lon0:     0,         // longitudeOfFirstGridPoint (0 degrees)
			lat1:     -90000000, // latitudeOfLastGridPoint (-90 degrees)
			lon1:     359750000, // longitudeOfLastGridPoint (359.75 degrees)
			dlat:     250000,    // jDirectionIncrement (0.25 degrees)
			dlon:     250000,    // iDirectionIncrement (0.25 degrees)
			wantI:    0,
			wantJ:    0,
			wantN:    0,
			wantGLat: 90.0,
			wantGLon: 0.0,
		},
		{
			name:     "global 0.25deg grid middle point",
			lat:      0.0,   // Equator
			lon:      180.0, // International Date Line
			lat0:     90000000,
			lon0:     0,
			lat1:     -90000000,
			lon1:     359750000,
			dlat:     250000,
			dlon:     250000,
			wantI:    360,    // (90 - 0) / 0.25 = 360
			wantJ:    720,    // (180 - 0) / 0.25 = 720
			wantN:    519120, // 360 * 1440 + 720
			wantGLat: 0.0,
			wantGLon: 180.0,
		},
		{
			name:     "global 0.25deg grid last point",
			lat:      -90.0,  // South Pole
			lon:      359.75, // Last longitude
			lat0:     90000000,
			lon0:     0,
			lat1:     -90000000,
			lon1:     359750000,
			dlat:     250000,
			dlon:     250000,
			wantI:    720,     // (90 - (-90)) / 0.25 + 1 = 721
			wantJ:    1439,    // (359.75 - 0) / 0.25 = 1439
			wantN:    1038239, // 720 * 1440 + 1439
			wantGLat: -90.0,
			wantGLon: 359.75,
		},
		{
			name:     "Beijing", // 北京
			lat:      39.9042,
			lon:      116.4074,
			lat0:     90000000,
			lon0:     0,
			lat1:     -90000000,
			lon1:     359750000,
			dlat:     250000,
			dlon:     250000,
			wantI:    200,
			wantJ:    466,
			wantN:    288466,
			wantGLat: 40.00,
			wantGLon: 116.50,
		},
		{
			name:     "New York", // 纽约
			lat:      40.7128,
			lon:      -74.0060,
			lat0:     90000000,
			lon0:     0,
			lat1:     -90000000,
			lon1:     359750000,
			dlat:     250000,
			dlon:     250000,
			wantI:    197,
			wantJ:    1144,
			wantN:    284824,
			wantGLat: 40.75,
			wantGLon: 286,
		},
		{
			name:     "London", // 伦敦
			lat:      51.5074,
			lon:      -0.1278,
			lat0:     90000000,
			lon0:     0,
			lat1:     -90000000,
			lon1:     359750000,
			dlat:     250000,
			dlon:     250000,
			wantI:    154,
			wantJ:    1439,
			wantN:    223199,
			wantGLat: 51.5,
			wantGLon: 359.75,
		},
		{
			name:     "Sydney", // 悉尼
			lat:      -33.8688,
			lon:      151.2093,
			lat0:     90000000,
			lon0:     0,
			lat1:     -90000000,
			lon1:     359750000,
			dlat:     250000,
			dlon:     250000,
			wantI:    495,
			wantJ:    605,
			wantN:    713405,
			wantGLat: -33.75,
			wantGLon: 151.25,
		},
		{
			name:     "Tokyo", // 东京
			lat:      35.6762,
			lon:      139.6503,
			lat0:     90000000,
			lon0:     0,
			lat1:     -90000000,
			lon1:     359750000,
			dlat:     250000,
			dlon:     250000,
			wantI:    217,
			wantJ:    559,
			wantN:    313039,
			wantGLat: 35.75,
			wantGLon: 139.75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotI, gotJ, gotN := GetRegularLLGridIndex(
				tt.lat, tt.lon,
				tt.lat0, tt.lon0,
				tt.lat1, tt.lon1,
				tt.dlat, tt.dlon,
			)

			assert.Equal(t, tt.wantI, gotI)
			assert.Equal(t, tt.wantJ, gotJ)
			assert.Equal(t, tt.wantN, gotN)

			gotGLat, gotGLon := GetRegularLLGridPointByIndex(
				gotI, gotJ,
				tt.lat0, tt.lon0,
				tt.lat1, tt.lon1,
				tt.dlat, tt.dlon,
			)
			assert.Equal(t, tt.wantGLat, gotGLat)
			assert.Equal(t, tt.wantGLon, gotGLon)
		})
	}
}

func BenchmarkGetRegularLLGridIndex(b *testing.B) {
	benchmarks := []struct {
		name       string
		lat, lon   float32
		lat0, lon0 int32
		lat1, lon1 int32
		dlat, dlon int32
	}{
		{
			name: "first_point",
			lat:  90.0,
			lon:  0.0,
			lat0: 90000000,
			lon0: 0,
			lat1: -90000000,
			lon1: 359750000,
			dlat: 250000,
			dlon: 250000,
		},
		{
			name: "equator_point",
			lat:  0.0,
			lon:  180.0,
			lat0: 90000000,
			lon0: 0,
			lat1: -90000000,
			lon1: 359750000,
			dlat: 250000,
			dlon: 250000,
		},
		{
			name: "last_point",
			lat:  -90.0,
			lon:  359.75,
			lat0: 90000000,
			lon0: 0,
			lat1: -90000000,
			lon1: 359750000,
			dlat: 250000,
			dlon: 250000,
		},
		{
			name: "random_point",
			lat:  45.5,
			lon:  120.5,
			lat0: 90000000,
			lon0: 0,
			lat1: -90000000,
			lon1: 359750000,
			dlat: 250000,
			dlon: 250000,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var i, j, n int
			b.ResetTimer()
			for k := 0; k < b.N; k++ {
				i, j, n = GetRegularLLGridIndex(
					bm.lat, bm.lon,
					bm.lat0, bm.lon0,
					bm.lat1, bm.lon1,
					bm.dlat, bm.dlon,
				)
			}
			_ = i
			_ = j
			_ = n
		})
	}
}
