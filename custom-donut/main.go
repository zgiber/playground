package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"

	"github.com/fogleman/gg"
	"github.com/llgcode/draw2d/draw2dimg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var (
	Grey = color.RGBA{
		R: 231,
		G: 231,
		B: 225,
		A: 255,
	}

	DarkGrey = color.RGBA{
		R: 139,
		G: 137,
		B: 137,
		A: 255,
	}

	Blue = color.RGBA{
		R: 55,
		G: 114,
		B: 199,
		A: 255,
	}

	DarkBlue = color.RGBA{
		R: 68,
		G: 69,
		B: 111,
		A: 255,
	}

	Cyan = color.RGBA{
		R: 84,
		G: 185,
		B: 212,
		A: 255,
	}

	LightBlue = color.RGBA{
		R: 231,
		G: 231,
		B: 225,
		A: 255,
	}
)

func main() {

	width := 400
	height := 400
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	gc := draw2dimg.NewGraphicContext(img)
	Donut(gc, img, DonutOptions{
		InnerRadius: 150,
		OuterRadius: 199,
		StartAngle:  20,
		EndAngle:    290,
		Color:       Blue,
	})

	dc := gg.NewContext(width, height)
	dc.DrawImage(img, 0, 0)

	if err := dc.SavePNG("donut.png"); err != nil {
		log.Println(err)
	}

	// // // Save to png
	// fileName := "donut.png"
	// err := draw2dimg.SaveToPngFile(fileName, img)
	// if err != nil {
	// 	log.Printf("Saving %q failed: %v", fileName, err)
	// 	return
	// }

	// log.Printf("Succesfully created %q", fileName)
}

func bgOptions(opts DonutOptions) DonutOptions {
	return DonutOptions{
		InnerRadius: opts.InnerRadius + opts.InnerRadius/1000,
		OuterRadius: opts.OuterRadius - opts.OuterRadius/1000,
		StartAngle:  0,
		EndAngle:    360,
		Color:       Grey,
	}
}

type DonutOptions struct {
	InnerRadius float64
	OuterRadius float64
	StartAngle  float64
	EndAngle    float64
	Color       color.RGBA
	Text        string
}

func Donut(gc *draw2dimg.GraphicContext, dest *image.RGBA, opts DonutOptions) {
	DounutSection(gc, dest, bgOptions(opts))
	DounutSection(gc, dest, opts)
}

func AsBase64PNG(img *image.RGBA) ([]byte, error) {
	pngOut := &bytes.Buffer{}
	err := png.Encode(pngOut, img)
	if err != nil {
		return nil, err
	}

	base64Out := &bytes.Buffer{}
	_, err = base64.NewEncoder(base64.StdEncoding, base64Out).Write(pngOut.Bytes())
	return base64Out.Bytes(), err
}

func DounutSection(gc *draw2dimg.GraphicContext, dest *image.RGBA, opts DonutOptions) {
	centerX := float64(dest.Rect.Dx()) / 2
	centerY := float64(dest.Rect.Dy()) / 2
	outerRadius := opts.OuterRadius
	innerRadius := opts.InnerRadius
	angle := math.Abs(opts.EndAngle - opts.StartAngle)

	startX, startY := pointOnCircle(centerX, centerY, outerRadius, opts.StartAngle)
	fmt.Println(startX, startY)
	gc.MoveTo(startX, startY)
	gc.ArcTo(centerX, centerY, outerRadius, outerRadius, d2r(270+opts.StartAngle), d2r(angle))
	gc.LineTo(pointOnCircle(centerX, centerY, innerRadius, angle+opts.StartAngle))
	gc.ArcTo(centerX, centerY, innerRadius, innerRadius, d2r(270+opts.StartAngle+angle), d2r(angle*-1))
	gc.LineTo(startX, startY)
	gc.Close()

	gc.SetFillColor(opts.Color)
	gc.Fill()
	// gc.
}

func d2r(degree float64) float64 {
	return degree * (math.Pi / 180)
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{200, 100, 0, 255}
	point := fixed.Point26_6{fixed.I(x), fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func pointOnCircle(cx, cy, r, angle float64) (float64, float64) {
	x := r * math.Sin(d2r(angle))
	y := r * math.Cos(d2r(angle))

	angle = angle - float64(int(angle/360)*360)
	// fmt.Println(angle)

	if d2r(180) < angle && angle < d2r(360) {
		x = -1 * x
	}

	if d2r(270) < angle || angle < d2r(90) {
		y = -1 * y
	}

	// fmt.Println(x, y)

	return cx + x, cy + y
}
