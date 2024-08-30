package gridpoint_test

import (
	"bytes"
	"math"
	"testing"

	gridpoint "github.com/scorix/grib-go/pkg/grib2/drt/grid_point"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleScale(t *testing.T) {
	var v uint16 = 0b_10100101_1110

	assert.Equal(t, uint64(0xa5e), uint64(v))
	assert.Equal(t, float64(1.3113e-320), math.Float64frombits(uint64(v)))

	sp := gridpoint.NewSimplePacking(gridpoint.DefSimplePacking{
		R:    0.0194875,
		E:    32786,
		D:    32772,
		Bits: 12,
		Type: 0,
	})

	b1, b2 := uint8(v>>4&0xff), uint8((v&0x0f)<<4)
	t.Logf("simple packing: %+v", sp.SimplePackingReader.DefSimplePacking)

	values, err := sp.ReadData(bytes.NewReader([]byte{b1, b2}))
	require.NoError(t, err)
	assert.Equal(t, float32(2.9611706734e+02), float32(values[0]))
}
