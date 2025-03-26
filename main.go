package main

import "github.com/fogleman/gg"

func main() {
	const W = 2000
	const H = 2000
	ctx := gg.NewContext(W, H)

	ctx.SetRGB(1, 1, 1) // White color
	ctx.Clear()

	ctx.SetRGB(0, 0, 0) // Black color

	hexs := GenerateHexagonRing(20, 1000, 1000, 100)
	for _, hex := range hexs {
		hex.Draw(ctx)
		ctx.SetLineWidth(10)
		ctx.Stroke()
	}

	err := ctx.SavePNG("hexagon.png")
	if err != nil {
		panic(err)
	}
}
