package grib

import (
	"encoding/binary"
	"io"
)

type Section6 struct {
	section6
}

func (s *Section6) SectionLength() int {
	return int(s.Section6Length)
}

func (s *Section6) SectionNumber() int {
	return int(s.NumberOfSection)
}

func (s *Section6) ReadFrom(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &s.section6)
}

type section6 struct {
	Section6Length  uint32 // Length of the section in octets (N)
	NumberOfSection uint8  // 6 - Number of the section
	BitMapIndicator uint8
}
