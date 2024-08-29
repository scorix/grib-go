package grib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

var (
	GribLiteral = [4]byte{'G', 'R', 'I', 'B'}
)

var (
	ErrNotWellFormed     = fmt.Errorf("grib is not well formed")
	ErrEditionNotMatched = fmt.Errorf("edition number is not matched")
)

type gribReader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

func NewGrib2(r gribReader) *Grib2 {
	return &Grib2{
		r:              r,
		offsetSection0: 0,
		offsetSection1: 16,
	}
}

type Grib2 struct {
	r gribReader

	offsetSection0 int64
	offsetSection1 int64
}

// ReadSection0
func (g *Grib2) ReadSection0() (*Section0, error) {
	if _, err := g.r.Seek(g.offsetSection0, io.SeekStart); err != nil {
		return nil, err
	}

	bs := make([]byte, 16)
	if _, err := g.r.Read(bs); err != nil {
		return nil, fmt.Errorf("reader: read section0: %w", err)
	}

	var sec Section0
	if err := binary.Read(bytes.NewReader(bs), binary.BigEndian, &sec); err != nil {
		return nil, fmt.Errorf("binary: read section0: %w", err)
	}

	if sec.GribLiteral != GribLiteral {
		return nil, ErrNotWellFormed
	}

	if sec.EditionNumber != 2 {
		return nil, fmt.Errorf("grib2: %w", ErrEditionNotMatched)
	}

	return &sec, nil
}

// ReadSection1
func (g *Grib2) ReadSection1() (*Section1, error) {
	if _, err := g.r.Seek(g.offsetSection1, io.SeekStart); err != nil {
		return nil, err
	}

	var sec section1

	length := make([]byte, 4)
	if _, err := g.r.Read(length); err != nil {
		return nil, fmt.Errorf("reader: read section 1 length: %w", err)
	}

	if err := binary.Read(bytes.NewReader(length), binary.BigEndian, &sec.Length); err != nil {
		return nil, fmt.Errorf("binary: read section 1 length: %w", err)
	}

	bs := make([]byte, sec.Length-4)
	if _, err := g.r.Read(bs); err != nil {
		return nil, fmt.Errorf("reader: read section 1: %w", err)
	}

	bs = append(length, bs...)

	if err := binary.Read(bytes.NewReader(bs), binary.BigEndian, &sec); err != nil {
		return nil, fmt.Errorf("binary: read section1: %w", err)
	}

	reserved := make([]byte, sec.Length-21)
	copy(bs[21:], reserved)

	return &Section1{section1: sec, reserved: reserved}, nil
}
