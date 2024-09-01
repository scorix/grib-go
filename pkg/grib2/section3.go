package grib

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/scorix/grib-go/pkg/grib2/gdt"
)

type Section3 struct {
	section3
	gdt.Template
}

func (s *Section3) SectionLength() int {
	return int(s.Section3Length)
}

func (s *Section3) SectionNumber() int {
	return int(s.NumberOfSection)
}

func (s *Section3) ReadFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.section3); err != nil {
		return err
	}

	switch n := s.section3.GridDefinitionTemplateNumber; n {
	case 0:
		s.Template = &gdt.Template0{}

	case 255:
		return nil

	default:
		return fmt.Errorf("unsupported grid definition template: %d", n)
	}

	return binary.Read(r, binary.BigEndian, s.Template)
}

type section3 struct {
	Section3Length                   uint32 // Length of the section in octets (N)
	NumberOfSection                  uint8  // 3 - Number of the section
	SourceOfGridDefinition           uint8  // Source of grid definition
	NumberOfDataPoints               uint32
	NumberOfOctectsForNumberOfPoints uint8
	InterpretationOfNumberOfPoints   uint8
	GridDefinitionTemplateNumber     uint16
}
