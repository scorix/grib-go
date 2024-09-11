package grib2

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/icza/bitio"
	"github.com/scorix/grib-go/pkg/grib2/drt"
	gridpoint "github.com/scorix/grib-go/pkg/grib2/drt/grid_point"
	"github.com/scorix/grib-go/pkg/grib2/gdt"
)

type Message interface {
	Parameter
	HasLevel

	GetProductDefinitionTemplateNumber() int
	GetDataRepresentationTemplateNumber() int
	GetDataRepresentationTemplate() drt.Template
	GetScanningMode() (gdt.ScanningMode, error)

	ReadData() ([]float64, error)
	Step() int

	GetGridPointLL(n int) (float32, float32, error)
	GetGridPointFromLL(float32, float32) (int, error)
	GetNi() int
	GetNj() int
	GetSize() int64
}

type IndexedMessage interface {
	Message

	GetOffset() int64
	GetDataOffset() int64
}

type Parameter interface {
	GetDiscipline() int
	GetParameterCategory() int
	GetParameterNumber() int
	GetTimestamp(loc *time.Location) time.Time
	GetForecastTime(loc *time.Location) time.Time
}

type HasLevel interface {
	GetLevel() int
	GetTypeOfFirstFixedSurface() int
	GetScaleFactorOfFirstFixedSurface() int
	GetScaledValueOfFirstFixedSurface() int
	GetTypeOfSecondFixedSurface() int
	GetScaleFactorOfSecondFixedSurface() int
	GetScaledValueOfSecondFixedSurface() int
}

type message struct {
	offset     int64
	dataOffset int64

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

func (m *message) GetDataRepresentationTemplateNumber() int {
	return int(m.sec5.DataRepresentationTemplateNumber)
}

func (m *message) GetDataRepresentationTemplate() drt.Template {
	return m.sec5.DataRepresentationTemplate
}

func (m *message) GetGridPointLL(n int) (float32, float32, error) {
	sm, err := m.sec3.GetGridDefinitionTemplate().GetScanningMode()
	if err != nil {
		return 0, 0, err
	}

	lat, lng := sm.GetGridPointLL(n)

	return lat, lng, nil
}

func (m *message) GetGridPointFromLL(lat float32, lon float32) (int, error) {
	sm, err := m.sec3.GetGridDefinitionTemplate().GetScanningMode()
	if err != nil {
		return 0, err
	}

	n := sm.GetGridPointFromLL(lat, lon)

	return n, nil
}

func (m *message) GetNi() int {
	return int(m.sec3.GridDefinitionTemplate.GetNi())
}

func (m *message) GetNj() int {
	return int(m.sec3.GridDefinitionTemplate.GetNj())
}

func (m *message) GetOffset() int64 {
	return m.offset
}

func (m *message) GetDataOffset() int64 {
	return m.dataOffset
}

func (m *message) GetSize() int64 {
	return int64(m.sec0.GribLength)
}

func (m *message) GetScanningMode() (gdt.ScanningMode, error) {
	return m.sec3.GridDefinitionTemplate.GetScanningMode()
}

type MessageReader interface {
	ReadLL(float32, float32) (float64, error)
}

type simplePackingMessageReader struct {
	sp  *gridpoint.SimplePacking
	spr *gridpoint.SimplePackingReader
	sm  gdt.ScanningMode
}

func NewSimplePackingMessageReader(r io.ReaderAt, m IndexedMessage) (MessageReader, error) {
	sp, ok := m.GetDataRepresentationTemplate().(*gridpoint.SimplePacking)
	if !ok {
		return nil, fmt.Errorf("%T is not supported", m.GetDataRepresentationTemplate())
	}

	sm, err := m.GetScanningMode()
	if err != nil {
		return nil, fmt.Errorf("get scanning mode: %w", err)
	}

	return &simplePackingMessageReader{
		spr: gridpoint.NewSimplePackingReader(r, m.GetDataOffset(), m.GetSize()+m.GetOffset(), sp),
		sp:  sp,
		sm:  sm,
	}, nil
}

func (r *simplePackingMessageReader) ReadLL(lat float32, lon float32) (float64, error) {
	grid := r.sm.GetGridPointFromLL(lat, lon)

	v, err := r.spr.ReadGridAt(grid)
	if err != nil {
		return 0, fmt.Errorf("grid %d: %w", grid, err)
	}

	return v, nil
}
