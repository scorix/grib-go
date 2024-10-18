package colormap

import (
	"image/color"
	"math"
)

// ColorPoint 定义了颜色映射中的一个点
type ColorPoint struct {
	Value float64
	Color color.Color
}

// ColorMap 定义了一个颜色映射
type ColorMap struct {
	Points []ColorPoint
}

// NewColorMap 创建一个新的颜色映射
func NewColorMap(points ...ColorPoint) *ColorMap {
	return &ColorMap{Points: points}
}

// GetColor 根据给定的值返回对应的颜色
func (cm *ColorMap) GetColor(value, min, max float64) color.Color {
	// 将值归一化到 0-1 范围
	normalized := (value - min) / (max - min)
	normalized = math.Max(0, math.Min(1, normalized))

	// 找到normalized所在的区间
	var i int
	for i = 0; i < len(cm.Points)-1; i++ {
		if normalized <= cm.Points[i+1].Value {
			break
		}
	}

	// 计算插值
	p1 := cm.Points[i]
	p2 := cm.Points[i+1]
	t := (normalized - p1.Value) / (p2.Value - p1.Value)

	return interpolateColor(p1.Color, p2.Color, t)
}

// interpolateColor 在两个颜色之间进行线性插值
func interpolateColor(c1, c2 color.Color, t float64) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	return color.RGBA{
		R: uint8(float64(r1) + t*float64(r2-r1)),
		G: uint8(float64(g1) + t*float64(g2-g1)),
		B: uint8(float64(b1) + t*float64(b2-b1)),
		A: uint8(float64(a1) + t*float64(a2-a1)),
	}
}

// 预定义一些常用的颜色映射
var (
	TemperatureColorMap = NewColorMap(
		ColorPoint{0.0, color.RGBA{0, 0, 128, 255}},   // 深蓝色 (极冷)
		ColorPoint{0.2, color.RGBA{0, 128, 255, 255}}, // 蓝色
		ColorPoint{0.4, color.RGBA{0, 255, 255, 255}}, // 青色
		ColorPoint{0.6, color.RGBA{255, 255, 0, 255}}, // 黄色
		ColorPoint{0.8, color.RGBA{255, 128, 0, 255}}, // 橙色
		ColorPoint{1.0, color.RGBA{128, 0, 0, 255}},   // 深红色 (极热)
	)

	RainbowColorMap = NewColorMap(
		ColorPoint{0.0, color.RGBA{255, 0, 0, 255}},   // 红色
		ColorPoint{0.2, color.RGBA{255, 127, 0, 255}}, // 橙色
		ColorPoint{0.4, color.RGBA{255, 255, 0, 255}}, // 黄色
		ColorPoint{0.6, color.RGBA{0, 255, 0, 255}},   // 绿色
		ColorPoint{0.8, color.RGBA{0, 0, 255, 255}},   // 蓝色
		ColorPoint{1.0, color.RGBA{139, 0, 255, 255}}, // 紫色
	)

	GrayColorMap = NewColorMap(
		ColorPoint{0.0, color.Gray16{0}},
		ColorPoint{0.2, color.Gray16{255 / 6}},
		ColorPoint{0.4, color.Gray16{255 / 3}},
		ColorPoint{0.6, color.Gray16{255 / 2}},
		ColorPoint{0.8, color.Gray16{255 * 2 / 3}},
		ColorPoint{1.0, color.Gray16{255}},
	)
)
