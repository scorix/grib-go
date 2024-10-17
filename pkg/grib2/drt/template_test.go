package drt_test

import (
	"testing"

	"github.com/scorix/grib-go/pkg/grib2/drt"
	gridpoint "github.com/scorix/grib-go/pkg/grib2/drt/grid_point"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateMarshaler_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  drt.TemplateMarshaler
		want    string
		wantErr bool
	}{
		{
			name: "simple packing",
			fields: drt.TemplateMarshaler{
				Template: &gridpoint.SimplePacking{
					ReferenceValue:     1.5,
					BinaryScaleFactor:  2,
					DecimalScaleFactor: 3,
					Bits:               16,
					NumVals:            721 * 1440,
				},
			},
			want: `{"number":0,"content":{"r":1.5,"b":2,"d":3,"l":16,"t":0},"vals":1038240}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.fields.MarshalJSON()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.JSONEq(t, tt.want, string(got), string(got))
			}
		})
	}
}

func TestTemplateMarshaler_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		json    string
		want    drt.Template
		wantErr bool
	}{
		{
			name: "simple packing",
			json: `{"number":0,"content":{"r":1.5,"b":2,"d":3,"l":16,"t":0},"vals":1038240}`,
			want: &gridpoint.SimplePacking{
				ReferenceValue:     1.5,
				BinaryScaleFactor:  2,
				DecimalScaleFactor: 3,
				Bits:               16,
				NumVals:            1038240,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var tm drt.TemplateMarshaler
			err := tm.UnmarshalJSON([]byte(tt.json))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, tm.Template)
			}
		})
	}
}
