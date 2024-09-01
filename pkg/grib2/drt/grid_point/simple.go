package gridpoint

import (
	"math"

	"github.com/scorix/grib-go/pkg/grib2/scale"
)

type SimplePacking struct {
	*DefSimplePacking
}

func NewSimplePacking(def DefSimplePacking) *SimplePacking {
	return &SimplePacking{
		DefSimplePacking: &def,
	}
}

type DefSimplePacking struct {
	R    float32      // Reference value (R) (IEEE 32-bit floating-point value)
	E    scale.Factor // Binary scale factor
	D    scale.Factor // Decimal scale factor
	Bits uint8        // Number of bits used for each packed value for simple packing, or for each group reference value for complex packing or spatial differencing
	Type uint8        // Type of original field values: https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table5-1.shtml
}

func (def DefSimplePacking) BinaryScaleFactor() float64 {
	return math.Pow(2, float64(def.E.Int16()))
}

func (def DefSimplePacking) DecimalScaleFactor() float64 {
	return math.Pow10(-int(def.D.Int16()))
}

func (def DefSimplePacking) ReferenceValue() float64 {
	return def.DecimalScaleFactor() * float64(def.R)
}

func (def DefSimplePacking) ScaleFactor() float64 {
	return def.BinaryScaleFactor() * def.DecimalScaleFactor()
}

func (def DefSimplePacking) ScaleFunc() func(uint64) float64 {
	ref, scale := def.ReferenceValue(), def.ScaleFactor()

	return func(v uint64) float64 {
		return SimpleScale(v, ref, scale)
	}
}

func (def DefSimplePacking) GetBits() uint8 {
	return def.Bits
}

func SimpleScale(v uint64, ref float64, scale float64) float64 {
	return ref + float64(v)*scale
}
