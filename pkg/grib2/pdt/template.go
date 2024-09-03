package pdt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type Template interface {
	GetParameterCategory() int
	GetParameterNumber() int
	GetForecastDuration() time.Duration
	GetLevel() int
}

type MissingTemplate struct{}

func (m MissingTemplate) GetParameterCategory() int          { return -1 }
func (m MissingTemplate) GetParameterNumber() int            { return -1 }
func (m MissingTemplate) GetForecastDuration() time.Duration { return 0 }
func (m MissingTemplate) GetLevel() int                      { return 0 }

func ReadTemplate(r io.Reader, n uint16) (Template, error) {
	switch n {
	case 0:
		t0, err := readTemplate0(r)
		if err != nil {
			return nil, err
		}

		return t0.Export(), nil

	case 8:
		t0, err := readTemplate0(r)
		if err != nil {
			return nil, err
		}

		t8, err := readTemplate8(r, t0)
		if err != nil {
			return nil, err
		}

		t8.template0 = t0

		return t8.Export(), nil

	case 255:
		return &MissingTemplate{}, nil

	default:
		return nil, fmt.Errorf("unsupported product definition template: %d", n)
	}
}

func readTemplate0(r io.Reader) (*template0, error) {
	var tpl template0
	if err := binary.Read(r, binary.BigEndian, &tpl); err != nil {
		return nil, err
	}

	return &tpl, nil
}

func readTemplate8(r io.Reader, t0 *template0) (*template8, error) {
	var tpl template8
	if err := binary.Read(r, binary.BigEndian, &tpl.template8fields); err != nil {
		return nil, err
	}

	tpl.template0 = t0

	bs := tpl.template8fields.GetAdditionalTimeRangeSpecifications()
	if _, err := io.CopyN(bytes.NewBuffer(bs), r, int64(len(bs))); err != nil {
		return nil, err
	}

	return &tpl, nil
}
