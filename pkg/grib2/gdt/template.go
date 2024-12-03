package gdt

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

type Template interface {
	GetNi() int32
	GetNj() int32
	GetGridIndex(lat, lon float32) (n int)
	GetGridPoint(n int) (float32, float32, bool)
}

type MissingTemplate struct{}

func (m MissingTemplate) GetGridPointFromLL(float32, float32) int { return 0 }
func (m MissingTemplate) GetNi() int32                            { return 0 }
func (m MissingTemplate) GetNj() int32                            { return 0 }
func (m MissingTemplate) GetGridIndex(lat, lon float32) (n int) {
	return 0
}
func (m MissingTemplate) GetGridPoint(n int) (float32, float32, bool) {
	return 0, 0, false
}

func ReadTemplate(r io.Reader, n uint16) (Template, error) {
	switch n {
	case 0:
		var tpl template0FixedPart
		if err := binary.Read(r, binary.BigEndian, &tpl); err != nil {
			return nil, err
		}

		return tpl.Export(), nil

	case 40:
		var tpl template40FixedPart
		if err := binary.Read(r, binary.BigEndian, &tpl); err != nil {
			return nil, err
		}

		return tpl.Export(), nil

	case 255:
		return &MissingTemplate{}, nil

	default:
		return nil, fmt.Errorf("unsupported grid definition template: %d", n)
	}
}

func UnMarshalJSONTemplate(data []byte) (Template, error) {
	var tpl struct {
		Template0  *Template0FixedPart  `json:"template0"`
		Template40 *Template40FixedPart `json:"template40"`
	}

	if err := json.Unmarshal(data, &tpl); err != nil {
		return nil, err
	}

	switch {
	case tpl.Template0 != nil:
		return tpl.Template0.AsTemplate(), nil
	case tpl.Template40 != nil:
		return tpl.Template40.AsTemplate(), nil
	}

	return nil, fmt.Errorf("unsupported grid definition template")
}
