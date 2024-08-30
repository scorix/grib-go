package grib

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/scorix/grib-go/pkg/grib2/pdt"
)

type Section4 struct {
	defSection4
}

func (s *Section4) SectionLength() int {
	return int(s.Section4Length)
}

func (s *Section4) SectionNumber() int {
	return int(s.NumberOfSection)
}

func (s *Section4) ReadFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.section4); err != nil {
		return err
	}

	switch n := s.section4.ProductDefinitionTemplateNumber; n {
	case 0:
		s.Template = &pdt.Template0{}

	case 255:
		return nil

	default:
		return fmt.Errorf("unsupported product definition template: %d", n)
	}

	return binary.Read(r, binary.BigEndian, s.Template)
}

type defSection4 struct {
	section4
	pdt.Template
}

type section4 struct {
	Section4Length                  uint32 // Length of the section in octets (N)
	NumberOfSection                 uint8  // 4 - Number of the section
	NV                              uint16
	ProductDefinitionTemplateNumber uint16
}
