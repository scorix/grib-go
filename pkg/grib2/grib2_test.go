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
	"github.com/scorix/grib-go/pkg/grib2/regulation"
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
	assert.EqualExportedValues(t, template, sec5.GetDataRepresentationTemplate())
}

func TestGrib_ReadSection_SimplePacking(t *testing.T) {
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
						ScaleFactorOfRadiusOfSphericalEarth:    -1,
						ScaledValueOfRadiusOfSphericalEarth:    -1,
						ScaleFactorOfEarthMajorAxis:            -1,
						ScaledValueOfEarthMajorAxis:            -1,
						ScaleFactorOfEarthMinorAxis:            -1,
						ScaledValueOfEarthMinorAxis:            -1,
						Ni:                                     363,
						Nj:                                     373,
						BasicAngleOfTheInitialProductionDomain: 0,
						SubdivisionsOfBasicAngle:               -1,
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
					BackgroundProcess:               -1,
					GeneratingProcessIdentifier:     -1,
					HoursAfterDataCutoff:            -1,
					MinutesAfterDataCutoff:          -1,
					IndicatorOfUnitOfTimeRange:      1,
					ForecastTime:                    0,
					TypeOfFirstFixedSurface:         1,
					ScaleFactorOfFirstFixedSurface:  -1,
					ScaledValueOfFirstFixedSurface:  -1,
					TypeOfSecondFixedSurface:        -1,
					ScaleFactorOfSecondFixedSurface: -1,
					ScaledValueOfSecondFixedSurface: -1,
				})
			},
		},
		{
			name: "section 5",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 5, 21)
				assertSection5(t, sec, &gridpoint.SimplePacking{
					R:    0.0194875,
					E:    -18,
					D:    -4,
					Bits: 12,
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

func TestGrib_ReadSection_ComplexPacking(t *testing.T) {
	f, err := os.Open("../testdata/grid_complex.grib2")
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
				assertSection0(t, sec, 2, 0, 81023)
			},
		},
		{
			name: "section 1",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 1, 21)
				assertSection1(t, sec, "2019-01-06T12:00:00Z")
			},
		},
		{
			name: "section 3",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 3, 72)

				tpl := gdt.Template0FixedPart{
					ShapeOfTheEarth:             6,
					Ni:                          360,
					Nj:                          181,
					SubdivisionsOfBasicAngle:    -1,
					LatitudeOfFirstGridPoint:    90000000,
					LongitudeOfFirstGridPoint:   0,
					ResolutionAndComponentFlags: 48,
					LatitudeOfLastGridPoint:     -90000000,
					LongitudeOfLastGridPoint:    359000000,
					IDirectionIncrement:         1000000,
					JDirectionIncrement:         1000000,
				}
				assertSection3(t, sec, &gdt.Template0{
					Template0FixedPart: tpl,
				})

				assert.Equal(t, true, regulation.IsMissingValue(tpl.SubdivisionsOfBasicAngle))
			},
		},
		{
			name: "section 4",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 4, 34)
				assertSection4(t, sec, &pdt.Template0{
					ParameterCategory:               2,
					ParameterNumber:                 2,
					TypeOfGeneratingProcess:         2,
					BackgroundProcess:               0,
					GeneratingProcessIdentifier:     96,
					HoursAfterDataCutoff:            0,
					MinutesAfterDataCutoff:          0,
					IndicatorOfUnitOfTimeRange:      1,
					ForecastTime:                    6,
					TypeOfFirstFixedSurface:         103,
					ScaleFactorOfFirstFixedSurface:  0,
					ScaledValueOfFirstFixedSurface:  10,
					TypeOfSecondFixedSurface:        -1,
					ScaleFactorOfSecondFixedSurface: 0,
					ScaledValueOfSecondFixedSurface: 0,
				})
			},
		},
		{
			name: "section 5",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 5, 47)
				assertSection5(t, sec, &gridpoint.ComplexPacking{
					SimplePacking:              &gridpoint.SimplePacking{R: -2023.1235, E: 0, D: 2, Bits: 12},
					GroupMethod:                1,
					MissingValue:               0,
					PrimaryMissingSubstitute:   1649987994,
					SecondaryMissingSubstitute: -1,
					NumberOfGroups:             2732,
					Group: &gridpoint.Group{
						Widths:            0,
						WidthsBits:        4,
						LengthsReference:  1,
						LengthIncrement:   1,
						LastLength:        17,
						ScaledLengthsBits: 7,
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
				assertSection(t, sec, 7, 80823)
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

func TestGrib_ReadSection_ComplexPackingAndSpatialDifferencing(t *testing.T) {
	f, err := os.Open("../testdata/hpbl.grib2")
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
				assertSection0(t, sec, 2, 0, 1476981)
			},
		},
		{
			name: "section 1",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 1, 21)
				assertSection1(t, sec, "2024-08-20T12:00:00Z")
			},
		},
		{
			name: "section 3",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 3, 72)

				tpl := gdt.Template0FixedPart{
					ShapeOfTheEarth:             6,
					Ni:                          1440,
					Nj:                          721,
					SubdivisionsOfBasicAngle:    -1,
					LatitudeOfFirstGridPoint:    90000000,
					LongitudeOfFirstGridPoint:   0,
					ResolutionAndComponentFlags: 48,
					LatitudeOfLastGridPoint:     -90000000,
					LongitudeOfLastGridPoint:    359750000,
					IDirectionIncrement:         250000,
					JDirectionIncrement:         250000,
				}
				assertSection3(t, sec, &gdt.Template0{
					Template0FixedPart: tpl,
				})

				assert.Equal(t, true, regulation.IsMissingValue(tpl.SubdivisionsOfBasicAngle))
			},
		},
		{
			name: "section 4",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 4, 34)
				assertSection4(t, sec, &pdt.Template0{
					ParameterCategory:               3,
					ParameterNumber:                 196,
					TypeOfGeneratingProcess:         2,
					BackgroundProcess:               0,
					GeneratingProcessIdentifier:     96,
					HoursAfterDataCutoff:            0,
					MinutesAfterDataCutoff:          0,
					IndicatorOfUnitOfTimeRange:      1,
					ForecastTime:                    44,
					TypeOfFirstFixedSurface:         1,
					ScaleFactorOfFirstFixedSurface:  0,
					ScaledValueOfFirstFixedSurface:  0,
					TypeOfSecondFixedSurface:        -1,
					ScaleFactorOfSecondFixedSurface: 0,
					ScaledValueOfSecondFixedSurface: 0,
				})
			},
		},
		{
			name: "section 5",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 5, 49)
				assertSection5(t, sec, &gridpoint.ComplexPackingAndSpatialDifferencing{
					ComplexPacking: &gridpoint.ComplexPacking{
						SimplePacking:              &gridpoint.SimplePacking{R: 772.85974, E: 3, D: 2, Bits: 17},
						GroupMethod:                1,
						MissingValue:               0,
						PrimaryMissingSubstitute:   1649987994,
						SecondaryMissingSubstitute: -1,
						NumberOfGroups:             30736,
						Group: &gridpoint.Group{
							Widths:            0,
							WidthsBits:        5,
							LengthsReference:  1,
							LengthIncrement:   1,
							LastLength:        41,
							ScaledLengthsBits: 7,
						},
					},
					SpatialOrderDifference: 2,
					OctetsNumber:           3,
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
				assertSection(t, sec, 7, 1476779)
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

func TestSection7_ReadData_SimplePacking(t *testing.T) {
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

func TestSection7_ReadData_ComplexPacking(t *testing.T) {
	f, err := os.Open("../testdata/grid_complex.grib2")
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

	// grib_dump -O pkg/testdata/hpbl.grib2
	exampleValues := []float64{
		1.1887646484e+00, 1.1687646484e+00, 1.1387646484e+00, 1.1187646484e+00, 1.0987646484e+00, 1.0687646484e+00, 1.0487646484e+00, 1.0287646484e+00,
		9.9876464844e-01, 9.6876464844e-01, 9.4876464844e-01, 9.1876464844e-01, 8.9876464844e-01, 8.6876464844e-01, 8.3876464844e-01, 8.1876464844e-01,
		7.8876464844e-01, 7.5876464844e-01, 7.3876464844e-01, 7.0876464844e-01, 6.7876464844e-01, 6.4876464844e-01, 6.1876464844e-01, 5.8876464844e-01,
		5.6876464844e-01, 5.3876464844e-01, 5.0876464844e-01, 4.7876464844e-01, 4.4876464844e-01, 4.1876464844e-01, 3.8876464844e-01, 3.5876464844e-01,
		3.2876464844e-01, 2.9876464844e-01, 2.6876464844e-01, 2.3876464844e-01, 2.0876464844e-01, 1.7876464844e-01, 1.4876464844e-01, 1.1876464844e-01,
		8.8764648437e-02, 5.8764648438e-02, 2.8764648438e-02, -1.2353515625e-03, -4.1235351563e-02, -7.1235351563e-02, -1.0123535156e-01, -1.3123535156e-01,
		-1.6123535156e-01, -1.9123535156e-01, -2.2123535156e-01, -2.5123535156e-01, -2.8123535156e-01, -3.1123535156e-01, -3.4123535156e-01, -3.7123535156e-01,
		-4.0123535156e-01, -4.3123535156e-01, -4.6123535156e-01, -4.9123535156e-01, -5.2123535156e-01, -5.5123535156e-01, -5.7123535156e-01, -6.0123535156e-01,
		-6.3123535156e-01, -6.6123535156e-01, -6.9123535156e-01, -7.2123535156e-01, -7.4123535156e-01, -7.7123535156e-01, -8.0123535156e-01, -8.3123535156e-01,
		-8.5123535156e-01, -8.8123535156e-01, -9.1123535156e-01, -9.3123535156e-01, -9.6123535156e-01, -9.8123535156e-01, -1.0112353516e+00, -1.0312353516e+00,
		-1.0612353516e+00, -1.0812353516e+00, -1.1012353516e+00, -1.1312353516e+00, -1.1512353516e+00, -1.1712353516e+00, -1.2012353516e+00, -1.2212353516e+00,
		-1.2412353516e+00, -1.2612353516e+00, -1.2812353516e+00, -1.3012353516e+00, -1.3212353516e+00, -1.3412353516e+00, -1.3612353516e+00, -1.3812353516e+00,
		-1.4012353516e+00, -1.4212353516e+00, -1.4412353516e+00, -1.4512353516e+00,
	}

	for i, v := range exampleValues {
		assert.InDelta(t, v, data[i], 1e-8)
	}
}

func TestSection7_ReadData_ComplexPackingAndSpatialDifferencing(t *testing.T) {
	f, err := os.Open("../testdata/hpbl.grib2")
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

	// grib_dump -O pkg/testdata/hpbl.grib2
	exampleValues := []float64{
		3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02,
		3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02,
		3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02,
		3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02,
		3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02,
		3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02,
		3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02,
		3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02,
		3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02,
		3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02,
		3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02,
		3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02,
		3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02, 3.2348859741e+02,
	}

	for i, v := range exampleValues {
		assert.InDelta(t, v, data[i], 1e-8)
	}
}
