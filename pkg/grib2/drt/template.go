package drt

import (
	"io"

	"github.com/scorix/grib-go/pkg/grib2/drt/datapacking"
)

type Template interface {
	NewUnpackReader(io.Reader) (datapacking.UnpackReader, error)
}
