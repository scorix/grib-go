package gridpoint

import (
	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
	"github.com/scorix/grib-go/pkg/grib2/drt/definition"
	"github.com/scorix/grib-go/pkg/grib2/regulation"
)

type ComplexPacking struct {
	*SimplePacking

	GroupMethod                int8
	MissingValue               int8
	PrimaryMissingSubstitute   int32
	SecondaryMissingSubstitute int32
	NumberOfGroups             int32
	GroupWidths                int8
	GroupWidthsBits            int8
	GroupLengthsReference      int32
	GroupLengthIncrement       int8
	GroupLastLength            int32
	GroupScaledLengthsBits     int8
}

func NewComplexPacking(def definition.ComplexPacking) *ComplexPacking {
	return &ComplexPacking{
		SimplePacking:              NewSimplePacking(def.SimplePacking),
		GroupMethod:                regulation.ToInt8(def.GroupMethod),
		MissingValue:               regulation.ToInt8(def.MissingValue),
		PrimaryMissingSubstitute:   regulation.ToInt32(def.PrimaryMissingSubstitute),
		SecondaryMissingSubstitute: regulation.ToInt32(def.SecondaryMissingSubstitute),
		NumberOfGroups:             regulation.ToInt32(def.NumberOfGroups),
		GroupWidths:                regulation.ToInt8(def.GroupWidths),
		GroupWidthsBits:            regulation.ToInt8(def.GroupWidthsBits),
		GroupLengthsReference:      regulation.ToInt32(def.GroupLengthsReference),
		GroupLengthIncrement:       regulation.ToInt8(def.GroupLengthIncrement),
		GroupLastLength:            regulation.ToInt32(def.GroupLastLength),
		GroupScaledLengthsBits:     regulation.ToInt8(def.GroupScaledLengthsBits),
	}
}

func (sp *ComplexPacking) ScaleFunc() func(uint32) float64 {
	return datapacking.SimpleScaleFunc(sp.E, sp.D, sp.R)
}
