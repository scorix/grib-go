package regulation_test

import (
	"testing"

	"github.com/scorix/grib-go/pkg/grib2/regulation"
	"github.com/stretchr/testify/assert"
)

func TestToInt32(t *testing.T) {
	assert.Equal(t, int32(-90000000), regulation.ToInt32(2237483648))
	assert.Equal(t, int32(359750000), regulation.ToInt32(359750000))
	assert.Equal(t, int8(-1), regulation.ToInt8(255))
}
