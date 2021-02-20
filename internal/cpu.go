package maybego

type Registers struct {
	A  byte // can be combined to AF
	B  byte // BC, B hi
	C  byte // BC, C lo
	D  byte // DE, D hi
	E  byte // DE, E lo
	H  byte // HL, H hi
	L  byte // HL, L lo
	SP uint16
	PC uint16
}

type Flags struct {
	Z bool // zero flag
	C bool // carry flag
	N bool // sub flag
	H bool // half carry
}

type CPU struct {
	reg *Registers
	flg *Flags
}

var currentOpcode uint16
var opcodes [256]func()

// dummy "constructor"
func NewCPU() *CPU {
	cpu := &CPU{reg: new(Registers), flg: new(Flags)}
	//    cpu.Registers := &Registers{}
	//    cpu.Flags := &Flags{}

	return cpu
}

func (cpu *CPU) cpu00() { // do I need parameters for args?
	cpu.reg.PC++
}

func (cpu *CPU) cpu01() int { // LD BC, u16
	cpu.reg.PC++
	cpu.reg.C = Read(cpu.reg.PC)
	cpu.reg.PC++
	cpu.reg.B = Read(cpu.reg.PC)
	cpu.reg.PC++

	return 0
}

func (cpu *CPU) cpu02() int { // LD (BC), A
	address := (uint16(cpu.reg.B) << 8) + uint16(cpu.reg.C)
	Write(address, cpu.reg.A)
	cpu.reg.PC++

	return 0
}
