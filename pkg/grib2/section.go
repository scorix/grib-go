package grib2

import (
	"io"
)

type Section interface {
	Number() int
	Length() int
	readFrom(r io.ReaderAt, offset int64, length int64) error
}

type Section0 interface {
	Section
	GetEditionNumber() int
	GetDiscipline() int
	GetGribLength() int
}
