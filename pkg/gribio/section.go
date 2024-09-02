package gribio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

type SectionReader interface {
	ReadSection() (*GribSection, error)
}

type gribSectionReader struct {
	raw io.Reader
	buf *bytes.Buffer
	m   sync.Mutex
}

func NewGribSectionReader(r io.Reader) SectionReader {
	buf := bytes.NewBuffer(nil)

	return &gribSectionReader{
		raw: io.TeeReader(r, buf),
		buf: buf,
	}
}

func (r *gribSectionReader) discernSection() (secN uint8, secLen uint32, err error) {
	bs := make([]byte, 16)

	if _, err := r.raw.Read(bs[:4]); err != nil {
		return 0, 0, fmt.Errorf("read section length: %w", err)
	}

	if bs[0] == 'G' && bs[1] == 'R' && bs[2] == 'I' && bs[3] == 'B' {
		return 0, 16, nil
	}

	if bs[0] == '7' && bs[1] == '7' && bs[2] == '7' && bs[3] == '7' {
		return 8, 4, nil
	}

	secLen = binary.BigEndian.Uint32(bs[:4])

	if _, err := r.raw.Read(bs[4:5]); err != nil {
		return 0, 0, fmt.Errorf("read section number: %w", err)
	}

	secN = bs[4]

	return
}

func (r *gribSectionReader) ReadSection() (*GribSection, error) {
	r.m.Lock()
	defer r.m.Unlock()

	n, l, err := r.discernSection()
	if err != nil {
		return nil, err
	}

	var section bytes.Buffer

	if _, err := io.Copy(&section, io.MultiReader(r.buf, io.LimitReader(r.raw, int64(int(l)-r.buf.Len())))); err != nil {
		return nil, err
	}

	r.buf.Reset()

	return &GribSection{
		number: n,
		length: l,
		body:   section.Bytes(),
	}, nil
}

type GribSection struct {
	number uint8
	length uint32
	body   []byte
}

func (gs GribSection) Number() int {
	return int(gs.number)
}

func (gs GribSection) Length() int {
	return int(gs.length)
}

func (gs GribSection) Reader() io.ReadSeeker {
	return bytes.NewReader(gs.body)
}
