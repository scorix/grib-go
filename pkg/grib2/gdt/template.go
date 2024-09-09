package gdt

import (
	"encoding/binary"
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
