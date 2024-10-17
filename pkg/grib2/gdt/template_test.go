package gdt_test

import (
	"encoding/json"
	"testing"

	"github.com/scorix/grib-go/pkg/grib2/gdt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanningModeJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   gdt.ScanningModeMarshaler
		want    string
		wantErr bool
	}{
		{
			name: "marshal scanning mode 0000",
			input: gdt.ScanningModeMarshaler{
				Template: &gdt.ScanningMode0000{
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
			want:    `{"mode":0,"content":{"ni":360,"nj":181,"latitudeOfFirstGridPoint":90000000,"longitudeOfFirstGridPoint":0,"resolutionAndComponentFlags":48,"latitudeOfLastGridPoint":-90000000,"longitudeOfLastGridPoint":359000000,"iDirectionIncrement":1000000,"jDirectionIncrement":1000000}}`,
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

func TestScanningMode0000_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    gdt.ScanningModeMarshaler
		wantErr bool
	}{
		{
			name:  "unmarshal scanning mode 0000",
			input: `{"mode":0,"content":{"ni":360,"nj":181,"latitudeOfFirstGridPoint":90000000,"longitudeOfFirstGridPoint":0,"resolutionAndComponentFlags":48,"latitudeOfLastGridPoint":-90000000,"longitudeOfLastGridPoint":359000000,"iDirectionIncrement":1000000,"jDirectionIncrement":1000000}}`,
			want: gdt.ScanningModeMarshaler{
				Template: &gdt.ScanningMode0000{
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
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got gdt.ScanningModeMarshaler
			err := json.Unmarshal([]byte(tt.input), &got)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
