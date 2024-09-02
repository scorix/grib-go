package gridpoint

import (
	"fmt"
	"io"

	"github.com/icza/bitio"
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
	*Group
}

func NewComplexPacking(def definition.ComplexPacking, numVals int) *ComplexPacking {
	return &ComplexPacking{
		SimplePacking:              NewSimplePacking(def.SimplePacking, numVals),
		GroupMethod:                regulation.ToInt8(def.GroupMethod),
		MissingValue:               regulation.ToInt8(def.MissingValue),
		PrimaryMissingSubstitute:   regulation.ToInt32(def.PrimaryMissingSubstitute),
		SecondaryMissingSubstitute: regulation.ToInt32(def.SecondaryMissingSubstitute),
		NumberOfGroups:             regulation.ToInt32(def.NumberOfGroups),
		Group: &Group{
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
	switch cp.MissingValue {
	case 0, -1:
		return 0, 0, nil
	case 1:
		return float64(cp.PrimaryMissingSubstitute), 0, nil
	case 2:
		return float64(cp.PrimaryMissingSubstitute), float64(cp.SecondaryMissingSubstitute), nil
	}

	return 0, 0, fmt.Errorf("unimplemented")
}

func (cp *ComplexPacking) unpackData(r *bitio.Reader, groups []group) ([]float64, error) {
	primary, secondary, err := cp.missingValueSubstitute()
	if err != nil {
		return nil, err
	}

	data := make([]float64, cp.numVals)
	ifldmiss := make([]uint32, 0, cp.numVals)
	s7i := 0

	scale := cp.ScaleFunc()

	for _, g := range groups {
		groupData, err := g.readData(r)
		if err != nil {
			return nil, fmt.Errorf("read (%d) data: %w", s7i, err)
		}

		if s7i+1 > cp.numVals {
			return nil, fmt.Errorf("got more than %d values", cp.numVals)
		}

		missingValueBits := g.width
		if missingValueBits == 0 {
			missingValueBits = cp.Bits
		}

		missingValues := []uint32{1<<missingValueBits - 1, 1<<missingValueBits - 2}

		switch cp.MissingValue {
		case 0:
			ifldmiss = append(ifldmiss, make([]uint32, len(groupData))...)
			for _, d := range groupData {
				data[s7i] = scale(g.ref + d)
				s7i++
			}

		case 1:
			for _, d := range groupData {
				if g.ref == missingValues[0] {
					data[s7i] = primary
					s7i++
					ifldmiss = append(ifldmiss, 1)
				} else {
					data[s7i] = scale(g.ref + d)
					s7i++
					ifldmiss = append(ifldmiss, 0)
				}
			}

		case 2:
			for _, d := range groupData {
				if g.ref == missingValues[0] || g.ref == missingValues[1] {
					data[s7i] = secondary
					s7i++

					if g.ref == missingValues[0] {
						ifldmiss = append(ifldmiss, 1)
					} else {
						ifldmiss = append(ifldmiss, 2)
					}
				} else {
					data[s7i] = scale(g.ref + d)
					s7i++
					ifldmiss = append(ifldmiss, 0)
				}
			}

		}
	}

	return data, nil
}

// ReadAllData parses data2 struct from the reader into the an array of floating-point values
func (cp *ComplexPacking) ReadAllData(r io.Reader) ([]float64, error) {
	br := bitio.NewReader(r)

	groups, err := cp.readGroups(br)
	if err != nil {
		return nil, fmt.Errorf("read groups: %w", err)
	}

	if len(groups) != int(cp.NumberOfGroups) {
		return nil, fmt.Errorf("expected groups: %d, got %d", cp.NumberOfGroups, len(groups))
	}

	return cp.unpackData(br, groups)
}

type bitGroupParameter struct {
	GroupLengthsReference uint64
	GroupWidths           uint64
	GroupLastLength       uint64
}

type Group struct {
	Widths            uint8
	WidthsBits        uint8
	LengthsReference  uint32
	LengthIncrement   uint8
	LastLength        uint32
	ScaledLengthsBits uint8
}

func (cp *ComplexPacking) readGroups(r *bitio.Reader) ([]group, error) {
	references := make([]uint32, cp.NumberOfGroups)
	for n := range cp.NumberOfGroups {
		b, err := r.ReadBits(cp.Bits)
		if err != nil {
			return nil, err
		}

		references[n] = uint32(b)
	}

	widths := make([]uint8, cp.NumberOfGroups)
	for n := range cp.NumberOfGroups {
		b, err := r.ReadBits(cp.Group.WidthsBits)
		if err != nil {
			return nil, err
		}

		if int8(b) < 0 {
			return nil, fmt.Errorf("invalid width: %d", b)
		}

		widths[n] = uint8(b) + cp.Group.Widths
	}

	lengths := make([]uint64, cp.NumberOfGroups)
	for n := range cp.NumberOfGroups {
		b, err := r.ReadBits(cp.Group.ScaledLengthsBits)
		if err != nil {
			return nil, err
		}

		lengths[n] = b*uint64(cp.Group.LengthIncrement) + uint64(cp.Group.LengthsReference)
	}

	lengths[cp.NumberOfGroups-1] = uint64(cp.Group.LastLength)

	groups := make([]group, cp.NumberOfGroups)
	for n := range cp.NumberOfGroups {
		g := group{
			ref:    references[n],
			width:  widths[n],
			length: lengths[n],
		}

		groups[n] = g
	}

	r.Align()

	return groups, nil
}

type group struct {
	ref    uint32
	length uint64
	width  uint8
}

func (g *group) readData(r *bitio.Reader) ([]uint32, error) {
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
