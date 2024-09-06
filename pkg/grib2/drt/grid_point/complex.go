package gridpoint

import (
	"fmt"
	"math"

	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
	"github.com/scorix/grib-go/pkg/grib2/drt/definition"
	"github.com/scorix/grib-go/pkg/grib2/regulation"
)

type ComplexPacking struct {
	*SimplePacking                   // 12-21
	GroupSplittingMethodUsed   int8  // 22
	MissingValueManagementUsed int8  // 23
	PrimaryMissingSubstitute   int32 // 24-27
	SecondaryMissingSubstitute int32 // 28-31
	*Grouping                        // 32-47
}

func NewComplexPacking(def definition.ComplexPacking, numVals int) *ComplexPacking {
	return &ComplexPacking{
		SimplePacking:              NewSimplePacking(def.SimplePacking, numVals),
		GroupSplittingMethodUsed:   regulation.ToInt8(def.GroupSplittingMethodUsed),
		MissingValueManagementUsed: regulation.ToInt8(def.MissingValueManagementUsed),
		PrimaryMissingSubstitute:   regulation.ToInt32(def.PrimaryMissingSubstitute),
		SecondaryMissingSubstitute: regulation.ToInt32(def.SecondaryMissingSubstitute),
		Grouping: &Grouping{
			NumberOfGroups:    regulation.ToInt32(def.NumberOfGroups),
			Widths:            def.GroupWidths,
			WidthsBits:        def.GroupWidthsBits,
			LengthsReference:  def.GroupLengthsReference,
			LengthIncrement:   def.GroupLengthIncrement,
			LastLength:        def.GroupLastLength,
			ScaledLengthsBits: def.GroupScaledLengthsBits,
		},
	}
}

func (cp *ComplexPacking) missingValueSubstitute() (float64, float64, error) {
	switch cp.MissingValueManagementUsed {
	case 0, -1:
		return 0, 0, nil
	case 1:
		return float64(cp.PrimaryMissingSubstitute), 0, nil
	case 2:
		return float64(cp.PrimaryMissingSubstitute), float64(cp.SecondaryMissingSubstitute), nil
	}

	return 0, 0, fmt.Errorf("unimplemented")
}

type scaleGroupDataFunc func(data []uint32, missing []uint32, primary float64, secondary float64, scaleFunc func(uint32) float64) ([]float64, error)

func (cp *ComplexPacking) unpackData(r datapacking.BitReader, groups []Group, f scaleGroupDataFunc) ([]float64, error) {
	data := make([]uint32, cp.numVals)
	miss := make([]uint32, 0, cp.numVals)
	idx := 0

	primary, secondary, err := cp.missingValueSubstitute()
	if err != nil {
		return nil, err
	}

	for _, g := range groups {
		groupData, err := g.ReadData(r)
		if err != nil {
			return nil, fmt.Errorf("read (%d) data: %w", idx, err)
		}

		if idx+1 > cp.numVals {
			return nil, fmt.Errorf("got more than %d values", cp.numVals)
		}

		missingValueBits := g.width
		if missingValueBits == 0 {
			missingValueBits = cp.Bits
		}

		missingValues := []uint32{1<<missingValueBits - 1, 1<<missingValueBits - 2}

		switch cp.MissingValueManagementUsed {
		case 0:
			miss = append(miss, make([]uint32, len(groupData))...)
			for i := range groupData {
				groupData[i] += g.ref
			}

		case 1:
			for i := range groupData {
				if g.ref == missingValues[0] {
					groupData[i] = math.MaxUint32
					miss = append(miss, 1)
				} else {
					groupData[i] += g.ref
					miss = append(miss, 0)
				}
			}

		case 2:
			for i := range groupData {
				if g.ref == missingValues[0] || g.ref == missingValues[1] {
					groupData[i] = math.MaxUint32

					if g.ref == missingValues[0] {
						miss = append(miss, 1)
					} else {
						miss = append(miss, 2)
					}
				} else {
					groupData[i] += g.ref
					miss = append(miss, 0)
				}
			}
		}

		idx += copy(data[idx:], groupData)
	}

	return f(data, miss, primary, secondary, cp.ScaleFunc())
}

func (cp *ComplexPacking) ReadAllData(r datapacking.BitReader) ([]float64, error) {
	groups, err := cp.ReadGroups(r, cp.Bits)
	if err != nil {
		return nil, fmt.Errorf("read groups: %w", err)
	}

	if len(groups) != int(cp.NumberOfGroups) {
		return nil, fmt.Errorf("expected groups: %d, got %d", cp.NumberOfGroups, len(groups))
	}

	return cp.unpackData(r, groups, cp.scaleValues)
}

func (cp *ComplexPacking) scaleValues(data []uint32, miss []uint32, primary float64, secondary float64, scaleFunc func(uint32) float64) ([]float64, error) {
	values := make([]float64, len(data))

	switch cp.MissingValueManagementUsed {
	case 0:
		// no missing values
		for n, dataValue := range data {
			values[n] = scaleFunc(dataValue)
		}

	case 1, 2:
		// missing values included
		for n, dataValue := range data {
			switch miss[n] {
			case 0:
				values[n] = scaleFunc(dataValue)
			case 1:
				values[n] = primary
			case 2:
				values[n] = secondary
			}
		}
	}

	return values, nil
}

type Grouping struct {
	NumberOfGroups    int32  // 32-35
	Widths            uint8  // 36
	WidthsBits        uint8  // 37
	LengthsReference  uint32 // 38-41
	LengthIncrement   uint8  // 42
	LastLength        uint32 // 43-46
	ScaledLengthsBits uint8  // 47
}

type Group struct {
	ref    uint32
	length uint64
	width  uint8
}

func (g Grouping) ReadGroups(r datapacking.BitReader, bits uint8) ([]Group, error) {
	references := make([]uint32, g.NumberOfGroups)
	for n := range g.NumberOfGroups {
		b, err := r.ReadBits(bits)
		if err != nil {
			return nil, err
		}

		references[n] = uint32(b)
	}

	r.Align()

	widths := make([]uint8, g.NumberOfGroups)
	for n := range g.NumberOfGroups {
		b, err := r.ReadBits(g.WidthsBits)
		if err != nil {
			return nil, err
		}

		if int(b)+int(g.Widths) < 0 {
			return nil, fmt.Errorf("invalid width: %d", b)
		}

		widths[n] = uint8(b) + g.Widths
	}

	r.Align()

	lengths := make([]uint64, g.NumberOfGroups)
	for n := range g.NumberOfGroups {
		b, err := r.ReadBits(g.ScaledLengthsBits)
		if err != nil {
			return nil, err
		}

		lengths[n] = b*uint64(g.LengthIncrement) + uint64(g.LengthsReference)
	}

	r.Align()

	lengths[g.NumberOfGroups-1] = uint64(g.LastLength)

	groups := make([]Group, g.NumberOfGroups)
	for n := range g.NumberOfGroups {
		g := Group{
			ref:    references[n],
			width:  widths[n],
			length: lengths[n],
		}

		groups[n] = g
	}

	return groups, nil
}

func (g *Group) ReadData(r datapacking.BitReader) ([]uint32, error) {
	data := make([]uint32, g.length)

	if g.width == 0 {
		return data, nil
	}

	for i := range g.length {
		b, err := r.ReadBits(g.width)
		if err != nil {
			return nil, err
		}

		data[i] = uint32(b)
	}

	return data, nil
}
