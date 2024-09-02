package gridpoint

import (
	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
	"github.com/scorix/grib-go/pkg/grib2/drt/definition"
	"github.com/scorix/grib-go/pkg/grib2/regulation"
)

type SimplePacking struct {
	R    float32
	E    int16
	D    int16
	Bits uint8
	Type int8
}

func NewSimplePacking(def definition.SimplePacking) *SimplePacking {
	return &SimplePacking{
		R:    def.R,
		E:    regulation.ToInt16(def.E),
		D:    regulation.ToInt16(def.D),
		Bits: def.Bits,
		Type: regulation.ToInt8(def.Type),
	}
}

func (sp *SimplePacking) ScaleFunc() func(uint32) float64 {
	return datapacking.SimpleScaleFunc(sp.E, sp.D, sp.R)
}

func (sp *SimplePacking) GetBits() uint8 {
	return sp.Bits
}
