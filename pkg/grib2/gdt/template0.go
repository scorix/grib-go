package gdt

import (
	"fmt"
	"math"

	"github.com/scorix/grib-go/pkg/grib2/regulation"
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
	Template0FixedPart
	// TODO: 73-nn: List of number of points along each meridian or parallel (These octets are only present for quasi-regular grids as described in notes 2 and 3)
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

func (t template0FixedPart) Export() Template0FixedPart {
	return Template0FixedPart{
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
}

type Template0FixedPart struct {
	ShapeOfTheEarth                        int8
	ScaleFactorOfRadiusOfSphericalEarth    int8
	ScaledValueOfRadiusOfSphericalEarth    int32
	ScaleFactorOfEarthMajorAxis            int8
	ScaledValueOfEarthMajorAxis            int32
	ScaleFactorOfEarthMinorAxis            int8
	ScaledValueOfEarthMinorAxis            int32
	Ni                                     int32
	Nj                                     int32
	BasicAngleOfTheInitialProductionDomain int32
	SubdivisionsOfBasicAngle               int32
	LatitudeOfFirstGridPoint               int32
	LongitudeOfFirstGridPoint              int32
	ResolutionAndComponentFlags            int8
	LatitudeOfLastGridPoint                int32
	LongitudeOfLastGridPoint               int32
	IDirectionIncrement                    int32
	JDirectionIncrement                    int32
	ScanningMode                           int8
}

type ScanningMode interface {
	GetGridPointLL(n int) (float32, float32)
	GetGridPointFromLL(lat float32, lng float32) int
}

func (t *Template0FixedPart) GetScanningMode() (ScanningMode, error) {
	switch t.ScanningMode {
	case 0:
		sm := ScanningMode0000{
			Ni:                          t.Ni,
			Nj:                          t.Nj,
			LatitudeOfFirstGridPoint:    t.LatitudeOfFirstGridPoint,
			LongitudeOfFirstGridPoint:   t.LongitudeOfFirstGridPoint,
			ResolutionAndComponentFlags: t.ResolutionAndComponentFlags,
			LatitudeOfLastGridPoint:     t.LatitudeOfLastGridPoint,
			LongitudeOfLastGridPoint:    t.LongitudeOfLastGridPoint,
			IDirectionIncrement:         t.IDirectionIncrement,
			JDirectionIncrement:         t.JDirectionIncrement,
		}

		return &sm, nil
	}

	return nil, fmt.Errorf("scanning mode %04b is not implemented", t.ScanningMode)
}

func (t *Template0FixedPart) GetNi() int32 {
	return t.Ni
}

func (t *Template0FixedPart) GetNj() int32 {
	return t.Nj
}

type ScanningMode0000 struct {
	Ni                          int32
	Nj                          int32
	LatitudeOfFirstGridPoint    int32
	LongitudeOfFirstGridPoint   int32
	ResolutionAndComponentFlags int8
	LatitudeOfLastGridPoint     int32
	LongitudeOfLastGridPoint    int32
	IDirectionIncrement         int32
	JDirectionIncrement         int32
}

func (sm *ScanningMode0000) GetLatitudeGridPoint(n int) int {
	latFirst, latLast, inc := sm.LatitudeOfFirstGridPoint, sm.LatitudeOfLastGridPoint, sm.IDirectionIncrement
	if (latFirst-latLast)/sm.IDirectionIncrement < 0 {
		inc = -inc
	}

	return int(latFirst) - (n/int(sm.Ni))*int(inc)
}

func (sm *ScanningMode0000) GetLatitudeGridIndex(lat float32) int {
	latFirst, latLast, inc := sm.LatitudeOfFirstGridPoint, sm.LatitudeOfLastGridPoint, sm.IDirectionIncrement
	if (latFirst-latLast)/sm.IDirectionIncrement < 0 {
		inc = -inc
	}

	return (int(latFirst) - toInt(lat)) / int(inc)
}

func (sm *ScanningMode0000) GetLongitudeGridPoint(n int) int {
	lonFirst, lonLast, inc := sm.LongitudeOfFirstGridPoint, sm.LongitudeOfLastGridPoint, sm.JDirectionIncrement
	if (lonFirst-lonLast)/sm.IDirectionIncrement < 0 {
		inc = -inc
	}

	return int(lonFirst) - (n%int(sm.Ni))*int(inc)
}

func (sm *ScanningMode0000) GetLongitudeGridIndex(lng float32) int {
	lonFirst, lonLast, inc := sm.LongitudeOfFirstGridPoint, sm.LongitudeOfLastGridPoint, sm.JDirectionIncrement
	if (lonFirst-lonLast)/sm.IDirectionIncrement < 0 {
		inc = -inc
	}

	return (int(lonFirst) - toInt(lng)) / int(inc)
}

func (sm *ScanningMode0000) GetGridPointLL(n int) (float32, float32) {
	return regulation.DegreedLatitudeLongitude(sm.GetLatitudeGridPoint(n)), regulation.DegreedLatitudeLongitude(sm.GetLongitudeGridPoint(n))
}

func (sm *ScanningMode0000) GetGridPointFromLL(lat float32, lng float32) int {
	return sm.GetLatitudeGridIndex(lat)*int(sm.Ni) + sm.GetLongitudeGridIndex(lng)
}

func toInt(v float32) int {
	i := math.Floor(float64(v) * 1e6)

	return int(i)
}
