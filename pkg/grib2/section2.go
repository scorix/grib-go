package grib2

import (
	"bytes"
	"encoding/binary"
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

func (s *section2) readFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.Section2.Section2FixedPart); err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)

	if _, err := io.Copy(buf, r); err != nil {
		return err
	}

	s.Section2.Local = buf.Bytes()

	return nil
}
