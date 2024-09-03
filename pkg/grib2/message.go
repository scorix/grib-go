package grib2

import (
	"bytes"
	"time"

	"github.com/icza/bitio"
)

type Message interface {
	GetDiscipline() int
	GetParameterCategory() int
	GetParameterNumber() int
	GetTimestamp(loc *time.Location) time.Time
	GetForecastTime(loc *time.Location) time.Time
	GetLevel() int
	ReadData() ([]float64, error)
}

type message struct {
	sec0 *section0
	sec1 *section1
	sec2 *section2
	sec3 *section3
	sec4 *section4
	sec5 *section5
	sec6 *section6
	sec7 *section7
	sec8 *section8
}

func (m message) GetDiscipline() int {
	return m.sec0.GetDiscipline()
}

func (m message) GetParameterCategory() int {
	return m.sec4.GetParameterCategory()
}

func (m message) GetParameterNumber() int {
	return m.sec4.GetParameterNumber()
}

func (m message) GetTimestamp(loc *time.Location) time.Time {
	return m.sec1.GetTime(loc)
}

func (m message) GetForecastTime(loc *time.Location) time.Time {
	return m.GetTimestamp(loc).Add(m.sec4.GetForecastDuration())
}

func (m *message) GetLevel() int {
	return m.sec4.GetLevel()
}

func (m *message) ReadData() ([]float64, error) {
	tpl := m.sec5.GetDataRepresentationTemplate()

	return tpl.ReadAllData(bitio.NewReader(bytes.NewReader(m.sec7.Data)))
}
