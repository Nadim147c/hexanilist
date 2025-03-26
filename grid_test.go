package main

import "testing"

func TestGrid(t *testing.T) {
	radius := 10.0
	grid := NewGrid(radius)
	point := Point{X: 15.3, Y: 25.7}

	if grid.IsOccupied(point) {
		t.Errorf("Expected IsOccupied() to return false, got true")
	}

	grid.MarkOccupied(point)

	if !grid.IsOccupied(point) {
		t.Errorf("Expected IsOccupied() to return true after MarkOccupied, got false")
	}
}
