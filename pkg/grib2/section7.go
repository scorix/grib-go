package grib

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
)

type Section7 struct {
	section7
	dataReader datapacking.UnpackReader
	data       []float64
}

func (s *Section7) SectionLength() int {
	return int(s.Section7Length)
}

func (s *Section7) SectionNumber() int {
	return int(s.NumberOfSection)
}

func (s *Section7) ReadFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.section7); err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	data, err := s.dataReader.ReadData(r)
	if err != nil {
		return fmt.Errorf("read data: %w", err)
	}

	s.data = data

	return nil
}

func (s *Section7) Data() []float64 {
	return s.data
}

type section7 struct {
	Section7Length  uint32 // Length of the section in octets (N)
	NumberOfSection uint8  // 7 - Number of the section
}
