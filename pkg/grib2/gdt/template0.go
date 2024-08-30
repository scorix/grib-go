package gdt

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
	// TODO: 73-nn: List of number of points along each meridian or parallel (These octets are only present for quasi-regular grids as described in notes 2 and 3)
}
