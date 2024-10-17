package gdt

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

type Template interface {
	GetScanningMode() (ScanningMode, error)
	GetNi() int32
	GetNj() int32
}

type MissingTemplate struct{}

func (m MissingTemplate) GetScanningMode() (ScanningMode, error) {
	return nil, fmt.Errorf("unknown scanning mode")
}
func (m MissingTemplate) GetGridPointFromLL(float32, float32) int { return 0 }
func (m MissingTemplate) GetNi() int32                            { return 0 }
func (m MissingTemplate) GetNj() int32                            { return 0 }

func ReadTemplate(r io.Reader, n uint16) (Template, error) {
	switch n {
	case 0:
		var tpl template0FixedPart
		if err := binary.Read(r, binary.BigEndian, &tpl); err != nil {
			return nil, err
		}

		return &Template0{Template0FixedPart: tpl.Export()}, nil

	case 255:
		return &MissingTemplate{}, nil

	default:
		return nil, fmt.Errorf("unsupported grid definition template: %d", n)
	}
}

type ScanningModeMarshaler struct {
	Template ScanningMode
}

type scanningModeMarshaler struct {
	Mode    int8            `json:"mode"`
	Content json.RawMessage `json:"content"`
}

func (sm ScanningModeMarshaler) MarshalJSON() ([]byte, error) {
	content, err := json.Marshal(sm.Template)
	if err != nil {
		return nil, err
	}

	m := scanningModeMarshaler{
		Mode:    sm.Template.GetScanMode(),
		Content: content,
	}

	return json.Marshal(m)
}

func (sm *ScanningModeMarshaler) UnmarshalJSON(data []byte) error {
	var m scanningModeMarshaler

	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	switch m.Mode {
	case 0:
		var mode ScanningMode0000
		if err := json.Unmarshal(m.Content, &mode); err != nil {
			return err
		}

		sm.Template = &mode
		return nil
	}

	return fmt.Errorf("unsupported scanning mode: %d", m.Mode)
}
