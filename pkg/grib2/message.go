package grib2

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"time"

	"github.com/scorix/grib-go/internal/pkg/bitio"
	"github.com/scorix/grib-go/pkg/grib2/cache"
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
	GetGridDefinitionTemplate() gdt.Template

	ReadData() ([]float32, error)
	Image() (image.Image, error)
	Step() int

	GetGridPointLL(n int) (float32, float32, bool)
	GetGridPointFromLL(lat float32, lon float32) int
	GetNi() int
	GetNj() int
	GetSize() int64
	DumpMessageIndex() (*MessageIndex, error)
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
	offset int64
	sec0   *section0
	sec1   *section1
	sec2   *section2
	sec3   *section3
	sec4   *section4
	sec5   *section5
	sec6   *section6
	sec7   *section7
	sec8   *section8
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

func (m *message) ReadData() ([]float32, error) {
	tpl := m.sec5.GetDataRepresentationTemplate()
	if err := m.sec7.LoadData(); err != nil {
		return nil, fmt.Errorf("load data from section 7: %w", err)
	}

	data, err := tpl.ReadAllData(bitio.NewReader(bytes.NewReader(m.sec7.Data)))
	if err != nil {
		return nil, fmt.Errorf("read data using template %T: %w", tpl, err)
	}
	return data, nil
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

func (m *message) GetGridPointLL(n int) (float32, float32, bool) {
	tpl := m.sec3.GetGridDefinitionTemplate()
	return tpl.GetGridPoint(n)
}

func (m *message) GetGridPointFromLL(lat float32, lon float32) int {
	tpl := m.sec3.GetGridDefinitionTemplate()
	return tpl.GetGridIndex(lat, lon)
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
	return m.sec7.GetDataOffset()
}

func (m *message) GetSize() int64 {
	return int64(m.sec0.GribLength)
}

func (m *message) GetGridDefinitionTemplate() gdt.Template {
	return m.sec3.GridDefinitionTemplate
}

func (m *message) assignSection(sec Section) error {
	switch sec.Number() {
	case 0:
		m.sec0 = sec.(*section0)
	case 1:
		m.sec1 = sec.(*section1)
	case 2:
		m.sec2 = sec.(*section2)
	case 3:
		m.sec3 = sec.(*section3)
	case 4:
		m.sec4 = sec.(*section4)
	case 5:
		m.sec5 = sec.(*section5)
	case 6:
		m.sec6 = sec.(*section6)
	case 7:
		m.sec7 = sec.(*section7)
	case 8:
		m.sec8 = sec.(*section8)
	default:
		return fmt.Errorf("unknown section number: %d", sec.Number())
	}
	return nil
}

func (m *message) DumpMessageIndex() (*MessageIndex, error) {
	return &MessageIndex{
		Offset:         m.offset,
		Size:           m.GetSize(),
		DataOffset:     m.GetDataOffset(),
		GridDefinition: m.GetGridDefinitionTemplate(),
		Packing:        m.GetDataRepresentationTemplate(),
	}, nil
}

func (m *message) Image() (image.Image, error) {
	tpl := m.GetDataRepresentationTemplate()

	switch t := tpl.(type) {
	case *gridpoint.PortableNetworkGraphics:
		if err := m.sec7.LoadData(); err != nil {
			return nil, fmt.Errorf("load data from section 7: %w", err)
		}
		return t.Image(bitio.NewReader(bytes.NewReader(m.sec7.Data)))
	default:
		return nil, fmt.Errorf("data is not an image: %T", tpl)
	}
}

type MessageReader interface {
	ReadLL(ctx context.Context, lat float32, lon float32) (float32, float32, float32, error)
	GetGridIndex(lat float32, lon float32) int
	GetGridPoint(n int) (float32, float32, bool)
}

type simplePackingMessageReader struct {
	sp    *gridpoint.SimplePacking
	spr   *gridpoint.SimplePackingReader
	gdt   gdt.Template
	cache cache.GridCache
}

func NewSimplePackingMessageReaderFromMessage(r io.ReaderAt, m IndexedMessage, opts ...SimplePackingMessageReaderOptions) (MessageReader, error) {
	sp, ok := m.GetDataRepresentationTemplate().(*gridpoint.SimplePacking)
	if !ok {
		return nil, fmt.Errorf("unsupported data representation template: %T", m.GetDataRepresentationTemplate())
	}

	gdt := m.GetGridDefinitionTemplate()

	return NewSimplePackingMessageReader(r, m.GetOffset(), m.GetSize(), m.GetDataOffset(), sp, gdt, opts...)
}

type SimplePackingMessageReaderOptions func(r *simplePackingMessageReader)

func WithBoundaryCache(minLat, maxLat, minLon, maxLon float32, store cache.Store) SimplePackingMessageReaderOptions {
	return func(r *simplePackingMessageReader) {
		r.cache = cache.NewBoundary(minLat, maxLat, minLon, maxLon, r.spr, store)
	}
}

func WithCustomCacheStrategy(inCache func(lat, lon float32) bool, store cache.Store) SimplePackingMessageReaderOptions {
	return func(r *simplePackingMessageReader) {
		r.cache = cache.NewCustom(inCache, r.spr, store)
	}
}

func NewSimplePackingMessageReader(r io.ReaderAt, messageOffset int64, messageSize int64, dataOffset int64, sp *gridpoint.SimplePacking, gdt gdt.Template, opts ...SimplePackingMessageReaderOptions) (MessageReader, error) {
	spr := gridpoint.NewSimplePackingReader(r, dataOffset, messageOffset+messageSize, sp)

	mr := &simplePackingMessageReader{
		spr:   spr,
		sp:    sp,
		gdt:   gdt,
		cache: cache.NewNoCache(spr),
	}

	for _, opt := range opts {
		opt(mr)
	}

	return mr, nil
}

func NewSimplePackingMessageReaderFromMessageIndex(r io.ReaderAt, mi *MessageIndex, opts ...SimplePackingMessageReaderOptions) (MessageReader, error) {
	sp, ok := mi.Packing.(*gridpoint.SimplePacking)
	if !ok {
		return nil, fmt.Errorf("unsupported packing: %T", mi.Packing)
	}

	return NewSimplePackingMessageReader(r, mi.Offset, mi.Size, mi.DataOffset, sp, mi.GridDefinition, opts...)
}

func (r *simplePackingMessageReader) ReadLL(ctx context.Context, lat float32, lon float32) (float32, float32, float32, error) {
	grid := r.gdt.GetGridIndex(lat, lon)
	lat, lng, _ := r.gdt.GetGridPoint(grid)

	v, err := r.cache.ReadGridAt(ctx, grid, lat, lng)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("read grid at point %d (lat: %f, lon: %f): %w", grid, lat, lng, err)
	}

	return lat, lng, v, nil
}

func (r *simplePackingMessageReader) GetGridIndex(lat float32, lon float32) int {
	return r.gdt.GetGridIndex(lat, lon)
}

func (r *simplePackingMessageReader) GetGridPoint(n int) (float32, float32, bool) {
	return r.gdt.GetGridPoint(n)
}
