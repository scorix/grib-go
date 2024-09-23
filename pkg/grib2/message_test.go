package grib2_test

import (
	"errors"
	"os"
	"testing"

	codes "github.com/scorix/go-eccodes"
	cio "github.com/scorix/go-eccodes/io"
	"github.com/scorix/grib-go/pkg/grib2"
	grib "github.com/scorix/grib-go/pkg/grib2"
	"github.com/scorix/grib-go/pkg/grib2/regulation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/mmap"
)

func TestMessageReader_ReadLL(t *testing.T) {
	t.Parallel()

	// grib_set -r -s packingType=grid_simple pkg/testdata/hpbl.grib2 pkg/testdata/hpbl.grib2.out
	const filename = "../testdata/hpbl.grib2.out"

	s, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		t.Skipf("%s not exist", filename)
	}

	t.Run(s.Name(), func(t *testing.T) {
		f, err := os.Open(filename)
		require.NoError(t, err)
		defer f.Close()

		cf, err := cio.OpenFile(f.Name(), "r")
		require.NoError(t, err)
		defer cf.Close()

		cgrib, err := codes.OpenFile(cf)
		require.NoError(t, err)
		defer cgrib.Close()

		handle, err := cgrib.Handle()
		require.NoError(t, err)
		defer handle.Close()

		cmsg := handle.Message()
		defer cmsg.Close()

		iter, err := cmsg.Iterator()
		require.NoError(t, err)
		defer iter.Close()

		g := grib.NewGrib2(f)

		msg, err := g.ReadMessage()
		require.NoError(t, err)
		require.NotNil(t, msg)

		assert.Equal(t, 1, msg.GetTypeOfFirstFixedSurface())
		assert.Equal(t, int64(2206439), msg.GetSize())
		assert.Equal(t, int64(0), msg.GetOffset())
		assert.Equal(t, int64(175), msg.GetDataOffset())

		mm, err := mmap.Open(f.Name())
		require.NoError(t, err)

		reader, err := grib2.NewSimplePackingMessageReader(mm, msg)
		require.NoError(t, err)

		for i := 0; iter.HasNext(); i++ {
			lat, lng, val, _ := iter.Next()
			lat32 := regulation.DegreedLatitudeLongitude(int(lat * 1e6))
			lng32 := regulation.DegreedLatitudeLongitude(int(lng * 1e6))

			grd, err := msg.GetGridPointFromLL(lat32, lng32)
			require.NoError(t, err)
			require.Equalf(t, i, grd, "expect: (%f,%f,%d), actual: (%f,%f,%d)", lat, lng, i, lat32, lng32, grd)

			{
				v, err := reader.ReadLL(lat32, lng32) // grib_get -l 90,0 pkg/testdata/hpbl.grib2.out
				require.NoError(t, err)
				assert.InDelta(t, float32(val), float32(v), 1e-5)
			}

			{
				// read again
				v, err := reader.ReadLL(lat32, lng32) // grib_get -l 90,0 pkg/testdata/hpbl.grib2.out
				require.NoError(t, err)
				assert.InDelta(t, float32(val), float32(v), 1e-5)
			}

		}
	})
}
