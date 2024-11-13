package bitio

import (
	"errors"
	"io"
	"sync"
)

type Reader struct {
	r        io.Reader
	buffer   byte       // Current byte buffer
	bitCount uint8      // Number of unread bits in buffer
	mu       sync.Mutex // Protects buffer and bitCount
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

func (r *Reader) ReadBits(n uint8) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if n > 64 {
		return 0, errors.New("cannot read more than 64 bits")
	}

	// Read initial byte if buffer is empty
	if r.bitCount == 0 {
		buf := make([]byte, 1)
		_, err := r.r.Read(buf)
		if err != nil {
			return 0, err
		}
		r.buffer = buf[0]
		r.bitCount = 8
	}

	var result uint64
	remainingBits := n

	for remainingBits > 0 {
		if r.bitCount == 0 {
			buf := make([]byte, 1)
			_, err := r.r.Read(buf)
			if err != nil {
				return 0, err
			}
			r.buffer = buf[0]
			r.bitCount = 8
		}

		bitsToRead := min(remainingBits, r.bitCount)

		// Left align the bits we want
		shift := 8 - r.bitCount
		// Then right align them
		bits := uint64((r.buffer << shift) >> (8 - bitsToRead))

		result = (result << bitsToRead) | bits

		r.bitCount -= bitsToRead
		remainingBits -= bitsToRead
	}

	return result, nil
}

func (r *Reader) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p); i++ {
		// Read 8 bits for each byte
		val, err := r.ReadBits(8)
		if err != nil {
			return i, err
		}
		p[i] = byte(val)
		n++
	}
	return n, nil
}

// Align discards any remaining bits in the current byte and aligns the reader
// to the next byte boundary. Returns the number of bits that were discarded.
func (r *Reader) Align() uint8 {
	r.mu.Lock()
	defer r.mu.Unlock()

	discarded := r.bitCount
	r.bitCount = 0
	r.buffer = 0
	return discarded
}
