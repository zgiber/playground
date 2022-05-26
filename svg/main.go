package main

import (
	"fmt"
	"image/color"

	"github.com/llgcode/draw2d/draw2dsvg"
)

func main() {
	// Initialize the graphic context on an RGBA image
	dest := draw2dsvg.NewSvg()
	gc := draw2dsvg.NewGraphicContext(dest)

	// Set some properties
	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
	gc.SetLineWidth(5)

	// Draw a closed shape
	gc.MoveTo(10, 10) // should always be called first for a new path
	gc.LineTo(100, 50)
	gc.QuadCurveTo(100, 10, 10, 10)
	gc.Close()
	gc.FillStroke()

	// gc.MoveTo(100, 100)
	gc.ArcTo(150, 150, 100, 100, 1, 1)
	gc.FillStroke()

	// Save to file
	if err := draw2dsvg.SaveToSvgFile("test.svg", dest); err != nil {
		fmt.Println(err)
	}
}
