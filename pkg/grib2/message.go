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
	Step() int

	GetTypeOfFirstFixedSurface() int
	GetScaleFactorOfFirstFixedSurface() int
	GetScaledValueOfFirstFixedSurface() int
	GetTypeOfSecondFixedSurface() int
	GetScaleFactorOfSecondFixedSurface() int
	GetScaledValueOfSecondFixedSurface() int
	GetProductDefinitionTemplateNumber() int
	GetGridPoint(n int) (float32, float32)
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
	return m.sec4.GetProductDefinitionTemplate().GetParameterCategory()
}

func (m message) GetParameterNumber() int {
	return m.sec4.GetProductDefinitionTemplate().GetParameterNumber()
}

func (m message) GetTimestamp(loc *time.Location) time.Time {
	return m.sec1.GetTime(loc)
}

func (m message) GetForecastTime(loc *time.Location) time.Time {
	return m.GetTimestamp(loc).Add(m.sec4.GetProductDefinitionTemplate().GetForecastDuration())
}

func (m *message) GetLevel() int {
	return m.sec4.GetProductDefinitionTemplate().GetLevel()
}

func (m *message) ReadData() ([]float64, error) {
	tpl := m.sec5.GetDataRepresentationTemplate()

	return tpl.ReadAllData(bitio.NewReader(bytes.NewReader(m.sec7.Data)))
}

func (m *message) Step() int {
	return m.sec4.GetProductDefinitionTemplate().GetForecast()
}

func (m *message) GetTypeOfFirstFixedSurface() int {
	return m.sec4.GetProductDefinitionTemplate().GetTypeOfFirstFixedSurface()
}

func (m *message) GetScaleFactorOfFirstFixedSurface() int {
	return m.sec4.GetProductDefinitionTemplate().GetScaleFactorOfFirstFixedSurface()
}

func (m *message) GetScaledValueOfFirstFixedSurface() int {
	return m.sec4.GetProductDefinitionTemplate().GetScaledValueOfFirstFixedSurface()
}

func (m *message) GetTypeOfSecondFixedSurface() int {
	return m.sec4.GetProductDefinitionTemplate().GetTypeOfSecondFixedSurface()
}

func (m *message) GetScaleFactorOfSecondFixedSurface() int {
	return m.sec4.GetProductDefinitionTemplate().GetScaleFactorOfSecondFixedSurface()
}

func (m *message) GetScaledValueOfSecondFixedSurface() int {
	return m.sec4.GetProductDefinitionTemplate().GetScaledValueOfSecondFixedSurface()
}

func (m *message) GetProductDefinitionTemplateNumber() int {
	return int(m.sec4.ProductDefinitionTemplateNumber)
}

func (m *message) GetGridPoint(n int) (float32, float32) {
	return m.sec3.GetGridDefinitionTemplate().GetGridPoint(n)
}
