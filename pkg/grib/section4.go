package grib

import (
	"encoding/binary"
	"io"
)

type Section4 struct {
	section section4
}

func (s *Section4) SectionLength() int {
	return s.section.SectionLength()
}

func (s *Section4) SectionNumber() int {
	return s.section.SectionNumber()
}

func (s *Section4) ReadFrom(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &s.section)
}

type section4 struct {
	Section4Length                  uint32 // Length of the section in octets (N)
	NumberOfSection                 uint8  // 4 - Number of the section
	NV                              uint16
	ProductDefinitionTemplateNumber uint16
	ParameterCategory               uint8
	ParameterNumber                 uint8
	TypeOfGeneratingProcess         uint8
	BackgroundProcess               uint8
	HoursAfterDataCutoff            uint16
	MinutesAfterDataCutoff          uint8
	IndicatorOfUnitForForecastTime  uint8
	ForecastTime                    uint32
	TypeOfFirstFixedSurface         uint8
	ScaleFactorOfFirstFixedSurface  uint8
	ScaledValueOfFirstFixedSurface  uint32
	TypeOfSecondFixedSurface        uint8
	ScaleFactorOfSecondFixedSurface uint8
	ScaledValueOfSecondFixedSurface uint32
}

func (s *section4) SectionLength() int {
	return int(s.Section4Length)
}

func (s *section4) SectionNumber() int {
	return int(s.NumberOfSection)
}
