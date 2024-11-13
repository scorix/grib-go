package grib2

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/scorix/grib-go/internal/pkg/bitio"
	"github.com/scorix/grib-go/pkg/grib2/definition"
	"github.com/scorix/grib-go/pkg/grib2/drt"
)

type Section7 interface {
	Section
	GetData(drt.Template) ([]float32, error)
	GetDataOffset() int64
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
	s.Section7.Section7FixedPart = definition.Section7FixedPart{
		Section7Length:  uint32(length),
		NumberOfSection: 7,
	}

	s.dataSize = int64(s.Section7.Section7Length - 5)
	s.dataOffset = offset + 5
	s.dataReader = r

	return nil
}

func (s *section7) GetDataOffset() int64 {
	return s.dataOffset
}

func (s *section7) LoadData() error {
	s.readDataOnce.Do(func() {
		data := make([]byte, s.dataSize)
		n, err := s.dataReader.ReadAt(data, s.dataOffset)
		if err != nil {
			s.readDataErr = fmt.Errorf("load data: total %d, read %d, offset %d: %w", n, s.dataSize, s.dataOffset, err)
		}

		s.Section7.Data = data
	})

	return s.readDataErr
}

func (s *section7) GetData(tpl drt.Template) ([]float32, error) {
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
