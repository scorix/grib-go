package gribio

import (
	"encoding/binary"
	"fmt"
	"io"
)

// DiscernSection identifies the section type and length
func DiscernSection(r io.ReaderAt, offset int64) (number uint8, length uint32, err error) {
	bs := make([]byte, 5)
	n, err := r.ReadAt(bs, offset)

	if n >= 4 {
		if bs[0] == '7' && bs[1] == '7' && bs[2] == '7' && bs[3] == '7' {
			return 8, 4, nil
		}
		if bs[0] == 'G' && bs[1] == 'R' && bs[2] == 'I' && bs[3] == 'B' {
			return 0, 16, nil
		}
	}

	if err != nil {
		return 0, 0, fmt.Errorf("section header: %w", err)
	}

	length = binary.BigEndian.Uint32(bs[:4])
	number = bs[4]

	return number, length, nil
}
