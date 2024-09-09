package grib2

import (
	"errors"
	"fmt"
	"io"

	"github.com/scorix/grib-go/pkg/gribio"
)

var (
	ErrNotWellFormed     = fmt.Errorf("grib is not well formed")
	ErrEditionNotMatched = fmt.Errorf("edition number is not matched")
	ErrUnknownSection    = fmt.Errorf("unknown section")
)

var (
	secMap = map[int]func() Section{
		0: func() Section { return &section0{} },
		1: func() Section { return &section1{} },
		2: func() Section { return &section2{} },
		3: func() Section { return &section3{} },
		4: func() Section { return &section4{} },
		5: func() Section { return &section5{} },
		6: func() Section { return &section6{} },
		7: func() Section { return &section7{} },
		8: func() Section { return &section8{} },
	}
)

func NewGrib2(r io.Reader) *Grib2 {
	return &Grib2{
		r: gribio.NewGribSectionReader(r),
	}
}

type Grib2 struct {
	r      gribio.SectionReader
	offset int64
}

func (g *Grib2) ReadSection() (Section, error) {
	sec, err := g.r.ReadSection()
	if err != nil {
		return nil, err
	}

	sectionFunc, ok := secMap[sec.Number()]
	if !ok {
		return nil, ErrUnknownSection
	}

	s := sectionFunc()
	if err := s.readFrom(sec.Reader()); err != nil {
		return nil, err
	}

	return s, nil
}

func (g *Grib2) ReadMessage() (IndexedMessage, error) {
	m, err := g.readIndexedMessage(g.offset)
	if err != nil {
		return nil, err
	}

	g.offset += m.GetSize()

	return m, nil
}

func (g *Grib2) readIndexedMessage(offset int64) (IndexedMessage, error) {
	m := &message{offset: offset}

	for {
		sec, err := g.ReadSection()
		if err != nil {
			return nil, fmt.Errorf("read section: %w", err)
		}

		switch sec.Number() {
		case 0:
			m.sec0 = sec.(*section0)
			m.dataOffset = offset + int64(m.sec0.Length())
		case 1:
			m.sec1 = sec.(*section1)
			m.dataOffset += int64(m.sec1.Length())
		case 2:
			m.sec2 = sec.(*section2)
			m.dataOffset += int64(m.sec2.Length())
		case 3:
			m.sec3 = sec.(*section3)
			m.dataOffset += int64(m.sec3.Length())
		case 4:
			m.sec4 = sec.(*section4)
			m.dataOffset += int64(m.sec4.Length())
		case 5:
			m.sec5 = sec.(*section5)
			m.dataOffset += int64(m.sec5.Length())
		case 6:
			m.sec6 = sec.(*section6)
			m.dataOffset += int64(m.sec6.Length())
		case 7:
			m.sec7 = sec.(*section7)
			m.dataOffset += int64(m.sec7.Length() - len(m.sec7.Data))
		case 8:
			m.sec8 = sec.(*section8)

			return m, nil

		default:
			return nil, fmt.Errorf("unknown section number: %d", sec.Number())
		}
	}
}

func (g *Grib2) EachMessage(f func(IndexedMessage) error) error {
	for {
		m, err := g.ReadMessage()

		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("read message: %w", err)
		}

		if err := f(m); err != nil {
			return err
		}
	}
}
