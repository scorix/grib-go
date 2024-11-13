package bitio_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/scorix/grib-go/internal/pkg/bitio"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	t.Run("read bits one by one", func(t *testing.T) {
		input := []byte{0b10110100}
		reader := bitio.NewReader(bytes.NewReader(input))

		expected := []byte{1, 0, 1, 1, 0, 1, 0, 0}
		for i, want := range expected {
			got, err := reader.ReadBits(1)
			assert.NoError(t, err)
			assert.Equal(t, uint64(want), got, "bit %d", i)
		}

		// Next read should return EOF
		_, err := reader.ReadBits(1)
		assert.Equal(t, io.EOF, err)
	})

	t.Run("read multiple bits", func(t *testing.T) {
		input := []byte{0b10110100, 0b11000000}
		reader := bitio.NewReader(bytes.NewReader(input))

		tests := []struct {
			bits uint8
			want uint64
		}{
			{4, 0b1011},
			{3, 0b010},
			{5, 0b01100},
		}

		for _, tt := range tests {
			got, err := reader.ReadBits(tt.bits)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got, "read %d bits, want %08b, got %08b", tt.bits, tt.want, got)
		}
	})

	t.Run("read full bytes", func(t *testing.T) {
		input := []byte{0b10110100, 0b11000000}
		reader := bitio.NewReader(bytes.NewReader(input))

		buf := make([]byte, 2)
		n, err := reader.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, 2, n)
		assert.Equal(t, input, buf)
	})

	t.Run("read partial bytes", func(t *testing.T) {
		input := []byte{0b10110100, 0b11000000}
		reader := bitio.NewReader(bytes.NewReader(input))

		// Read 12 bits (1.5 bytes)
		got, err := reader.ReadBits(12)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0b101101001100), got)
	})

	t.Run("error on too many bits", func(t *testing.T) {
		reader := bitio.NewReader(bytes.NewReader([]byte{0xFF}))
		_, err := reader.ReadBits(65)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read more than 64 bits")
	})

	t.Run("read after partial bits", func(t *testing.T) {
		input := []byte{0b10110100, 0b11000000, 0b10101010}
		reader := bitio.NewReader(bytes.NewReader(input))

		// First read 4 bits
		got, err := reader.ReadBits(4)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0b1011), got)

		// Now read full bytes
		buf := make([]byte, 2)
		n, err := reader.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, 2, n)
		// Should get remaining 4 bits from first byte + second byte
		assert.Equal(t, []byte{0b01001100, 0b00001010}, buf)
	})

	t.Run("align to next byte", func(t *testing.T) {
		input := []byte{0b10110100, 0b11000000}
		reader := bitio.NewReader(bytes.NewReader(input))

		// Read 3 bits
		got, err := reader.ReadBits(3)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0b101), got)

		// Align to next byte
		discarded := reader.Align()
		assert.Equal(t, uint8(5), discarded) // 5 bits were remaining in the buffer

		// Read next byte - should start from second byte
		got, err = reader.ReadBits(8)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0b11000000), got)
	})
}
