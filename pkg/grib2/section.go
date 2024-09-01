package grib2

import (
	"io"
)

type Section interface {
	Number() int
	Length() int
	readFrom(r io.Reader) error
	// Bytes() []byte
}

type Section0 interface {
	Section
	GetEditionNumber() int
	GetDiscipline() int
	GetGribLength() int
}
