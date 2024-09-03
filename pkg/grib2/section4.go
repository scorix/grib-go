package grib2

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/scorix/grib-go/pkg/grib2/definition"
	"github.com/scorix/grib-go/pkg/grib2/pdt"
)

type Section4 interface {
	Section
	GetProductDefinitionTemplate() pdt.Template
	pdt.Template
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
		return err
	}

	s.Section4.ProductDefinitionTemplate = tpl

	return nil
}

func (s *section4) GetProductDefinitionTemplate() pdt.Template {
	return s.Section4.ProductDefinitionTemplate
}

func (s *section4) GetParameterCategory() int {
	return s.GetProductDefinitionTemplate().GetParameterCategory()
}

func (s *section4) GetParameterNumber() int {
	return s.GetProductDefinitionTemplate().GetParameterNumber()
}

func (s *section4) GetForecastDuration() time.Duration {
	return s.GetProductDefinitionTemplate().GetForecastDuration()
}

func (s *section4) GetLevel() int {
	return s.GetProductDefinitionTemplate().GetLevel()
}
