package grib

import (
	"encoding/binary"
	"io"
)

type Section2 struct {
	section2
	local []byte // 6-N Local Use
}

func (s *Section2) SectionLength() int {
	return int(s.Section2Length)
}

func (s *Section2) SectionNumber() int {
	return int(s.NumberOfSection)
}

func (s *Section2) ReadFrom(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &s.section2)
}

type section2 struct {
	Section2Length  uint32 // Length of the section in octets (N)
	NumberOfSection uint8  // 2 - Number of the section
}
