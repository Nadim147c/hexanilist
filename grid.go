package main

import (
	"fmt"
	"math"
	"slices"
)

type Grid struct {
	occupied map[string]bool
	radius   float64
}

func NewGrid(radius float64) *Grid {
	return &Grid{occupied: make(map[string]bool), radius: radius}
}

func (hg Grid) key(x, y float64) string {
	sx := math.Round(x / hg.radius)
	sy := math.Round(y / hg.radius)
	return fmt.Sprintf("%.0f,%.0f", sx, sy)
}

func (hg Grid) IsOccupied(p Point) bool {
	return hg.occupied[hg.key(p.X, p.Y)]
}

func (hg *Grid) MarkOccupied(p Point) {
	hg.occupied[hg.key(p.X, p.Y)] = true
}

func GenerateHexagonRing(n int, x, y, radius float64) []Hexagon {
	if n <= 1 {
		return nil
	}

	full := NewGrid(radius)
	empty := make(map[Point]bool)

	centerHex := NewHexagon(x, y, radius, 0)
	full.MarkOccupied(centerHex.Center)
	hexagons := []Hexagon{centerHex}

	for _, p := range centerHex.Neiboors() {
		empty[p] = true
	}

	for range n - 2 {
		d := math.Inf(1)
		hexCenter := centerHex.Center
		for point := range empty {
			newDist := point.Distance(centerHex.Center)
			if newDist < d {
				hexCenter = point
				d = newDist
			}
		}

		hex := NewHexagon(hexCenter.X, hexCenter.Y, radius, 0)
		hexagons = append(hexagons, hex)

		full.MarkOccupied(hexCenter)
		for _, p := range hex.Neiboors() {
			if !full.IsOccupied(p) {
				empty[p] = true
			}
		}

		for p := range empty {
			if full.IsOccupied(p) {
				delete(empty, p)
			}
		}

	}

	center := NewPoint(x, y)

	slices.SortFunc(hexagons, func(i, j Hexagon) int {
		if i.Center.Distance(center) < j.Center.Distance(center) {
			return -1
		}
		if i.Center.Distance(center) > j.Center.Distance(center) {
			return 1
		}
		return 0
	})

	return hexagons
}
