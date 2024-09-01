package grib2

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/scorix/grib-go/pkg/grib2/definition"
	"github.com/scorix/grib-go/pkg/grib2/gdt"
)

type Section3 interface {
	Section
	GetGridDefinitionTemplate() gdt.Template
}

type section3 struct {
	definition.Section3
}

func (s *section3) Length() int {
	return int(s.Section3Length)
}

func (s *section3) Number() int {
	return int(s.NumberOfSection)
}

func (s *section3) readFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.Section3.Section3FixedPart); err != nil {
		return err
	}

	switch n := s.Section3.GridDefinitionTemplateNumber; n {
	case 0:
		var tpl gdt.Template0
		if err := binary.Read(r, binary.BigEndian, &tpl.Template0FixedPart); err != nil {
			return err
		}

		s.GridDefinitionTemplate = &tpl

		return nil

	case 255:
		return nil

	default:
		return fmt.Errorf("unsupported grid definition template: %d", n)
	}
}

func (s *section3) GetGridDefinitionTemplate() gdt.Template {
	return s.GridDefinitionTemplate
}
