package gdt

import (
	"fmt"

	"github.com/scorix/grib-go/pkg/grib2/regulation"
)

type Template40 struct {
	Template40FixedPart
	// TODO: 73-nn: List of number of points along each meridian or parallel (These octets are only present for quasi-regular grids as described in notes 2 and 3)
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

func (t template40FixedPart) Export() Template40FixedPart {
	return Template40FixedPart{
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
}

type Template40FixedPart struct {
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
	N                                      int32
	ScanningMode                           int8
}

func (t *Template40FixedPart) GetScanningMode() (ScanningMode, error) {
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
			N:                           t.N,
			getGridIndexFunc:            t.GetGridIndex,
			getGridPointByIndexFunc:     t.GetGridPointByIndex,
		}

		return &sm, nil
	}

	return nil, fmt.Errorf("scanning mode %04b is not implemented", t.ScanningMode)
}

func (t *Template40FixedPart) GetNi() int32 {
	return t.Ni
}

func (t *Template40FixedPart) GetNj() int32 {
	return t.Nj
}

func (t *Template40FixedPart) GetGridIndex(lat, lon float32) (i, j, n int) {
	return GetRegularGGGridIndex(lat, lon, t.LatitudeOfFirstGridPoint, t.LongitudeOfFirstGridPoint, t.LatitudeOfLastGridPoint, t.LongitudeOfLastGridPoint, t.N, t.Ni)
}

func (t *Template40FixedPart) GetGridPointByIndex(i, j int) (lat, lon float32) {
	return GetRegularGGGridPointByIndex(i, j, t.LatitudeOfFirstGridPoint, t.LongitudeOfFirstGridPoint, t.LatitudeOfLastGridPoint, t.LongitudeOfLastGridPoint, t.N, t.Ni)
}
