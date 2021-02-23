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

func FlagToBit(flag bool) byte {
	if flag {
		return 1
	}
	return 0
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

func (cpu *CPU) cpu03() int { // INC BC
	cpu.reg.C++
	if cpu.reg.C == 0 {
		cpu.reg.B++
	}
	cpu.reg.PC++

	return 0
}

func (cpu *CPU) cpu04() int { // INC B
	cpu.flg.N = false
	cpu.reg.B++
	if cpu.reg.B == 0 {
		cpu.flg.Z = true
	}
	if cpu.reg.B == 0x10 {
		cpu.flg.H = true
	}

	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu05() int { // DEC B
	cpu.flg.N = true
	cpu.reg.B--
	if cpu.reg.B == 0 {
		cpu.flg.Z = true
	}
	if cpu.reg.B == 0xF {
		cpu.flg.H = true
	}

	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu06() int { // LD B, u8
	cpu.reg.B = Read(cpu.reg.PC + 1)

	cpu.reg.PC += 2
	return 0
}

func (cpu *CPU) cpu07() int { // RLCA
	if (cpu.reg.A & 0x80) == 0x80 {
		cpu.flg.C = true
	} else {
		cpu.flg.C = false
	}

	cpu.reg.A <<= 1
	cpu.reg.A += FlagToBit(cpu.flg.C)

	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu08() int { // LD (u16),SP
	address := uint16(Read(cpu.reg.PC+1)) + (uint16(Read(cpu.reg.PC+2)) << 8)
	Write(address, byte(cpu.reg.SP&0xFF))
	Write(address+1, byte((cpu.reg.SP&0xFF00)>>8))

	cpu.reg.PC += 3
	return 0
}

func (cpu *CPU) cpu09() int { // ADD HL, BC
	cpu.flg.N = false
	sum := int(cpu.reg.L) + int(cpu.reg.C)
	cpu.reg.L = byte(sum & 0xFF)
	cpu.flg.H = (byte(sum>>8)+(cpu.reg.H&0xf)+(cpu.reg.B&0xf))&0x10 == 0x10
	sum = (sum >> 8) + int(cpu.reg.H) + int(cpu.reg.B)
	cpu.reg.H = byte(sum & 0xFF)
	cpu.flg.C = sum > 0xFF

	cpu.reg.PC++

	return 0
}

func (cpu *CPU) cpu0A() int { // LD A, (BC)
	address := uint16(cpu.reg.B)<<8 + uint16(cpu.reg.C)
	cpu.reg.A = Read(address)

	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu0B() int { // DEC BC
	if cpu.reg.C == 0 && cpu.reg.B > 0 {
		cpu.reg.C = 0xFF
		cpu.reg.B--
	} else {
		cpu.reg.C--
	}

	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu0C() int { // INC C
	cpu.flg.N = false
	cpu.reg.C++
	if cpu.reg.C == 0 {
		cpu.flg.Z = true
	}
	if cpu.reg.C == 0x10 {
		cpu.flg.H = true
	}

	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu0D() int { // DEC C
	cpu.flg.N = true
	cpu.reg.C--
	if cpu.reg.C == 0 {
		cpu.flg.Z = true
	}
	if cpu.reg.C == 0xF {
		cpu.flg.H = true
	}

	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu0E() int {
	cpu.reg.PC++
	cpu.reg.C = Read(cpu.reg.PC)

	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu0F() int {
	cpu.flg.Z = false
	cpu.flg.N = false
	cpu.flg.H = false
	cpu.flg.C = cpu.reg.A&0x01 == 0x01
	cpu.reg.A >>= 1
	cpu.reg.A += FlagToBit(cpu.flg.C) << 7

	cpu.reg.PC++
	return 0
}
