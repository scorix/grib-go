package bitio

import (
	"bufio"
	"io"
)

// readerAndByteReader combines io.Reader and io.ByteReader interfaces
// for efficient bit-level reading operations
type readerAndByteReader interface {
	io.Reader
	io.ByteReader
}

// Reader provides bit-level reading capabilities over an io.Reader.
// It buffers partial bytes and handles unaligned bit reads efficiently.
type Reader struct {
	in    readerAndByteReader
	cache byte // Holds partially read bits from the last byte
	bits  byte // Number of valid bits remaining in cache (0-8)

	// TryError stores the first error from TryXXX() methods
	// to allow for chainable error handling
	TryError error
}

// NewReader returns a new Reader using the specified io.Reader as the input (source).
func NewReader(in io.Reader) *Reader {
	bin, ok := in.(readerAndByteReader)
	if !ok {
		bin = bufio.NewReader(in)
	}
	return &Reader{in: bin}
}

// Read reads up to len(p) bytes (8 * len(p) bits) from the underlying reader.
//
// Read implements io.Reader, and gives a byte-level view of the bit stream.
// This will give best performance if the underlying io.Reader is aligned
// to a byte boundary (else all the individual bytes are assembled from multiple bytes).
// Byte boundary can be ensured by calling Align().
func (r *Reader) Read(p []byte) (n int, err error) {
	// r.bits will be the same after reading 8 bits, so we don't need to update that.
	if r.bits == 0 {
		return r.in.Read(p)
	}

	for ; n < len(p); n++ {
		if p[n], err = r.readUnalignedByte(); err != nil {
			return
		}
	}

	return
}

// ReadBits reads exactly n bits (0-64) and returns them as the lowest bits of uint64.
// The bits are read from most significant to least significant position.
func (r *Reader) ReadBits(n uint8) (u uint64, err error) {
	// Fast path: all bits available in cache
	if n < r.bits {
		shift := r.bits - n
		u = uint64(r.cache >> shift)
		r.cache &= 1<<shift - 1 // Clear read bits
		r.bits = shift
		return
	}

	if n > r.bits {
		// all cache bits needed, and it's not even enough so more will be read
		if r.bits > 0 {
			u = uint64(r.cache)
			n -= r.bits
		}
		// Read whole bytes
		for n >= 8 {
			b, err2 := r.in.ReadByte()
			if err2 != nil {
				return 0, err2
			}
			u = u<<8 + uint64(b)
			n -= 8
		}
		// Read last fraction, if any
		if n > 0 {
			if r.cache, err = r.in.ReadByte(); err != nil {
				return 0, err
			}
			shift := 8 - n
			u = u<<n + uint64(r.cache>>shift)
			r.cache &= 1<<shift - 1
			r.bits = shift
		} else {
			r.bits = 0
		}
		return u, nil
	}

	// cache has exactly as many as needed
	r.bits = 0 // no need to clear cache, will be overwritten on next read
	return uint64(r.cache), nil
}

// ReadByte reads the next 8 bits and returns them as a byte.
//
// ReadByte implements io.ByteReader.
func (r *Reader) ReadByte() (b byte, err error) {
	// r.bits will be the same after reading 8 bits, so we don't need to update that.
	if r.bits == 0 {
		return r.in.ReadByte()
	}
	return r.readUnalignedByte()
}

// readUnalignedByte handles reading 8 bits that span byte boundaries.
// It combines cached bits with bits from the next byte.
func (r *Reader) readUnalignedByte() (b byte, err error) {
	bits := r.bits
	b = r.cache << (8 - bits) // Place cached bits in their final position

	r.cache, err = r.in.ReadByte()
	if err != nil {
		return 0, err
	}

	b |= r.cache >> bits   // Combine with bits from new byte
	r.cache &= 1<<bits - 1 // Keep remaining bits in cache
	return
}

// ReadBool reads the next bit, and returns true if it is 1.
func (r *Reader) ReadBool() (b bool, err error) {
	if r.bits == 0 {
		r.cache, err = r.in.ReadByte()
		if err != nil {
			return
		}
		b = (r.cache & 0x80) != 0
		r.cache, r.bits = r.cache&0x7f, 7
		return
	}

	r.bits--
	b = (r.cache & (1 << r.bits)) != 0
	r.cache &= 1<<r.bits - 1
	return
}

// Align aligns the bit stream to a byte boundary,
// so next read will read/use data from the next byte.
// Returns the number of unread / skipped bits.
func (r *Reader) Align() (skipped uint8) {
	skipped = r.bits
	r.bits = 0 // no need to clear cache, will be overwritten on next read
	return
}
