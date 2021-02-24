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

	return cpu
}

func FlagToBit(flag bool) byte {
	if flag {
		return 1
	}
	return 0
}

// LD r8, r8/n8
func (cpu *CPU) ld8(dest *byte, src byte) {
	*dest = src
}

// LD r16, r16/n16
func ld16(destLo *byte, destHi *byte, srcLo byte, srcHi byte) {
	*destLo = srcLo
	*destHi = srcHi
}

// LD r16, r16/n16
func (cpu *CPU) ld16reg(dest uint16, srcLo byte, srcHi byte) {
	dest = uint16(srcHi)<<8 + uint16(srcLo)
}

// LD [r16], r8/n8
func (cpu *CPU) ldToAddress(adrLo byte, adrHi byte, val byte) {
	address := uint16(adrHi)<<8 + uint16(adrLo)
	Write(address, val)
}

// LD [r16], r16
func (cpu *CPU) ldToAddress16(adrLo byte, adrHi byte, valLo byte, valHi byte) {
	address := uint16(adrHi)<<8 + uint16(adrLo)
	Write(address, valLo)
	Write(address+1, valHi)

}

// LD r8, [r16]
func (cpu *CPU) ldFromAddress(dest *byte, adrLo byte, adrHi byte) {
	address := uint16(adrHi)<<8 + uint16(adrLo)
	*dest = Read(address)
}

func (cpu *CPU) inc8(reg *byte, flags bool) {
	*reg++
	if flags {
		cpu.flg.N = false
		cpu.flg.Z = *reg == 0
		cpu.flg.H = *reg == 0x10
	}
}

func (cpu *CPU) inc16(destLo *byte, destHi *byte) {
	cpu.inc8(destLo, false)
	if *destLo == 0 { // increase if overflow in low byte
		cpu.inc8(destHi, false)
	}
}

func (cpu *CPU) dec8(reg *byte, flags bool) {
	*reg--
	if flags {
		cpu.flg.N = true
		cpu.flg.Z = *reg == 0
		cpu.flg.H = *reg == 0xF
	}
}

func (cpu *CPU) dec16(destLo *byte, destHi *byte) {
	cpu.dec8(destLo, false)
	if *destLo == 0xFF {
		cpu.dec8(destHi, false)
	}
}

func (cpu *CPU) add16(destLo *byte, destHi *byte, srcLo byte, srcHi byte) {
	cpu.flg.N = false

	sum := int(*destLo) + int(srcLo)
	*destLo = byte(sum & 0xFF)
	cpu.flg.H = (byte(sum>>8)+(*destHi&0xf)+(srcHi&0xf))&0x10 == 0x10
	sum = (sum >> 8) + int(*destHi) + int(srcHi)
	*destHi = byte(sum & 0xFF)
	cpu.flg.C = sum > 0xFF

}

func (cpu *CPU) rl8(reg *byte, carry bool) {
	lsb := byte(0)
	if carry {
		lsb = FlagToBit(cpu.flg.C)
	} else {
		lsb = *reg >> 7
	}
	cpu.flg.C = *reg>>7 == 1

	*reg = *reg<<1 + lsb

	cpu.flg.N = false
	cpu.flg.H = false
	cpu.flg.Z = *reg == 0
}

func (cpu *CPU) rr8(reg *byte, carry bool) {
	msb := byte(0)
	if carry {
		msb = FlagToBit(cpu.flg.C)
	} else {
		msb = *reg & 0x01
	}
	cpu.flg.C = *reg&0x01 == 1

	*reg = *reg>>1 + (msb << 7)

	cpu.flg.N = false
	cpu.flg.H = false
	cpu.flg.Z = *reg == 0
}

func (cpu *CPU) cpu00() { // do I need parameters for args?
	cpu.reg.PC++
}

func (cpu *CPU) cpu01() int { // LD BC, u16
	ld16(&cpu.reg.C, &cpu.reg.B, Read(cpu.reg.PC+1), Read(cpu.reg.PC+2))
	cpu.reg.PC += 3

	return 0
}

func (cpu *CPU) cpu02() int { // LD (BC), A
	cpu.ldToAddress(cpu.reg.C, cpu.reg.B, cpu.reg.A)
	cpu.reg.PC++

	return 0
}

func (cpu *CPU) cpu03() int { // INC BC
	cpu.inc16(&cpu.reg.C, &cpu.reg.B)
	cpu.reg.PC++

	return 0
}

func (cpu *CPU) cpu04() int { // INC B
	cpu.inc8(&cpu.reg.B, true)
	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu05() int { // DEC B
	cpu.dec8(&cpu.reg.B, true)
	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu06() int { // LD B, u8
	cpu.ld8(&cpu.reg.B, Read(cpu.reg.PC+1))

	cpu.reg.PC += 2
	return 0
}

func (cpu *CPU) cpu07() int { // RLCA
	cpu.rl8(&cpu.reg.A, false)

	cpu.flg.Z = false

	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu08() int { // LD (u16),SP
	cpu.ldToAddress16(Read(cpu.reg.PC+1), Read(cpu.reg.PC+2),
		byte(cpu.reg.SP&0xFF), byte(cpu.reg.SP>>8))

	cpu.reg.PC += 3
	return 0
}

func (cpu *CPU) cpu09() int { // ADD HL, BC
	cpu.add16(&cpu.reg.L, &cpu.reg.H, cpu.reg.C, cpu.reg.B)

	cpu.reg.PC++

	return 0
}

func (cpu *CPU) cpu0A() int { // LD A, (BC)
	cpu.ldFromAddress(&cpu.reg.A, cpu.reg.C, cpu.reg.B)
	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu0B() int { // DEC BC
	cpu.dec16(&cpu.reg.C, &cpu.reg.B)

	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu0C() int { // INC C
	cpu.inc8(&cpu.reg.C, true)

	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu0D() int { // DEC C
	cpu.dec8(&cpu.reg.C, true)

	cpu.reg.PC++
	return 0
}

func (cpu *CPU) cpu0E() int { // LD C, u8
	cpu.ld8(&cpu.reg.C, Read(cpu.reg.PC+1))

	cpu.reg.PC += 2
	return 0
}

func (cpu *CPU) cpu0F() int { // RRCA
	cpu.rr8(&cpu.reg.A, false)
	cpu.flg.Z = false

	cpu.reg.PC++
	return 0
}
