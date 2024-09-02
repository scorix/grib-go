package grib2

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/scorix/grib-go/pkg/grib2/definition"
	"github.com/scorix/grib-go/pkg/grib2/drt"
)

type Section5 interface {
	Section
	GetDataRepresentationTemplate() drt.Template
	GetNumberOfValues() int
}

type section5 struct {
	definition.Section5
}

func (s *section5) Length() int {
	return int(s.Section5.Section5Length)
}

func (s *section5) Number() int {
	return int(s.Section5.NumberOfSection)
}

func (s *section5) readFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.Section5.Section5FixedPart); err != nil {
		return fmt.Errorf("read: %w", err)
	}

	tpl, err := drt.ReadTemplate(r, s.Section5.DataRepresentationTemplateNumber)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}

	s.Section5.DataRepresentationTemplate = tpl

	return nil
}

func (s *section5) GetDataRepresentationTemplate() drt.Template {
	return s.Section5.DataRepresentationTemplate
}

func (s *section5) GetNumberOfValues() int {
	return int(s.Section5.NumberOfValues)
}
