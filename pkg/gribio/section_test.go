package gribio_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/scorix/grib-go/pkg/gribio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSectionReader_ReadSection(t *testing.T) {
	t.Parallel()

	f, err := os.Open("../testdata/temp.grib2")
	require.NoError(t, err)
	defer f.Close()

	sections := []struct {
		number int
		length int
		offset int64
		body   []byte
	}{
		{number: 0, length: 16, offset: 0, body: []byte{0x47, 0x52, 0x49, 0x42, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x1a, 0xe}},
		{number: 1, length: 21, offset: 16, body: []byte{0x0, 0x0, 0x0, 0x15, 0x1, 0x0, 0x4a, 0x0, 0x5, 0x1d, 0x1, 0x1, 0x7, 0xe7, 0x7, 0xb, 0x0, 0x0, 0x0, 0x0, 0x1}},
		{number: 3, length: 72, offset: 37},
		{number: 4, length: 34, offset: 109},
		{number: 5, length: 21, offset: 143},
		{number: 6, length: 6, offset: 164},
		{number: 7, length: 203_104, offset: 170},
		{number: 8, length: 4, offset: 203_274, body: []byte{'7', '7', '7', '7'}},
	}

	for _, section := range sections {
		t.Run(fmt.Sprintf("section %d", section.number), func(t *testing.T) {
			num, length, err := gribio.DiscernSection(f, section.offset)
			require.NoError(t, err)

			assert.Equal(t, section.number, int(num))
			assert.Equal(t, section.length, int(length))

			if section.body != nil {
				body := make([]byte, section.length)
				n, err := f.ReadAt(body, section.offset)
				require.NoError(t, err)
				assert.Equal(t, section.length, int(n))
				assert.Equal(t, section.body, body)
			}
		})
	}
}
