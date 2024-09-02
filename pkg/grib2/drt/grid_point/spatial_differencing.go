package gridpoint

import (
	"github.com/scorix/grib-go/pkg/grib2/drt/definition"
	"github.com/scorix/grib-go/pkg/grib2/regulation"
)

type ComplexPackingAndSpatialDifferencing struct {
	*ComplexPacking

	SpatialOrderDifference int8
	OctetsNumber           int8
}

func NewComplexPackingAndSpatialDifferencing(def definition.ComplexPackingAndSpatialDifferencing, numVals int) *ComplexPackingAndSpatialDifferencing {
	return &ComplexPackingAndSpatialDifferencing{
		ComplexPacking:         NewComplexPacking(def.ComplexPacking, numVals),
		SpatialOrderDifference: regulation.ToInt8(def.SpatialOrderDifference),
		OctetsNumber:           regulation.ToInt8(def.OctetsNumber),
	}
}
