package gribimg

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/scorix/grib-go/pkg/colormap"
	"github.com/scorix/grib-go/pkg/grib2"
	gridpoint "github.com/scorix/grib-go/pkg/grib2/drt/grid_point"
	"golang.org/x/exp/mmap"
	"golang.org/x/image/tiff"
	"golang.org/x/sync/errgroup"
)

func ExportPNG(gribFilename string, outDir string, min, max float64) ([]image.Image, error) {
	return ExportImage(gribFilename, outDir, min, max, func(basename string, img image.Image) error {
		f, err := os.Create(filepath.Join(outDir, fmt.Sprintf("%s.png", basename)))
		if err != nil {
			return err
		}

		return png.Encode(f, img)
	})
}

func ExportTIFF(gribFilename string, outDir string, min, max float64) ([]image.Image, error) {
	return ExportImage(gribFilename, outDir, min, max, func(basename string, img image.Image) error {
		f, err := os.Create(filepath.Join(outDir, fmt.Sprintf("%s.tiff", basename)))
		if err != nil {
			return err
		}

		return tiff.Encode(f, img, &tiff.Options{Compression: tiff.Uncompressed})
	})
}

func ExportImage(gribFilename string, outDir string, min, max float64, encode func(string, image.Image) error) ([]image.Image, error) {
	mm, err := mmap.Open(gribFilename)
	if err != nil {
		return nil, err
	}

	g := grib2.NewGrib2(mm)
	var images []image.Image
	var eg errgroup.Group
	var i int

	err = g.EachMessage(func(msg grib2.IndexedMessage) (next bool, err error) {
		tpl, ok := msg.GetDataRepresentationTemplate().(*gridpoint.SimplePacking)
		if !ok {
			return false, fmt.Errorf("unsupported data representation template: %T", msg.GetDataRepresentationTemplate())
		}

		filename := filepath.Join(outDir, fmt.Sprintf("%s.%d", filepath.Base(gribFilename), i))
		r := gridpoint.NewSimplePackingReader(mm, msg.GetDataOffset(), msg.GetSize()-msg.GetDataOffset(), tpl)
		bounds := image.Rect(0, 0, msg.GetNi(), msg.GetNj())

		eg.Go(func() error {
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

			images = append(images, rgbaImg)

			return encode(filename, rgbaImg)
		})

		i++

		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return images, eg.Wait()
}
