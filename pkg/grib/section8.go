package grib

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Section8 struct {
	section section8
}

func (s *Section8) SectionLength() int {
	return s.section.SectionLength()
}

func (s *Section8) SectionNumber() int {
	return s.section.SectionNumber()
}

func (s *Section8) ReadFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.section); err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	magicNumber := [4]byte{'7', '7', '7', '7'}

	if s.section.MagicNumber != magicNumber {
		return fmt.Errorf("malformed section8: %d", s.section.MagicNumber)
	}

	return nil
}

type section8 struct {
	MagicNumber [4]byte
}

func (s *section8) SectionLength() int {
	return 4
}

func (s *section8) SectionNumber() int {
	return 8
}
