package grib2

import (
	"errors"
	"fmt"
	"io"
	"sync/atomic"

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

type grib2reader interface {
	io.Reader
	io.ReaderAt
}

func NewGrib2(r grib2reader) *Grib2 {
	return &Grib2{
		ReaderAt: r,
		r:        gribio.NewGribSectionReader(r),
	}
}

type Grib2 struct {
	io.ReaderAt
	r      gribio.SectionReader
	offset int64
	cursor int64
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

	if err := s.readFrom(g.ReaderAt, g.cursor, int64(sec.Length())); err != nil {
		return nil, fmt.Errorf("section %d: %w", sec.Number(), err)
	}

	atomic.AddInt64(&g.cursor, int64(s.Length()))

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

func (g *Grib2) ReadMessageAt(offset int64) (IndexedMessage, error) {
	m, err := g.readIndexedMessage(offset)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (g *Grib2) readIndexedMessage(offset int64) (IndexedMessage, error) {
	m := &message{offset: offset}

ReadAllSections:
	for {
		sec, err := g.ReadSection()
		if err != nil {
			return nil, fmt.Errorf("read section: %w", err)
		}

		switch sec.Number() {
		case 0:
			m.sec0 = sec.(*section0)
		case 1:
			m.sec1 = sec.(*section1)
		case 2:
			m.sec2 = sec.(*section2)
		case 3:
			m.sec3 = sec.(*section3)
		case 4:
			m.sec4 = sec.(*section4)
		case 5:
			m.sec5 = sec.(*section5)
		case 6:
			m.sec6 = sec.(*section6)
		case 7:
			m.sec7 = sec.(*section7)
		case 8:
			m.sec8 = sec.(*section8)

			break ReadAllSections

		default:
			return nil, fmt.Errorf("unknown section number: %d", sec.Number())
		}
	}

	return m, nil
}

func (g *Grib2) EachMessage(f func(m IndexedMessage) (next bool, err error)) error {
	for {
		m, err := g.ReadMessage()

		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("read message: %w", err)
		}

		next, err := f(m)
		if err != nil {
			return err
		}

		if next {
			continue
		}

		return nil
	}
}
