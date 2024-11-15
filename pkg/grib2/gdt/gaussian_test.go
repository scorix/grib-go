package gdt_test

import (
	"testing"

	"github.com/scorix/grib-go/pkg/grib2/gdt"
	"github.com/stretchr/testify/assert"
)

func TestGetRegularGGGridIndex(t *testing.T) {
	tests := []struct {
		name         string
		lat, lon     float32
		lat0, lon0   int32
		lat1, lon1   int32
		n            int32
		ni           int32
		wantI, wantJ int
		wantN        int
		wantGLat     float32
		wantGLon     float32
	}{
		{
			name:     "N768 grid first point",
			lat:      89.910324,
			lon:      0.0,
			lat0:     89910324,
			lon0:     0,
			lat1:     -89910324,
			lon1:     359882813,
			n:        768,
			ni:       3072,
			wantI:    0,
			wantJ:    0,
			wantN:    0,
			wantGLat: 89.910324,
			wantGLon: 0.0,
		},
		{
			name:     "N768 grid equator point",
			lat:      0.0,
			lon:      180.0,
			lat0:     89910324,
			lon0:     0,
			lat1:     -89910324,
			lon1:     359882813,
			n:        768,
			ni:       3072,
			wantI:    768,
			wantJ:    1536,
			wantN:    2360832,
			wantGLat: -0.06,
			wantGLon: 180.0,
		},
		{
			name:     "N768 grid last point",
			lat:      -89.910324,
			lon:      359.882813,
			lat0:     89910324,
			lon0:     0,
			lat1:     -89910324,
			lon1:     359882813,
			n:        768,
			ni:       3072,
			wantI:    1535,
			wantJ:    3071,
			wantN:    4718591,
			wantGLat: -89.91,
			wantGLon: 359.882813,
		},
		{
			name:     "Beijing",
			lat:      39.9042,
			lon:      116.4074,
			lat0:     89910324,
			lon0:     0,
			lat1:     -89910324,
			lon1:     359882813,
			n:        768,
			ni:       3072,
			wantI:    427,
			wantJ:    993,
			wantN:    1312737,
			wantGLat: 39.89,
			wantGLon: 116.37,
		},
		// {
		// 	name:     "New York 4.78km",
		// 	lat:      40.7128,
		// 	lon:      -74.0060,
		// 	lat0:     89910324,
		// 	lon0:     0,
		// 	lat1:     -89910324,
		// 	lon1:     359882813,
		// 	n:        768,
		// 	ni:       3072,
		// 	wantI:    384,
		// 	wantJ:    2447,
		// 	wantN:    1292680,
		// 	wantGLat: 40.71,
		// 	wantGLon: 285.94,
		// },
		{
			name:     "New York 5.13km",
			lat:      40.7128,
			lon:      -74.0060,
			lat0:     89910324,
			lon0:     0,
			lat1:     -89910324,
			lon1:     359882813,
			n:        768,
			ni:       3072,
			wantI:    420,
			wantJ:    2441,
			wantN:    1292681,
			wantGLat: 40.71,
			wantGLon: 286.05,
		},
		{
			name:     "London",
			lat:      51.5074,
			lon:      -0.1278,
			lat0:     89910324,
			lon0:     0,
			lat1:     -89910324,
			lon1:     359882813,
			n:        768,
			ni:       3072,
			wantI:    328,
			wantJ:    3071,
			wantN:    1010687,
			wantGLat: 51.49,
			wantGLon: 359.88,
		},
		{
			name:     "Sydney",
			lat:      -33.8688,
			lon:      151.2093,
			lat0:     89910324,
			lon0:     0,
			lat1:     -89910324,
			lon1:     359882813,
			n:        768,
			ni:       3072,
			wantI:    1057,
			wantJ:    1290,
			wantN:    3248394,
			wantGLat: -33.91,
			wantGLon: 151.17,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotI, gotJ, gotN := gdt.GetRegularGGGridIndex(
				tt.lat, tt.lon,
				tt.lat0, tt.lon0,
				tt.lat1, tt.lon1,
				tt.n, tt.ni,
			)

			assert.Equal(t, tt.wantI, gotI)
			assert.Equal(t, tt.wantJ, gotJ)
			assert.Equal(t, tt.wantN, gotN)

			gotGLat, gotGLon := gdt.GetRegularGGGridPointByIndex(
				gotI, gotJ,
				tt.lat0, tt.lon0,
				tt.lat1, tt.lon1,
				tt.n, tt.ni,
			)
			assert.InDelta(t, tt.wantGLat, gotGLat, 1.2e-1)
			assert.InDelta(t, tt.wantGLon, gotGLon, 1.2e-1)
		})
	}
}

func BenchmarkGetRegularGGGridIndex(b *testing.B) {
	benchmarks := []struct {
		name       string
		lat, lon   float32
		lat0, lon0 int32
		lat1, lon1 int32
		n          int32
		ni         int32
	}{
		{
			name: "first_point",
			lat:  89.910324,
			lon:  0.0,
			lat0: 89910324,
			lon0: 0,
			lat1: -89910324,
			lon1: 359882813,
			n:    768,
			ni:   3072,
		},
		{
			name: "equator_point",
			lat:  0.0,
			lon:  180.0,
			lat0: 89910324,
			lon0: 0,
			lat1: -89910324,
			lon1: 359882813,
			n:    768,
			ni:   3072,
		},
		{
			name: "last_point",
			lat:  -89.910324,
			lon:  359.882813,
			lat0: 89910324,
			lon0: 0,
			lat1: -89910324,
			lon1: 359882813,
			n:    768,
			ni:   3072,
		},
		{
			name: "random_point",
			lat:  45.5,
			lon:  120.5,
			lat0: 89910324,
			lon0: 0,
			lat1: -89910324,
			lon1: 359882813,
			n:    768,
			ni:   3072,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var i, j, n int
			b.ResetTimer()
			for k := 0; k < b.N; k++ {
				i, j, n = gdt.GetRegularGGGridIndex(
					bm.lat, bm.lon,
					bm.lat0, bm.lon0,
					bm.lat1, bm.lon1,
					bm.n, bm.ni,
				)
			}
			_ = i
			_ = j
			_ = n
		})
	}
}
