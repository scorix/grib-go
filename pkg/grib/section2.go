package grib

import (
	"encoding/binary"
	"io"
)

type Section2 struct {
	section section2
	local   []byte // 6-N Local Use
}

func (s *Section2) SectionLength() int {
	return s.section.SectionLength()
}

func (s *Section2) SectionNumber() int {
	return s.section.SectionNumber()
}

func (s *Section2) ReadFrom(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &s.section)
}

type section2 struct {
	Section2Length  uint32 // Length of the section in octets (N)
	NumberOfSection uint8  // 2 - Number of the section
}

func (s *section2) SectionLength() int {
	return int(s.Section2Length)
}

func (s *section2) SectionNumber() int {
	return int(s.NumberOfSection)
}
