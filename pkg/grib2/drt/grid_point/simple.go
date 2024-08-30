package gridpoint

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/icza/bitio"
	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
)

type SimplePacking struct {
	*DefSimplePacking
	*SimplePackingReader
}

func NewSimplePacking(def DefSimplePacking) *SimplePacking {
	return &SimplePacking{
		DefSimplePacking: &def,
		SimplePackingReader: &SimplePackingReader{
			DefSimplePacking: &def,
		},
	}
}

type DefSimplePacking struct {
	R    float32     // Reference value (R) (IEEE 32-bit floating-point value)
	E    scaleFactor // Binary scale factor
	D    scaleFactor // Decimal scale factor
	Bits uint8       // Number of bits used for each packed value for simple packing, or for each group reference value for complex packing or spatial differencing
	Type uint8       // Type of original field values: https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table5-1.shtml
}

func (def DefSimplePacking) BinaryScaleFactor() float64 {
	return math.Pow(2, float64(def.E.Int16()))
}

func (def DefSimplePacking) DecimalScaleFactor() float64 {
	return math.Pow10(-int(def.D.Int16()))
}

func (def DefSimplePacking) ReferenceValue() float64 {
	return def.DecimalScaleFactor() * float64(def.R)
}

func (def DefSimplePacking) ScaleFactor() float64 {
	return def.BinaryScaleFactor() * def.DecimalScaleFactor()
}

func (def DefSimplePacking) ScaleFunc() func(uint64) float64 {
	ref, scale := def.ReferenceValue(), def.ScaleFactor()

	return func(v uint64) float64 {
		return SimpleScale(v, ref, scale)
	}
}

func SimpleScale(v uint64, ref float64, scale float64) float64 {
	return ref + float64(v)*scale
}

func (sp *SimplePacking) NewUnpackReader(r io.Reader) (datapacking.UnpackReader, error) {
	pr := SimplePackingReader{
		DefSimplePacking: sp.DefSimplePacking,
	}

	return &pr, nil
}

type SimplePackingReader struct {
	*DefSimplePacking
}

func (pr *SimplePackingReader) ReadData(r io.Reader) ([]float64, error) {
	br := bitio.NewReader(r)

	var (
		values    []float64
		scaleFunc = pr.ScaleFunc()
	)

	for {
		bitsVal, err := br.ReadBits(pr.Bits)
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

type scaleFactor uint16

func readScaleFactor(r io.Reader) (scaleFactor, error) {
	var f scaleFactor

	if err := binary.Read(r, binary.BigEndian, &f); err != nil {
		return f, err
	}

	return f, nil
}

// https://codes.ecmwf.int/grib/format/grib2/regulations/
// 92.1.5 If applicable, negative values shall be indicated by setting the most significant bit to “1”.
func (s scaleFactor) Int16() int16 {
	negtive := s&0x8000 > 0
	i := int16(s & 0x7fff)

	if negtive {
		return -i
	}

	return i
}
