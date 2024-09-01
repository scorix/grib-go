package drt

import (
	"errors"
	"io"

	"github.com/icza/bitio"
	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
)

type Template interface {
	datapacking.UnpackReader
}

func ScaleData(pr Template, r io.Reader) ([]float64, error) {
	var (
		br        = bitio.NewReader(r)
		values    []float64
		scaleFunc = pr.ScaleFunc()
	)

	for {
		bitsVal, err := br.ReadBits(pr.GetBits())
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, err
		}

		values = append(values, scaleFunc(bitsVal))
	}

	return values, nil
}
