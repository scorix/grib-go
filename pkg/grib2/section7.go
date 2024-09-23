package grib2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/icza/bitio"
	"github.com/scorix/grib-go/pkg/grib2/definition"
	"github.com/scorix/grib-go/pkg/grib2/drt"
)

type Section7 interface {
	Section
	GetData(drt.Template) ([]float64, error)
}

type section7 struct {
	definition.Section7

	dataReader   io.ReaderAt
	dataOffset   int64
	dataSize     int64
	readDataOnce sync.Once
	readDataErr  error
}

func (s *section7) Length() int {
	return int(s.Section7.Section7Length)
}

func (s *section7) Number() int {
	return int(s.Section7.NumberOfSection)
}

func (s *section7) readFrom(r io.ReaderAt, offset int64, length int64) error {
	p := make([]byte, length)
	if _, err := r.ReadAt(p, offset); err != nil {
		return fmt.Errorf("read %d bytes at %d: %w", length, offset, err)
	}

	n, err := binary.Decode(p, binary.BigEndian, &s.Section7.Section7FixedPart)
	if err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	s.dataSize = int64(s.Section7.Section7Length) - int64(n)
	s.dataOffset = offset + int64(n)
	s.dataReader = r

	return nil
}

func (s *section7) LoadData() error {
	s.readDataOnce.Do(func() {
		data := make([]byte, s.dataSize)
		_, err := s.dataReader.ReadAt(data, s.dataOffset)
		if err != nil {
			s.readDataErr = err
		}

		s.Section7.Data = data
	})

	return s.readDataErr
}

func (s *section7) GetData(tpl drt.Template) ([]float64, error) {
	if err := s.LoadData(); err != nil {
		return nil, err
	}

	br := bitio.NewReader(bytes.NewReader(s.Section7.Data))

	data, err := tpl.ReadAllData(br)
	if err != nil {
		return nil, fmt.Errorf("read data from %T: %w", tpl, err)
	}

	return data, nil
}
