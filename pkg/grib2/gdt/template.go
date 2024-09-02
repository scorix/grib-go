package gdt

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Template interface{}

type MissingTemplate struct{}

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
