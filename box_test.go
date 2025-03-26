package main

import (
	"image"
	"testing"
)

func TestBox(t *testing.T) {
	start := Point{X: 10.3, Y: 20.7}
	end := Point{X: 50.9, Y: 70.2}
	box := NewBox(start, end)

	expectedX, expectedY, expectedW, expectedH := 10, 21, 41, 50
	x, y, w, h := box.Values()
	if int(x) != expectedX || int(y) != expectedY || int(w) != expectedW || int(h) != expectedH {
		t.Errorf("Expected Values() to return (%d, %d, %d, %d), got (%d, %d, %d, %d)", expectedX, expectedY, expectedW, expectedH, int(x), int(y), int(w), int(h))
	}

	sx, sy := box.Start()
	if sx != expectedX || sy != expectedY {
		t.Errorf("Expected Start() to return (%d, %d), got (%d, %d)", expectedX, expectedY, sx, sy)
	}

	ex, ey := box.End()
	if ex != expectedX+expectedW || ey != expectedY+expectedH {
		t.Errorf("Expected End() to return (%d, %d), got (%d, %d)", expectedX+expectedW, expectedY+expectedH, ex, ey)
	}

	sw, sh := box.Size()
	if sw != expectedW || sh != expectedH {
		t.Errorf("Expected Size() to return (%d, %d), got (%d, %d)", expectedW, expectedH, sw, sh)
	}

	rect := box.Rect()
	expectedRect := image.Rect(expectedX, expectedY, expectedX+expectedW, expectedY+expectedH)
	if rect != expectedRect {
		t.Errorf("Expected Rect() to return %v, got %v", expectedRect, rect)
	}
}
