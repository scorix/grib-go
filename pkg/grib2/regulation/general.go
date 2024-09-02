package regulation

import "math"

/*
92.1 General

92.1.1 The GRIB code shall be used for the exchange and storage of general regularly-distributed information expressed in binary form.

92.1.2 The beginning and the end of the code shall be identified by 4 octets coded according to the International Alphabet No. 5 to represent the indicators "GRIB" and "7777" in Indicator Section 0 and End Section 8, respectively. All other octets included in the code shall represent data in binary form.

92.1.3 Each section included in the code shall always end on an octet boundary. This rule shall be applied by appending bits set to zero to the section where necessary.

92.1.4 All bits set to “1" for any value indicates that value is missing. This rule shall not apply to packed data.

92.1.5 If applicable, negative values shall be indicated by setting the most significant bit to “1”.

92.1.6 Latitude, longitude, and angle values shall be in units of 10-6 degree, except for specific cases explicitly stated in some grid definitions.

92.1.7 The latitude values shall be limited to the range 0 to 90 degrees inclusive. Orientation shall be north latitude positive, south latitude negative. Bit 1 is set to 1 to indicate south latitude.

92.1.8 The longitude values shall be limited to the range 0 to 360 degrees inclusive. Orientation shall be east longitude positive, with only positive values being used.

92.1.9 The latitude and longitude of the first grid point and the last grid point shall always be given for regular grids.

92.1.10 Vector components at the North and South Poles shall be coded according to the following conventions.

92.1.10.1 If the resolution and component flags in section 3 (Flag table 3.3) indicate that the vector components are relative to the defined grid, the vector components at the Pole shall be resolved relative to the grid.

92.1.10.2 Otherwise, for projections where there are multiple points at a given pole, the vector components shall be resolved as if measured an infinitesimal distance from the Pole at the longitude corresponding to each grid point. At the North Pole, the West to East (x direction) component at a grid point with longitude L shall be resolved along the meridian 90 degrees East of L, and the South to North (y direction) component shall be resolved along the meridian 180 degrees from L. At the South Pole the West to East component at a grid point with longitude L shall be resolved along the meridian 90 degrees East of L and the South to North component shall be resolved along L.

92.1.10.3 Otherwise, if there is only one Pole point, either on a cylindrical projection with all but one Pole point deleted, or on any projection (such as polar stereographic) where the Pole maps to a unique point, the West to East and South to North components shall be resolved along longitudes 270 and 0 respectively at the North Pole and along longitudes 270 and 180 respectively at the South Pole.

Note: (1) This differs from the treatment of the Poles in the WMO traditional alphanumeric codes.

92.1.11 The first and last grid points shall not necessarily correspond to the first and last data points, respectively, if the bit-map is used.

92.1.12 Items in sections 3 and 4 which consist of a scale factor F and a scaled value V are related to the original value L as follows:

L * 10F = V
*/

// 92.1.5
func ToInt8(v uint8) int8 {
	if IsMissingValue(v) {
		return int8(v)
	}

	i := int8(v & 0x7f)

	if negtive := v&0x80 > 0; negtive {
		return -i
	}

	return i
}

func ToInt16(v uint16) int16 {
	if IsMissingValue(v) {
		return int16(v)
	}

	i := int16(v & 0x7fff)

	if negtive := v&0x8000 > 0; negtive {
		return -i
	}

	return i
}

func ToInt32(v uint32) int32 {
	if IsMissingValue(v) {
		return int32(v)
	}

	i := int32(v & 0x7fffffff)

	if negtive := v&0x80000000 > 0; negtive {
		return -i
	}

	return i
}

// 92.1.4
func IsMissingValue(value any) bool {
	switch v := value.(type) {
	case int8:
		return uint8(v) == math.MaxUint8
	case uint8:
		return v == math.MaxUint8
	case int16:
		return uint16(v) == math.MaxUint16
	case uint16:
		return v == math.MaxUint16
	case int32:
		return uint32(v) == math.MaxUint32
	case uint32:
		return v == math.MaxUint32
	case int64:
		return uint64(v) == math.MaxUint64
	case uint64:
		return v == math.MaxUint64
	}

	return false
}

// 92.1.6
func DegreedLatitudeLongitude(v float64) float64 {
	return v / float64(10e6)
}