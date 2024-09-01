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
