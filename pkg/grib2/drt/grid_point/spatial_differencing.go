package gridpoint

import (
	"encoding/json"
	"fmt"

	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
	"github.com/scorix/grib-go/pkg/grib2/drt/definition"
	"github.com/scorix/grib-go/pkg/grib2/regulation"
)

type ComplexPackingAndSpatialDifferencing struct {
	*ComplexPacking              // 12-47
	SpatialOrderDifference int8  // 48
	OctetsNumber           uint8 // 49
}

func NewComplexPackingAndSpatialDifferencing(def definition.ComplexPackingAndSpatialDifferencing, numVals int) *ComplexPackingAndSpatialDifferencing {
	return &ComplexPackingAndSpatialDifferencing{
		ComplexPacking:         NewComplexPacking(def.ComplexPacking, numVals),
		SpatialOrderDifference: regulation.ToInt8(def.SpatialOrderDifference),
		OctetsNumber:           def.OctetsNumber,
	}
}

func (cpsd *ComplexPackingAndSpatialDifferencing) ReadAllData(r datapacking.BitReader) ([]float64, error) {
	sd, err := cpsd.ReadSpacingDifferential(r)
	if err != nil {
		return nil, fmt.Errorf("read spacing differential value: %w", err)
	}

	groups, err := cpsd.ReadGroups(r, cpsd.Bits)
	if err != nil {
		return nil, fmt.Errorf("read groups: %w", err)
	}

	if len(groups) != int(cpsd.NumberOfGroups) {
		return nil, fmt.Errorf("expected groups: %d, got %d", cpsd.NumberOfGroups, len(groups))
	}

	data, err := cpsd.unpackData(r, groups, func(data, miss []uint32, primary, secondary float64, scaleFunc func(uint32) float64) ([]float64, error) {
		sd.Apply(data)

		return cpsd.scaleValues(data, miss, primary, secondary, cpsd.ScaleFunc())
	})
	if err != nil {
		return nil, fmt.Errorf("unpack data: %w", err)
	}

	return data, nil
}

type complexPackingAndSpatialDifferencing struct {
	ComplexPacking         *ComplexPacking `json:"complex_packing"`
	SpatialOrderDifference int8            `json:"spatial_order_difference"`
	OctetsNumber           uint8           `json:"octets_number"`
	NumVals                int             `json:"num_vals"`
}

func (cpsd *ComplexPackingAndSpatialDifferencing) MarshalJSON() ([]byte, error) {
	return json.Marshal(complexPackingAndSpatialDifferencing{
		ComplexPacking:         cpsd.ComplexPacking,
		SpatialOrderDifference: cpsd.SpatialOrderDifference,
		OctetsNumber:           cpsd.OctetsNumber,
		NumVals:                cpsd.NumVals,
	})
}

func (cpsd *ComplexPackingAndSpatialDifferencing) UnmarshalJSON(data []byte) error {
	var temp complexPackingAndSpatialDifferencing

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	cpsd.ComplexPacking = temp.ComplexPacking
	cpsd.SpatialOrderDifference = temp.SpatialOrderDifference
	cpsd.OctetsNumber = temp.OctetsNumber
	cpsd.NumVals = temp.NumVals
	return nil
}

func (cpsd *ComplexPackingAndSpatialDifferencing) Definition() any {
	return definition.ComplexPackingAndSpatialDifferencing{
		ComplexPacking:         cpsd.ComplexPacking.Definition().(definition.ComplexPacking),
		SpatialOrderDifference: regulation.ToUint8(cpsd.SpatialOrderDifference),
		OctetsNumber:           cpsd.OctetsNumber,
	}
}

type spacingDifferential struct {
	vals []uint32
	min  uint32
}

// Spatial differencing is a pre-processing before group splitting at encoding time.
// It is intended to reduce the size of sufficiently smooth fields, when combined with
// a splitting scheme as described in Data Representation Template 5.2.
func (sd *spacingDifferential) Apply(data []uint32) {
	copy(data[0:len(sd.vals)], sd.vals)

	switch len(sd.vals) {
	case 1:
		// At order 1, an initial field of values f is replaced by a new field of values g,
		// where g1 = f1, g2 = f2 - f1, ..., gn = fn - fn-1.

		for n := int(1); n < len(data); n++ {
			data[n] = data[n] + data[n-1] + sd.min
		}
	case 2:
		// At order 2, the field of values g is itself replaced by a new field of values h,
		// where h1 = f1, h2 = f2, h3 = g3 - g2, ..., hn = gn - gn-1.

		for n := int(2); n < len(data); n++ {
			data[n] = data[n] + (2 * data[n-1]) - data[n-2] + sd.min
		}
	}

	// To keep values positive, the overall minimum of the resulting field (either gmin or
	// hmin) is removed. At decoding time, after bit string unpacking, the original scaled
	// values are recovered by adding the overall minimum and summing up recursively.
}

func (cpsd *ComplexPackingAndSpatialDifferencing) ReadSpacingDifferential(r datapacking.BitReader) (*spacingDifferential, error) {
	if cpsd.OctetsNumber == 0 {
		return nil, nil
	}

	rc := cpsd.OctetsNumber * 8
	var v spacingDifferential

	val, err := r.ReadBits(rc)
	if err != nil {
		return nil, fmt.Errorf("Spacial differencing Value 1: %w", err)
	}

	v.vals = append(v.vals, uint32(regulation.ToInt(int(val), int(rc))))

	if cpsd.SpatialOrderDifference == 2 {
		val, err := r.ReadBits(rc)
		if err != nil {
			return nil, fmt.Errorf("Spacial differencing Value 2: %w", err)
		}

		v.vals = append(v.vals, uint32(regulation.ToInt(int(val), int(rc))))
	}

	minVal, err := r.ReadBits(rc)
	if err != nil {
		return nil, fmt.Errorf("Spacial differencing Reference: %w", err)
	}

	v.min = uint32(regulation.ToInt(int(minVal), int(rc)))

	return &v, nil
}
