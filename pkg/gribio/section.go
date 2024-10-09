package gribio

import (
	"encoding/binary"
	"fmt"
	"io"
)

type GribSection interface {
	Number() int
	Length() int
	Offset() int64
	Reader() io.ReaderAt
}

type SectionReader interface {
	ReadSectionAt(offset int64) (GribSection, error)
}

type gribSectionReader struct {
	io.ReaderAt
}

func NewGribSectionReader(r io.ReaderAt) SectionReader {
	return &gribSectionReader{
		ReaderAt: r,
	}
}

func discernSection(r io.ReaderAt, offset int64) (uint8, uint32, error) {
	bs := make([]byte, 16)

	n, err := r.ReadAt(bs, offset)
	if n >= 4 && bs[0] == '7' && bs[1] == '7' && bs[2] == '7' && bs[3] == '7' {
		return 8, 4, nil
	}

	if n == 16 && bs[0] == 'G' && bs[1] == 'R' && bs[2] == 'I' && bs[3] == 'B' {
		return 0, 16, nil
	}

	if err != nil {
		return 0, 0, fmt.Errorf("read section header: %w", err)
	}

	secLen := binary.BigEndian.Uint32(bs[:4])
	secN := bs[4]

	if secLen <= uint32(len(bs)) {
		return secN, secLen, nil
	}

	return secN, secLen, nil
}

func (r *gribSectionReader) ReadSectionAt(offset int64) (GribSection, error) {
	secNumber, l, err := discernSection(r, offset)
	if err != nil {
		return nil, err
	}

	return &gribSection{
		number: secNumber,
		length: l,
		offset: offset,
		body:   r.ReaderAt,
	}, nil
}

type gribSection struct {
	number uint8
	length uint32
	offset int64
	body   io.ReaderAt
}

func (gs *gribSection) Number() int {
	return int(gs.number)
}

func (gs *gribSection) Length() int {
	return int(gs.length)
}

func (gs *gribSection) Reader() io.ReaderAt {
	return gs.body
}

func (gs *gribSection) Offset() int64 {
	return gs.offset
}
