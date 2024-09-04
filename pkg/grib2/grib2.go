package grib2

import (
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
	r gribio.SectionReader
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

func (g *Grib2) ReadMessage() (Message, error) {
	m := &message{}

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

			return m, nil

		default:
			return nil, fmt.Errorf("unknown section number: %d", sec.Number())
		}
	}
}
