package grib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

var (
	GribLiteral = [4]byte{'G', 'R', 'I', 'B'}
)

var (
	ErrNotWellFormed     = fmt.Errorf("grib is not well formed")
	ErrEditionNotMatched = fmt.Errorf("edition number is not matched")
	ErrSectionNotMatched = fmt.Errorf("section not matched")
)

type gribReader interface {
	io.ReadSeeker
	io.ReaderAt
}

func NewGrib2(r gribReader) *Grib2 {
	return &Grib2{
		r: r,
	}
}

type Grib2 struct {
	r gribReader

	m sync.Mutex
}

// ReadSection0
func (g *Grib2) ReadSection0() (*Section0, error) {
	g.m.Lock()
	defer g.m.Unlock()

	if _, err := g.r.Seek(0, io.SeekStart); err != nil {
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

func (g *Grib2) NextSection() (Section, error) {
	g.m.Lock()
	defer g.m.Unlock()

	var (
		secHead struct {
			SectionLength uint32
			Number        uint8
		}
		bs  = make([]byte, 5)
		buf = bytes.NewBuffer(nil)
		tee = io.TeeReader(g.r, buf)
	)

	if _, err := tee.Read(bs); err != nil {
		return nil, fmt.Errorf("reader: next section: %w", err)
	}

	switch buf.Len() {
	case 5:
		if err := binary.Read(bytes.NewReader(buf.Bytes()), binary.BigEndian, &secHead); err != nil {
			return nil, fmt.Errorf("binary: next section length: %w", err)
		}

	case 4:
		sec := Section8{section: section8{}}

		if err := sec.ReadFrom(bytes.NewReader(buf.Bytes())); err != nil {
			return nil, fmt.Errorf("section8: %w", err)
		}

		return &sec, nil
	}

	if moreBytesLen := secHead.SectionLength - uint32(buf.Len()); moreBytesLen > 0 {
		bs = make([]byte, moreBytesLen)

		if _, err := tee.Read(bs); err != nil {
			return nil, fmt.Errorf("section %d: next section: %w", secHead.Number, err)
		}

		var sec Section

		switch secHead.Number {
		case 1:
			sec = &Section1{section: section1{}}

		case 2:
			sec = &Section2{section: section2{}}

		case 3:
			sec = &Section3{section: section3{}}

		case 4:
			sec = &Section4{section: section4{}}

		case 5:
			sec = &Section5{section: section5{}}

		case 6:
			sec = &Section6{section: section6{}}

		case 7:
			sec = &Section7{section: section7{}}

		default:
			return nil, fmt.Errorf("section %d: %w", secHead.Number, ErrNotWellFormed)
		}

		if err := sec.ReadFrom(buf); err != nil {
			return nil, fmt.Errorf("binary: next section: %w", err)
		}

		return sec, nil
	}

	return nil, ErrNotWellFormed
}
