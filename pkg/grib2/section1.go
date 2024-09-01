package grib2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/scorix/grib-go/pkg/grib2/definition"
)

type Section1 interface {
	Section
	GetTime(loc *time.Location) time.Time
	GetReferenceTime() definition.ReferenceTime
}

type section1 struct {
	definition.Section1
}

func (s *section1) Length() int {
	return int(s.Section1Length)
}

func (s *section1) Number() int {
	return int(s.NumberOfSection)
}

func (s *section1) readFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.Section1.Section1FixedPart); err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	n := int64(s.Section1FixedPart.Section1Length - 21)

	if _, err := io.CopyN(bytes.NewBuffer(s.Reserved), r, n); err != nil {
		return fmt.Errorf("copy reserved: %w", err)
	}

	return nil
}

func (s *section1) GetTime(loc *time.Location) time.Time {
	return time.Date(int(s.Year), time.Month(s.Month), int(s.Day), int(s.Hour), int(s.Minute), int(s.Second), 0, loc)
}

func (s *section1) GetReferenceTime() definition.ReferenceTime {
	return s.SignificanceOfReferenceTime
}
