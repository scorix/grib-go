package grib2_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/scorix/grib-go/pkg/grib2"
	grib "github.com/scorix/grib-go/pkg/grib2"
	"github.com/scorix/grib-go/pkg/grib2/drt"
	gridpoint "github.com/scorix/grib-go/pkg/grib2/drt/grid_point"
	"github.com/scorix/grib-go/pkg/grib2/gdt"
	"github.com/scorix/grib-go/pkg/grib2/pdt"
	"github.com/scorix/grib-go/pkg/grib2/scale"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertSection(t testing.TB, sec grib2.Section, number int, length int) {
	t.Helper()

	assert.Equal(t, number, sec.Number())
	assert.Equal(t, length, sec.Length())
}

func assertSection0(t testing.TB, sec grib2.Section, editionNumber int, discipline int, gribLen int) {
	t.Helper()

	require.Implements(t, (*grib2.Section0)(nil), sec)

	sec0 := sec.(grib2.Section0)
	assert.Equal(t, editionNumber, sec0.GetEditionNumber())
	assert.Equal(t, discipline, sec0.GetDiscipline())
	assert.Equal(t, gribLen, sec0.GetGribLength())
}

func assertSection1(t testing.TB, sec grib.Section, rfc3339 string) {
	t.Helper()

	require.Implements(t, (*grib2.Section1)(nil), sec)

	sec1 := sec.(grib2.Section1)
	assert.Equal(t, rfc3339, sec1.GetTime(time.UTC).Format(time.RFC3339))
}

func assertSection3(t testing.TB, sec grib.Section, template gdt.Template) {
	t.Helper()

	require.Implements(t, (*grib2.Section3)(nil), sec)

	sec3 := sec.(grib2.Section3)
	assert.Equal(t, template, sec3.GetGridDefinitionTemplate())
}

func assertSection4(t testing.TB, sec grib.Section, template pdt.Template) {
	t.Helper()

	require.Implements(t, (*grib2.Section4)(nil), sec)

	sec4 := sec.(grib2.Section4)
	assert.Equal(t, template, sec4.GetProductDefinitionTemplate())
}

func assertSection5(t testing.TB, sec grib.Section, template drt.Template) {
	t.Helper()

	require.Implements(t, (*grib2.Section5)(nil), sec)

	sec5 := sec.(grib2.Section5)
	assert.Equal(t, template, sec5.GetDataRepresentationTemplate())
}

func TestGrib_ReadSection(t *testing.T) {
	f, err := os.Open("../testdata/temp.grib2")
	require.NoError(t, err)

	g := grib.NewGrib2(f)

	tests := []struct {
		name string
		test func(t *testing.T, sec grib2.Section)
		err  error
	}{
		{
			name: "section 0",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 0, 16)
				assertSection0(t, sec, 2, 0, 203278)
			},
		},
		{
			name: "section 1",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 1, 21)
				assertSection1(t, sec, "2023-07-11T00:00:00Z")
			},
		},
		{
			name: "section 3",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 3, 72)
				assertSection3(t, sec, &gdt.Template0{
					Template0FixedPart: gdt.Template0FixedPart{
						ShapeOfTheEarth:                        6,
						ScaleFactorOfRadiusOfSphericalEarth:    255,
						ScaledValueOfRadiusOfSphericalEarth:    4294967295,
						ScaleFactorOfEarthMajorAxis:            255,
						ScaledValueOfEarthMajorAxis:            4294967295,
						ScaleFactorOfEarthMinorAxis:            255,
						ScaledValueOfEarthMinorAxis:            4294967295,
						Ni:                                     363,
						Nj:                                     373,
						BasicAngleOfTheInitialProductionDomain: 0,
						SubdivisionsOfBasicAngle:               4294967295,
						LatitudeOfFirstGridPoint:               33046875,
						LongitudeOfFirstGridPoint:              346007813,
						ResolutionAndComponentFlags:            48,
						LatitudeOfLastGridPoint:                67921875,
						LongitudeOfLastGridPoint:               36914063,
						IDirectionIncrement:                    140625,
						JDirectionIncrement:                    93750,
						ScanningMode:                           64,
					},
				})
			},
		},
		{
			name: "section 4",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 4, 34)
				assertSection4(t, sec, &pdt.Template0{
					ParameterCategory:               0,
					ParameterNumber:                 0,
					TypeOfGeneratingProcess:         2,
					BackgroundProcess:               255,
					GeneratingProcessIdentifier:     255,
					HoursAfterDataCutoff:            65535,
					MinutesAfterDataCutoff:          255,
					IndicatorOfUnitOfTimeRange:      1,
					ForecastTime:                    0,
					TypeOfFirstFixedSurface:         1,
					ScaleFactorOfFirstFixedSurface:  255,
					ScaledValueOfFirstFixedSurface:  4294967295,
					TypeOfSecondFixedSurface:        255,
					ScaleFactorOfSecondFixedSurface: 255,
					ScaledValueOfSecondFixedSurface: 4294967295,
				})
			},
		},
		{
			name: "section 5",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 5, 21)
				assertSection5(t, sec, &gridpoint.SimplePacking{
					DefSimplePacking: &gridpoint.DefSimplePacking{
						R:    float32(0.0194875),
						E:    scale.Factor(32786),
						D:    scale.Factor(32772),
						Bits: 12,
					},
				})
			},
		},
		{
			name: "section 6",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 6, 6)
			},
		},
		{
			name: "section 7",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 7, 203104)
			},
		},
		{
			name: "section 8",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 8, 4)
			},
		},
		{
			name: "eof",
			test: func(t *testing.T, sec grib2.Section) {
				assert.Nil(t, sec)
			},
			err: io.EOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t2 *testing.T) {
			sec, err := g.ReadSection()
			if tt.err == nil {
				require.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.err)
			}

			tt.test(t2, sec)
		})
	}
}

func TestSection7_ReadData(t *testing.T) {
	f, err := os.Open("../testdata/temp.grib2")
	require.NoError(t, err)

	g := grib.NewGrib2(f)

	var sec7 grib2.Section7
	var tpl drt.Template
	var dataLen int

	for {
		sec, err := g.ReadSection()
		require.NoError(t, err)

		if sec.Number() == 5 {
			tpl = sec.(grib2.Section5).GetDataRepresentationTemplate()
			dataLen = sec.(grib2.Section5).GetNumberOfValues()
		}

		if sec.Number() == 7 {
			sec7 = sec.(grib2.Section7)

			break
		}
	}

	data, err := sec7.GetData(tpl)
	require.NoError(t, err)
	require.NotZero(t, dataLen)
	require.Equal(t, dataLen, len(data))

	// grib_dump -O pkg/testdata/temp.grib2
	exampleValues := []float64{
		2.9611706734e+02, 2.9600262642e+02, 2.9588818550e+02, 2.9562115669e+02, 2.9562115669e+02, 2.9550671577e+02, 2.9550671577e+02, 2.9562115669e+02,
		2.9562115669e+02, 2.9573559761e+02, 2.9573559761e+02, 2.9588818550e+02, 2.9588818550e+02, 2.9588818550e+02, 2.9600262642e+02, 2.9611706734e+02,
		2.9638409615e+02, 2.9649853706e+02, 2.9661297798e+02, 2.9676556587e+02, 2.9676556587e+02, 2.9661297798e+02, 2.9661297798e+02, 2.9661297798e+02,
		2.9661297798e+02, 2.9661297798e+02, 2.9661297798e+02, 2.9661297798e+02, 2.9676556587e+02, 2.9676556587e+02, 2.9699444771e+02, 2.9699444771e+02,
		2.9688000679e+02, 2.9649853706e+02, 2.9600262642e+02, 2.9512524605e+02, 2.9485821724e+02, 2.9512524605e+02, 2.9424786568e+02, 2.9424786568e+02,
		2.9462933540e+02, 2.9401898384e+02, 2.9287457466e+02, 2.9523968697e+02, 2.9401898384e+02, 2.9501080513e+02, 2.9386639595e+02, 2.9363751411e+02,
		2.9588818550e+02, 2.9688000679e+02, 2.9787182808e+02, 2.9661297798e+02, 2.9726147652e+02, 2.9848217964e+02, 2.9726147652e+02, 2.9924511909e+02,
		2.9913067818e+02, 2.9485821724e+02, 2.9825329781e+02, 2.9798626900e+02, 2.9787182808e+02, 2.9523968697e+02, 2.9287457466e+02, 2.9073834419e+02,
		2.8810620308e+02, 2.9413342476e+02, 2.9462933540e+02, 2.9539227486e+02, 2.9649853706e+02, 2.9611706734e+02, 2.9638409615e+02, 2.9737591743e+02,
		2.9764294624e+02, 3.0000805855e+02, 2.9951214790e+02, 2.9798626900e+02, 2.9501080513e+02, 2.9375195503e+02, 2.9337048531e+02, 2.9523968697e+02,
		2.9649853706e+02, 2.9688000679e+02, 2.9649853706e+02, 2.9699444771e+02, 2.9562115669e+02, 2.9550671577e+02, 2.9588818550e+02, 2.9989361763e+02,
		2.9848217964e+02, 3.0023694038e+02, 3.0126690865e+02, 3.0088543892e+02, 3.0023694038e+02, 2.9962658882e+02, 2.9825329781e+02, 2.9726147652e+02,
		2.9699444771e+02, 2.9699444771e+02, 2.9726147652e+02, 2.9649853706e+02,
	}

	for i, v := range exampleValues {
		assert.InDelta(t, v, data[i], 1e-8)
	}
}
