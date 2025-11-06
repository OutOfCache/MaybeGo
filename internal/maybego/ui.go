package maybego

import (
	"image/color"
	"math/rand"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	// "fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
)

const (
	LCDWidth  float32 = 160
	LCDHeight float32 = 144
)

type Interface struct {
	app      fyne.App
	window   fyne.Window
	display   *canvas.Raster
}

func NewUI() *Interface {
	a := app.New()
	w := a.NewWindow("MaybeGo")
	display := canvas.NewRasterWithPixels(
		func(x, y, w, h int) color.Color {
			return color.RGBA {
				R: uint8(rand.Intn(255)),
				G: uint8(rand.Intn(255)),
				B: uint8(rand.Intn(255)),
				A: 0xff,
			}
		})
	// TODO: scaling factor
	display.SetMinSize(fyne.NewSize(LCDWidth, LCDHeight))
	ui := &Interface{app: a, window: w, display: display}

	return ui
}

func (ui *Interface) Update(cycles byte) {
	ui.window.SetContent(ui.display)
	ui.window.Show()
}
