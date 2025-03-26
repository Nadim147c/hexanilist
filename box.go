package main

import (
	"image"
	"math"
)

type Box struct {
	X, Y, W, H float64
}

func NewBox(start, end Point) Box {
	return Box{start.X, start.Y, end.X - start.X, end.Y - start.Y}
}

func (b Box) Values() (float64, float64, float64, float64) {
	return math.Round(b.X), math.Round(b.Y), math.Round(b.W), math.Round(b.H)
}

func (b Box) Start() (int, int) {
	x, y, _, _ := b.Values()
	return int(x), int(y)
}

func (b Box) End() (int, int) {
	x, y, w, h := b.Values()
	return int(x + w), int(y + h)
}

func (b Box) Size() (int, int) {
	_, _, w, h := b.Values()
	return int(w), int(h)
}

func (b Box) Rect() image.Rectangle {
	x, y, w, h := b.Values()
	return image.Rect(int(x), int(y), int(x+w), int(y+h))
}
