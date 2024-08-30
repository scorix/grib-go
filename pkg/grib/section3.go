package grib

import (
	"encoding/binary"
	"io"
)

type Section3 struct {
	section section3
}

func (s *Section3) SectionLength() int {
	return s.section.SectionLength()
}

func (s *Section3) SectionNumber() int {
	return s.section.SectionNumber()
}

func (s *Section3) ReadFrom(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &s.section)
}

type section3 struct {
	Section3Length                         uint32 // Length of the section in octets (N)
	NumberOfSection                        uint8  // 3 - Number of the section
	SourceOfGridDefinition                 uint8  // Source of grid definition
	NumberOfDataPoints                     uint32
	NumberOfOctectsForNumberOfPoints       uint8
	InterpretationOfNumberOfPoints         uint8
	GridDefinitionTemplateNumber           uint16
	ShapeOfTheEarth                        uint8
	ScaleFactorOfRadiusOfSphericalEarth    uint8
	ScaledValueOfRadiusOfSphericalEarth    uint32
	ScaleFactorOfEarthMajorAxis            uint8
	ScaledValueOfEarthMajorAxis            uint32
	ScaleFactorOfEarthMinorAxis            uint8
	ScaledValueOfEarthMinorAxis            uint32
	Ni                                     uint32
	Nj                                     uint32
	BasicAngleOfTheInitialProductionDomain uint32
	SubdivisionsOfBasicAngle               uint32
	LatitudeOfFirstGridPoint               uint32
	LongitudeOfFirstGridPoint              uint32
	ResolutionAndComponentFlags            uint8
	LatitudeOfLastGridPoint                uint32
	LongitudeOfLastGridPoint               uint32
	IDirectionIncrement                    uint32
	JDirectionIncrement                    uint32
	ScanningMode                           uint8
}

func (s *section3) SectionLength() int {
	return int(s.Section3Length)
}

func (s *section3) SectionNumber() int {
	return int(s.NumberOfSection)
}
