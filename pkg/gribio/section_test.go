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
	f, err := os.Open("../testdata/temp.grib2")
	require.NoError(t, err)

	sr := gribio.NewGribSectionReader(f)

	sections := []struct {
		number int
		length int
	}{
		{number: 0, length: 16},
		{number: 1, length: 21},
		{number: 3, length: 72},
		{number: 4, length: 34},
		{number: 5, length: 21},
		{number: 6, length: 6},
		{number: 7, length: 203_104},
		{number: 8, length: 4},
	}

	for _, section := range sections {
		t.Run(fmt.Sprintf("section %d", section.number), func(t *testing.T) {
			sec, err := sr.ReadSection()
			require.NoError(t, err)

			assert.Equal(t, section.number, sec.Number())
			assert.Equal(t, section.length, sec.Length())
		})
	}
}
