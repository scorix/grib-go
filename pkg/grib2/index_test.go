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
				GridDefinition: &gdt.Template0{
					Template0FixedPart: gdt.Template0FixedPart{
						LatitudeOfFirstGridPoint:  90000000,
						LongitudeOfFirstGridPoint: 0,
						LatitudeOfLastGridPoint:   -90000000,
						LongitudeOfLastGridPoint:  359750000,
						IDirectionIncrement:       250000,
						JDirectionIncrement:       250000,
						ScanningMode:              1,
					},
				},
				Packing: drt.Template(&gridpoint.SimplePacking{
					ReferenceValue:     1.5,
					BinaryScaleFactor:  2,
					DecimalScaleFactor: 3,
					Bits:               16,
				}),
			},
			want: `{"offset":100,"size":1000,"data_offset":200,"grid_definition":{"template0":{"latitudeOfFirstGridPoint":90000000,"longitudeOfFirstGridPoint":0,"latitudeOfLastGridPoint":-90000000,"longitudeOfLastGridPoint":359750000,"iDirectionIncrement":250000,"jDirectionIncrement":250000,"scanningMode":1}},"packing":{"number":0,"content":{"r":1.5,"b":2,"d":3,"l":16,"t":0},"vals":0}}`,
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
			json: `{"offset":100,"size":1000,"data_offset":200,"grid_definition":{"template0":{"latitudeOfFirstGridPoint":90000000,"longitudeOfFirstGridPoint":0,"latitudeOfLastGridPoint":-90000000,"longitudeOfLastGridPoint":359750000,"iDirectionIncrement":250000,"jDirectionIncrement":250000,"scanningMode":1}},"packing":{"number":0,"content":{"r":1.5,"b":2,"d":3,"l":16,"t":0},"vals":0}}`,
			want: grib2.MessageIndex{
				Offset:     100,
				Size:       1000,
				DataOffset: 200,
				GridDefinition: &gdt.Template0{
					Template0FixedPart: gdt.Template0FixedPart{
						LatitudeOfFirstGridPoint:  90000000,
						LongitudeOfFirstGridPoint: 0,
						LatitudeOfLastGridPoint:   -90000000,
						LongitudeOfLastGridPoint:  359750000,
						IDirectionIncrement:       250000,
						JDirectionIncrement:       250000,
						ScanningMode:              1,
					},
				},
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
				assert.EqualExportedValues(t, tt.want, got)
			}
		})
	}
}
