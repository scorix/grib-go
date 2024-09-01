package grib2

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/scorix/grib-go/pkg/grib2/definition"
	"github.com/scorix/grib-go/pkg/grib2/pdt"
)

type Section4 interface {
	Section
	GetProductDefinitionTemplate() pdt.Template
}

type section4 struct {
	definition.Section4
}

func (s *section4) Length() int {
	return int(s.Section4Length)
}

func (s *section4) Number() int {
	return int(s.NumberOfSection)
}

func (s *section4) readFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.Section4.Section4FixedPart); err != nil {
		return err
	}

	switch n := s.Section4.ProductDefinitionTemplateNumber; n {
	case 0:
		s.ProductDefinitionTemplate = &pdt.Template0{}

	case 255:
		return nil

	default:
		return fmt.Errorf("unsupported product definition template: %d", n)
	}

	return binary.Read(r, binary.BigEndian, s.ProductDefinitionTemplate)
}

func (s *section4) GetProductDefinitionTemplate() pdt.Template {
	return s.ProductDefinitionTemplate
}
