package gribimg_test

import (
	"testing"

	"github.com/scorix/grib-go/pkg/gribimg"
	"github.com/stretchr/testify/require"
)

func TestExportPNG(t *testing.T) {
	t.Parallel()

	outDir := "../testdata"

	images, err := gribimg.ExportPNG("../testdata/temp.grib2", outDir, 273.15, 323.15)
	require.NoError(t, err)
	require.Len(t, images, 1)
}

func TestExportTIFF(t *testing.T) {
	t.Parallel()

	outDir := "../testdata"

	images, err := gribimg.ExportTIFF("../testdata/temp.grib2", outDir, 273.15, 323.15)
	require.NoError(t, err)
	require.Len(t, images, 1)
}
