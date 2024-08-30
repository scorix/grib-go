package grib

import (
	"encoding/binary"
	"io"
)

type Section6 struct {
	section section6
}

func (s *Section6) SectionLength() int {
	return s.section.SectionLength()
}

func (s *Section6) SectionNumber() int {
	return s.section.SectionNumber()
}

func (s *Section6) ReadFrom(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &s.section)
}

type section6 struct {
	Section6Length  uint32 // Length of the section in octets (N)
	NumberOfSection uint8  // 6 - Number of the section
	BitMapIndicator uint8
}

func (s *section6) SectionLength() int {
	return int(s.Section6Length)
}

func (s *section6) SectionNumber() int {
	return int(s.NumberOfSection)
}
