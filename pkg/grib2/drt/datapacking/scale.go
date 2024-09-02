package datapacking

import (
	"math"
)

func BinaryScaleFactor(e int16) float64 {
	return math.Pow(2, float64(e))
}

func DecimalScaleFactor(d int16) float64 {
	return math.Pow10(int(d))
}

// The original data value Y (in the units of code table 4.2) can be recovered with the formula:
//
// Y * 10D= R + (X1+X2) * 2E
//
// For simple packing and all spectral data
// E = Binary scale factor,
// D = Decimal scale factor
// R = Reference value of the whole field,
// X1 = 0,
// X2 = Scaled (encoded) value.
//
// # For complex grid point packing schemes, E, D, and R are as above, but
//
// X1 = Reference value (scaled integer) of the group the data value belongs to,
// X2 = Scaled (encoded) value with the group reference value (XI) removed..
func SimpleScaleFunc(e int16, d int16, r float32) func(x2 uint32) float64 {
	var (
		dec   = DecimalScaleFactor(d)
		ref   = float64(r) / dec
		scale = BinaryScaleFactor(e) / dec
	)

	return func(x2 uint32) float64 {
		return ref + float64(x2)*scale
	}
}
