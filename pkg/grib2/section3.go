package grib2

import (
	"bytes"
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
	return int(s.Section3.Section3Length)
}

func (s *section3) Number() int {
	return int(s.Section3.NumberOfSection)
}

func (s *section3) readFrom(r io.ReaderAt, offset int64, length int64) error {
	p := make([]byte, length)
	if _, err := r.ReadAt(p, offset); err != nil {
		return fmt.Errorf("read %d bytes at %d: %w", length, offset, err)
	}

	n, err := binary.Decode(p, binary.BigEndian, &s.Section3.Section3FixedPart)
	if err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	tpl, err := gdt.ReadTemplate(bytes.NewBuffer(p[n:]), s.Section3.GridDefinitionTemplateNumber)
	if err != nil {
		return err
	}

	s.Section3.GridDefinitionTemplate = tpl

	return nil
}

func (s *section3) GetGridDefinitionTemplate() gdt.Template {
	return s.Section3.GridDefinitionTemplate
}
