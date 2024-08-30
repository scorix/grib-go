package grib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
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
	r          gribReader
	dataReader datapacking.UnpackReader

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

	type secHead struct {
		SectionLength uint32
		Number        uint8
	}

	var (
		head secHead
		bs   = make([]byte, 5)
	)

	n, err := g.r.Read(bs)
	if err != nil {
		return nil, fmt.Errorf("reader: next section: %w", err)
	}

	switch n {
	case 5:
		if err := binary.Read(bytes.NewReader(bs), binary.BigEndian, &head); err != nil {
			return nil, fmt.Errorf("binary: next section length: %w", err)
		}

	case 4:
		sec := Section8{section8: section8{}}

		if err := sec.ReadFrom(bytes.NewReader(bs)); err != nil {
			return nil, fmt.Errorf("section8: %w", err)
		}

		return &sec, nil
	}

	if moreBytesLen := head.SectionLength - uint32(n); moreBytesLen > 0 {
		bs = append(bs, make([]byte, moreBytesLen)...)

		if _, err := g.r.Read(bs[n:]); err != nil {
			return nil, fmt.Errorf("section %d: next section: %w", head.Number, err)
		}

		var section Section

		switch head.Number {
		case 1:
			sec := &Section1{section1: section1{}}
			if err := sec.ReadFrom(bytes.NewReader(bs)); err != nil {
				return nil, fmt.Errorf("section %d: %w", head.Number, err)
			}

			section = sec

		case 2:
			sec := &Section2{section2: section2{}}
			if err := sec.ReadFrom(bytes.NewReader(bs)); err != nil {
				return nil, fmt.Errorf("section %d: %w", head.Number, err)
			}

			section = sec

		case 3:
			sec := &Section3{section3: section3{}}
			if err := sec.ReadFrom(bytes.NewReader(bs)); err != nil {
				return nil, fmt.Errorf("section %d: %w", head.Number, err)
			}

			section = sec

		case 4:
			sec := &Section4{defSection4: defSection4{}}
			if err := sec.ReadFrom(bytes.NewReader(bs)); err != nil {
				return nil, fmt.Errorf("section %d: %w", head.Number, err)
			}

			section = sec

		case 5:
			sec := &Section5{defSection5: defSection5{}}
			if err := sec.ReadFrom(io.MultiReader(bytes.NewReader(bs), g.r)); err != nil {
				return nil, fmt.Errorf("section %d: %w", head.Number, err)
			}

			g.dataReader = sec.DataPackingReader
			section = sec

		case 6:
			sec := &Section6{section6: section6{}}
			if err := sec.ReadFrom(bytes.NewReader(bs)); err != nil {
				return nil, fmt.Errorf("section %d: %w", head.Number, err)
			}

			section = sec

		case 7:
			sec := &Section7{section7: section7{}, dataReader: g.dataReader}
			if err := sec.ReadFrom(bytes.NewReader(bs)); err != nil {
				return nil, fmt.Errorf("section %d: %w", head.Number, err)
			}

			section = sec

		default:
			return nil, fmt.Errorf("section %d: %w", head.Number, ErrNotWellFormed)
		}

		return section, nil
	}

	return nil, ErrNotWellFormed
}
