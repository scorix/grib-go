package grib

import (
	"encoding/binary"
	"io"
)

type Section7 struct {
	section section7
}

func (s *Section7) SectionLength() int {
	return s.section.SectionLength()
}

func (s *Section7) SectionNumber() int {
	return s.section.SectionNumber()
}

func (s *Section7) ReadFrom(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &s.section)
}

type section7 struct {
	Section7Length  uint32 // Length of the section in octets (N)
	NumberOfSection uint8  // 7 - Number of the section
	BitMapIndicator uint8
}

func (s *section7) SectionLength() int {
	return int(s.Section7Length)
}

func (s *section7) SectionNumber() int {
	return int(s.NumberOfSection)
}
