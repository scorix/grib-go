package datapacking_test

import (
	"testing"

	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
	"github.com/stretchr/testify/assert"
)

func TestSimpleScaleFunc(t *testing.T) {
	tests := []struct {
		r   float32
		d   int16
		e   int16
		x2  uint32
		y   float64
		dsf float64
		bsf float64
	}{
		{r: 0.0194875, e: -18, d: -4, x2: 2654, y: 296.11706733703613, dsf: 0.0001, bsf: 3.814697265625e-06},
	}

	for _, tt := range tests {
		t.Run("simple", func(t *testing.T) {
			t.Log(tt)

			assert.Equal(t, tt.dsf, datapacking.DecimalScaleFactor(tt.d))
			assert.Equal(t, tt.bsf, datapacking.BinaryScaleFactor(tt.e))

			f := datapacking.SimpleScaleFunc(tt.e, tt.d, tt.r)
			assert.Equal(t, float32(tt.y), float32(f(tt.x2)))
		})
	}
}
