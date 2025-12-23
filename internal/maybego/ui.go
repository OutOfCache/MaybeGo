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

type cpuStateWindow struct {
	container *fyne.Container
	state     *cpu_state_bindings
}

type debugView struct {
	cpu_win    *cpuStateWindow
	disasm_win *disasmWindow
	halt       bool
	step       bool
}

type Interface struct {
	app     fyne.App
	window  fyne.Window
	display *canvas.Raster
	vram    *fyne.Container
	emu     *Emulator
	debug   *debugView
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
	cpu := createCpuStateWindow()
	cpu_state_visibility := fyne.NewMenuItem("CPU state", func() {
		if cpu.container.Hidden {
			cpu.container.Refresh()
			cpu.container.Show()
		} else {
			cpu.container.Hide()
		}
	})
	// ============= Debugger: CPU State =============
	// ============= Debugger: Disassembler =============
	debug := createDebugView(cpu)
	debug_container := createDebugContainer(e, display, debug)

	var debug_visibility *fyne.MenuItem
	debug_visibility = fyne.NewMenuItem("Debugger", func() {
		if debug_container.Hidden {
			debug_container.Refresh()
			debug_container.Show()

			cpu.container.Show()
		} else {
			debug_container.Hide()
			cpu.container.Hide()
		}
		debug_visibility.Checked = debug_container.Visible()
	})
	debug_visibility.Checked = debug_container.Visible()
	// ============= Debugger: Disassembler =============

	vram := createVramView()
	// TODO: scaling factor
	display.SetMinSize(fyne.NewSize(160, 144))
	content := container.New(layout.NewHBoxLayout(), debug_container, layout.NewSpacer(), cpu.container, layout.NewSpacer(), display, layout.NewSpacer(), vram)

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

	ui := &Interface{app: a, window: w, display: display, vram: vram, emu: e, debug: debug}
	ui.debug.disasm_win.ExtendBaseWidget(debug.disasm_win)

	return ui
}

func (ui *Interface) LoadRom(rom *[]byte) {
	for i, buffer := range *rom {
		Write(uint16(i), buffer)
	}

	ui.debug.disasm_win.disasm.SetFile(rom)

	go func() {
		ui.debug.disasm_win.disasm.Disassemble()

		for _, line := range ui.debug.disasm_win.disasm.lines {
			ui.debug.disasm_win.Append(fmt.Sprintf("%04X|\t%s", line.offset, line.disasm))
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
	ui.debug.cpu_win.state.cycles.Set(int(current_state.cycles))
	ui.debug.cpu_win.state.registers.a.Set(int(current_state.registers.A))
	ui.debug.cpu_win.state.registers.b.Set(int(current_state.registers.B))
	ui.debug.cpu_win.state.registers.c.Set(int(current_state.registers.C))
	ui.debug.cpu_win.state.registers.d.Set(int(current_state.registers.D))
	ui.debug.cpu_win.state.registers.e.Set(int(current_state.registers.E))
	ui.debug.cpu_win.state.registers.h.Set(int(current_state.registers.H))
	ui.debug.cpu_win.state.registers.l.Set(int(current_state.registers.L))
	ui.debug.cpu_win.state.registers.pc.Set(int(current_state.registers.PC))

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

	ui.debug.cpu_win.state.flagstring.Set(renderFlags())
}

func (ui *Interface) Run() {
	go func() {
		frame_time := 16 * time.Millisecond // for 60 fps
		for range time.NewTicker(frame_time).C {
			if ui.debug.halt {
				continue
			}
			fyne.DoAndWait(func() {

				frame_ready := false
				max_render_time := (456 /* dots */ * 153 /* lines */ / 4 /* cpu cyc */)
				if ui.debug.step {
					max_render_time = 1
					ui.debug.halt = true
				}
				for _ = range max_render_time {
					frame_ready = ui.emu.Run()
					if frame_ready {
						break
					}
					next_pc := ui.emu.GetCPUState().registers.PC
					if slices.Contains(ui.debug.disasm_win.breakpoints, uint(next_pc)) {
						ui.debug.halt = true
						break
					}
				}
				if frame_ready {
					ui.display.Refresh()
				}

				if ui.debug.cpu_win.container.Visible() {
					ui.SetCPUState()
				}
			})
		}
	}()
	ui.window.ShowAndRun()

}

func generateVramTile(tileID int, scale int) func(x, y, w, h int) color.Color {
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

func createVramView() *fyne.Container {
	vram := container.New(layout.NewGridLayout(16))
	scale := 2
	tile_size := float32(scale * 8)
	for i := 0; i < 384; i++ {
		tile := canvas.NewRasterWithPixels(generateVramTile(i, scale))
		tile.SetMinSize(fyne.NewSize(tile_size, tile_size))
		vram.Add(tile)
	}

	return vram
}

func createCpuStateWindow() *cpuStateWindow {
	cpu := &cpuStateWindow{
		container: container.New(layout.NewVBoxLayout()),
		state:     &cpu_state_bindings{cycles: binding.NewInt()},
	}
	cpu_state_label := widget.NewLabel("CPU State")
	cpu_state_label.TextStyle.Bold = true
	cpu.container.Add(cpu_state_label)
	cpu.container.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu.state.cycles, "cpu cycle: %d")))

	register_container := container.NewHBox()
	register_container_lo := container.NewVBox()
	register_container_hi := container.NewVBox()
	cpu.state.registers.a = binding.NewInt()
	cpu.state.registers.b = binding.NewInt()
	cpu.state.registers.c = binding.NewInt()
	cpu.state.registers.d = binding.NewInt()
	cpu.state.registers.e = binding.NewInt()
	cpu.state.registers.h = binding.NewInt()
	cpu.state.registers.l = binding.NewInt()
	cpu.state.registers.pc = binding.NewInt()
	cpu.state.registers.sp = binding.NewInt()
	registers_label := widget.NewLabel("Registers")
	registers_label.TextStyle.Bold = true
	cpu.container.Add(registers_label)

	register_container_lo.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu.state.registers.a, "A: %X")))
	register_container_hi.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu.state.registers.pc, "PC: %X")))
	register_container_lo.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu.state.registers.b, "B: %X")))
	register_container_hi.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu.state.registers.c, "C: %X")))
	register_container_lo.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu.state.registers.d, "D: %X")))
	register_container_hi.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu.state.registers.e, "E: %X")))
	register_container_lo.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu.state.registers.h, "H: %X")))
	register_container_hi.Add(widget.NewLabelWithData(binding.IntToStringWithFormat(cpu.state.registers.l, "L: %X")))
	register_container.Add(register_container_lo)
	register_container.Add(register_container_hi)
	register_container.Resize(fyne.NewSize(160, 160))
	cpu.container.Add(register_container)

	cpu.state.flagstring = binding.NewString()

	flag_container := container.NewVBox()
	flag_label := widget.NewLabel("Flags")
	flag_label.TextStyle.Bold = true
	flagstring_label := widget.NewLabelWithData(cpu.state.flagstring)
	flagstring_label.TextStyle.Monospace = true
	flag_container.Add(flag_label)
	flag_container.Add(flagstring_label)
	cpu.container.Add(flag_container)

	return cpu
}

func createDisasmView() *disasmWindow {
	disasm_container := &disasmWindow{
		TextGrid: &widget.TextGrid{},
		disasm:   NewDisasm(),
	}
	disasm_container.Scroll = fyne.ScrollVerticalOnly
	return disasm_container
}

func createDebugView(cpu *cpuStateWindow) *debugView {
	disasm_container := createDisasmView()
	debug := &debugView{
		disasm_win: disasm_container,
		cpu_win:    cpu,
		halt:       false,
		step:       false,
	}

	return debug
}

func createDebugContainer(emu *Emulator, display *canvas.Raster, debug *debugView) *fyne.Container {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.MediaPauseIcon(), func() {
			debug.halt = true
		}),
		widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
			debug.halt = false
			debug.step = true
		}),
		widget.NewToolbarAction(theme.MediaFastForwardIcon(), func() {
			debug.halt = false
			debug.step = false
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.MediaReplayIcon(), func() {
			emu.Reset()
			debug.step = false
			display.Refresh()
		}),
	)

	return container.NewBorder(toolbar, nil, nil, nil, debug.disasm_win)
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
