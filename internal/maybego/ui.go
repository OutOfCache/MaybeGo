package maybego

import (
	// "fmt"
	"image/color"
	// "math/rand"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

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
	vram    *fyne.Container
	emu     *Emulator
}

func GenerateVramTile(tileID int, scale int) func(x, y, w, h int) color.Color {
	return func(x, y, w, h int) color.Color {
		x_conv := (x / scale)
		y_conv := (y / scale)
		if x_conv > 7 || y_conv > 7 {
			return color.RGBA{R: 0, G: 0, B: 0, A: 0}
		}

		address := uint16(0x8000 + (tileID * 16) + (y_conv * 2))
		pixelcolor := (Read(address) >> (7 - x_conv) & 0x1) + (Read(address+1)>>(7-x_conv)&0x1)*2

		return Palette[pixelcolor]
	}
}

func NewUI(logger *Logger) *Interface {
	a := app.New()
	w := a.NewWindow("MaybeGo")
	e := NewEmulator(logger)
	w.Resize(fyne.NewSize(160, 144))
	display := canvas.NewRasterWithPixels(
		func(x, y, w, h int) color.Color {
			// TODO: wtf out-of-bounds??
			// TODO: center layout to set to minimum size?
			if x > 159 || y > 143 {
				return color.RGBA{R: 0, G: 0, B: 0, A: 0}
			}
			return Palette[e.ppu.GetCurrentFrame()[(160*y)+x]]
		})

	vram := container.New(layout.NewGridLayout(16))
	scale := 2
	tile_size := float32(scale * 8)
	for i := 0; i < 384; i++ {
		tile := canvas.NewRasterWithPixels(GenerateVramTile(i, scale))
		tile.SetMinSize(fyne.NewSize(tile_size, tile_size))
		vram.Add(tile)
	}
	// TODO: scaling factor
	display.SetMinSize(fyne.NewSize(160, 144))
	// display_content := container.New(layout.NewCenterLayout(), display)
	content := container.New(layout.NewHBoxLayout(), display, layout.NewSpacer(), vram)
	w.SetContent(content)

	ui := &Interface{app: a, window: w, display: display, vram: vram, emu: e}

	return ui
}

func (ui *Interface) Update() {
	frame_ready := ui.emu.FetchDecodeExec()
	if frame_ready {
		ui.vram.Refresh() // TODO: refresh vram data only if there was a write to tiledata memory
		ui.display.Refresh()
		ui.window.Show()
	}
}

func (ui *Interface) LoadRom(rom *[]byte) {
	for i, buffer := range *rom {
		Write(uint16(i), buffer)
	}

	// TODO: option to skip boot rom or not?

	// for i, buffer := range (*rom)[0x100:] {
	// 	Write(uint16(i+0x100), buffer)
	// }

	// for i, buffer := range *rom {
	// 	Write(uint16(i+0x100), buffer)
	// }

}
