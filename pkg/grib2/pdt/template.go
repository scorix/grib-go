package pdt

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type Template interface {
	GetParameterCategory() int
	GetParameterNumber() int
	GetForecastDuration() time.Duration
}

type MissingTemplate struct{}

func (m MissingTemplate) GetParameterCategory() int          { return -1 }
func (m MissingTemplate) GetParameterNumber() int            { return -1 }
func (m MissingTemplate) GetForecastDuration() time.Duration { return 0 }

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
