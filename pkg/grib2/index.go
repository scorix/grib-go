package grib2

import (
	"encoding/json"
	"fmt"

	"github.com/scorix/grib-go/pkg/grib2/drt"
	"github.com/scorix/grib-go/pkg/grib2/gdt"
)

type MessageIndex struct {
	Offset       int64            `json:"offset"`
	Size         int64            `json:"size"`
	DataOffset   int64            `json:"data_offset"`
	ScanningMode gdt.ScanningMode `json:"scanning_mode"`
	Packing      drt.Template     `json:"packing"`
}

func (mi MessageIndex) MarshalJSON() ([]byte, error) {
	tm := drt.TemplateMarshaler{
		Template: mi.Packing,
	}

	sm := gdt.ScanningModeMarshaler{
		Template: mi.ScanningMode,
	}

	return json.Marshal(struct {
		Offset       int64                     `json:"offset"`
		Size         int64                     `json:"size"`
		DataOffset   int64                     `json:"data_offset"`
		ScanningMode gdt.ScanningModeMarshaler `json:"scanning_mode"`
		Packing      drt.TemplateMarshaler     `json:"packing"`
	}{
		Offset:       mi.Offset,
		Size:         mi.Size,
		DataOffset:   mi.DataOffset,
		ScanningMode: sm,
		Packing:      tm,
	})
}

func (mi *MessageIndex) UnmarshalJSON(data []byte) error {
	var temp struct {
		Offset       int64                     `json:"offset"`
		Size         int64                     `json:"size"`
		DataOffset   int64                     `json:"data_offset"`
		ScanningMode gdt.ScanningModeMarshaler `json:"scanning_mode"`
		Packing      drt.TemplateMarshaler     `json:"packing"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal message index: %w, data: %s", err, data)
	}

	mi.Offset = temp.Offset
	mi.Size = temp.Size
	mi.DataOffset = temp.DataOffset
	mi.Packing = temp.Packing.Template
	mi.ScanningMode = temp.ScanningMode.Template

	return nil
}
