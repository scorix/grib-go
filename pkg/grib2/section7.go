package grib2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/scorix/grib-go/pkg/grib2/definition"
	"github.com/scorix/grib-go/pkg/grib2/drt"
)

type Section7 interface {
	Section
	GetData(drt.Template) ([]float64, error)
}

type section7 struct {
	definition.Section7
}

func (s *section7) Length() int {
	return int(s.Section7Length)
}

func (s *section7) Number() int {
	return int(s.NumberOfSection)
}

func (s *section7) readFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.Section7.Section7FixedPart); err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	data := bytes.NewBuffer(nil)

	if _, err := io.Copy(data, r); err != nil {
		return err
	}

	s.Data = data.Bytes()

	return nil
}

func (s *section7) GetData(tpl drt.Template) ([]float64, error) {
	data, err := drt.ScaleData(tpl, bytes.NewReader(s.Data))
	if err != nil {
		return nil, fmt.Errorf("read data: %w", err)
	}

	return data, nil
}
