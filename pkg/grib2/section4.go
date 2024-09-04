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
	return int(s.Section4.Section4Length)
}

func (s *section4) Number() int {
	return int(s.Section4.NumberOfSection)
}

func (s *section4) readFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.Section4.Section4FixedPart); err != nil {
		return err
	}

	tpl, err := pdt.ReadTemplate(r, s.Section4.ProductDefinitionTemplateNumber)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}

	s.Section4.ProductDefinitionTemplate = tpl

	return nil
}

func (s *section4) GetProductDefinitionTemplate() pdt.Template {
	return s.Section4.ProductDefinitionTemplate
}
