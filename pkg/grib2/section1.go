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
	return int(s.Section1.Section1Length)
}

func (s *section1) Number() int {
	return int(s.Section1.NumberOfSection)
}

func (s *section1) readFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.Section1.Section1FixedPart); err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	n := int64(s.Section1.Section1FixedPart.Section1Length - 21)

	if _, err := io.CopyN(bytes.NewBuffer(s.Section1.Reserved), r, n); err != nil {
		return fmt.Errorf("copy reserved: %w", err)
	}

	return nil
}

func (s *section1) GetTime(loc *time.Location) time.Time {
	return time.Date(int(s.Section1.Year), time.Month(s.Section1.Month), int(s.Section1.Day), int(s.Section1.Hour), int(s.Section1.Minute), int(s.Section1.Second), 0, loc)
}

func (s *section1) GetReferenceTime() definition.ReferenceTime {
	return s.Section1.SignificanceOfReferenceTime
}
