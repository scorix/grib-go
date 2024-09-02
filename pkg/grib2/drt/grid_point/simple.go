package gridpoint

import (
	"errors"
	"fmt"
	"io"

	"github.com/icza/bitio"
	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
	"github.com/scorix/grib-go/pkg/grib2/drt/definition"
	"github.com/scorix/grib-go/pkg/grib2/regulation"
)

type SimplePacking struct {
	R       float32
	E       int16
	D       int16
	Bits    uint8
	Type    int8
	numVals int
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

func (sp *SimplePacking) ReadAllData(r io.Reader) ([]float64, error) {
	var (
		br        = bitio.NewReader(r)
		values    []float64
		scaleFunc = sp.ScaleFunc()
	)

	for {
		bitsVal, err := br.ReadBits(sp.Bits)
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
