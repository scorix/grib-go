package pdt

import "time"

type IndicatorOfUnitForTime uint8

const (
	IndicatorOfUnitForTimeMinute  IndicatorOfUnitForTime = 0
	IndicatorOfUnitForTimeHour    IndicatorOfUnitForTime = 1
	IndicatorOfUnitForTimeDay     IndicatorOfUnitForTime = 2
	IndicatorOfUnitForTimeMonth   IndicatorOfUnitForTime = 3
	IndicatorOfUnitForTimeYear    IndicatorOfUnitForTime = 4
	IndicatorOfUnitForTimeDecade  IndicatorOfUnitForTime = 5 // 10 years
	IndicatorOfUnitForTimeNormal  IndicatorOfUnitForTime = 6 // 30 years
	IndicatorOfUnitForTimeCentury IndicatorOfUnitForTime = 7 // 100 years
	IndicatorOfUnitForTime3Hours  IndicatorOfUnitForTime = 10
	IndicatorOfUnitForTime6Hours  IndicatorOfUnitForTime = 11
	IndicatorOfUnitForTime12Hours IndicatorOfUnitForTime = 12
	IndicatorOfUnitForTimeSecond  IndicatorOfUnitForTime = 13
)

func (u IndicatorOfUnitForTime) AsDuration(i int) time.Duration {
	var t time.Duration

	switch u {
	case IndicatorOfUnitForTimeSecond:
		t = time.Second

	case IndicatorOfUnitForTimeMinute:
		t = time.Minute

	case IndicatorOfUnitForTimeHour:
		t = time.Hour

	case IndicatorOfUnitForTime3Hours:
		t = time.Hour * 3

	case IndicatorOfUnitForTime6Hours:
		t = time.Hour * 6

	case IndicatorOfUnitForTime12Hours:
		t = time.Hour * 12

	case IndicatorOfUnitForTimeDay:
		t = time.Hour * 24

	case IndicatorOfUnitForTimeMonth:
		now := time.Now()
		t = now.AddDate(0, 1, 0).Sub(now)

	case IndicatorOfUnitForTimeYear:
		now := time.Now()
		t = now.AddDate(1, 0, 0).Sub(now)

	case IndicatorOfUnitForTimeDecade:
		now := time.Now()
		t = now.AddDate(10, 0, 0).Sub(now)

	case IndicatorOfUnitForTimeNormal:
		now := time.Now()
		t = now.AddDate(30, 0, 0).Sub(now)

	case IndicatorOfUnitForTimeCentury:
		now := time.Now()
		t = now.AddDate(100, 0, 0).Sub(now)
	}

	return time.Duration(i) * t
}
