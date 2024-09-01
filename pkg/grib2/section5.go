package grib2

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/scorix/grib-go/pkg/grib2/definition"
	"github.com/scorix/grib-go/pkg/grib2/drt"
	gridpoint "github.com/scorix/grib-go/pkg/grib2/drt/grid_point"
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
	return int(s.Section5Length)
}

func (s *section5) Number() int {
	return int(s.NumberOfSection)
}

func (s *section5) readFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.Section5.Section5FixedPart); err != nil {
		return fmt.Errorf("read: %w", err)
	}

	if err := s.readTemplate(r); err != nil {
		return fmt.Errorf("read template: %w", err)
	}

	return nil
}

func (s *section5) readTemplate(r io.Reader) error {
	switch s.DataRepresentationTemplateNumber {
	case definition.GridPointDataSimplePacking:
		var tplDef gridpoint.DefSimplePacking

		if err := binary.Read(r, binary.BigEndian, &tplDef); err != nil {
			return err
		}

		s.DataRepresentationTemplate = gridpoint.NewSimplePacking(tplDef)

		return nil

	case definition.MatrixValueAtGridPointSimplePacking:

	case definition.GridPointDataComplexPacking:

	}

	return fmt.Errorf("unsupported data template: %d", s.DataRepresentationTemplateNumber)
}

func (s *section5) GetDataRepresentationTemplate() drt.Template {
	return s.DataRepresentationTemplate
}

func (s *section5) GetNumberOfValues() int {
	return int(s.NumberOfValues)
}
