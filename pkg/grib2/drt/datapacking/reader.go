package datapacking

import "io"

type UnpackReader interface {
	ReadData(io.Reader) ([]float64, error)
}
