package grib2

import (
	"encoding/binary"
	"io"

	"github.com/scorix/grib-go/pkg/grib2/definition"
)

type section0 struct {
	definition.Section0
}

// Edition number
//
// - 1 for GRIB1
// - 2 for GRIB2
func (s *section0) GetEditionNumber() int {
	return int(s.EditionNumber)
}

// Discipline (From [Table 0.0](https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table0-0.shtml))
func (s *section0) GetDiscipline() int {
	return int(s.Discipline)
}

// Total length of GRIB message in octets (All sections)
func (s *section0) GetGribLength() int {
	return int(s.GribLength)
}

func (s *section0) Number() int {
	return 0
}

func (s *section0) Length() int {
	return 16
}

func (s *section0) readFrom(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &s.Section0)
}
