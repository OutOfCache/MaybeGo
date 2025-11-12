package maybego

import "fmt"

type Emulator struct {
	cpu    *CPU
	ppu    *PPU
	logger *Logger
	rom_loaded bool
}

func NewEmulator(logger *Logger) *Emulator {
	// TODO: no logger in CPU or PPU
	cpu := NewCPU(logger)
	ppu := NewPPU(logger)
	e := &Emulator{cpu: cpu, ppu: ppu, logger: logger}

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
	frame_ready := emu.ppu.Render(cycles)

	// FIX: belongs in UI
	// for blarggs tests:
	if Read(0xff02) == 0x81 {
		c := Read(0xff01)
		fmt.Printf("%c", c)
		Write(0xff02, 0)
	}
	return frame_ready
}
