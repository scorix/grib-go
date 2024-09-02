package datapacking

import "io"

type BitReader interface {
	io.Reader
	ReadBits(uint8) (uint64, error)
	Align() uint8
}

type CountBitReader interface {
	BitReader
	Count() int
}

type countBitReader struct {
	BitReader
	count int
}

func NewCountBitReader(r BitReader) CountBitReader {
	return &countBitReader{BitReader: r}
}

func (r *countBitReader) ReadBits(n uint8) (uint64, error) {
	r.count += int(n)

	return r.BitReader.ReadBits(n)
}

func (r *countBitReader) Read(p []byte) (n int, err error) {
	r.count += len(p) * 8

	return r.BitReader.Read(p)
}

func (r *countBitReader) Count() int {
	return r.count / 8
}
