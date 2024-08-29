package grib_test

import (
	"os"
	"testing"
	"time"

	"github.com/scorix/grib-go/pkg/grib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrib_ReadSection0(t *testing.T) {
	f, err := os.Open("testdata/temp.grib2")
	require.NoError(t, err)

	g := grib.NewGrib2(f)
	sec, err := g.ReadSection0()
	require.NoError(t, err)

	assert.Equal(t, "GRIB", string(sec.GribLiteral[:]))
	assert.Equal(t, 2, sec.GetEditionNumber())
	assert.Equal(t, 0, sec.GetDiscipline())
	assert.Equal(t, 203278, sec.GetGribLength())
}

func TestGrib_ReadSection1(t *testing.T) {
	f, err := os.Open("testdata/temp.grib2")
	require.NoError(t, err)

	g := grib.NewGrib2(f)
	sec, err := g.ReadSection1()
	require.NoError(t, err)

	assert.Equal(t, 1, sec.GetSectionNumber())
	assert.Equal(t, "2023-07-11T00:00:00Z", sec.GetTime().Format(time.RFC3339))
	assert.Equal(t, grib.ReferenceTimeStartOfForecast, sec.SignificanceOfReferenceTime)
	assert.Equal(t, 29, sec.GetMasterTableVersion())
	assert.Equal(t, 1, sec.GetLocalTableVersion())
}
