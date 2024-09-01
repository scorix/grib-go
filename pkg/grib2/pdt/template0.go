package pdt

type Template0 struct {
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
