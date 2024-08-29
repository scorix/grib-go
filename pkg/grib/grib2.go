package grib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

var (
	GribLiteral = [4]byte{'G', 'R', 'I', 'B'}
)

var (
	ErrNotWellFormed = fmt.Errorf("grib is not well formed")
)

type gribReader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

type Grib struct {
	r gribReader
}

func New(r gribReader) *Grib {
	return &Grib{r: r}
}

// ReadSection0
func (g *Grib) ReadSection0() (*Section0, error) {
	if _, err := g.r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	bs := make([]byte, 16)

	n, err := g.r.Read(bs)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	if n != 16 {
		return nil, fmt.Errorf("should read 16 bytes, but read %d", n)
	}

	var sec Section0
	if err := binary.Read(bytes.NewReader(bs), binary.BigEndian, &sec); err != nil {
		return nil, fmt.Errorf("read section0: %w", err)
	}

	if sec.GribLiteral != GribLiteral {
		return nil, ErrNotWellFormed
	}

	return &sec, nil
}
