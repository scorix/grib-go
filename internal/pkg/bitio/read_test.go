package bitio_test

import (
	"testing"

	"github.com/scorix/grib-go/internal/pkg/bitio"
	"github.com/stretchr/testify/assert"
)

func TestReadBits(t *testing.T) {
	tests := []struct {
		name       string
		data       []byte
		offset     uint8
		numBits    uint8
		want       uint64
		wantErr    bool
		errMessage string
	}{
		{
			name:    "read single bit - 1",
			data:    []byte{0b10000000},
			offset:  0,
			numBits: 1,
			want:    1,
		},
		{
			name:    "read single bit - 0",
			data:    []byte{0b00000000},
			offset:  0,
			numBits: 1,
			want:    0,
		},
		{
			name:    "read multiple bits from single byte",
			data:    []byte{0b10110100},
			offset:  0,
			numBits: 4,
			want:    0b1011,
		},
		{
			name:    "read with offset in byte",
			data:    []byte{0b10110100},
			offset:  2,
			numBits: 3,
			want:    0b110,
		},
		{
			name:    "read across byte boundary",
			data:    []byte{0b10110100, 0b11000000},
			offset:  6,
			numBits: 4,
			want:    0b0011,
		},
		{
			name:    "read full byte",
			data:    []byte{0b10110100},
			offset:  0,
			numBits: 8,
			want:    0b10110100,
		},
		{
			name:    "read multiple bytes",
			data:    []byte{0xFF, 0xFF},
			offset:  0,
			numBits: 16,
			want:    0xFFFF,
		},
		{
			name:       "error on too many bits",
			data:       []byte{0xFF},
			offset:     0,
			numBits:    65,
			wantErr:    true,
			errMessage: "cannot read more than 64 bits at once",
		},
		{
			name:       "error on EOF",
			data:       []byte{0xFF},
			offset:     0,
			numBits:    16,
			wantErr:    true,
			errMessage: "EOF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bitio.ReadBits(tt.data, tt.offset, tt.numBits)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMessage != "" {
					assert.Equal(t, tt.errMessage, err.Error())
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
