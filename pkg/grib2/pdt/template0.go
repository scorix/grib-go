package pdt

import (
	"time"

	"github.com/scorix/grib-go/pkg/grib2/regulation"
)

type template0 struct {
	ParameterCategory                             uint8  // 10
	ParameterNumber                               uint8  // 11
	TypeOfGeneratingProcess                       uint8  // 12
	BackgroundProcess                             uint8  // 13
	AnalysisOrForecastGeneratingProcessIdentified uint8  // 14
	HoursAfterDataCutoff                          uint16 // 15-16
	MinutesAfterDataCutoff                        uint8  // 17
	IndicatorOfUnitForForecastTime                uint8  // 18
	ForecastTime                                  uint32 // 19-22
	TypeOfFirstFixedSurface                       uint8  // 23
	ScaleFactorOfFirstFixedSurface                int8   // 24
	ScaledValueOfFirstFixedSurface                int32  // 25-28
	TypeOfSecondFixedSurface                      uint8  // 29
	ScaleFactorOfSecondFixedSurface               int8   // 30
	ScaledValueOfSecondFixedSurface               int32  // 31-34
}

func (t template0) Export() *Template0 {
	return &Template0{
		ParameterCategory:       t.ParameterCategory,
		ParameterNumber:         t.ParameterNumber,
		TypeOfGeneratingProcess: regulation.ToInt8(t.TypeOfGeneratingProcess),
		BackgroundProcess:       regulation.ToInt8(t.BackgroundProcess),
		AnalysisOrForecastGeneratingProcessIdentified: regulation.ToInt8(t.AnalysisOrForecastGeneratingProcessIdentified),
		HoursAfterDataCutoff:                          regulation.ToInt16(t.HoursAfterDataCutoff),
		MinutesAfterDataCutoff:                        regulation.ToInt8(t.MinutesAfterDataCutoff),
		IndicatorOfUnitForForecastTime:                IndicatorOfUnitForTime(regulation.ToInt8(t.IndicatorOfUnitForForecastTime)),
		ForecastTime:                                  regulation.ToInt32(t.ForecastTime),
		TypeOfFirstFixedSurface:                       t.TypeOfFirstFixedSurface,
		ScaleFactorOfFirstFixedSurface:                t.ScaleFactorOfFirstFixedSurface,
		ScaledValueOfFirstFixedSurface:                t.ScaledValueOfFirstFixedSurface,
		TypeOfSecondFixedSurface:                      t.TypeOfSecondFixedSurface,
		ScaleFactorOfSecondFixedSurface:               t.ScaleFactorOfSecondFixedSurface,
		ScaledValueOfSecondFixedSurface:               t.ScaledValueOfSecondFixedSurface,
	}
}

type Template0 struct {
	ParameterCategory                             uint8
	ParameterNumber                               uint8
	TypeOfGeneratingProcess                       int8
	BackgroundProcess                             int8
	AnalysisOrForecastGeneratingProcessIdentified int8
	HoursAfterDataCutoff                          int16
	MinutesAfterDataCutoff                        int8
	IndicatorOfUnitForForecastTime                IndicatorOfUnitForTime
	ForecastTime                                  int32
	TypeOfFirstFixedSurface                       uint8
	ScaleFactorOfFirstFixedSurface                int8
	ScaledValueOfFirstFixedSurface                int32
	TypeOfSecondFixedSurface                      uint8
	ScaleFactorOfSecondFixedSurface               int8
	ScaledValueOfSecondFixedSurface               int32
}

func (t Template0) GetParameterCategory() int { return int(t.ParameterCategory) }
func (t Template0) GetParameterNumber() int   { return int(t.ParameterNumber) }

func (t Template0) GetForecastDuration() time.Duration {
	return t.IndicatorOfUnitForForecastTime.AsDuration(int(t.ForecastTime))
}

func (t Template0) GetLevel() int {
	return regulation.CalculateLevel(t.GetScaledValueOfFirstFixedSurface(), t.GetScaleFactorOfFirstFixedSurface())
}

func (t Template0) GetForecast() int {
	return int(t.ForecastTime)
}

func (t Template0) GetTypeOfFirstFixedSurface() int {
	return int(t.TypeOfFirstFixedSurface)
}

func (t Template0) GetScaleFactorOfFirstFixedSurface() int {
	return int(t.ScaleFactorOfFirstFixedSurface)
}

func (t Template0) GetScaledValueOfFirstFixedSurface() int {
	return int(t.ScaledValueOfFirstFixedSurface)
}

func (t Template0) GetTypeOfSecondFixedSurface() int {
	return int(t.TypeOfSecondFixedSurface)
}

func (t Template0) GetScaleFactorOfSecondFixedSurface() int {
	return int(t.ScaleFactorOfSecondFixedSurface)
}

func (t Template0) GetScaledValueOfSecondFixedSurface() int {
	return int(t.ScaledValueOfSecondFixedSurface)
}
