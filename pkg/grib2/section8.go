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

func (s *section8) readFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.Section8); err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	magicNumber := [4]byte{'7', '7', '7', '7'}

	if s.MagicNumber != magicNumber {
		return fmt.Errorf("malformed section8: %d", s.MagicNumber)
	}

	return nil
}
