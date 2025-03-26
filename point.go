package main

import "math"

type Point struct {
	X, Y float64
}

func (p Point) Value() (float64, float64) {
	return p.X, p.Y
}

func NewPoint(x, y float64) Point {
	return Point{x, y}
}

// Rotate rotates the point around a base point by a given angle (in radians)
func (p *Point) Rotate(base Point, angle float64) {
	cosTheta := math.Cos(angle)
	sinTheta := math.Sin(angle)
	dx := float64(p.X - base.X)
	dy := float64(p.Y - base.Y)

	newX := cosTheta*dx - sinTheta*dy + float64(base.X)
	newY := sinTheta*dx + cosTheta*dy + float64(base.Y)

	p.X = math.Round(newX)
	p.Y = math.Round(newY)
}

func (p Point) Distance(d Point) float64 {
	return math.Sqrt(math.Pow(p.X-d.X, 2) + math.Pow(p.Y-d.Y, 2))
}
