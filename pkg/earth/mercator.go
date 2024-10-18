package earth

import (
	"math"
)

func LatLonToMercator(lat, lon float64) (float64, float64) {
	x := lon * 20037508.34 / 180
	y := math.Log(math.Tan((90+lat)*math.Pi/360)) / (math.Pi / 180)
	y = (y * 20037508.34) / 180
	return x, y
}

func MercatorToLatLon(x, y float64) (float64, float64) {
	x = x * 180 / 20037508.34
	y = y * 180 / 20037508.34
	y = (math.Atan(math.Pow(math.E, y*(math.Pi/180)))*360)/math.Pi - 90
	return x, y
}
