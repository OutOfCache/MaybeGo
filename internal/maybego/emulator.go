package maybego

type Emulator struct {
	cpu        *CPU
	ppu        *PPU
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
	e := &Emulator{cpu: cpu, ppu: ppu, logger: logger}
	InitMemory()

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
	return emu.ppu.Render(cycles)
}
