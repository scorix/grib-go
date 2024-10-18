package grib2_test

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	codes "github.com/scorix/go-eccodes"
	cio "github.com/scorix/go-eccodes/io"
	"github.com/scorix/grib-go/pkg/colormap"
	"github.com/scorix/grib-go/pkg/grib2"
	grib "github.com/scorix/grib-go/pkg/grib2"
	gridpoint "github.com/scorix/grib-go/pkg/grib2/drt/grid_point"
	"github.com/scorix/grib-go/pkg/grib2/gdt"
	"github.com/scorix/grib-go/pkg/grib2/regulation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/mmap"
	"golang.org/x/sync/errgroup"
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

			grd, err := msg.GetGridPointFromLL(lat32, lng32)
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

func TestMessage_Image(t *testing.T) {
	t.Parallel()

	const (
		f1 = "/Users/scorix/Downloads/2t_heightAboveGround_2_0_0.grib2"
		f2 = "/Users/scorix/Downloads/2t_heightAboveGround_2_0_0.png.grib2"
	)

	mm1, err := mmap.Open(f1)
	if err != nil {
		t.Skipf("skip %s: %v", f1, err)
	}

	mm2, err := mmap.Open(f2)
	if err != nil {
		t.Skipf("skip %s: %v", f2, err)
	}

	g1 := grib.NewGrib2(mm1)
	g2 := grib.NewGrib2(mm2)

	t.Run("compare data", func(t *testing.T) {
		t.Parallel()

		msg1, err := g1.ReadMessageAt(0)
		require.NoError(t, err)

		msg2, err := g2.ReadMessageAt(0)
		require.NoError(t, err)

		data1, err := msg1.ReadData()
		require.NoError(t, err)

		data2, err := msg2.ReadData()
		require.NoError(t, err)

		assert.Equal(t, data1, data2)
	})

	t.Run("export image from grid_png", func(t *testing.T) {
		t.Parallel()

		var i int
		var eg errgroup.Group

		_ = g2.EachMessage(func(msg grib.IndexedMessage) (next bool, err error) {
			require.IsType(t, &gridpoint.PortableNetworkGraphics{}, msg.GetDataRepresentationTemplate())
			filename := fmt.Sprintf("../../tmp/2t_heightAboveGround_2_0_0.%d.png", i)

			eg.Go(func() error {
				img, err := msg.Image()
				if err != nil {
					return err
				}

				bounds := img.Bounds()

				var max, min uint32
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					for x := bounds.Min.X; x < bounds.Max.X; x++ {
						r, _, _, _ := img.At(x, y).RGBA()
						if r > max {
							max = r
						}
						if r < min {
							min = r
						}
					}
				}

				rgbaImg := image.NewRGBA(bounds)
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					for x := bounds.Min.X; x < bounds.Max.X; x++ {
						r, _, _, _ := img.At(x, y).RGBA()
						r = (r - min) * 255 / (max - min)
						rgbaImg.Set(x, y, color.RGBA{R: uint8(r), G: uint8(r), B: uint8(r), A: 255})
					}
				}

				tmp, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o644)
				if err != nil {
					return err
				}
				defer tmp.Close()

				return png.Encode(tmp, rgbaImg)
			})

			i++

			return true, nil
		})
		require.NoError(t, eg.Wait())
	})

	t.Run("export image from grid_simple", func(t *testing.T) {
		t.Parallel()

		var i int
		const celsiusZero = 273.15
		var max, min float64 = celsiusZero + 50, celsiusZero - 50
		var eg errgroup.Group

		_ = g1.EachMessage(func(msg grib.IndexedMessage) (next bool, err error) {
			require.Equal(t, 1440, msg.GetNi())
			require.Equal(t, 721, msg.GetNj())

			require.IsType(t, &gridpoint.SimplePacking{}, msg.GetDataRepresentationTemplate())
			filename := fmt.Sprintf("../../tmp/2t_heightAboveGround_2_0_0.%d.png", i)
			r := gridpoint.NewSimplePackingReader(mm1, msg.GetDataOffset(), msg.GetSize()-msg.GetDataOffset(), msg.GetDataRepresentationTemplate().(*gridpoint.SimplePacking))
			bounds := image.Rect(0, 0, msg.GetNi(), msg.GetNj())

			tmp, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o644)
			if err != nil {
				return false, err
			}

			eg.Go(func() error {
				defer tmp.Close()
				rgbaImg := image.NewGray16(bounds)
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					for x := bounds.Min.X; x < bounds.Max.X; x++ {
						grd := y*msg.GetNi() + x
						val, err := r.ReadGridAt(grd)
						if err != nil {
							return err
						}

						gray := colormap.GrayColorMap.GetColor(val, min, max)
						rgbaImg.Set(x, y, gray)
					}
				}

				return png.Encode(tmp, rgbaImg)
			})

			i++

			return true, nil
		})
		require.NoError(t, eg.Wait())
	})
}
