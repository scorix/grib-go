package cache_test

import (
	"testing"

	"github.com/scorix/grib-go/pkg/grib2/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockGridDataSource struct {
	gridValue float32
	readCount int
}

func (m *mockGridDataSource) ReadGridAt(grid int) (float32, error) {
	m.readCount++
	return m.gridValue, nil
}

func TestNoCache(t *testing.T) {
	ds := &mockGridDataSource{gridValue: 100}
	nc := cache.NewNoCache(ds)

	// first read should be from source
	v, err := nc.ReadGridAt(1, 1, 1)
	require.NoError(t, err)
	assert.Equal(t, float32(100), v)
	assert.Equal(t, 1, ds.readCount)

	// second read should not be cached
	v, err = nc.ReadGridAt(1, 1, 1)
	require.NoError(t, err)
	assert.Equal(t, float32(100), v)
	assert.Equal(t, 2, ds.readCount)
}
