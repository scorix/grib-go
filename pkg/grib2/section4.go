package grib2

import (
	"bytes"
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
	return int(s.Section4.Section4Length)
}

func (s *section4) Number() int {
	return int(s.Section4.NumberOfSection)
}

func (s *section4) readFrom(r io.ReaderAt, offset int64, length int64) error {
	p := make([]byte, length)
	if _, err := r.ReadAt(p, offset); err != nil {
		return fmt.Errorf("read %d bytes at %d: %w", length, offset, err)
	}

	n, err := binary.Decode(p, binary.BigEndian, &s.Section4.Section4FixedPart)
	if err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	tpl, err := pdt.ReadTemplate(bytes.NewBuffer(p[n:]), s.Section4.ProductDefinitionTemplateNumber)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}

	s.Section4.ProductDefinitionTemplate = tpl

	return nil
}

func (s *section4) GetProductDefinitionTemplate() pdt.Template {
	return s.Section4.ProductDefinitionTemplate
}
