package grib2_test

import (
	"errors"
	"os"
	"testing"

	codes "github.com/scorix/go-eccodes"
	cio "github.com/scorix/go-eccodes/io"
	"github.com/scorix/grib-go/pkg/grib2"
	grib "github.com/scorix/grib-go/pkg/grib2"
	gridpoint "github.com/scorix/grib-go/pkg/grib2/drt/grid_point"
	"github.com/scorix/grib-go/pkg/grib2/gdt"
	"github.com/scorix/grib-go/pkg/grib2/regulation"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/mmap"
)

func TestMessageReader_ReadLL(t *testing.T) {
	t.Parallel()

	// grib_set -r -s packingType=grid_simple pkg/testdata/hpbl.grib2 pkg/testdata/hpbl.grib2.out
	t.Run("regular_ll", func(t *testing.T) {
		const filename = "../testdata/hpbl.grib2.out"

		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			t.Skipf("%s not exist", filename)
		}

		t.Parallel()

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

		msg, err := g.ReadMessageAt(0)
		require.NoError(t, err)
		require.NotNil(t, msg)

		require.Equal(t, 1, msg.GetTypeOfFirstFixedSurface())
		require.Equal(t, int64(2206439), msg.GetSize())
		require.Equal(t, int64(0), msg.GetOffset())
		require.Equal(t, int64(175), msg.GetDataOffset())

		mm, err := mmap.Open(f.Name())
		require.NoError(t, err)

		sm, err := msg.GetScanningMode()
		require.NoError(t, err)

		reader, err := grib2.NewSimplePackingMessageReader(mm, msg.GetOffset(), msg.GetSize(), msg.GetDataOffset(), msg.GetDataRepresentationTemplate().(*gridpoint.SimplePacking), sm)
		require.NoError(t, err)

		for i := 0; iter.HasNext(); i++ {
			lat, lng, val, _ := iter.Next()
			lat32 := regulation.DegreedLatitudeLongitude(int(lat * 1e6))
			lng32 := regulation.DegreedLatitudeLongitude(int(lng * 1e6))

			_, _, grd, err := msg.GetGridPointFromLL(lat32, lng32)
			require.NoError(t, err)
			require.Equalf(t, i, grd, "expect: (%f,%f,%d), actual: (%f,%f,%d)", lat, lng, i, lat32, lng32, grd)

			{
				_, _, v, err := reader.ReadLL(lat32, lng32) // grib_get -l 90,0 pkg/testdata/hpbl.grib2.out
				require.NoError(t, err)
				require.InDelta(t, float32(val), float32(v), 1e-5)
			}

			{
				// read again
				_, _, v, err := reader.ReadLL(lat32, lng32) // grib_get -l 90,0 pkg/testdata/hpbl.grib2.out
				require.NoError(t, err)
				require.InDelta(t, float32(val), float32(v), 1e-5)
			}

		}
	})

	// grib_set -r -s packingType=grid_simple pkg/testdata/regular_gg.grib2 pkg/testdata/regular_gg.grib2.out
	t.Run("regular_gg", func(t *testing.T) {
		const filename = "../testdata/regular_gg.grib2.out"

		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			t.Skipf("%s not exist", filename)
		}

		t.Parallel()

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

		msg, err := g.ReadMessageAt(0)
		require.NoError(t, err)
		require.NotNil(t, msg)

		require.Equal(t, 1, msg.GetTypeOfFirstFixedSurface())
		require.Equal(t, int64(10027187), msg.GetSize())
		require.Equal(t, int64(0), msg.GetOffset())
		require.Equal(t, int64(175), msg.GetDataOffset())

		mm, err := mmap.Open(f.Name())
		require.NoError(t, err)

		sm, err := msg.GetScanningMode()
		require.NoError(t, err)

		reader, err := grib2.NewSimplePackingMessageReader(mm, msg.GetOffset(), msg.GetSize(), msg.GetDataOffset(), msg.GetDataRepresentationTemplate().(*gridpoint.SimplePacking), sm)
		require.NoError(t, err)

		errorItemsCount := 0
		totalItemsCount := 0

		for i := 0; iter.HasNext(); i++ {
			lat, lng, val, _ := iter.Next()
			lat32 := regulation.DegreedLatitudeLongitude(int(lat * 1e6))
			lng32 := regulation.DegreedLatitudeLongitude(int(lng * 1e6))

			_, _, grd, err := msg.GetGridPointFromLL(lat32, lng32)
			require.NoError(t, err)

			totalItemsCount++
			if i != grd {
				errorItemsCount++
				continue
			}

			require.Equalf(t, i, grd, "expect: (%f,%f,%d), actual: (%f,%f,%d)", lat, lng, i, lat32, lng32, grd)

			{
				_, _, v, err := reader.ReadLL(lat32, lng32) // grib_get -l 90,0 pkg/testdata/hpbl.grib2.out
				require.NoError(t, err)
				require.InDelta(t, float32(val), float32(v), 1e-5)
			}

			{
				// read again
				_, _, v, err := reader.ReadLL(lat32, lng32) // grib_get -l 90,0 pkg/testdata/hpbl.grib2.out
				require.NoError(t, err)
				require.InDelta(t, float32(val), float32(v), 1e-5)
			}
		}

		errorRate := float64(errorItemsCount) / float64(totalItemsCount)
		t.Logf("error rate: %f, total items: %d, error items: %d", errorRate, totalItemsCount, errorItemsCount)
		require.Less(t, errorRate, 0.01)
	})
}

func TestMessage_DumpMessageIndex(t *testing.T) {
	t.Parallel()

	// grib_set -r -s packingType=grid_simple pkg/testdata/hpbl.grib2 pkg/testdata/hpbl.grib2.out
	const filename = "../testdata/hpbl.grib2.out"

	s, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		t.Skipf("%s not exist", filename)
	}

	mm, err := mmap.Open(filename)
	require.NoError(t, err)
	defer mm.Close()

	g := grib.NewGrib2(mm)

	msg, err := g.ReadMessageAt(0)
	require.NoError(t, err)
	require.NotNil(t, msg)

	tests := []struct {
		name    string
		message grib.Message
		want    *grib.MessageIndex
		wantErr bool
	}{
		{
			name:    s.Name(),
			message: msg,
			want: &grib.MessageIndex{
				Offset:     0,
				Size:       2206439,
				DataOffset: 175,
				ScanningMode: gdt.ScanningMode(&gdt.ScanningMode0000{
					Ni:                          1440,
					Nj:                          721,
					LatitudeOfFirstGridPoint:    90000000,
					LongitudeOfFirstGridPoint:   0,
					ResolutionAndComponentFlags: 48,
					LatitudeOfLastGridPoint:     -90000000,
					LongitudeOfLastGridPoint:    359750000,
					IDirectionIncrement:         250000,
					JDirectionIncrement:         250000,
				}),
				Packing: &gridpoint.SimplePacking{
					ReferenceValue:     7.728597,
					BinaryScaleFactor:  -4,
					DecimalScaleFactor: 0,
					Bits:               17,
					Type:               0,
					NumVals:            1038240,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.message.DumpMessageIndex()
			require.NoError(t, err)
			require.EqualExportedValues(t, tt.want, got)
		})
	}
}
