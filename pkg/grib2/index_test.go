package grib2_test

import (
	"encoding/json"
	"testing"

	"github.com/scorix/grib-go/pkg/grib2"
	"github.com/scorix/grib-go/pkg/grib2/drt"
	gridpoint "github.com/scorix/grib-go/pkg/grib2/drt/grid_point"
	"github.com/scorix/grib-go/pkg/grib2/gdt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageIndex_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  grib2.MessageIndex
		want    string
		wantErr bool
	}{
		{
			name: "marshal simple packing",
			fields: grib2.MessageIndex{
				Offset:     100,
				Size:       1000,
				DataOffset: 200,
				ScanningMode: gdt.ScanningMode(&gdt.ScanningMode0000{
					Ni: 10,
					Nj: 20,
				}),
				Packing: drt.Template(&gridpoint.SimplePacking{
					ReferenceValue:     1.5,
					BinaryScaleFactor:  2,
					DecimalScaleFactor: 3,
					Bits:               16,
				}),
			},
			want: `{"offset":100,"size":1000,"data_offset":200,"scanning_mode":{"mode":0,"content":{"ni":10,"nj":20,"latitudeOfFirstGridPoint":0,"longitudeOfFirstGridPoint":0,"resolutionAndComponentFlags":0,"latitudeOfLastGridPoint":0,"longitudeOfLastGridPoint":0,"iDirectionIncrement":0,"jDirectionIncrement":0}},"packing":{"number":0,"content":{"r":1.5,"b":2,"d":3,"l":16,"t":0},"vals":0}}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := json.Marshal(tt.fields)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.JSONEq(t, tt.want, string(got), string(got))
			}
		})
	}
}

func TestMessageIndex_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		json    string
		want    grib2.MessageIndex
		wantErr bool
	}{
		{
			name: "unmarshal simple packing",
			json: `{"offset":100,"size":1000,"data_offset":200,"scanning_mode":{"mode":0,"content":{"ni":10,"nj":20,"latitudeOfFirstGridPoint":0,"longitudeOfFirstGridPoint":0,"resolutionAndComponentFlags":0,"latitudeOfLastGridPoint":0,"longitudeOfLastGridPoint":0,"iDirectionIncrement":0,"jDirectionIncrement":0}},"packing":{"number":0,"content":{"r":1.5,"b":2,"d":3,"l":16,"t":0},"vals":0}}`,
			want: grib2.MessageIndex{
				Offset:     100,
				Size:       1000,
				DataOffset: 200,
				ScanningMode: gdt.ScanningMode(&gdt.ScanningMode0000{
					Ni: 10,
					Nj: 20,
				}),
				Packing: drt.Template(&gridpoint.SimplePacking{
					ReferenceValue:     1.5,
					BinaryScaleFactor:  2,
					DecimalScaleFactor: 3,
					Bits:               16,
				}),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got grib2.MessageIndex
			err := json.Unmarshal([]byte(tt.json), &got)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
