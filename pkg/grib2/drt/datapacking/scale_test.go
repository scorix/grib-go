package datapacking_test

import (
	"testing"

	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
	"github.com/stretchr/testify/assert"
)

func TestSimpleScaleFunc(t *testing.T) {
	tests := []struct {
		r  float32
		d  int16
		e  int16
		x2 uint32
		y  float64
	}{
		{r: 0.0194875, e: -18, d: -4, x2: 0b_10100101_1110, y: 296.11706733703613},
	}

	for _, tt := range tests {
		t.Run("simple", func(t *testing.T) {
			t.Log(tt)

			assert.Equal(t, float64(0.0001), datapacking.DecimalScaleFactor(tt.d))
			assert.Equal(t, float64(3.814697265625e-06), datapacking.BinaryScaleFactor(tt.e))

			f := datapacking.SimpleScaleFunc(tt.e, tt.d, tt.r)
			assert.Equal(t, tt.y, f(tt.x2))
		})
	}
}
