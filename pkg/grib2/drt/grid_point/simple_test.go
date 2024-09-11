package gridpoint_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/icza/bitio"
	"github.com/scorix/grib-go/pkg/grib2/drt/definition"
	gridpoint "github.com/scorix/grib-go/pkg/grib2/drt/grid_point"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleScale(t *testing.T) {
	var v uint16 = 0b_10100101_1110

	assert.Equal(t, uint64(0xa5e), uint64(v))
	assert.Equal(t, float64(1.3113e-320), math.Float64frombits(uint64(v)))

	sp := gridpoint.NewSimplePacking(definition.SimplePacking{}, 1)
	sp.R = 0.0194875
	sp.E = -18
	sp.D = -4
	sp.Bits = 12
	t.Logf("simple packing: %+v", sp)

	b1, b2 := uint8(v>>4&0xff), uint8((v&0x0f)<<4)

	values, err := sp.ReadAllData(bitio.NewCountReader(bytes.NewReader([]byte{b1, b2})))
	require.NoError(t, err)
	assert.Equal(t, float32(2.9611706734e+02), float32(values[0]))
}

func TestSimpleScaleReader(t *testing.T) {
	bs := []byte{0xff, 0xff, 0xff, 0b1010_0101, 0b1110_0000}
	bf := bytes.NewReader(bs)

	sp := gridpoint.NewSimplePacking(definition.SimplePacking{}, 3)
	sp.R = 0.0194875
	sp.E = -18
	sp.D = -4
	sp.Bits = 12
	t.Logf("simple packing: %+v", sp)

	r := gridpoint.NewSimplePackingReader(bf, 0, 5, sp)
	f, err := r.ReadGridAt(2)
	require.NoError(t, err)
	assert.InDelta(t, float32(2.9611706734e+02), f, 1e-5)
}
