package regulation_test

import (
	"math"
	"testing"

	"github.com/scorix/grib-go/pkg/grib2/regulation"
	"github.com/stretchr/testify/assert"
)

func TestToInt32(t *testing.T) {
	assert.Equal(t, int32(-90000000), regulation.ToInt32(2237483648))
	assert.Equal(t, int32(359750000), regulation.ToInt32(359750000))

	assert.Equal(t, int32(-1), regulation.ToInt32(math.MaxUint32))
	assert.Equal(t, int32(math.MaxInt32), regulation.ToInt32(math.MaxInt32))
}

func TestToInt16(t *testing.T) {
	assert.Equal(t, int16(-1), regulation.ToInt16(math.MaxUint16))
	assert.Equal(t, int16(math.MaxInt16), regulation.ToInt16(math.MaxInt16))
}

func TestToInt8(t *testing.T) {
	assert.Equal(t, int8(103), regulation.ToInt8(103))

	assert.Equal(t, int8(-1), regulation.ToInt8(math.MaxUint8))
	assert.Equal(t, int8(math.MaxInt8), regulation.ToInt8(math.MaxInt8))
}

func TestToInt(t *testing.T) {
	assert.Equal(t, -64651, regulation.ToInt(0b100000001111110010001011, 24))
}

func TestIsMissingValue(t *testing.T) {
	assert.Equal(t, true, regulation.IsMissingValue(255, 8))
	assert.Equal(t, true, regulation.IsMissingValue(65535, 16))

	i := -1
	assert.Equal(t, true, regulation.IsMissingValue(uint(i), 8))
	assert.Equal(t, true, regulation.IsMissingValue(uint(i), 16))
	assert.Equal(t, true, regulation.IsMissingValue(uint(i), 24))
	assert.Equal(t, true, regulation.IsMissingValue(uint(i), 32))
}
