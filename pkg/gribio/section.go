package gribio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

type GribSection interface {
	Number() int
	Length() int
	Reader() io.Reader
}

type SectionReader interface {
	ReadSection() (GribSection, error)
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

func discernSection(r io.Reader) (secN uint8, secLen uint32, err error) {
	bs := make([]byte, 16)

	if _, err := r.Read(bs[:4]); err != nil {
		return 0, 0, fmt.Errorf("read section length: %w", err)
	}

	if bs[0] == 'G' && bs[1] == 'R' && bs[2] == 'I' && bs[3] == 'B' {
		return 0, 16, nil
	}

	if bs[0] == '7' && bs[1] == '7' && bs[2] == '7' && bs[3] == '7' {
		return 8, 4, nil
	}

	secLen = binary.BigEndian.Uint32(bs[:4])

	if _, err := r.Read(bs[4:5]); err != nil {
		return 0, 0, fmt.Errorf("read section number: %w", err)
	}

	secN = bs[4]

	return
}

func (r *gribSectionReader) ReadSection() (GribSection, error) {
	r.m.Lock()
	defer r.m.Unlock()

	n, l, err := discernSection(r.raw)
	if err != nil {
		return nil, err
	}

	var section bytes.Buffer

	if _, err := io.Copy(&section, io.MultiReader(r.buf, io.LimitReader(r.raw, int64(int(l)-r.buf.Len())))); err != nil {
		return nil, err
	}

	r.buf.Reset()

	return &gribSection{
		number: n,
		length: l,
		body:   section.Bytes(),
	}, nil
}

type gribSection struct {
	number uint8
	length uint32
	body   []byte
}

func (gs *gribSection) Number() int {
	return int(gs.number)
}

func (gs *gribSection) Length() int {
	return int(gs.length)
}

func (gs *gribSection) Reader() io.Reader {
	return bytes.NewReader(gs.body)
}

// lazy load section reader
type lazyGribSectionReader struct {
	offset int64
	raw    io.ReadSeeker
	m      sync.Mutex
}

func NewLazyGribSectionReader(r io.ReadSeeker) SectionReader {
	return &lazyGribSectionReader{
		raw: r,
	}
}

type lazyGribSection struct {
	offset     int64
	number     uint8
	length     uint32
	bodyOffset int64
	bodySize   int64
	body       io.ReadSeeker
}

func (gs *lazyGribSection) Number() int {
	return int(gs.number)
}

func (gs *lazyGribSection) Length() int {
	return int(gs.length)
}

func (gs *lazyGribSection) Reader() io.Reader {
	if _, err := gs.body.Seek(gs.offset, io.SeekStart); err != nil {
		r, w := io.Pipe()
		w.CloseWithError(err)

		return r
	}

	return io.LimitReader(gs.body, int64(gs.length))
}

func (r *lazyGribSectionReader) ReadSection() (GribSection, error) {
	r.m.Lock()
	defer r.m.Unlock()

	if _, err := r.raw.Seek(r.offset, io.SeekStart); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)

	n, l, err := discernSection(io.TeeReader(r.raw, buf))
	if err != nil {
		return nil, err
	}

	gs := lazyGribSection{
		offset:     r.offset,
		number:     n,
		length:     l,
		bodyOffset: int64(buf.Len()),
		bodySize:   int64(int(l) - buf.Len()),
		body:       r.raw,
	}

	r.offset += int64(l)

	return &gs, nil
}
