package maybego

import (
	"fmt"
	"image/color"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
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
	cycles    binding.Int
	registers struct {
		a  binding.Int
		b  binding.Int
		c  binding.Int
		d  binding.Int
		e  binding.Int
		h  binding.Int
		l  binding.Int
		sp binding.Int
		pc binding.Int
	}
	flagstring binding.String
}

type disasmWindow struct {
	*widget.TextGrid
	disasm      *Disasm
	breakpoints []uint
}

type debugView struct {
	disasm_win *disasmWindow
	halt       bool
	step       bool
}

type Interface struct {
	app        fyne.App
	window     fyne.Window
	display    *canvas.Raster
	vram       *fyne.Container
	cpu        *fyne.Container
	cpu_state  *cpu_state_bindings
	emu        *Emulator
	debug_view *debugView
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
			if x > 159 || y > 143 {
				return color.RGBA{R: 0, G: 0, B: 0, A: 0}
			}
			return Palette[e.ppu.GetCurrentFrame()[(160*y)+x]]
		})

	// ============= Debugger: CPU State =============
	cpu_state_container := container.New(layout.NewVBoxLayout())
	cpu_state_label := widget.NewLabel("CPU State")
	cpu_state_label.TextStyle.Bold = true
	cpu_state := &cpu_state_bindings{cycles: binding.NewInt()}
	cpu_state_container.Add(cpu_state_label)
	register_container := container.NewHBox()
	register_container_lo := container.NewVBox()
	register_container_hi := container.NewVBox()
	cpu_state_container.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu_state.cycles, "cpu cycle: %d")))
	cpu_state.registers.a = binding.NewInt()
	cpu_state.registers.b = binding.NewInt()
	cpu_state.registers.c = binding.NewInt()
	cpu_state.registers.d = binding.NewInt()
	cpu_state.registers.e = binding.NewInt()
	cpu_state.registers.h = binding.NewInt()
	cpu_state.registers.l = binding.NewInt()
	cpu_state.registers.pc = binding.NewInt()
	cpu_state.registers.sp = binding.NewInt()
	registers_label := widget.NewLabel("Registers")
	registers_label.TextStyle.Bold = true
	cpu_state_container.Add(registers_label)
	// registers := binding.BindIntList(
	// 	&[]int{},
	// )
	// list := widget.NewListWithData(registers,
	// 	func() fyne.CanvasObject {
	// 		return canvas.NewText("template")
	// 	},
	// 	func(i binding.DataItem, o fyne.CanvasObject) {
	// 		o.(*canvas.Text).Text = binding.IntToStringWithFormat(i.(binding.Int), "A: %X")
	// 	})

	// cpu_state_container.Add(widget.NewRichTextWithText(binding.IntToStringWithFormat(cpu_state.registers.a, "A: %X")))
	register_container_lo.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu_state.registers.a, "A: %X")))
	register_container_hi.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu_state.registers.pc, "PC: %X")))
	register_container_lo.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu_state.registers.b, "B: %X")))
	register_container_hi.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu_state.registers.c, "C: %X")))
	register_container_lo.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu_state.registers.d, "D: %X")))
	register_container_hi.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu_state.registers.e, "E: %X")))
	register_container_lo.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu_state.registers.h, "H: %X")))
	register_container_hi.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu_state.registers.l, "L: %X")))
	register_container.Add(register_container_lo)
	register_container.Add(register_container_hi)
	register_container.Resize(fyne.NewSize(160, 160))
	cpu_state_container.Add(register_container)

	cpu_state.flagstring = binding.NewString()

	flag_container := container.NewVBox()
	flag_label := widget.NewLabel("Flags")
	flag_label.TextStyle.Bold = true
	flagstring_label := widget.NewLabelWithData(cpu_state.flagstring)
	flagstring_label.TextStyle.Monospace = true
	flag_container.Add(flag_label)
	flag_container.Add(flagstring_label)
	cpu_state_container.Add(flag_container)

	// cpu_state_container.Hide()
	cpu_state_visibility := fyne.NewMenuItem("CPU state", func() {
		if cpu_state_container.Hidden {
			cpu_state_container.Refresh()
			cpu_state_container.Show()
		} else {
			cpu_state_container.Hide()
		}
	})
	// ============= Debugger: CPU State =============
	// ============= Debugger: Disassembler =============
	disasm_container := &disasmWindow{
		TextGrid: &widget.TextGrid{},
		disasm:   NewDisasm(),
	}
	disasm_container.Scroll = fyne.ScrollVerticalOnly

	debug_view := &debugView{
		disasm_win: disasm_container,
		halt:       false,
	}
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.MediaPauseIcon(), func() {
			debug_view.halt = true
			fmt.Printf("halt: %t\n", debug_view.halt)
		}),
		widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
			debug_view.halt = false
			debug_view.step = true
		}),
		widget.NewToolbarAction(theme.MediaFastForwardIcon(), func() {
			debug_view.halt = false
			debug_view.step = false
			fmt.Printf("halt: %t\n", debug_view.halt)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.MediaReplayIcon(), func() {
			e.Reset()
			debug_view.halt = false
			debug_view.step = false
		}),
	)

	disasm_content := container.NewBorder(toolbar, nil, nil, nil, disasm_container)

	// ============= Debugger: Disassembler =============

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
	content := container.New(layout.NewHBoxLayout(), disasm_content, layout.NewSpacer(), cpu_state_container, layout.NewSpacer(), display, layout.NewSpacer(), vram)

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

	ui := &Interface{app: a, window: w, display: display, vram: vram, cpu: cpu_state_container, cpu_state: cpu_state, emu: e, debug_view: debug_view}
	ui.debug_view.disasm_win.ExtendBaseWidget(disasm_container)

	return ui
}

func (ui *Interface) LoadRom(rom *[]byte) {
	for i, buffer := range *rom {
		Write(uint16(i), buffer)
	}

	ui.debug_view.disasm_win.disasm.SetFile(rom)

	go func() {
		ui.debug_view.disasm_win.disasm.Disassemble()

		for _, line := range ui.debug_view.disasm_win.disasm.lines {
			ui.debug_view.disasm_win.Append(fmt.Sprintf("%04X|\t%s", line.offset, line.disasm))
		}
	}()

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
	ui.cpu_state.registers.a.Set(int(current_state.registers.A))
	ui.cpu_state.registers.b.Set(int(current_state.registers.B))
	ui.cpu_state.registers.c.Set(int(current_state.registers.C))
	ui.cpu_state.registers.d.Set(int(current_state.registers.D))
	ui.cpu_state.registers.e.Set(int(current_state.registers.E))
	ui.cpu_state.registers.h.Set(int(current_state.registers.H))
	ui.cpu_state.registers.l.Set(int(current_state.registers.L))
	ui.cpu_state.registers.pc.Set(int(current_state.registers.PC))

	renderFlags := func() string {
		flags := ""

		if current_state.flags.C {
			flags += "C"
		} else {
			flags += "-"
		}
		if current_state.flags.H {
			flags += "H"
		} else {
			flags += "-"
		}
		if current_state.flags.N {
			flags += "N"
		} else {
			flags += "-"
		}
		if current_state.flags.Z {
			flags += "Z"
		} else {
			flags += "-"
		}

		flags += " "

		if current_state.flags.IME {
			flags += "IME"
		} else {
			flags += "---"
		}

		flags += " "

		if current_state.flags.HALT {
			flags += "HALT"
		} else {
			flags += "----"
		}

		return flags
	}

	ui.cpu_state.flagstring.Set(renderFlags())
}

func (ui *Interface) Run() {
	go func() {
		frame_time := 16 * time.Millisecond // for 60 fps
		for range time.NewTicker(frame_time).C {
			if ui.debug_view.halt {
				continue
			}
			fyne.DoAndWait(func() {

				frame_ready := false
				max_render_time := (456 /* dots */ * 153 /* lines */ / 4 /* cpu cyc */)
				if ui.debug_view.step {
					max_render_time = 1
					ui.debug_view.halt = true
				}
				for _ = range max_render_time {
					frame_ready = ui.emu.Run()
					if frame_ready {
						break
					}
					next_pc := ui.emu.GetCPUState().registers.PC
					if slices.Contains(ui.debug_view.disasm_win.breakpoints, uint(next_pc)) {
						fmt.Printf("halted. next_pc: %X", next_pc)
						ui.debug_view.halt = true
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

func (dw *disasmWindow) Tapped(ev *fyne.PointEvent) {
	xpos, _ := dw.CursorLocationForPosition(ev.Position)

	selectedStyle := widget.CustomTextGridStyle{}
	selectedStyle.BGColor = theme.Color(theme.ColorNameFocus)
	// TODO visual indication that it is selected instantly
	dw.breakpoints = append(dw.breakpoints, dw.disasm.lines[xpos].offset)
	dw.SetRowStyle(xpos, &selectedStyle)
	dw.BaseWidget.Refresh()
}
