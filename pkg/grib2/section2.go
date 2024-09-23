package grib2

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/scorix/grib-go/pkg/grib2/definition"
)

type Section2 interface {
	Section
}

type section2 struct {
	definition.Section2
}

func (s *section2) Length() int {
	return int(s.Section2.Section2Length)
}

func (s *section2) Number() int {
	return int(s.Section2.NumberOfSection)
}

func (s *section2) readFrom(r io.ReaderAt, offset int64, length int64) error {
	p := make([]byte, length)
	if _, err := r.ReadAt(p, offset); err != nil {
		return fmt.Errorf("read %d bytes at %d: %w", length, offset, err)
	}

	n, err := binary.Decode(p, binary.BigEndian, &s.Section2.Section2FixedPart)
	if err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	s.Section2.Local = p[n:]

	return nil
}
