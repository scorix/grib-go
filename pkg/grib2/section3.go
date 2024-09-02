package grib2

import (
	"encoding/binary"
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
	return int(s.Section3.Section3Length)
}

func (s *section3) Number() int {
	return int(s.Section3.NumberOfSection)
}

func (s *section3) readFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.Section3.Section3FixedPart); err != nil {
		return err
	}

	tpl, err := gdt.ReadTemplate(r, s.Section3.GridDefinitionTemplateNumber)
	if err != nil {
		return err
	}

	s.Section3.GridDefinitionTemplate = tpl

	return nil
}

func (s *section3) GetGridDefinitionTemplate() gdt.Template {
	return s.Section3.GridDefinitionTemplate
}
