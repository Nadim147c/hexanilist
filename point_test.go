package main

import (
	"math"
	"testing"
)

// Unit test for Rotate function
func TestRotate(t *testing.T) {
	tests := []struct {
		point    Point
		base     Point
		angle    float64
		expected Point
	}{
		{Point{10, 0}, Point{0, 0}, math.Pi / 2, Point{0, 10}},
		{Point{0, 10}, Point{0, 0}, math.Pi / 2, Point{-10, 0}},
		{Point{-10, 0}, Point{0, 0}, math.Pi / 2, Point{0, -10}},
		{Point{0, -10}, Point{0, 0}, math.Pi / 2, Point{10, 0}},
	}

	for _, tt := range tests {
		p := tt.point
		p.Rotate(tt.base, tt.angle)
		if p != tt.expected {
			t.Errorf("Rotate(%v, %v, %v) = %v; want %v", tt.point, tt.base, tt.angle, p, tt.expected)
		}
	}
}
