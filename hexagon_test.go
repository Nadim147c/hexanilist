package main

import (
	"math"
	"testing"
)

func TestNewHexagon(t *testing.T) {
	h := NewHexagon(0, 0, 10, 0)

	if h.Center.X != 0 || h.Center.Y != 0 {
		t.Errorf("Expected center (0,0), got (%f,%f)", h.Center.X, h.Center.Y)
	}

	if h.Radius != 10 {
		t.Errorf("Expected radius 10, got %f", h.Radius)
	}
}

func TestHexagonSide(t *testing.T) {
	h := NewHexagon(0, 0, 10, 0)
	side := h.Side()
	expected := 10.0 // Since it's an equilateral hexagon

	if math.Abs(side-expected) > 1e-6 {
		t.Errorf("Expected side length %f, got %f", expected, side)
	}
}

func TestHexagonNeighbors(t *testing.T) {
	h := NewHexagon(0, 0, 10, 0)
	n := h.Neiboors()

	if len(n) != 6 {
		t.Errorf("Expected 6 neighbors, got %d", len(n))
	}

	expectedDistance := (h.Radius * math.Sin(math.Pi/6)) * 2

	for _, p := range n {
		d := h.Center.Distance(p)
		if math.Abs(d-expectedDistance) > 1e-6 {
			t.Errorf("Expected neighbor distance %f, got %f", expectedDistance, d)
		}
	}
}
