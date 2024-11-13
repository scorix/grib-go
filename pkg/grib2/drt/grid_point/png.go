package gridpoint

import (
	"fmt"
	"image"
	"image/png"
	"io"

	"github.com/scorix/grib-go/internal/pkg/bitio"
	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
	"github.com/scorix/grib-go/pkg/grib2/drt/definition"
	"github.com/scorix/grib-go/pkg/grib2/regulation"
)

type PortableNetworkGraphics struct {
	ReferenceValue     float32
	BinaryScaleFactor  int16
	DecimalScaleFactor int16
	Bits               uint8
	Type               int8
	NumVals            int
}

func (p *PortableNetworkGraphics) Definition() any {
	return definition.PNG{
		R: p.ReferenceValue,
		B: regulation.ToUint16(p.BinaryScaleFactor),
		D: regulation.ToUint16(p.DecimalScaleFactor),
		L: p.Bits,
		T: regulation.ToUint8(p.Type),
	}
}

func NewPortableNetworkGraphics(def definition.PNG, numVals int) *PortableNetworkGraphics {
	return &PortableNetworkGraphics{
		ReferenceValue:     def.R,
		BinaryScaleFactor:  regulation.ToInt16(def.B),
		DecimalScaleFactor: regulation.ToInt16(def.D),
		Bits:               def.L,
		Type:               regulation.ToInt8(def.T),
		NumVals:            numVals,
	}
}

func (p *PortableNetworkGraphics) GetNumVals() int {
	return p.NumVals
}

func (p *PortableNetworkGraphics) ScaleFunc() func(uint32) float32 {
	return datapacking.SimpleScaleFunc(p.BinaryScaleFactor, p.DecimalScaleFactor, p.ReferenceValue)
}

func (p *PortableNetworkGraphics) ReadAllData(r *bitio.Reader) ([]float32, error) {
	var err error
	values := make([]float32, p.NumVals)
	scaleFunc := p.ScaleFunc()

	// Special case: if bits per value is 0, all values are equal to the reference value
	if p.Bits == 0 {
		for i := range values {
			values[i] = p.ReferenceValue
		}
		return values, nil
	}

	// Decode PNG
	img, err := png.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if i >= p.NumVals {
				break
			}
			r, g, b, a := img.At(x, y).RGBA()
			values[i] = scaleFunc(p.rgbaToUint32(r, g, b, a))
			i++
		}
	}

	if i != p.NumVals {
		return nil, fmt.Errorf("expected %d values, got %d", p.NumVals, i)
	}

	return values, nil
}

func (p *PortableNetworkGraphics) rgbaToUint32(r, g, b, a uint32) uint32 {
	switch p.Bits {
	case 1, 2, 4, 8, 16:
		// 1, 2, 4, 8, or 16 - treat as a grayscale image
		return r
	case 24:
		// 24 - treat as RGB color image (each component having 8 bit depth)
		return r<<16 | g<<8 | b
	case 32:
		// 32 - treat as RGB w/alpha sample color image (each component having 8 bit depth)
		return r<<24 | g<<16 | b<<8 | a
	}

	return r
}

func (p *PortableNetworkGraphics) Image(r datapacking.BitReader) (image.Image, error) {
	img, err := png.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}

	return img, nil
}

type PortableNetworkGraphicsReader struct {
	r      io.ReaderAt
	p      *PortableNetworkGraphics
	sf     func(uint32) float32
	offset int64
	length int64
}

func NewPortableNetworkGraphicsReader(r io.ReaderAt, start, end int64, p *PortableNetworkGraphics) *PortableNetworkGraphicsReader {
	return &PortableNetworkGraphicsReader{
		r:      r,
		p:      p,
		sf:     p.ScaleFunc(),
		offset: start,
		length: end - start,
	}
}
