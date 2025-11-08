package maybego

import (
	"image/color"
	// "math/rand"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	// "fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
)

const (
	LCDWidth  float32 = 160.0
	LCDHeight float32 = 144.0
)

var defaultColor = color.RGBA{R: 0xFF, G: 0x80, B: 0x80, A: 0xFF}
var Palette = []color.RGBA{
	{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
	{R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
	{R: 0x08, G: 0x08, B: 0x08, A: 0xFF},
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
}

type Interface struct {
	app     fyne.App
	window  fyne.Window
	display *canvas.Raster
	emu     *Emulator
}

func NewUI(logger *Logger) *Interface {
	a := app.New()
	w := a.NewWindow("MaybeGo")
	e := NewEmulator(logger)
	w.Resize(fyne.NewSize(160, 144))
	display := canvas.NewRasterWithPixels(
		func(x, y, w, h int) color.Color {
			// TODO: wtf out-of-bounds??
			if x > 159 || y > 143 {
				return color.RGBA{R: 0, G: 0, B: 0, A: 0}
			}
			return Palette[e.ppu.GetCurrentFrame()[(160*y)+x]]
		})
	// TODO: scaling factor
	display.Resize(fyne.NewSize(160, 144))
	w.SetContent(display)

	ui := &Interface{app: a, window: w, display: display, emu: e}

	return ui
}

func (ui *Interface) Update() {
	frame_ready := ui.emu.FetchDecodeExec()
	if frame_ready {
		ui.display.Refresh()
		ui.window.Show()
	}
}

func (ui *Interface) LoadRom(rom *[]byte) {
	for i, buffer := range *rom {
		Write(uint16(i), buffer)
	}

	// TODO: option to skip boot rom or not?

	// for i, buffer := range rom[0x100:] {
	// 	maybego.Write(uint16(i + 0x100), buffer)
	// }

	// for i, buffer := range rom {
	// 	maybego.Write(uint16(i + 0x100), buffer)
	// }

}
