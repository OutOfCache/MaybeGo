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

var defaultColor = color.RGBA{R: 0xFF, G: 0x80, B: 0x80, A: 0xFF}
var Palette = []color.RGBA {
		{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF },
		{R: 0x80, G: 0x80, B: 0x80, A: 0xFF },
		{R: 0x08, G: 0x08, B: 0x08, A: 0xFF },
		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF },
	}

type Interface struct {
	app      fyne.App
	window   fyne.Window
	display  *canvas.Raster
	ppu      *PPU
}

func NewUI() *Interface {
	a := app.New()
	w := a.NewWindow("MaybeGo")
	display := canvas.NewRasterWithPixels(
		func(x, y, w, h int) color.Color {
			return Palette[rand.Intn(3)]
		})
	// TODO: scaling factor
	display.SetMinSize(fyne.NewSize(LCDWidth, LCDHeight))
	ui := &Interface{app: a, window: w, display: display}

	return ui
}

func (ui *Interface) Update(frame [160 * 144]byte) {
	ui.window.SetContent(ui.display)
	ui.window.Show()
}
