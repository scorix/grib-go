package gdt_test

import (
	"encoding/json"
	"testing"

	"github.com/scorix/grib-go/pkg/grib2/gdt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   gdt.Template
		want    string
		wantErr bool
	}{
		{
			name: "marshal template 0",
			input: &gdt.Template0{
				Template0FixedPart: gdt.Template0FixedPart{
					Ni:                          360,
					Nj:                          181,
					LatitudeOfFirstGridPoint:    90000000,
					LongitudeOfFirstGridPoint:   0,
					ResolutionAndComponentFlags: 48,
					LatitudeOfLastGridPoint:     -90000000,
					LongitudeOfLastGridPoint:    359000000,
					IDirectionIncrement:         1000000,
					JDirectionIncrement:         1000000,
				},
			},
			want:    `{"template0":{"latitudeOfFirstGridPoint":90000000,"longitudeOfFirstGridPoint":0,"latitudeOfLastGridPoint":-90000000,"longitudeOfLastGridPoint":359000000,"iDirectionIncrement":1000000,"jDirectionIncrement":1000000,"scanningMode":0}}`,
			wantErr: false,
		},
		{
			name: "marshal template 40",
			input: &gdt.Template40{
				Template40FixedPart: gdt.Template40FixedPart{
					N:            768,
					ScanningMode: 0,
				},
			},
			want:    `{"template40":{"n":768,"scanningMode":0}}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.JSONEq(t, tt.want, string(got), string(got))
			}
		})
	}
}

func TestTemplate0_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    gdt.Template
		wantErr bool
	}{
		{
			name:  "unmarshal template 0",
			input: `{"template0":{"latitudeOfFirstGridPoint":90000000,"longitudeOfFirstGridPoint":0,"latitudeOfLastGridPoint":-90000000,"longitudeOfLastGridPoint":359000000,"iDirectionIncrement":1000000,"jDirectionIncrement":1000000,"scanningMode":0}}`,
			want: &gdt.Template0{
				Template0FixedPart: gdt.Template0FixedPart{
					LatitudeOfFirstGridPoint:  90000000,
					LongitudeOfFirstGridPoint: 0,
					LatitudeOfLastGridPoint:   -90000000,
					LongitudeOfLastGridPoint:  359000000,
					IDirectionIncrement:       1000000,
					JDirectionIncrement:       1000000,
				},
			},
			wantErr: false,
		},
		{
			name:  "unmarshal template 40",
			input: `{"template40":{"n":768,"scanningMode":0}}`,
			want: &gdt.Template40{
				Template40FixedPart: gdt.Template40FixedPart{
					N:            768,
					ScanningMode: 0,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gdt.UnMarshalJSONTemplate([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.EqualExportedValues(t, tt.want, got)
			}
		})
	}
}
