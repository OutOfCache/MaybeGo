package maybego

import "fmt"

type Signals struct {
	exec_cy    chan int
	exec_ready chan bool
	frame_done chan bool
}

type Emulator struct {
	cpu        *CPU
	ppu        *PPU
	rom_loaded bool
	logger     *Logger
	signals    *Signals
}

func NewEmulator(logger *Logger) *Emulator {
	// TODO: no logger in CPU or PPU
	cpu := NewCPU(logger)
	ppu := NewPPU(logger)
	// signals := &Signals{exec_cy: make(chan int, 1)}
	// signals := &Signals{exec_ready: make(chan bool, 1)}
	signals := &Signals{}
	e := &Emulator{cpu: cpu, ppu: ppu, logger: logger, signals: signals}

	return e
}

func (emu *Emulator) GetPPU() *PPU {
	return emu.ppu
}

func (emu *Emulator) GetCPU() *CPU {
	return emu.cpu
}

// TODO: for loading roms during runtime
func (emu *Emulator) Reset() {
	// cpu.Reset()
	// ppu.Reset()
}

func (emu *Emulator) FetchDecodeExec() bool {
	if !emu.rom_loaded {
		return false
	}
	emu.cpu.Fetch()
	cycles := emu.cpu.Decode()

	emu.cpu.Handle_timer(cycles)
	fmt.Println("after handle timer")
	frame_ready := emu.ppu.Render(cycles)
	fmt.Println("after render: ", frame_ready)

	// FIX: belongs in UI
	// for blarggs tests:
	if Read(0xff02) == 0x81 {
		c := Read(0xff01)
		fmt.Printf("%c", c)
		Write(0xff02, 0)
	}
	emu.signals.exec_cy <- int(cycles)
	fmt.Println("exec_done")
	// emu.signals.exec_ready <- true
	// emu.signals.frame_done <- frame_ready
	// fmt.Println("frame ready")
	return frame_ready
}

func (emu *Emulator) Run() {}
