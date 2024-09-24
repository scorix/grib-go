package grib2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/icza/bitio"
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

func (s *section5) readFrom(r io.ReaderAt, offset int64, length int64) error {
	p := make([]byte, length)
	if _, err := r.ReadAt(p, offset); err != nil {
		return fmt.Errorf("read %d bytes at %d: %w", length, offset, err)
	}

	n, err := binary.Decode(p, binary.BigEndian, &s.Section5.Section5FixedPart)
	if err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	tpl, err := drt.ReadTemplate(bitio.NewReader(bytes.NewBuffer(p[n:])), s.Section5.DataRepresentationTemplateNumber, int(s.Section5.NumberOfValues))
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
