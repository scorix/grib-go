package grib2

import (
	"errors"
	"fmt"
	"io"

	"github.com/scorix/grib-go/pkg/gribio"
)

var (
	ErrNotWellFormed     = errors.New("grib file is not well-formed")
	ErrEditionNotMatched = errors.New("grib edition number does not match expected value")
	ErrUnknownSection    = errors.New("encountered an unknown grib section")
)

// SectionFactory uses the factory pattern to create Section instances
type SectionFactory interface {
	CreateSection(number int) (Section, error)
}

type DefaultSectionFactory struct{}

func (f *DefaultSectionFactory) CreateSection(number int) (Section, error) {
	sections := map[int]Section{
		0: &section0{},
		1: &section1{},
		2: &section2{},
		3: &section3{},
		4: &section4{},
		5: &section5{},
		6: &section6{},
		7: &section7{},
		8: &section8{},
	}

	if section, ok := sections[number]; ok {
		return section, nil
	}
	return nil, fmt.Errorf("%w: section number %d", ErrUnknownSection, number)
}

// Grib2Reader uses the strategy pattern to define an interface for reading Grib2 data
type Grib2Reader interface {
	Reader() io.ReaderAt
	ReadSectionAt(offset int64) (Section, error)
	ReadMessageAt(offset int64) (IndexedMessage, error)
	EachMessage(f func(m IndexedMessage) (next bool, err error)) error
}

type grib2 struct {
	io.ReaderAt
	sectionFactory SectionFactory
}

func NewGrib2(r io.ReaderAt) Grib2Reader {
	return &grib2{
		ReaderAt:       r,
		sectionFactory: &DefaultSectionFactory{},
	}
}

func (g *grib2) ReadSectionAt(offset int64) (Section, error) {
	num, length, err := gribio.DiscernSection(g.ReaderAt, offset)
	if err != nil {
		return nil, fmt.Errorf("discern section at offset %d: %w", offset, err)
	}

	s, err := g.sectionFactory.CreateSection(int(num))
	if err != nil {
		return nil, fmt.Errorf("create section %d: %w", num, err)
	}

	readLen := int(length)
	if s.Number() == 7 {
		readLen = 5
	}

	if err := s.readFrom(g.ReaderAt, offset, int64(readLen)); err != nil {
		return nil, fmt.Errorf("read section %d at offset %d: %w", num, offset, err)
	}

	return s, nil
}

func (g *grib2) ReadMessageAt(offset int64) (IndexedMessage, error) {
	m, err := g.readIndexedMessageAt(offset)
	if err != nil {
		return nil, fmt.Errorf("read indexed message at offset %d: %w", offset, err)
	}

	return m, nil
}

func (g *grib2) readIndexedMessageAt(offset int64) (IndexedMessage, error) {
	m := &message{offset: offset}
	cursor := offset

	for {
		sec, err := g.ReadSectionAt(cursor)
		if err != nil {
			return nil, fmt.Errorf("read section at offset %d: %w", cursor, err)
		}

		if err := m.assignSection(sec); err != nil {
			return nil, fmt.Errorf("assign section %d: %w", sec.Number(), err)
		}

		if sec.Number() == 8 {
			break
		}

		cursor += int64(sec.Length())
	}

	return m, nil
}

func (g *grib2) EachMessage(f func(m IndexedMessage) (next bool, err error)) error {
	var offset int64

	for {
		m, err := g.ReadMessageAt(offset)

		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("read message at offset %d: %w", offset, err)
		}

		next, err := f(m)
		if err != nil {
			return fmt.Errorf("process message at offset %d: %w", offset, err)
		}

		if !next {
			return nil
		}

		offset += m.GetSize()
	}
}

func (g *grib2) Reader() io.ReaderAt {
	return g.ReaderAt
}
