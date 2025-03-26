package main

import (
	"math"

	"github.com/fogleman/gg"
)

// Define Hexagon as an array of 6 Points
type Hexagon struct {
	Points [6]Point
	Center Point
	Angle  float64
	Radius float64
}

// NewHexagon creates a hexagon centered at (x, y) with a given radius and rotation angle
func NewHexagon(x, y, radius, angle float64) Hexagon {
	center := Point{x, y}

	hex := Hexagon{
		Center: center,
		Radius: radius,
		Angle:  angle,
	}
	// Compute 6 vertices of the hexagon
	for i := range 6 {
		// 60-degree increments (π/3 radians)
		theta := angle + float64(i)*(math.Pi/3)
		p := NewPoint(x+radius, y)
		p.Rotate(center, theta)
		hex.Points[i] = p
	}

	return hex
}

func (h Hexagon) Draw(ctx *gg.Context) {
	ctx.ClearPath()
	ctx.MoveTo(h.Points[0].Value())
	for i := range 5 {
		idx := i + 1
		ctx.LineTo(h.Points[idx].Value())
	}
	ctx.ClosePath()
}

func (h Hexagon) Side() float64 {
	return h.Points[0].Distance(h.Points[1])
}

func (h Hexagon) Box() Box {
	xSlice := make([]float64, 6)
	ySlice := make([]float64, 6)
	for i, p := range h.Points {
		xSlice[i] = p.X
		ySlice[i] = p.Y
	}

	start := NewPoint(min(xSlice[0], xSlice[1:]...), min(ySlice[0], ySlice[1:]...))
	end := NewPoint(max(xSlice[0], xSlice[1:]...), max(ySlice[0], ySlice[1:]...))

	return NewBox(start, end)
}

func (h Hexagon) Neiboors() []Point {
	ringRadius := (h.Radius * math.Cos(math.Pi/6)) * 2
	n := make([]Point, 6)

	for i := range 6 {
		angle := h.Angle + (math.Pi / 6) + (float64(i) * (math.Pi / 3))
		point := NewPoint(h.Center.X+ringRadius, h.Center.Y)
		point.Rotate(h.Center, angle)
		n[i] = point
	}
	return n
}
