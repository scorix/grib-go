package gdt

import (
	"math"

	"github.com/scorix/grib-go/pkg/grib2/regulation"
	"github.com/scorix/walg/pkg/geo/grids"
	"github.com/scorix/walg/pkg/geo/grids/latlon"
)

/*
Notes:
1.  Basic angle of the initial production domain and subdivisions of this basic angle are provided to manage cases where the recommended unit of 10-6 degrees is not applicable to describe the extreme longitudes and latitudes, and direction increments. For these last six descriptors, the unit is equal to the ratio of the basic angle and the subdivisions number. For ordinary cases, zero and missing values should be coded, equivalent to respective values of 1 and 106  (10-6  degrees unit).

2.  For data on a quasi-regular grid, in which all the rows or columns do not necessarily have the same number of grid points either Ni (octets 31-34) of Nj (octets 35-38) and the corresponding Di (octets 64-67) or Dj (octets 68-71) shall be coded with all bits set to 1 (missing). The actual number of points along each parallel or meridian shall be coded in the octets immediately following the grid definition template (octets [xx+1]-nn), as described in the description of the grid definition section.

3.  A quasi-regular grid is only defined for appropriate grid scanning modes. Either rows or columns, but not both simultaneously, may have variable numbers of points or variable spacing. The first point in each row (column) shall be positioned at the meridian (parallel) indicted by octets 47-54. The grid points shall be evenly spaced in latitude (longitude).

4.  A scale value of radius of spherical Earth, or major axis of oblate spheroid Earth is delivered from applying appropriate scale factor to the value expressed in meters.

5.  It is recommended to use unsigned direction increments.

6.  In most cases, multiplying Ni (octets 31-34) by Nj (octets 35-38) yields the total number of points
in the grid. However, this may not be true if bit 8 of the scanning mode flags (octet 72) is set to 1.
*/
type Template0 struct {
	Template0FixedPart `json:"template0"`
	// TODO: 73-nn: List of number of points along each meridian or parallel (These octets are only present for quasi-regular grids as described in notes 2 and 3)
	grids grids.Grid `json:"-"`
}

type template0FixedPart struct {
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
	JDirectionIncrement                    uint32
	ScanningMode                           uint8
}

func (t template0FixedPart) Export() Template {
	t0 := Template0FixedPart{
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
		JDirectionIncrement:                    regulation.ToInt32(t.JDirectionIncrement),
		ScanningMode:                           regulation.ToInt8(t.ScanningMode),
	}

	return t0.AsTemplate()
}

type Template0FixedPart struct {
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
	LatitudeOfFirstGridPoint               int32 `json:"latitudeOfFirstGridPoint"`
	LongitudeOfFirstGridPoint              int32 `json:"longitudeOfFirstGridPoint"`
	ResolutionAndComponentFlags            int8  `json:"-"`
	LatitudeOfLastGridPoint                int32 `json:"latitudeOfLastGridPoint"`
	LongitudeOfLastGridPoint               int32 `json:"longitudeOfLastGridPoint"`
	IDirectionIncrement                    int32 `json:"iDirectionIncrement"`
	JDirectionIncrement                    int32 `json:"jDirectionIncrement"`
	ScanningMode                           int8  `json:"scanningMode"`
}

func (t *Template0FixedPart) AsTemplate() Template {
	firstLat := float64(t.LatitudeOfFirstGridPoint) / 1e6
	lastLat := float64(t.LatitudeOfLastGridPoint) / 1e6
	firstLon := float64(t.LongitudeOfFirstGridPoint) / 1e6
	lastLon := float64(t.LongitudeOfLastGridPoint) / 1e6
	minLat := math.Min(firstLat, lastLat)
	maxLat := math.Max(firstLat, lastLat)
	minLon := math.Min(firstLon, lastLon)
	maxLon := math.Max(firstLon, lastLon)

	return &Template0{
		Template0FixedPart: *t,
		grids: latlon.NewLatLonGrid(
			minLat,
			maxLat,
			minLon,
			maxLon,
			float64(t.IDirectionIncrement)/1e6,
			float64(t.JDirectionIncrement)/1e6,
			latlon.WithScanMode(grids.ScanMode(t.ScanningMode)),
		),
	}
}

func (t *Template0FixedPart) GetNi() int32 {
	return t.Ni
}

func (t *Template0FixedPart) GetNj() int32 {
	return t.Nj
}

func (t *Template0) GetGridIndex(lat, lon float32) (n int) {
	return grids.GridIndex(t.grids, float64(lat), float64(lon))
}

func (t *Template0) GetGridPoint(n int) (float32, float32) {
	lat, lon := grids.GridPoint(t.grids, n)
	return float32(lat), float32(lon)
}
