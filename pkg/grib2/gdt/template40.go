package gdt

import (
	"github.com/scorix/grib-go/pkg/grib2/regulation"
	"github.com/scorix/walg/pkg/geo/grids"
	"github.com/scorix/walg/pkg/geo/grids/gaussian"
)

/*
Notes:
( 1) Basic angle of the initial production domain and subdivisions of this basic angle are provided to manage cases where the recommended unit of 10-6 degrees is not applicable to describe the extreme longitudes and latitudes, and direction increments. For these last six descriptors, unit is equal to the ratio of the basic angle and the subdivisions number. For ordinary cases, zero and missing values should be coded, equivalent to respective values of 1 and 106 (10-6 degrees unit).

( 2) The number of parallels between a pole and the equator is used to establish the variable (Gaussian) spacing of the parallels; this value must always be given.

( 3) A scaled value of radius of spherical Earth, or major or minor axis of oblate spheroid Earth is derived from applying appropriate scale factor to the value expressed in metres.

( 4) A quasi-regular grid is only defined for appropriate grid scanning modes. Either rows or columns, but not both simultaneously, may have variable numbers of points. The first point in each row (column) shall be positioned at the meridian (parallel) indicated by Octets 47-54. The grid points shall be evenly spaced in latitude (longitude).

( 5) It is recommended to use unsigned direction increments
*/

type Template40 struct {
	Template40FixedPart `json:"template40"`
	// TODO: 73-nn: List of number of points along each meridian or parallel (These octets are only present for quasi-regular grids as described in note 4)
	grids grids.Grid `json:"-"`
}

// https://codes.ecmwf.int/grib/format/grib2/templates/3/40/
type template40FixedPart struct {
	ShapeOfTheEarth                        uint8
	ScaleFactorOfRadiusOfSphericalEarth    uint8
	ScaledValueOfRadiusOfSphericalEarth    uint32
	ScaleFactorOfEarthMajorAxis            uint8
	ScaledValueOfEarthMajorAxis            uint32
	ScaleFactorOfEarthMinorAxis            uint8
	ScaledValueOfEarthMinorAxis            uint32
	Ni                                     uint32
	Nj                                     uint32
	BasicAngleOfTheInitialProductionDomain uint32
	SubdivisionsOfBasicAngle               uint32
	LatitudeOfFirstGridPoint               uint32
	LongitudeOfFirstGridPoint              uint32
	ResolutionAndComponentFlags            uint8
	LatitudeOfLastGridPoint                uint32
	LongitudeOfLastGridPoint               uint32
	IDirectionIncrement                    uint32
	N                                      uint32
	ScanningMode                           uint8
}

func (t template40FixedPart) Export() Template {
	t40 := Template40FixedPart{
		ShapeOfTheEarth:                        regulation.ToInt8(t.ShapeOfTheEarth),
		ScaleFactorOfRadiusOfSphericalEarth:    regulation.ToInt8(t.ScaleFactorOfRadiusOfSphericalEarth),
		ScaledValueOfRadiusOfSphericalEarth:    regulation.ToInt32(t.ScaledValueOfRadiusOfSphericalEarth),
		ScaleFactorOfEarthMajorAxis:            regulation.ToInt8(t.ScaleFactorOfEarthMajorAxis),
		ScaledValueOfEarthMajorAxis:            regulation.ToInt32(t.ScaledValueOfEarthMajorAxis),
		ScaleFactorOfEarthMinorAxis:            regulation.ToInt8(t.ScaleFactorOfEarthMinorAxis),
		ScaledValueOfEarthMinorAxis:            regulation.ToInt32(t.ScaledValueOfEarthMinorAxis),
		Ni:                                     regulation.ToInt32(t.Ni),
		Nj:                                     regulation.ToInt32(t.Nj),
		BasicAngleOfTheInitialProductionDomain: regulation.ToInt32(t.BasicAngleOfTheInitialProductionDomain),
		SubdivisionsOfBasicAngle:               regulation.ToInt32(t.SubdivisionsOfBasicAngle),
		LatitudeOfFirstGridPoint:               regulation.ToInt32(t.LatitudeOfFirstGridPoint),
		LongitudeOfFirstGridPoint:              regulation.ToInt32(t.LongitudeOfFirstGridPoint),
		ResolutionAndComponentFlags:            regulation.ToInt8(t.ResolutionAndComponentFlags),
		LatitudeOfLastGridPoint:                regulation.ToInt32(t.LatitudeOfLastGridPoint),
		LongitudeOfLastGridPoint:               regulation.ToInt32(t.LongitudeOfLastGridPoint),
		IDirectionIncrement:                    regulation.ToInt32(t.IDirectionIncrement),
		N:                                      regulation.ToInt32(t.N),
		ScanningMode:                           regulation.ToInt8(t.ScanningMode),
	}

	return t40.AsTemplate()
}

type Template40FixedPart struct {
	ShapeOfTheEarth                        int8  `json:"-"`
	ScaleFactorOfRadiusOfSphericalEarth    int8  `json:"-"`
	ScaledValueOfRadiusOfSphericalEarth    int32 `json:"-"`
	ScaleFactorOfEarthMajorAxis            int8  `json:"-"`
	ScaledValueOfEarthMajorAxis            int32 `json:"-"`
	ScaleFactorOfEarthMinorAxis            int8  `json:"-"`
	ScaledValueOfEarthMinorAxis            int32 `json:"-"`
	Ni                                     int32 `json:"-"`
	Nj                                     int32 `json:"-"`
	BasicAngleOfTheInitialProductionDomain int32 `json:"-"`
	SubdivisionsOfBasicAngle               int32 `json:"-"`
	LatitudeOfFirstGridPoint               int32 `json:"-"`
	LongitudeOfFirstGridPoint              int32 `json:"-"`
	ResolutionAndComponentFlags            int8  `json:"-"`
	LatitudeOfLastGridPoint                int32 `json:"-"`
	LongitudeOfLastGridPoint               int32 `json:"-"`
	IDirectionIncrement                    int32 `json:"-"`
	N                                      int32 `json:"n"`
	ScanningMode                           int8  `json:"scanningMode"`
}

func (t *Template40FixedPart) AsTemplate() Template {
	return &Template40{
		Template40FixedPart: *t,
		grids: gaussian.NewRegular(
			int(t.N),
		),
	}
}

func (t *Template40FixedPart) GetNi() int32 {
	return t.Ni
}

func (t *Template40FixedPart) GetNj() int32 {
	return t.Nj
}

func (t *Template40) GetGridIndex(lat, lon float32) (n int) {
	return grids.GuessGridIndex(t.grids, float64(lat), float64(lon), grids.ScanMode(t.ScanningMode))
}

func (t *Template40) GetGridPoint(n int) (float32, float32) {
	lat, lon := grids.GridPoint(t.grids, n, grids.ScanMode(t.ScanningMode))
	return float32(lat), float32(lon)
}
