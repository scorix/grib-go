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

func NewGrib2(r io.ReaderAt) *Grib2 {
	return &Grib2{
		ReaderAt: r,
		r:        gribio.NewGribSectionReader(r),
	}
}

type Grib2 struct {
	io.ReaderAt
	r gribio.SectionReader
}

func (g *Grib2) ReadSectionAt(offset int64) (Section, error) {
	r := gribio.NewGribSectionReader(g.ReaderAt)

	sec, err := r.ReadSectionAt(offset)
	if err != nil {
		return nil, err
	}

	s, err := g.wrapSection(sec)
	if err != nil {
		return nil, fmt.Errorf("wrap: %w", err)
	}

	return s, nil
}

func (g *Grib2) wrapSection(sec gribio.GribSection) (Section, error) {
	sectionFunc, ok := secMap[sec.Number()]
	if !ok {
		return nil, ErrUnknownSection
	}

	s := sectionFunc()
	readLen := sec.Length()

	if s.Number() == 7 {
		readLen = 5
	}

	if err := s.readFrom(sec.Reader(), sec.Offset(), int64(readLen)); err != nil {
		return nil, fmt.Errorf("section %d: %w", sec.Number(), err)
	}

	return s, nil
}

func (g *Grib2) ReadMessageAt(offset int64) (IndexedMessage, error) {
	m, err := g.readIndexedMessageAt(offset)
	if err != nil {
		return nil, fmt.Errorf("read message at %d: %w", offset, err)
	}

	return m, nil
}

func (g *Grib2) readIndexedMessageAt(offset int64) (IndexedMessage, error) {
	m := &message{offset: offset}
	cursor := offset

ReadAllSections:
	for {
		sec, err := g.ReadSectionAt(cursor)
		if err != nil {
			return nil, fmt.Errorf("read section at %d: %w", cursor, err)
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

		cursor += int64(sec.Length())
	}

	return m, nil
}

func (g *Grib2) EachMessage(f func(m IndexedMessage) (next bool, err error)) error {
	var offset int64

	for {
		m, err := g.ReadMessageAt(offset)

		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("read message: %w", err)
		}

		offset += m.GetSize()

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
