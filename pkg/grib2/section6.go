package grib2

import (
	"encoding/binary"
	"io"

	"github.com/scorix/grib-go/pkg/grib2/definition"
)

type Section6 interface {
	Section
}

type section6 struct {
	definition.Section6
}

func (s *section6) Length() int {
	return int(s.Section6.Section6Length)
}

func (s *section6) Number() int {
	return int(s.Section6.NumberOfSection)
}

func (s *section6) readFrom(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &s.Section6)
}
