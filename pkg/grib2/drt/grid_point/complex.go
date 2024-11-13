package gridpoint

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/scorix/grib-go/internal/pkg/bitio"
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

func (cp *ComplexPacking) unpackData(r *bitio.Reader, groups []Group, f scaleGroupDataFunc) ([]float64, error) {
	data := make([]uint32, cp.NumVals)
	miss := make([]uint32, 0, cp.NumVals)
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

		if idx+1 > cp.NumVals {
			return nil, fmt.Errorf("got more than %d values", cp.NumVals)
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

func (cp *ComplexPacking) ReadAllData(r *bitio.Reader) ([]float64, error) {
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

func (cp *ComplexPacking) Definition() any {
	return definition.ComplexPacking{
		SimplePacking:              cp.SimplePacking.Definition().(definition.SimplePacking),
		GroupSplittingMethodUsed:   regulation.ToUint8(cp.GroupSplittingMethodUsed),
		MissingValueManagementUsed: regulation.ToUint8(cp.MissingValueManagementUsed),
		PrimaryMissingSubstitute:   regulation.ToUint32(cp.PrimaryMissingSubstitute),
		SecondaryMissingSubstitute: regulation.ToUint32(cp.SecondaryMissingSubstitute),
		NumberOfGroups:             regulation.ToUint32(cp.NumberOfGroups),
		GroupWidths:                cp.Grouping.Widths,
		GroupWidthsBits:            cp.Grouping.WidthsBits,
		GroupLengthsReference:      cp.Grouping.LengthsReference,
		GroupLengthIncrement:       cp.Grouping.LengthIncrement,
		GroupLastLength:            cp.Grouping.LastLength,
		GroupScaledLengthsBits:     cp.Grouping.ScaledLengthsBits,
	}
}

type complexPacking struct {
	SimplePacking              *SimplePacking `json:"simple_packing"`
	GroupSplittingMethodUsed   int8           `json:"group_splitting_method_used"`
	MissingValueManagementUsed int8           `json:"missing_value_management_used"`
	PrimaryMissingSubstitute   int32          `json:"primary_missing_substitute"`
	SecondaryMissingSubstitute int32          `json:"secondary_missing_substitute"`
	Grouping                   *Grouping      `json:"grouping"`
	NumVals                    int            `json:"num_vals"`
}

func (cp *ComplexPacking) MarshalJSON() ([]byte, error) {
	return json.Marshal(complexPacking{
		SimplePacking:              cp.SimplePacking,
		GroupSplittingMethodUsed:   cp.GroupSplittingMethodUsed,
		MissingValueManagementUsed: cp.MissingValueManagementUsed,
		PrimaryMissingSubstitute:   cp.PrimaryMissingSubstitute,
		SecondaryMissingSubstitute: cp.SecondaryMissingSubstitute,
		Grouping:                   cp.Grouping,
		NumVals:                    cp.NumVals,
	})
}

func (cp *ComplexPacking) UnmarshalJSON(data []byte) error {
	var temp complexPacking

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	cp.SimplePacking = temp.SimplePacking
	cp.GroupSplittingMethodUsed = temp.GroupSplittingMethodUsed
	cp.MissingValueManagementUsed = temp.MissingValueManagementUsed
	cp.PrimaryMissingSubstitute = temp.PrimaryMissingSubstitute
	cp.SecondaryMissingSubstitute = temp.SecondaryMissingSubstitute
	cp.Grouping = temp.Grouping
	cp.NumVals = temp.NumVals
	return nil
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

func (g Grouping) ReadGroups(r *bitio.Reader, bits uint8) ([]Group, error) {
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

func (g *Group) ReadData(r *bitio.Reader) ([]uint32, error) {
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

type grouping struct {
	NumberOfGroups    int32  `json:"number_of_groups"`
	Widths            uint8  `json:"widths"`
	WidthsBits        uint8  `json:"widths_bits"`
	LengthsReference  uint32 `json:"lengths_reference"`
	LengthIncrement   uint8  `json:"length_increment"`
	LastLength        uint32 `json:"last_length"`
	ScaledLengthsBits uint8  `json:"scaled_lengths_bits"`
}

func (g Grouping) MarshalJSON() ([]byte, error) {
	return json.Marshal(grouping{
		NumberOfGroups:    g.NumberOfGroups,
		Widths:            g.Widths,
		WidthsBits:        g.WidthsBits,
		LengthsReference:  g.LengthsReference,
		LengthIncrement:   g.LengthIncrement,
		LastLength:        g.LastLength,
		ScaledLengthsBits: g.ScaledLengthsBits,
	})
}

func (g *Grouping) UnmarshalJSON(data []byte) error {
	var gg grouping

	if err := json.Unmarshal(data, &gg); err != nil {
		return err
	}

	g.NumberOfGroups = gg.NumberOfGroups
	g.Widths = gg.Widths
	g.WidthsBits = gg.WidthsBits
	g.LengthsReference = gg.LengthsReference
	g.LengthIncrement = gg.LengthIncrement
	g.LastLength = gg.LastLength
	g.ScaledLengthsBits = gg.ScaledLengthsBits

	return nil
}
