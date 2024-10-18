package grib2_test

import (
	"context"
	"errors"
	"image/png"
	"io"
	"os"
	"testing"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/icza/bitio"
	ossio "github.com/scorix/aliyun-oss-io"
	"github.com/scorix/grib-go/pkg/grib2"
	grib "github.com/scorix/grib-go/pkg/grib2"
	"github.com/scorix/grib-go/pkg/grib2/drt"
	gridpoint "github.com/scorix/grib-go/pkg/grib2/drt/grid_point"
	"github.com/scorix/grib-go/pkg/grib2/gdt"
	"github.com/scorix/grib-go/pkg/grib2/pdt"
	"github.com/scorix/grib-go/pkg/grib2/regulation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/mmap"
	"golang.org/x/sync/errgroup"
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

func assertSection7(t testing.TB, sec grib.Section, dataOffset int64) {
	t.Helper()

	require.Implements(t, (*grib2.Section7)(nil), sec)

	sec7 := sec.(grib2.Section7)
	assert.Equal(t, dataOffset, sec7.GetDataOffset())
}

func TestGrib_ReadSectionAt_SimplePacking(t *testing.T) {
	t.Parallel()

	f, err := os.Open("../testdata/temp.grib2")
	require.NoError(t, err)
	defer f.Close()

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
					ParameterCategory:       0,
					ParameterNumber:         0,
					TypeOfGeneratingProcess: 2,
					BackgroundProcess:       -1,
					AnalysisOrForecastGeneratingProcessIdentified: -1,
					HoursAfterDataCutoff:                          -1,
					MinutesAfterDataCutoff:                        -1,
					IndicatorOfUnitForForecastTime:                1,
					ForecastTime:                                  0,
					TypeOfFirstFixedSurface:                       1,
					ScaleFactorOfFirstFixedSurface:                -1,
					ScaledValueOfFirstFixedSurface:                -1,
					TypeOfSecondFixedSurface:                      255,
					ScaleFactorOfSecondFixedSurface:               -1,
					ScaledValueOfSecondFixedSurface:               -1,
				})
			},
		},
		{
			name: "section 5",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 5, 21)
				assertSection5(t, sec, &gridpoint.SimplePacking{
					ReferenceValue:     0.0194875,
					BinaryScaleFactor:  -18,
					DecimalScaleFactor: -4,
					Bits:               12,
					NumVals:            135399,
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
				assertSection7(t, sec, 175)
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

	var offset int64

	for _, tt := range tests {
		t.Run(tt.name, func(t2 *testing.T) {
			sec, err := g.ReadSectionAt(offset)
			if tt.err == nil {
				require.NoError(t, err)
				offset += int64(sec.Length())
			} else {
				assert.ErrorIs(t, err, tt.err)
			}

			tt.test(t2, sec)
		})
	}
}

func TestGrib_ReadSectionAt_ComplexPacking(t *testing.T) {
	t.Parallel()

	f, err := os.Open("../testdata/grid_complex.grib2")
	require.NoError(t, err)
	defer f.Close()

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

				assert.Equal(t, true, regulation.IsMissingValue(uint(tpl.SubdivisionsOfBasicAngle), 32))
			},
		},
		{
			name: "section 4",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 4, 34)
				assertSection4(t, sec, &pdt.Template0{
					ParameterCategory:       2,
					ParameterNumber:         2,
					TypeOfGeneratingProcess: 2,
					BackgroundProcess:       0,
					AnalysisOrForecastGeneratingProcessIdentified: 96,
					HoursAfterDataCutoff:                          0,
					MinutesAfterDataCutoff:                        0,
					IndicatorOfUnitForForecastTime:                1,
					ForecastTime:                                  6,
					TypeOfFirstFixedSurface:                       103,
					ScaleFactorOfFirstFixedSurface:                0,
					ScaledValueOfFirstFixedSurface:                10,
					TypeOfSecondFixedSurface:                      255,
					ScaleFactorOfSecondFixedSurface:               0,
					ScaledValueOfSecondFixedSurface:               0,
				})
			},
		},
		{
			name: "section 5",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 5, 47)
				assertSection5(t, sec, &gridpoint.ComplexPacking{
					SimplePacking:              &gridpoint.SimplePacking{ReferenceValue: -2023.1235, BinaryScaleFactor: 0, DecimalScaleFactor: 2, Bits: 12, NumVals: 65160},
					GroupSplittingMethodUsed:   1,
					MissingValueManagementUsed: 0,
					PrimaryMissingSubstitute:   1649987994,
					SecondaryMissingSubstitute: -1,
					Grouping: &gridpoint.Grouping{
						NumberOfGroups:    2732,
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
				assertSection7(t, sec, 201)
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

	var offset int64

	for _, tt := range tests {
		t.Run(tt.name, func(t2 *testing.T) {
			sec, err := g.ReadSectionAt(offset)
			if tt.err == nil {
				require.NoError(t, err)
				offset += int64(sec.Length())
			} else {
				assert.ErrorIs(t, err, tt.err)
			}

			tt.test(t2, sec)
		})
	}
}

func TestGrib_ReadSectionAt_ComplexPackingAndSpatialDifferencing(t *testing.T) {
	t.Parallel()

	f, err := os.Open("../testdata/hpbl.grib2")
	require.NoError(t, err)
	defer f.Close()

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

				assert.Equal(t, true, regulation.IsMissingValue(uint(tpl.SubdivisionsOfBasicAngle), 32))
			},
		},
		{
			name: "section 4",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 4, 34)
				assertSection4(t, sec, &pdt.Template0{
					ParameterCategory:       3,
					ParameterNumber:         196,
					TypeOfGeneratingProcess: 2,
					BackgroundProcess:       0,
					AnalysisOrForecastGeneratingProcessIdentified: 96,
					HoursAfterDataCutoff:                          0,
					MinutesAfterDataCutoff:                        0,
					IndicatorOfUnitForForecastTime:                1,
					ForecastTime:                                  44,
					TypeOfFirstFixedSurface:                       1,
					ScaleFactorOfFirstFixedSurface:                0,
					ScaledValueOfFirstFixedSurface:                0,
					TypeOfSecondFixedSurface:                      255,
					ScaleFactorOfSecondFixedSurface:               0,
					ScaledValueOfSecondFixedSurface:               0,
				})
			},
		},
		{
			name: "section 5",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 5, 49)
				assertSection5(t, sec, &gridpoint.ComplexPackingAndSpatialDifferencing{
					ComplexPacking: &gridpoint.ComplexPacking{
						SimplePacking:              &gridpoint.SimplePacking{ReferenceValue: 772.85974, BinaryScaleFactor: 3, DecimalScaleFactor: 2, Bits: 17, NumVals: 1038240},
						GroupSplittingMethodUsed:   1,
						MissingValueManagementUsed: 0,
						PrimaryMissingSubstitute:   1649987994,
						SecondaryMissingSubstitute: -1,
						Grouping: &gridpoint.Grouping{
							NumberOfGroups:    30736,
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
				assertSection7(t, sec, 203)
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

	var offset int64

	for _, tt := range tests {
		t.Run(tt.name, func(t2 *testing.T) {
			sec, err := g.ReadSectionAt(offset)
			if tt.err == nil {
				require.NoError(t, err)
				offset += int64(sec.Length())
			} else {
				assert.ErrorIs(t, err, tt.err)
			}

			tt.test(t2, sec)
		})
	}
}

func TestGrib_ReadSectionAt_PNG(t *testing.T) {
	t.Parallel()

	f, err := os.Open("../testdata/grid_png.grib2")
	require.NoError(t, err)
	defer f.Close()

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
				assertSection0(t, sec, 2, 0, 78374)
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
				assertSection3(t, sec, &gdt.Template0{
					Template0FixedPart: gdt.Template0FixedPart{
						ShapeOfTheEarth:                        6,
						ScaleFactorOfRadiusOfSphericalEarth:    0,
						ScaledValueOfRadiusOfSphericalEarth:    0,
						ScaleFactorOfEarthMajorAxis:            0,
						ScaledValueOfEarthMajorAxis:            0,
						ScaleFactorOfEarthMinorAxis:            0,
						ScaledValueOfEarthMinorAxis:            0,
						Ni:                                     360,
						Nj:                                     181,
						BasicAngleOfTheInitialProductionDomain: 0,
						SubdivisionsOfBasicAngle:               -1,
						LatitudeOfFirstGridPoint:               90000000,
						LongitudeOfFirstGridPoint:              0,
						ResolutionAndComponentFlags:            48,
						LatitudeOfLastGridPoint:                -90000000,
						LongitudeOfLastGridPoint:               359000000,
						IDirectionIncrement:                    1000000,
						JDirectionIncrement:                    1000000,
						ScanningMode:                           0,
					},
				})
			},
		},
		{
			name: "section 4",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 4, 34)
				assertSection4(t, sec, &pdt.Template0{
					ParameterCategory:       2,
					ParameterNumber:         2,
					TypeOfGeneratingProcess: 2,
					BackgroundProcess:       0,
					AnalysisOrForecastGeneratingProcessIdentified: 96,
					HoursAfterDataCutoff:                          0,
					MinutesAfterDataCutoff:                        0,
					IndicatorOfUnitForForecastTime:                1,
					ForecastTime:                                  6,
					TypeOfFirstFixedSurface:                       103,
					ScaleFactorOfFirstFixedSurface:                0,
					ScaledValueOfFirstFixedSurface:                10,
					TypeOfSecondFixedSurface:                      255,
					ScaleFactorOfSecondFixedSurface:               0,
					ScaledValueOfSecondFixedSurface:               0,
				})
			},
		},
		{
			name: "section 5",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 5, 21)
				assertSection5(t, sec, &gridpoint.PortableNetworkGraphics{
					ReferenceValue:     -2023.1235,
					BinaryScaleFactor:  1,
					DecimalScaleFactor: 2,
					Bits:               12,
					NumVals:            65160,
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
				assertSection(t, sec, 7, 78200)
				assertSection7(t, sec, 175)
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

	var offset int64

	for _, tt := range tests {
		t.Run(tt.name, func(t2 *testing.T) {
			sec, err := g.ReadSectionAt(offset)
			if tt.err == nil {
				require.NoError(t, err)
				offset += int64(sec.Length())
			} else {
				assert.ErrorIs(t, err, tt.err)
			}

			tt.test(t2, sec)
		})
	}
}

func TestGrib_ReadSectionAt_tmax(t *testing.T) {
	t.Parallel()

	f, err := os.Open("../testdata/tmax.grib2")
	require.NoError(t, err)
	defer f.Close()

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
				assertSection0(t, sec, 2, 0, 499330)
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

				assert.Equal(t, true, regulation.IsMissingValue(uint(tpl.SubdivisionsOfBasicAngle), 32))
			},
		},
		{
			name: "section 4",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 4, 58)
				assertSection4(t, sec, &pdt.Template8{
					Template0: &pdt.Template0{
						ParameterCategory:       0,
						ParameterNumber:         4,
						TypeOfGeneratingProcess: 2,
						BackgroundProcess:       0,
						AnalysisOrForecastGeneratingProcessIdentified: 96,
						HoursAfterDataCutoff:                          0,
						MinutesAfterDataCutoff:                        0,
						IndicatorOfUnitForForecastTime:                1,
						ForecastTime:                                  42,
						TypeOfFirstFixedSurface:                       103,
						ScaleFactorOfFirstFixedSurface:                0,
						ScaledValueOfFirstFixedSurface:                2,
						TypeOfSecondFixedSurface:                      255,
						ScaleFactorOfSecondFixedSurface:               0,
						ScaledValueOfSecondFixedSurface:               0,
					},
				})
			},
		},
		{
			name: "section 5",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 5, 49)
				assertSection5(t, sec, &gridpoint.ComplexPackingAndSpatialDifferencing{
					ComplexPacking: &gridpoint.ComplexPacking{
						SimplePacking:              &gridpoint.SimplePacking{ReferenceValue: 2125.6357, BinaryScaleFactor: 0, DecimalScaleFactor: 1, Bits: 9, NumVals: 1038240},
						GroupSplittingMethodUsed:   1,
						MissingValueManagementUsed: 0,
						PrimaryMissingSubstitute:   1649987994,
						SecondaryMissingSubstitute: -1,
						Grouping: &gridpoint.Grouping{
							NumberOfGroups:    32100,
							Widths:            0,
							WidthsBits:        4,
							LengthsReference:  1,
							LengthIncrement:   1,
							LastLength:        41,
							ScaledLengthsBits: 7,
						},
					},
					SpatialOrderDifference: 2,
					OctetsNumber:           2,
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
				assertSection(t, sec, 7, 499104)
				assertSection7(t, sec, 227)
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

	var offset int64

	for _, tt := range tests {
		t.Run(tt.name, func(t2 *testing.T) {
			sec, err := g.ReadSectionAt(offset)
			if tt.err == nil {
				require.NoError(t, err)
				offset += int64(sec.Length())
			} else {
				assert.ErrorIs(t, err, tt.err)
			}

			tt.test(t2, sec)
		})
	}
}

func TestGrib_ReadSectionAt_cwat(t *testing.T) {
	t.Parallel()

	f, err := os.Open("../testdata/cwat.grib2")
	require.NoError(t, err)
	defer f.Close()

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
				assertSection0(t, sec, 2, 0, 402221)
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

				assert.Equal(t, true, regulation.IsMissingValue(uint(tpl.SubdivisionsOfBasicAngle), 32))
			},
		},
		{
			name: "section 4",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 4, 34)
				assertSection4(t, sec, &pdt.Template0{
					ParameterCategory:       6,
					ParameterNumber:         6,
					TypeOfGeneratingProcess: 2,
					BackgroundProcess:       0,
					AnalysisOrForecastGeneratingProcessIdentified: 81,
					HoursAfterDataCutoff:                          0,
					MinutesAfterDataCutoff:                        0,
					IndicatorOfUnitForForecastTime:                1,
					ForecastTime:                                  0,
					TypeOfFirstFixedSurface:                       200,
					ScaleFactorOfFirstFixedSurface:                0,
					ScaledValueOfFirstFixedSurface:                0,
					TypeOfSecondFixedSurface:                      255,
					ScaleFactorOfSecondFixedSurface:               0,
					ScaledValueOfSecondFixedSurface:               0,
				})
			},
		},
		{
			name: "section 5",
			test: func(t *testing.T, sec grib2.Section) {
				assertSection(t, sec, 5, 49)
				assertSection5(t, sec, &gridpoint.ComplexPackingAndSpatialDifferencing{
					ComplexPacking: &gridpoint.ComplexPacking{
						SimplePacking:              &gridpoint.SimplePacking{ReferenceValue: 0, BinaryScaleFactor: 0, DecimalScaleFactor: 2, Bits: 10, NumVals: 1038240},
						GroupSplittingMethodUsed:   1,
						MissingValueManagementUsed: 0,
						PrimaryMissingSubstitute:   1649987994,
						SecondaryMissingSubstitute: -1,
						Grouping: &gridpoint.Grouping{
							NumberOfGroups:    31329,
							Widths:            0,
							WidthsBits:        4,
							LengthsReference:  1,
							LengthIncrement:   1,
							LastLength:        23,
							ScaledLengthsBits: 7,
						},
					},
					SpatialOrderDifference: 2,
					OctetsNumber:           2,
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
				assertSection(t, sec, 7, 402019)
				assertSection7(t, sec, 203)
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

	var offset int64

	for _, tt := range tests {
		t.Run(tt.name, func(t2 *testing.T) {
			sec, err := g.ReadSectionAt(offset)
			if tt.err == nil {
				require.NoError(t, err)
				offset += int64(sec.Length())
			} else {
				assert.ErrorIs(t, err, tt.err)
			}

			tt.test(t2, sec)
		})
	}
}

func TestSection7_ReadData_SimplePacking(t *testing.T) {
	t.Parallel()

	f, err := os.Open("../testdata/temp.grib2")
	require.NoError(t, err)
	defer f.Close()

	g := grib.NewGrib2(f)

	var sec7 grib2.Section7
	var tpl drt.Template
	var dataLen int

	var offset int64

	for {
		sec, err := g.ReadSectionAt(offset)
		require.NoError(t, err)

		offset += int64(sec.Length())

		if sec.Number() == 5 {
			tpl = sec.(grib2.Section5).GetDataRepresentationTemplate()
			dataLen = sec.(grib2.Section5).GetNumberOfValues()
			require.Equal(t, 135399, dataLen)
		}

		if sec.Number() == 7 {
			sec7 = sec.(grib2.Section7)

			break
		}
	}

	data, err := sec7.GetData(tpl)
	require.NoError(t, err)
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
	t.Parallel()

	f, err := os.Open("../testdata/grid_complex.grib2")
	require.NoError(t, err)
	defer f.Close()

	g := grib.NewGrib2(f)

	var sec7 grib2.Section7
	var tpl drt.Template
	var dataLen int

	var offset int64

	for {
		sec, err := g.ReadSectionAt(offset)
		require.NoError(t, err)

		offset += int64(sec.Length())

		if sec.Number() == 5 {
			tpl = sec.(grib2.Section5).GetDataRepresentationTemplate()
			dataLen = sec.(grib2.Section5).GetNumberOfValues()
			require.Equal(t, 65160, dataLen)
		}

		if sec.Number() == 7 {
			sec7 = sec.(grib2.Section7)
			require.Equal(t, 80823, sec7.Length())

			break
		}
	}

	data, err := sec7.GetData(tpl)
	require.NoError(t, err)
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
	t.Parallel()

	f, err := os.Open("../testdata/hpbl.grib2")
	require.NoError(t, err)
	defer f.Close()

	g := grib.NewGrib2(f)

	var sec7 grib2.Section7
	var tpl drt.Template
	var dataLen int

	var offset int64

	for {
		sec, err := g.ReadSectionAt(offset)
		require.NoError(t, err)

		offset += int64(sec.Length())

		if sec.Number() == 5 {
			tpl = sec.(grib2.Section5).GetDataRepresentationTemplate()
			dataLen = sec.(grib2.Section5).GetNumberOfValues()
			require.Equal(t, 1038240, dataLen)
		}

		if sec.Number() == 7 {
			sec7 = sec.(grib2.Section7)

			break
		}
	}

	data, err := sec7.GetData(tpl)
	require.NoError(t, err)
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

func TestSection7_ReadData_PortableNetworkGraphics(t *testing.T) {
	t.Parallel()

	f, err := os.Open("../testdata/grid_png.grib2")
	require.NoError(t, err)
	defer f.Close()

	g := grib.NewGrib2(f)

	var sec7 grib2.Section7
	var tpl drt.Template
	var dataLen int

	var offset int64

	for {
		sec, err := g.ReadSectionAt(offset)
		require.NoError(t, err)

		offset += int64(sec.Length())

		if sec.Number() == 5 {
			tpl = sec.(grib2.Section5).GetDataRepresentationTemplate()
			dataLen = sec.(grib2.Section5).GetNumberOfValues()
			require.Equal(t, 65160, dataLen)
		}

		if sec.Number() == 7 {
			sec7 = sec.(grib2.Section7)

			break
		}
	}

	data, err := sec7.GetData(tpl)
	require.NoError(t, err)
	require.Equal(t, dataLen, len(data))

	// grib_dump -O pkg/testdata/hpbl.grib2
	exampleValues := []float64{
		1.18876, 1.16876, 1.14876, 1.12876, 1.10876,
		1.06876, 1.04876, 1.02876, 1.00876, 0.968765,
		0.948765, 0.928765, 0.908765, 0.868765, 0.848765,
		0.828765, 0.788765, 0.768765, 0.748765, 0.708765,
		0.688765, 0.648765, 0.628765, 0.588765, 0.568765,
		0.548765, 0.508765, 0.488765, 0.448765, 0.428765,
		0.388765, 0.368765, 0.328765, 0.308765, 0.268765,
		0.248765, 0.208765, 0.188765, 0.148765, 0.128765,
		0.0887646, 0.0687646, 0.0287646, 0.00876465, -0.0312354,
		-0.0712354, -0.0912354, -0.131235, -0.151235, -0.191235,
		-0.211235, -0.251235, -0.271235, -0.311235, -0.331235,
		-0.371235, -0.391235, -0.431235, -0.451235, -0.491235,
		-0.511235, -0.551235, -0.571235, -0.591235, -0.631235,
		-0.651235, -0.691235, -0.711235, -0.731235, -0.771235,
		-0.791235, -0.831235, -0.851235, -0.871235, -0.911235,
		-0.931235, -0.951235, -0.971235, -1.01124, -1.03124,
		-1.05124, -1.07124, -1.09124, -1.13124, -1.15124,
		-1.17124, -1.19124, -1.21124, -1.23124, -1.25124,
		-1.27124, -1.29124, -1.31124, -1.33124, -1.35124,
		-1.37124, -1.39124, -1.41124, -1.43124, -1.45124,
	}

	for i, v := range exampleValues {
		assert.InDelta(t, v, data[i], 1e-5)
	}
}

func TestGrib2_ReadMessages(t *testing.T) {
	t.Parallel()
	// aws s3 cp --no-sign-request s3://noaa-gfs-bdp-pds/gfs.20240820/12/atmos/gfs.t12z.pgrb2.0p25.f044 pkg/testdata/gfs.t12z.pgrb2.0p25.f044
	const filename = "../testdata/gfs.t12z.pgrb2.0p25.f044"

	s, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		t.Skipf("%s not exist", filename)
	}

	t.Run(s.Name(), func(t *testing.T) {
		t.Parallel()

		f, err := os.Open(filename)
		require.NoError(t, err)
		defer f.Close()

		g := grib.NewGrib2(f)
		var msgs []grib2.Message

		var offset int64

		for i := 0; ; i++ {
			msg, err := g.ReadMessageAt(offset)
			if i < 743 && err != nil {
				t.Fatal(err)
			}

			if errors.Is(err, io.EOF) {
				break
			}
			require.NoError(t, err)
			require.NotNil(t, msg)

			msgs = append(msgs, msg)
			offset += msg.GetSize()
		}

		assert.Equal(t, 743, len(msgs))

		eg, _ := errgroup.WithContext(context.TODO())
		eg.SetLimit(1)
		for i, msg := range msgs {
			eg.Go(func() error {
				data, err := msg.ReadData()
				require.NoError(t, err)

				t.Logf("count: %d, discipline: %d, category: %d, number: %d, forecast: %s, dataLen: %d", i+1, msg.GetDiscipline(), msg.GetParameterCategory(), msg.GetParameterNumber(), msg.GetForecastTime(time.UTC), len(data))

				return nil
			})
		}

		require.NoError(t, eg.Wait())
	})
}

func TestGrib2_ReadMessageAt_cwat(t *testing.T) {
	t.Parallel()

	const filename = "../testdata/cwat.grib2"

	s, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		t.Skipf("%s not exist", filename)
	}

	t.Run(s.Name(), func(t *testing.T) {
		t.Parallel()

		f, err := os.Open(filename)
		require.NoError(t, err)
		defer f.Close()

		g := grib.NewGrib2(f)

		msg, err := g.ReadMessageAt(0)
		require.NoError(t, err)
		require.NotNil(t, msg)

		assert.Equal(t, 200, msg.GetTypeOfFirstFixedSurface())
	})

	t.Run("mmap", func(t *testing.T) {
		t.Parallel()

		mm, err := mmap.Open(filename)
		require.NoError(t, err)
		defer mm.Close()

		g := grib.NewGrib2(mm)

		msg, err := g.ReadMessageAt(0)
		require.NoError(t, err)
		require.NotNil(t, msg)

		assert.Equal(t, 200, msg.GetTypeOfFirstFixedSurface())
	})
}

func TestGrib2_ReadMessageAt(t *testing.T) {
	t.Parallel()

	t.Run("file", func(t *testing.T) {
		t.Parallel()

		const filename = "../testdata/cwat.grib2"

		f, err := os.Open(filename)
		if errors.Is(err, os.ErrNotExist) {
			t.Skipf("%s not exist", filename)
		}

		require.NoError(t, err)
		defer f.Close()

		g := grib.NewGrib2(f)

		sec0, err := g.ReadSectionAt(0)
		require.NoError(t, err)
		require.Equal(t, 0, sec0.Number())
		require.Equal(t, 16, sec0.Length())

		sec1, err := g.ReadSectionAt(16)
		require.NoError(t, err)
		require.Equal(t, 1, sec1.Number())
		require.Equal(t, 21, sec1.Length())

		sec2, err := g.ReadSectionAt(37)
		require.NoError(t, err)
		require.Equal(t, 3, sec2.Number())
		require.Equal(t, 72, sec2.Length())

		sec4, err := g.ReadSectionAt(109)
		require.NoError(t, err)
		require.Equal(t, 4, sec4.Number())
		require.Equal(t, 34, sec4.Length())

		sec5, err := g.ReadSectionAt(143)
		require.NoError(t, err)
		require.Equal(t, 5, sec5.Number())
		require.Equal(t, 49, sec5.Length())

		sec6, err := g.ReadSectionAt(192)
		require.NoError(t, err)
		require.Equal(t, 6, sec6.Number())
		require.Equal(t, 6, sec6.Length())

		sec7, err := g.ReadSectionAt(198)
		require.NoError(t, err)
		require.Equal(t, 7, sec7.Number())
		require.Equal(t, 402_019, sec7.Length())

		sec8, err := g.ReadSectionAt(402_217)
		require.NoError(t, err)
		require.Equal(t, 8, sec8.Number())
		require.Equal(t, 4, sec8.Length())

		msg1, err := g.ReadMessageAt(0)
		require.NoError(t, err)
		require.NotNil(t, msg1)

		{
			_, err := g.ReadMessageAt(msg1.GetOffset() + msg1.GetSize())
			require.ErrorIs(t, err, io.EOF)
		}
	})

	t.Run("mmap", func(t *testing.T) {
		t.Parallel()

		const filename = "../testdata/cwat.grib2"

		f, err := mmap.Open(filename)
		if errors.Is(err, os.ErrNotExist) {
			t.Skipf("%s not exist", filename)
		}

		require.NoError(t, err)
		defer f.Close()

		g := grib.NewGrib2(f)

		sec0, err := g.ReadSectionAt(0)
		require.NoError(t, err)
		require.Equal(t, 0, sec0.Number())
		require.Equal(t, 16, sec0.Length())

		sec1, err := g.ReadSectionAt(16)
		require.NoError(t, err)
		require.Equal(t, 1, sec1.Number())
		require.Equal(t, 21, sec1.Length())

		sec2, err := g.ReadSectionAt(37)
		require.NoError(t, err)
		require.Equal(t, 3, sec2.Number())
		require.Equal(t, 72, sec2.Length())

		sec4, err := g.ReadSectionAt(109)
		require.NoError(t, err)
		require.Equal(t, 4, sec4.Number())
		require.Equal(t, 34, sec4.Length())

		sec5, err := g.ReadSectionAt(143)
		require.NoError(t, err)
		require.Equal(t, 5, sec5.Number())
		require.Equal(t, 49, sec5.Length())

		sec6, err := g.ReadSectionAt(192)
		require.NoError(t, err)
		require.Equal(t, 6, sec6.Number())
		require.Equal(t, 6, sec6.Length())

		sec7, err := g.ReadSectionAt(198)
		require.NoError(t, err)
		require.Equal(t, 7, sec7.Number())
		require.Equal(t, 402_019, sec7.Length())

		sec8, err := g.ReadSectionAt(402_217)
		require.NoError(t, err)
		require.Equal(t, 8, sec8.Number())
		require.Equal(t, 4, sec8.Length())

		msg1, err := g.ReadMessageAt(0)
		require.NoError(t, err)
		require.NotNil(t, msg1)

		{
			_, err := g.ReadMessageAt(msg1.GetOffset() + msg1.GetSize())
			require.ErrorIs(t, err, io.EOF)
		}
	})

	t.Run("oss", func(t *testing.T) {
		t.Parallel()

		const (
			bucketName = "cy-meteorology"
			key        = "noaa-gfs/develop/2024/09/30/18/atmos/0p25/2t_heightAboveGround_2_0_0.grib2"
			msgOffset  = 0
		)

		var (
			endpoint        = os.Getenv("ALIYUN_OSS_ENDPOINT")
			accessKeyId     = os.Getenv("ALIYUN_OSS_ACCESS_KEY_ID")
			accessKeySecret = os.Getenv("ALIYUN_OSS_ACCESS_KEY_SECRET")
		)

		ctx := context.TODO()
		cli, err := oss.New(
			endpoint,
			accessKeyId,
			accessKeySecret,
		)
		if err != nil {
			t.Skip(err.Error())
		}

		bucket, err := cli.Bucket(bucketName)
		require.NoError(t, err)

		r, err := ossio.NewReader(ctx, bucket, key)
		if err != nil {
			t.Skip(err.Error())
		}

		p := make([]byte, 16)
		n, err := r.ReadAt(p, msgOffset)
		require.Equal(t, 16, n)
		require.NoError(t, err)

		g := grib.NewGrib2(r)

		sec0, err := g.ReadSectionAt(0)
		require.NoError(t, err)
		require.Equal(t, 0, sec0.Number())
		require.Equal(t, 16, sec0.Length())

		sec1, err := g.ReadSectionAt(16)
		require.NoError(t, err)
		require.Equal(t, 1, sec1.Number())
		require.Equal(t, 21, sec1.Length())

		sec2, err := g.ReadSectionAt(37)
		require.NoError(t, err)
		require.Equal(t, 3, sec2.Number())
		require.Equal(t, 72, sec2.Length())

		sec4, err := g.ReadSectionAt(109)
		require.NoError(t, err)
		require.Equal(t, 4, sec4.Number())
		require.Equal(t, 34, sec4.Length())

		sec5, err := g.ReadSectionAt(143)
		require.NoError(t, err)
		require.Equal(t, 5, sec5.Number())
		require.Equal(t, 21, sec5.Length())

		sec6, err := g.ReadSectionAt(164)
		require.NoError(t, err)
		require.Equal(t, 6, sec6.Number())
		require.Equal(t, 6, sec6.Length())

		sec7, err := g.ReadSectionAt(170)
		require.NoError(t, err)
		require.Equal(t, 7, sec7.Number())
		require.Equal(t, 1_427_585, sec7.Length())

		sec8, err := g.ReadSectionAt(1_427_755)
		require.NoError(t, err)
		require.Equal(t, 8, sec8.Number())
		require.Equal(t, 4, sec8.Length())

		msg1, err := g.ReadMessageAt(303462693)
		require.NoError(t, err)
		require.NotNil(t, msg1)

		offset := msg1.GetOffset() + msg1.GetSize()

		{
			sec0, err := g.ReadSectionAt(offset)
			require.NoError(t, err)
			require.Equal(t, 0, sec0.Number())
			require.Equal(t, 16, sec0.Length())

			sec1, err := g.ReadSectionAt(offset + 16)
			require.NoError(t, err)
			require.Equal(t, 1, sec1.Number())
			require.Equal(t, 21, sec1.Length())

			sec2, err := g.ReadSectionAt(offset + 16 + 21)
			require.NoError(t, err)
			require.Equal(t, 3, sec2.Number())
			require.Equal(t, 72, sec2.Length())

			sec4, err := g.ReadSectionAt(offset + 16 + 21 + 72)
			require.NoError(t, err)
			require.Equal(t, 4, sec4.Number())
			require.Equal(t, 34, sec4.Length())

			sec5, err := g.ReadSectionAt(offset + 16 + 21 + 72 + 34)
			require.NoError(t, err)
			require.Equal(t, 5, sec5.Number())
			require.Equal(t, 21, sec5.Length())

			sec6, err := g.ReadSectionAt(offset + 16 + 21 + 72 + 34 + 21)
			require.NoError(t, err)
			require.Equal(t, 6, sec6.Number())
			require.Equal(t, 6, sec6.Length())

			sec7, err := g.ReadSectionAt(offset + 16 + 21 + 72 + 34 + 21 + 6)
			require.NoError(t, err)
			require.Equal(t, 7, sec7.Number())
			require.Equal(t, 1_427_585, sec7.Length())

			sec8, err := g.ReadSectionAt(offset + 16 + 21 + 72 + 34 + 21 + 6 + 1_427_585)
			require.NoError(t, err)
			require.Equal(t, 8, sec8.Number())
			require.Equal(t, 4, sec8.Length())

			_, err = g.ReadSectionAt(offset + 16 + 21 + 72 + 34 + 21 + 6 + 1_427_585 + 1)
			require.ErrorIs(t, err, io.EOF)
		}

		{
			lastMsg, err := g.ReadMessageAt(offset)
			require.NoError(t, err)
			require.NotNil(t, lastMsg)

			_, err = g.ReadMessageAt(lastMsg.GetOffset() + lastMsg.GetSize())
			require.ErrorIs(t, err, io.EOF)
		}
	})
}

func TestGrib2_EachMessage(t *testing.T) {
	t.Parallel()

	t.Run("file", func(t *testing.T) {
		t.Parallel()

		const filename = "../testdata/cwat.grib2"

		f, err := os.Open(filename)
		if errors.Is(err, os.ErrNotExist) {
			t.Skipf("%s not exist", filename)
		}

		require.NoError(t, err)
		defer f.Close()

		g := grib.NewGrib2(f)

		count := 0
		assert.NoError(t, g.EachMessage(func(msg grib2.IndexedMessage) (bool, error) {
			count++
			require.NotNil(t, msg)
			return true, nil
		}))
		assert.Equal(t, 1, count)
	})

	t.Run("mmap", func(t *testing.T) {
		t.Parallel()

		const filename = "../testdata/cwat.grib2"

		f, err := mmap.Open(filename)
		if errors.Is(err, os.ErrNotExist) {
			t.Skipf("%s not exist", filename)
		}

		require.NoError(t, err)
		defer f.Close()

		g := grib.NewGrib2(f)

		count := 0
		assert.NoError(t, g.EachMessage(func(msg grib2.IndexedMessage) (bool, error) {
			count++
			require.NotNil(t, msg)
			return true, nil
		}))
		assert.Equal(t, 1, count)
	})

	t.Run("oss", func(t *testing.T) {
		t.Parallel()

		const (
			bucketName = "cy-meteorology"
			key        = "noaa-gfs/develop/2024/09/30/18/atmos/0p25/2t_heightAboveGround_2_0_0.grib2"
			msgOffset  = 0
		)

		var (
			endpoint        = os.Getenv("ALIYUN_OSS_ENDPOINT")
			accessKeyId     = os.Getenv("ALIYUN_OSS_ACCESS_KEY_ID")
			accessKeySecret = os.Getenv("ALIYUN_OSS_ACCESS_KEY_SECRET")
		)

		ctx := context.TODO()
		cli, err := oss.New(
			endpoint,
			accessKeyId,
			accessKeySecret,
		)
		if err != nil {
			t.Skip(err.Error())
		}

		bucket, err := cli.Bucket(bucketName)
		require.NoError(t, err)

		r, err := ossio.NewReader(ctx, bucket, key)
		if err != nil {
			t.Skip(err.Error())
		}

		g := grib.NewGrib2(r)

		count := 0
		assert.NoError(t, g.EachMessage(func(msg grib2.IndexedMessage) (bool, error) {
			count++
			require.NotNil(t, msg)
			return true, nil
		}))
		assert.Equal(t, 209, count)
	})
}

func TestGrib2_ReadImage(t *testing.T) {
	t.Parallel()

	const filename = "../testdata/grid_png.grib2"

	p := &gridpoint.PortableNetworkGraphics{
		ReferenceValue:     -2023.1235,
		BinaryScaleFactor:  1,
		DecimalScaleFactor: 2,
		Bits:               12,
		NumVals:            65160,
	}

	t.Run("read all data", func(t *testing.T) {
		f, err := os.Open(filename)
		if errors.Is(err, os.ErrNotExist) {
			t.Skipf("%s not exist", filename)
		}
		defer f.Close()

		r := io.NewSectionReader(f, 175, 78200)

		vals, err := p.ReadAllData(bitio.NewReader(r))
		require.NoError(t, err)
		require.Equal(t, 65160, len(vals))
		assert.InDelta(t, 1.18876, vals[0], 1e-5)
	})

	t.Run("read image", func(t *testing.T) {
		f, err := os.Open(filename)
		if errors.Is(err, os.ErrNotExist) {
			t.Skipf("%s not exist", filename)
		}
		defer f.Close()

		r := io.NewSectionReader(f, 175, 78200)
		img, err := p.Image(bitio.NewReader(r))
		require.NoError(t, err)
		require.NotNil(t, img)

		tmp, err := os.CreateTemp(os.TempDir(), "img.png")
		require.NoError(t, err)
		defer os.Remove(tmp.Name())
		defer tmp.Close()

		assert.Equal(t, 360, img.Bounds().Max.X)
		assert.Equal(t, 181, img.Bounds().Max.Y)

		require.NoError(t, png.Encode(tmp, img))
	})

	// t.Run("read grid image", func(t *testing.T) {
	// 	f, err := os.Open(filename)
	// 	if errors.Is(err, os.ErrNotExist) {
	// 		t.Skipf("%s not exist", filename)
	// 	}
	// 	defer f.Close()

	// 	pr := gridpoint.NewPortableNetworkGraphicsReader(f, 175, 78200, p)

	// 	img, err := pr.ImageGridAt(0)
	// 	require.NoError(t, err)
	// 	require.NotNil(t, img)
	// })
}
