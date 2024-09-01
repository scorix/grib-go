package grib

import (
	"io"
)

type Section interface {
	fixedLengthSection
	ReadFrom(r io.Reader) error
}

type fixedLengthSection interface {
	SectionLength() int
	SectionNumber() int
}
