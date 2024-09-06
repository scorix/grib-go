package gridpoint

import (
	"errors"
	"fmt"
	"io"

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
