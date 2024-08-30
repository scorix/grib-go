package grib

import (
	"encoding/binary"
	"io"
)

type Section5 struct {
	section section5
}

func (s *Section5) SectionLength() int {
	return s.section.SectionLength()
}

func (s *Section5) SectionNumber() int {
	return s.section.SectionNumber()
}

func (s *Section5) ReadFrom(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &s.section)
}

type section5 struct {
	Section5Length                   uint32 // Length of the section in octets (N)
	NumberOfSection                  uint8  // 5 - Number of the section
	NumberOfValues                   uint32
	DataRepresentationTemplateNumber uint16
	ReferenceValue                   uint32
	BinaryScaleFactor                uint16
	DecimalScaleFactor               uint16
	BitsPerValue                     uint8
	TypeOfOriginalFieldValues        uint8
}

func (s *section5) SectionLength() int {
	return int(s.Section5Length)
}

func (s *section5) SectionNumber() int {
	return int(s.NumberOfSection)
}
