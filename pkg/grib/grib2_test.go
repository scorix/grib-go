package grib_test

import (
	"os"
	"testing"

	"github.com/scorix/grib-go/pkg/grib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrib_ReadSection0(t *testing.T) {
	f, err := os.Open("testdata/temp.grib2")
	require.NoError(t, err)

	g := grib.New(f)
	sec, err := g.ReadSection0()
	require.NoError(t, err)

	assert.Equal(t, "GRIB", string(sec.GribLiteral[:]))
	assert.Equal(t, 2, sec.GetEditionNumber())
	assert.Equal(t, 0, sec.GetDiscipline())
	assert.Equal(t, 203278, sec.GetGribLength())
}
