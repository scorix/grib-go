package pdt

import "github.com/scorix/grib-go/pkg/grib2/regulation"

type template0 struct {
	ParameterCategory               uint8
	ParameterNumber                 uint8
	TypeOfGeneratingProcess         uint8
	BackgroundProcess               uint8
	GeneratingProcessIdentifier     uint8
	HoursAfterDataCutoff            uint16
	MinutesAfterDataCutoff          uint8
	IndicatorOfUnitOfTimeRange      uint8
	ForecastTime                    uint32
	TypeOfFirstFixedSurface         uint8
	ScaleFactorOfFirstFixedSurface  uint8
	ScaledValueOfFirstFixedSurface  uint32
	TypeOfSecondFixedSurface        uint8
	ScaleFactorOfSecondFixedSurface uint8
	ScaledValueOfSecondFixedSurface uint32
}

func (t template0) Export() *Template0 {
	return &Template0{
		ParameterCategory:               t.ParameterCategory,
		ParameterNumber:                 t.ParameterNumber,
		TypeOfGeneratingProcess:         regulation.ToInt8(t.TypeOfGeneratingProcess),
		BackgroundProcess:               regulation.ToInt8(t.BackgroundProcess),
		GeneratingProcessIdentifier:     regulation.ToInt8(t.GeneratingProcessIdentifier),
		HoursAfterDataCutoff:            regulation.ToInt16(t.HoursAfterDataCutoff),
		MinutesAfterDataCutoff:          regulation.ToInt8(t.MinutesAfterDataCutoff),
		IndicatorOfUnitOfTimeRange:      regulation.ToInt8(t.IndicatorOfUnitOfTimeRange),
		ForecastTime:                    regulation.ToInt32(t.ForecastTime),
		TypeOfFirstFixedSurface:         regulation.ToInt8(t.TypeOfFirstFixedSurface),
		ScaleFactorOfFirstFixedSurface:  regulation.ToInt8(t.ScaleFactorOfFirstFixedSurface),
		ScaledValueOfFirstFixedSurface:  regulation.ToInt32(t.ScaledValueOfFirstFixedSurface),
		TypeOfSecondFixedSurface:        regulation.ToInt8(t.TypeOfSecondFixedSurface),
		ScaleFactorOfSecondFixedSurface: regulation.ToInt8(t.ScaleFactorOfSecondFixedSurface),
		ScaledValueOfSecondFixedSurface: regulation.ToInt32(t.ScaledValueOfSecondFixedSurface),
	}
}

type Template0 struct {
	ParameterCategory               uint8
	ParameterNumber                 uint8
	TypeOfGeneratingProcess         int8
	BackgroundProcess               int8
	GeneratingProcessIdentifier     int8
	HoursAfterDataCutoff            int16
	MinutesAfterDataCutoff          int8
	IndicatorOfUnitOfTimeRange      int8
	ForecastTime                    int32
	TypeOfFirstFixedSurface         int8
	ScaleFactorOfFirstFixedSurface  int8
	ScaledValueOfFirstFixedSurface  int32
	TypeOfSecondFixedSurface        int8
	ScaleFactorOfSecondFixedSurface int8
	ScaledValueOfSecondFixedSurface int32
}
