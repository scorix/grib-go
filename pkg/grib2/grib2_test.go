package grib_test

import (
	"io"
	"os"
	"testing"
	"time"

	grib "github.com/scorix/grib-go/pkg/grib2"
	"github.com/scorix/grib-go/pkg/grib2/drt"
	gridpoint "github.com/scorix/grib-go/pkg/grib2/drt/grid_point"
	"github.com/scorix/grib-go/pkg/grib2/gdt"
	"github.com/scorix/grib-go/pkg/grib2/pdt"
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

		s := sec.(*grib.Section1)
		assert.Equal(t, "2023-07-11T00:00:00Z", s.GetTime(time.UTC).Format(time.RFC3339))
		assert.Equal(t, uint16(74), s.Center)
		assert.Equal(t, uint16(5), s.SubCenter)
		assert.Equal(t, uint8(29), s.TableVersion)
		assert.Equal(t, uint8(1), s.LocalTableVersion)
		assert.Equal(t, grib.ReferenceTime(1), s.SignificanceOfReferenceTime)
		assert.Equal(t, uint8(0), s.ProductionStatusOfProcessedData)
		assert.Equal(t, uint8(1), s.TypeOfProcessedData)
	}

	{
		sec, err := g.NextSection()
		require.NoError(t, err)

		assert.IsType(t, &grib.Section3{}, sec)
		require.Equal(t, 3, sec.SectionNumber())
		require.Equal(t, 72, sec.SectionLength())

		s := sec.(*grib.Section3)
		assert.Equal(t, uint16(0), s.GridDefinitionTemplateNumber)

		require.IsType(t, &gdt.Template0{}, s.Template)
		tpl := s.Template.(*gdt.Template0)
		assert.Equal(t, uint8(6), tpl.ShapeOfTheEarth)
		assert.Equal(t, uint8(0xff), tpl.ScaleFactorOfRadiusOfSphericalEarth)
		assert.Equal(t, uint32(0xffffffff), tpl.ScaledValueOfRadiusOfSphericalEarth)
		assert.Equal(t, uint8(0xff), tpl.ScaleFactorOfEarthMajorAxis)
		assert.Equal(t, uint32(0xffffffff), tpl.ScaledValueOfEarthMajorAxis)
		assert.Equal(t, uint8(0xff), tpl.ScaleFactorOfEarthMinorAxis)
		assert.Equal(t, uint32(0xffffffff), tpl.ScaledValueOfEarthMinorAxis)
		assert.Equal(t, uint32(363), tpl.Ni)
		assert.Equal(t, uint32(373), tpl.Nj)
		assert.Equal(t, uint32(0), tpl.BasicAngleOfTheInitialProductionDomain)
		assert.Equal(t, uint32(0xffffffff), tpl.SubdivisionsOfBasicAngle)
		assert.Equal(t, uint32(33046875), tpl.LatitudeOfFirstGridPoint)
		assert.Equal(t, uint32(346007813), tpl.LongitudeOfFirstGridPoint)
		assert.Equal(t, uint8(48), tpl.ResolutionAndComponentFlags)
		assert.Equal(t, uint32(67921875), tpl.LatitudeOfLastGridPoint)
		assert.Equal(t, uint32(36914063), tpl.LongitudeOfLastGridPoint)
		assert.Equal(t, uint32(140625), tpl.IDirectionIncrement)
		assert.Equal(t, uint32(93750), tpl.JDirectionIncrement)
		assert.Equal(t, uint8(64), tpl.ScanningMode)
	}

	{
		sec, err := g.NextSection()
		require.NoError(t, err)

		assert.IsType(t, &grib.Section4{}, sec)
		require.Equal(t, 4, sec.SectionNumber())
		require.Equal(t, 34, sec.SectionLength())

		s := sec.(*grib.Section4)
		assert.Equal(t, uint16(0), s.ProductDefinitionTemplateNumber)

		require.IsType(t, &pdt.Template0{}, s.Template)
		tpl := s.Template.(*pdt.Template0)
		assert.Equal(t, uint8(0), tpl.ParameterCategory)
		assert.Equal(t, uint8(0), tpl.ParameterNumber)
		assert.Equal(t, uint8(2), tpl.TypeOfGeneratingProcess)
		assert.Equal(t, uint8(255), tpl.BackgroundProcess)
		assert.Equal(t, uint8(255), tpl.GeneratingProcessIdentifier)
		assert.Equal(t, uint16(65535), tpl.HoursAfterDataCutoff)
		assert.Equal(t, uint8(255), tpl.MinutesAfterDataCutoff)
		assert.Equal(t, uint8(1), tpl.IndicatorOfUnitOfTimeRange)
		assert.Equal(t, uint32(0), tpl.ForecastTime)
		assert.Equal(t, uint8(1), tpl.TypeOfFirstFixedSurface)
		assert.Equal(t, uint8(255), tpl.ScaleFactorOfFirstFixedSurface)
		assert.Equal(t, uint32(0xffffffff), tpl.ScaledValueOfFirstFixedSurface)
		assert.Equal(t, uint8(255), tpl.TypeOfSecondFixedSurface)
		assert.Equal(t, uint8(255), tpl.ScaleFactorOfFirstFixedSurface)
		assert.Equal(t, uint32(0xffffffff), tpl.ScaledValueOfFirstFixedSurface)
	}

	var tpl drt.Template
	var numVals int
	{
		sec, err := g.NextSection()
		require.NoError(t, err)

		assert.IsType(t, &grib.Section5{}, sec)
		require.Equal(t, 5, sec.SectionNumber())
		require.Equal(t, 21, sec.SectionLength())

		s := sec.(*grib.Section5)
		assert.IsType(t, &gridpoint.SimplePacking{}, s.DataRepresentationTemplate)
		assert.IsType(t, &gridpoint.SimplePackingReader{}, s.DataPackingReader)
		numVals = int(s.NumberOfValues)

		tpl = s.DataRepresentationTemplate
		t.Logf("data representation template: %+v", tpl.(*gridpoint.SimplePacking).DefSimplePacking)

		assert.Equal(t, float32(0.0194875), tpl.(*gridpoint.SimplePacking).R)
		assert.Equal(t, int16(-18), tpl.(*gridpoint.SimplePacking).E.Int16())
		assert.Equal(t, int16(-4), tpl.(*gridpoint.SimplePacking).D.Int16())
		assert.Equal(t, uint8(12), tpl.(*gridpoint.SimplePacking).Bits)
		assert.Equal(t, uint8(0), tpl.(*gridpoint.SimplePacking).Type)
	}

	{
		sec, err := g.NextSection()
		require.NoError(t, err)

		assert.IsType(t, &grib.Section6{}, sec)
		require.Equal(t, 6, sec.SectionNumber())
		require.Equal(t, 6, sec.SectionLength())

		s := sec.(*grib.Section6)
		assert.Equal(t, uint8(255), s.BitMapIndicator)
	}

	{
		sec, err := g.NextSection()
		require.NoError(t, err)

		assert.IsType(t, &grib.Section7{}, sec)
		require.Equal(t, 7, sec.SectionNumber())
		require.Equal(t, 203104, sec.SectionLength())

		data := sec.(*grib.Section7).Data()
		assert.Equal(t, numVals, len(data))
		assert.InDelta(t, 2.9611706734e+02, data[0], 1e-8)
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
