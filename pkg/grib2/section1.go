package grib2

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/scorix/grib-go/pkg/grib2/definition"
)

type Section1 interface {
	Section
	GetTime(loc *time.Location) time.Time
}

type section1 struct {
	definition.Section1
}

func (s *section1) Length() int {
	return int(s.Section1.Section1Length)
}

func (s *section1) Number() int {
	return int(s.Section1.NumberOfSection)
}

func (s *section1) readFrom(r io.ReaderAt, offset int64, length int64) error {
	p := make([]byte, length)
	if _, err := r.ReadAt(p, offset); err != nil {
		return fmt.Errorf("read %d bytes at %d: %w", length, offset, err)
	}

	n, err := binary.Decode(p, binary.BigEndian, &s.Section1.Section1FixedPart)
	if err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	s.Section1.Reserved = p[n:]

	return nil
}

func (s *section1) GetTime(loc *time.Location) time.Time {
	return time.Date(int(s.Section1.Year), time.Month(s.Section1.Month), int(s.Section1.Day), int(s.Section1.Hour), int(s.Section1.Minute), int(s.Section1.Second), 0, loc)
}
