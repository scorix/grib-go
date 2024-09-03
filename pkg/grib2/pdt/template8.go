package pdt

import (
	"math"
	"time"
)

type template8 struct {
	*template0 // 10-34
	template8fields
}

type template8fields struct {
	Year                                               uint16 // 35-36
	Month                                              uint8  // 37
	Day                                                uint8  // 38
	Hour                                               uint8  // 39
	Minute                                             uint8  // 40
	Second                                             uint8  // 41
	NumberOfTimeRanges                                 uint8  // 42 n
	TotalNumberOfDataValuesMissingInStatisticalProcess uint32 // 43-46
	// 47 - 58 Specification of the outermost (or only) time range over which statistical processing is done
	StatisticalProcess                uint8  // 47
	TypeOfTimeIncrement               uint8  // 48
	IndicatorOfUnitOfTimeForTimeRange uint8  // 49
	LengthOfTimeRange                 uint32 // 50-53
	IndicatorOfUnitOfTimeForIncrement uint8  // 54
	TimeIncrement                     uint32 // 55-58
	// 59-70 As octets 47 to 58, next innermost step of processing
}

func (fields template8fields) GetAdditionalTimeRangeSpecifications() []byte {
	if fields.NumberOfTimeRanges == 0 {
		return nil
	}

	// 59 - nn These octets are included only if n>1, where nn = 46 + 12 x n
	var nn = 46 + 12*fields.NumberOfTimeRanges

	return make([]byte, nn)
}

func (t template8) Export() *Template8 {
	return &Template8{
		Template0: t.template0.Export(),
	}
}

type Template8 struct {
	*Template0
}

func (t Template8) GetParameterCategory() int { return int(t.ParameterCategory) }
func (t Template8) GetParameterNumber() int   { return int(t.ParameterNumber) }
func (t Template8) GetForecastDuration() time.Duration {
	return t.IndicatorOfUnitForForecastTime.AsDuration(int(t.ForecastTime))
}
func (t Template8) GetLevel() int {
	return int(t.ScaledValueOfFirstFixedSurface) * int(math.Pow10(int(t.ScaleFactorOfFirstFixedSurface)))
}
