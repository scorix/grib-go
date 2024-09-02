package pdt

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
		var tpl template0
		if err := binary.Read(r, binary.BigEndian, &tpl); err != nil {
			return nil, err
		}

		return tpl.Export(), nil

	case 255:
		return &MissingTemplate{}, nil

	default:
		return nil, fmt.Errorf("unsupported product definition template: %d", n)
	}
}
