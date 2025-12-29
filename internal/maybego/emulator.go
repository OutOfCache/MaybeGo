package maybego

type Emulator struct {
	cpu        *CPU
	ppu        *PPU
	joypad     *Joypad
	rom_loaded bool
	logger     *Logger
}

type cpu_state struct {
	cycles    uint
	registers *Registers
	flags     *Flags
}

func NewEmulator(logger *Logger) *Emulator {
	// TODO: no logger in CPU or PPU
	cpu := NewCPU(logger)
	ppu := NewPPU(logger)
	InitMemory()
	joy := NewJoypad()
	e := &Emulator{cpu: cpu, ppu: ppu, joypad: joy, logger: logger}

	return e
}

func (emu *Emulator) GetPPU() *PPU {
	return emu.ppu
}

func (emu *Emulator) GetCPU() *CPU {
	return emu.cpu
}

func (emu *Emulator) GetCPUState() cpu_state {
	return cpu_state{cycles: emu.cpu.clk.cycles, registers: emu.cpu.reg, flags: emu.cpu.flg}
}

// TODO: for loading roms during runtime
func (emu *Emulator) Reset() {
	emu.cpu.Reset()
	emu.ppu.Reset()
}

func (emu *Emulator) FetchDecodeExec() byte {
	emu.cpu.Fetch()
	cycles := emu.cpu.Decode()

	emu.cpu.Handle_timer(cycles)
	return cycles
}

func (emu *Emulator) Run() bool {
	if !emu.rom_loaded {
		return false
	}

	cycles := emu.FetchDecodeExec()
	emu.joypad.updateControls()
	return emu.ppu.Render(cycles)
}

func (emu *Emulator) PressButton(key string) {
	// keycode := byte(0)
	switch key {
	case "V":
		emu.joypad.setButton(ButtonA)
		return
	case "C":
		emu.joypad.setButton(ButtonB)
		return
	case "X":
		emu.joypad.setButton(ButtonSelect)
		return
	case "Z":
		emu.joypad.setButton(ButtonStart)
		return
	}

}

func (emu *Emulator) ReleaseButton(key string) {
	// keycode := byte(0)
	switch key {
	case "V":
		emu.joypad.resetButton(ButtonA)
		return
	case "C":
		emu.joypad.resetButton(ButtonB)
		return
	case "X":
		emu.joypad.resetButton(ButtonSelect)
		return
	case "Z":
		emu.joypad.resetButton(ButtonStart)
		return
	}

}
