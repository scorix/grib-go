package gridpoint

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/icza/bitio"
	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
	"github.com/scorix/grib-go/pkg/grib2/drt/definition"
	"github.com/scorix/grib-go/pkg/grib2/regulation"
)

type SimplePacking struct {
	R       float32 // 12-15
	E       int16   // 16-17
	D       int16   // 18-19
	Bits    uint8   // 20
	Type    int8    // 21
	numVals int     // 22
}

func NewSimplePacking(def definition.SimplePacking, numVals int) *SimplePacking {
	return &SimplePacking{
		R:       def.R,
		E:       regulation.ToInt16(def.E),
		D:       regulation.ToInt16(def.D),
		Bits:    def.Bits,
		Type:    regulation.ToInt8(def.Type),
		numVals: numVals,
	}
}

func (sp *SimplePacking) ScaleFunc() func(uint32) float64 {
	return datapacking.SimpleScaleFunc(sp.E, sp.D, sp.R)
}

func (sp *SimplePacking) ReadAllData(r datapacking.BitReader) ([]float64, error) {
	var (
		values    []float64
		scaleFunc = sp.ScaleFunc()
	)

	if sp.Bits == 0 {
		for range sp.numVals {
			values = append(values, scaleFunc(0))
		}
	}

	for sp.Bits > 0 {
		bitsVal, err := r.ReadBits(sp.Bits)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, err
		}

		values = append(values, scaleFunc(uint32(bitsVal)))
	}

	if len(values) != sp.numVals {
		return nil, fmt.Errorf("expected %d values, got %d", sp.numVals, len(values))
	}

	return values, nil
}

type SimplePackingReader struct {
	r      io.ReaderAt
	sp     *SimplePacking
	sf     func(uint32) float64
	offset int64
	length int64
}

func NewSimplePackingReader(r io.ReaderAt, start, end int64, sp *SimplePacking) *SimplePackingReader {
	return &SimplePackingReader{
		r:      r,
		sp:     sp,
		sf:     sp.ScaleFunc(),
		offset: start,
		length: end - start,
	}
}

func (r *SimplePackingReader) ReadGridAt(n int) (float64, error) {
	if n >= r.sp.numVals {
		return 0, fmt.Errorf("requesting[%d] is out of range, total number of values is %d", n, r.sp.numVals)
	}

	bitsOffset := n * int(r.sp.Bits)
	skipBits := bitsOffset % 8
	needBytes := int(math.Ceil(float64(int(r.sp.Bits)+skipBits) / float64(8.0)))

	bs := make([]byte, needBytes)
	if _, err := r.r.ReadAt(bs, r.offset+int64(bitsOffset/8)); err != nil {
		return 0, fmt.Errorf("read at: %w, range: %d - %d", err, r.offset, r.offset+r.length)
	}

	br := bitio.NewReader(bytes.NewReader(bs))

	if skipBits > 0 {
		if _, err := br.ReadBits(uint8(skipBits)); err != nil {
			return 0, fmt.Errorf("skip %d bits: %w", skipBits, err)
		}
	}

	u, err := br.ReadBits(r.sp.Bits)
	if err != nil {
		return 0, fmt.Errorf("read %d bits: %w", r.sp.Bits, err)
	}

	return r.sf(uint32(u)), nil
}
