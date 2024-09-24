package grib2

import (
	"encoding/binary"
	"fmt"
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

func (s *section6) readFrom(r io.ReaderAt, offset int64, length int64) error {
	p := make([]byte, length)
	if _, err := r.ReadAt(p, offset); err != nil {
		return fmt.Errorf("read %d bytes at %d: %w", length, offset, err)
	}

	_, err := binary.Decode(p, binary.BigEndian, &s.Section6)
	if err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	return nil
}
