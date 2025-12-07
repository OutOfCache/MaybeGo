package maybego

import (
	"image/color"
	"time"

	// "math/rand"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

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

type cpu_state_bindings struct {
	cycles binding.Int
}

type Interface struct {
	app       fyne.App
	window    fyne.Window
	display   *canvas.Raster
	vram      *fyne.Container
	cpu       *fyne.Container
	cpu_state *cpu_state_bindings
	emu       *Emulator
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
	display := canvas.NewRasterWithPixels(
		func(x, y, w, h int) color.Color {
			// frame_ready := e.FetchDecodeExec()
			// if frame_ready {
			// if ui.vram.Hidden {
			// 	ui.vram.Refresh()
			// 	ui.vram.Show()
			// }
			// ui.vram.Hide()
			// ui.vram.Refresh() // TODO: refresh vram data only if there was a write to tiledata memory
			// ui.display.Refresh()
			// ui.window.Show()
			// TODO: wtf out-of-bounds??
			// TODO: center layout to set to minimum size?
			// }
			if x > 159 || y > 143 {
				return color.RGBA{R: 0, G: 0, B: 0, A: 0}
			}
			return Palette[e.ppu.GetCurrentFrame()[(160*y)+x]]
		})

	// ============= Debugger: CPU State =============
	cpu_state_container := container.New(layout.NewVBoxLayout())
	cpu_state_label := widget.NewLabel("CPU State")
	cpu_state := &cpu_state_bindings{cycles: binding.NewInt()}
	cpu_state_container.Add(cpu_state_label)
	cpu_state_container.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu_state.cycles, "cpu cycle: %d")))
	cpu_state_container.Hide()
	cpu_state_visibility := fyne.NewMenuItem("CPU state", func() {
		if cpu_state_container.Hidden {
			cpu_state_container.Refresh()
			cpu_state_container.Show()
		} else {
			cpu_state_container.Hide()
		}
	})
	// ============= Debugger: CPU State =============

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
	content := container.New(layout.NewHBoxLayout(), cpu_state_container, layout.NewSpacer(), display, layout.NewSpacer(), vram)

	vram.Hide()
	vram_visibility := fyne.NewMenuItem("VRAM viewer", func() {
		if vram.Hidden {
			vram.Refresh()
			vram.Show()
		} else {
			vram.Hide()
		}
	})

	debug_menu := fyne.NewMenu("Debug", cpu_state_visibility, vram_visibility)
	main_menu := fyne.NewMainMenu(debug_menu)
	w.SetMainMenu(main_menu)
	w.SetContent(content)

	ui := &Interface{app: a, window: w, display: display, vram: vram, cpu: cpu_state_container, cpu_state: cpu_state, emu: e}

	return ui
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
	ui.emu.rom_loaded = true

}

func (ui *Interface) SetCPUState() {
	current_state := ui.emu.GetCPUState()
	ui.cpu_state.cycles.Set(int(current_state.cycles))
}

func (ui *Interface) Run() {
	go func() {
		frame_time := 16 * time.Millisecond // for 60 fps
		for range time.NewTicker(frame_time).C {
			fyne.DoAndWait(func() {

				frame_ready := false
				max_render_time := (456 /* dots */ * 153 /* lines */ / 4 /* cpu cyc */)
				for _ = range max_render_time {
					frame_ready = ui.emu.Run()
					if frame_ready {
						break
					}
				}
				if frame_ready {
					ui.display.Refresh()
				}

				if ui.cpu.Visible() {
					ui.SetCPUState()
				}
			})
		}
	}()
	ui.window.ShowAndRun()

}
