package grib

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/scorix/grib-go/pkg/grib2/drt"
	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
	gridpoint "github.com/scorix/grib-go/pkg/grib2/drt/grid_point"
)

type Section5 struct {
	defSection5

	DataPackingReader datapacking.UnpackReader
}

func (s *Section5) SectionLength() int {
	return int(s.Section5Length)
}

func (s *Section5) SectionNumber() int {
	return int(s.NumberOfSection)
}

func (s *Section5) ReadFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.section5); err != nil {
		return fmt.Errorf("read: %w", err)
	}

	if err := s.ReadTemplate(r); err != nil {
		return fmt.Errorf("read template: %w", err)
	}

	pr, err := s.DataRepresentationTemplate.NewUnpackReader(r)
	if err != nil {
		return fmt.Errorf("unpack reader: %w", err)
	}

	s.DataPackingReader = pr

	return nil
}

type defSection5 struct {
	section5
	DataRepresentationTemplate drt.Template
}

type section5 struct {
	Section5Length                   uint32 // Length of the section in octets (N)
	NumberOfSection                  uint8  // 5 - Number of the section
	NumberOfValues                   uint32
	DataRepresentationTemplateNumber drt.DataRepresentationTemplateNumber
}

func (s *Section5) ReadTemplate(r io.Reader) error {
	switch s.DataRepresentationTemplateNumber {
	case drt.GridPointDataSimplePacking:
		var tplDef gridpoint.DefSimplePacking

		if err := binary.Read(r, binary.BigEndian, &tplDef); err != nil {
			return err
		}

		s.defSection5.DataRepresentationTemplate = gridpoint.NewSimplePacking(tplDef)

		return nil

	case drt.MatrixValueAtGridPointSimplePacking:

	case drt.GridPointDataComplexPacking:

	}

	return fmt.Errorf("unsupported data template: %d", s.DataRepresentationTemplateNumber)
}
