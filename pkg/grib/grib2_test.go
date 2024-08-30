package grib_test

import (
	"io"
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
	{
		sec, err := g.ReadSection0()
		require.NoError(t, err)

		assert.Equal(t, "GRIB", string(sec.GribLiteral[:]))
		assert.Equal(t, 2, sec.GetEditionNumber())
		assert.Equal(t, 0, sec.GetDiscipline())
		assert.Equal(t, 203278, sec.GetGribLength())
	}

	{
		sec, err := g.ReadSection0()
		require.NoError(t, err)

		assert.Equal(t, "GRIB", string(sec.GribLiteral[:]))
		assert.Equal(t, 2, sec.GetEditionNumber())
		assert.Equal(t, 0, sec.GetDiscipline())
		assert.Equal(t, 203278, sec.GetGribLength())
	}
}

func TestGrib_NextSection(t *testing.T) {
	f, err := os.Open("testdata/temp.grib2")
	require.NoError(t, err)

	g := grib.NewGrib2(f)
	{
		sec, err := g.ReadSection0()
		require.NoError(t, err)

		assert.Equal(t, "GRIB", string(sec.GribLiteral[:]))
		assert.Equal(t, 2, sec.GetEditionNumber())
		assert.Equal(t, 0, sec.GetDiscipline())
		assert.Equal(t, 203278, sec.GetGribLength())
	}

	{
		sec, err := g.NextSection()
		require.NoError(t, err)

		assert.IsType(t, &grib.Section1{}, sec)
		require.Equal(t, 1, sec.SectionNumber())
		require.Equal(t, 21, sec.SectionLength())

		assert.Equal(t, "2023-07-11T00:00:00Z", sec.(*grib.Section1).GetTime(time.UTC).Format(time.RFC3339))
	}

	{
		sec, err := g.NextSection()
		require.NoError(t, err)

		assert.IsType(t, &grib.Section3{}, sec)
		require.Equal(t, 3, sec.SectionNumber())
		require.Equal(t, 72, sec.SectionLength())
	}

	{
		sec, err := g.NextSection()
		require.NoError(t, err)

		assert.IsType(t, &grib.Section4{}, sec)
		require.Equal(t, 4, sec.SectionNumber())
		require.Equal(t, 34, sec.SectionLength())
	}

	{
		sec, err := g.NextSection()
		require.NoError(t, err)

		assert.IsType(t, &grib.Section5{}, sec)
		require.Equal(t, 5, sec.SectionNumber())
		require.Equal(t, 21, sec.SectionLength())
	}

	{
		sec, err := g.NextSection()
		require.NoError(t, err)

		assert.IsType(t, &grib.Section6{}, sec)
		require.Equal(t, 6, sec.SectionNumber())
		require.Equal(t, 6, sec.SectionLength())
	}

	{
		sec, err := g.NextSection()
		require.NoError(t, err)

		assert.IsType(t, &grib.Section7{}, sec)
		require.Equal(t, 7, sec.SectionNumber())
		require.Equal(t, 203104, sec.SectionLength())
	}

	{
		sec, err := g.NextSection()
		require.NoError(t, err)

		assert.IsType(t, &grib.Section8{}, sec)
		require.Equal(t, 8, sec.SectionNumber())
		require.Equal(t, 4, sec.SectionLength())
	}

	{
		sec, err := g.NextSection()
		require.ErrorIs(t, err, io.EOF)
		require.Nil(t, sec)
	}
}
