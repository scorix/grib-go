package grib2

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/scorix/grib-go/pkg/grib2/definition"
)

type Section8 interface {
	Section
}

type section8 struct {
	definition.Section8
}

func (s *section8) Length() int {
	return 4
}

func (s *section8) Number() int {
	return 8
}

func (s *section8) readFrom(r io.ReaderAt, offset int64, length int64) error {
	p := make([]byte, length)
	if _, err := r.ReadAt(p, offset); err != nil {
		return fmt.Errorf("read %d bytes at %d: %w", length, offset, err)
	}

	_, err := binary.Decode(p, binary.BigEndian, &s.Section8)
	if err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	magicNumber := [4]byte{'7', '7', '7', '7'}

	if s.Section8.MagicNumber != magicNumber {
		return fmt.Errorf("malformed section8: %d", s.Section8.MagicNumber)
	}

	return nil
}
