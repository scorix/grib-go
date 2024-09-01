package datapacking

type UnpackReader interface {
	ScaleFunc() func(uint64) float64
	GetBits() uint8
}
