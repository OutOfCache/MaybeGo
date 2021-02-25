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
func (cpu *CPU) ld16(destLo *byte, destHi *byte, srcLo byte, srcHi byte) {
	*destLo = srcLo
	*destHi = srcHi
}

// LD r16, r16/n16
func (cpu *CPU) ld16reg(dest *uint16, srcLo byte, srcHi byte) {
	*dest = uint16(srcHi)<<8 + uint16(srcLo)
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

func (cpu *CPU) addA(reg byte, carry bool) {
	if carry {
		carry = cpu.flg.C
	}
	cpu.flg.H = uint16(cpu.reg.A&0xF)+uint16(reg&0xF)&0x10 == 0x10

	sum := uint16(cpu.reg.A) + uint16(reg) + uint16(FlagToBit(carry))
	cpu.reg.A = byte(sum)

	cpu.flg.Z = cpu.reg.A == 0
	cpu.flg.C = sum&0x10 == 0x10
	cpu.flg.N = false
}

func (cpu *CPU) subA(reg byte, carry bool) {
	if carry {
		carry = cpu.flg.C
	}
	cpu.flg.H = cpu.reg.A&0xF < reg&0xF
	cpu.flg.C = cpu.reg.A < reg

	cpu.reg.A -= reg + FlagToBit(carry)

	cpu.flg.Z = cpu.reg.A == 0
	cpu.flg.N = true
}

func (cpu *CPU) andA(reg byte) {
	cpu.reg.A &= reg

	cpu.flg.Z = cpu.reg.A == 0
	cpu.flg.N = false
	cpu.flg.H = true
	cpu.flg.C = false
}

func (cpu *CPU) xorA(reg byte) {
	cpu.reg.A ^= reg

	cpu.flg.Z = cpu.reg.A == 0
	cpu.flg.N = false
	cpu.flg.H = false
	cpu.flg.C = false
}

func (cpu *CPU) orA(reg byte) {
	cpu.reg.A |= reg

	cpu.flg.Z = cpu.reg.A == 0
	cpu.flg.N = false
	cpu.flg.H = false
	cpu.flg.C = false
}

func (cpu *CPU) cpA(reg byte) {
	cpu.flg.H = cpu.reg.A&0xF < reg&0xF
	cpu.flg.C = cpu.reg.A < reg

	result := cpu.reg.A - reg

	cpu.flg.Z = result == 0
	cpu.flg.N = true
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

func (cpu *CPU) jr(flag bool) int {
	if flag {
		cpu.reg.PC += uint16(2 + int8(Read(cpu.reg.PC+1)))
		return 3
	}
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpu00() int { // do I need parameters for args?
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu01() int { // LD BC, u16
	cpu.ld16(&cpu.reg.C, &cpu.reg.B, Read(cpu.reg.PC+1), Read(cpu.reg.PC+2))
	cpu.reg.PC += 3

	return 3
}

func (cpu *CPU) cpu02() int { // LD (BC), A
	cpu.ldToAddress(cpu.reg.C, cpu.reg.B, cpu.reg.A)
	cpu.reg.PC++

	return 2
}

func (cpu *CPU) cpu03() int { // INC BC
	cpu.inc16(&cpu.reg.C, &cpu.reg.B)
	cpu.reg.PC++

	return 2
}

func (cpu *CPU) cpu04() int { // INC B
	cpu.inc8(&cpu.reg.B, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu05() int { // DEC B
	cpu.dec8(&cpu.reg.B, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu06() int { // LD B, u8
	cpu.ld8(&cpu.reg.B, Read(cpu.reg.PC+1))

	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpu07() int { // RLCA
	cpu.rl8(&cpu.reg.A, false)

	cpu.flg.Z = false

	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu08() int { // LD (u16),SP
	cpu.ldToAddress16(Read(cpu.reg.PC+1), Read(cpu.reg.PC+2),
		byte(cpu.reg.SP&0xFF), byte(cpu.reg.SP>>8))

	cpu.reg.PC += 3
	return 5
}

func (cpu *CPU) cpu09() int { // ADD HL, BC
	cpu.add16(&cpu.reg.L, &cpu.reg.H, cpu.reg.C, cpu.reg.B)

	cpu.reg.PC++

	return 2
}

func (cpu *CPU) cpu0A() int { // LD A, (BC)
	cpu.ldFromAddress(&cpu.reg.A, cpu.reg.C, cpu.reg.B)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu0B() int { // DEC BC
	cpu.dec16(&cpu.reg.C, &cpu.reg.B)

	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu0C() int { // INC C
	cpu.inc8(&cpu.reg.C, true)

	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu0D() int { // DEC C
	cpu.dec8(&cpu.reg.C, true)

	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu0E() int { // LD C, u8
	cpu.ld8(&cpu.reg.C, Read(cpu.reg.PC+1))

	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpu0F() int { // RRCA
	cpu.rr8(&cpu.reg.A, false)
	cpu.flg.Z = false

	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu10() int { // TODO: STOP
	cpu.reg.PC += 2
	return 1
}

func (cpu *CPU) cpu11() int { // LD DE, u16
	cpu.ld16(&cpu.reg.E, &cpu.reg.D,
		Read(cpu.reg.PC+1), Read(cpu.reg.PC+2))
	cpu.reg.PC += 3

	return 3
}

func (cpu *CPU) cpu12() int { // LD (DE), A
	cpu.ldToAddress(cpu.reg.E, cpu.reg.D, cpu.reg.A)

	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu13() int { // INC DE
	cpu.inc16(&cpu.reg.E, &cpu.reg.D)

	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu14() int { // INC D
	cpu.inc8(&cpu.reg.D, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu15() int { // DEC D
	cpu.dec8(&cpu.reg.D, true)

	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu16() int { // LD D, u8
	cpu.ld8(&cpu.reg.D, Read(cpu.reg.PC+1))
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpu17() int { // RLA
	cpu.rl8(&cpu.reg.A, true)
	cpu.flg.Z = false

	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu18() int { //  JR i8
	return cpu.jr(true)
}

func (cpu *CPU) cpu19() int { // ADD HL, DE
	cpu.add16(&cpu.reg.L, &cpu.reg.H, cpu.reg.E, cpu.reg.D)

	cpu.reg.PC++

	return 2
}

func (cpu *CPU) cpu1A() int { // LD A, (DE)
	cpu.ldFromAddress(&cpu.reg.A, cpu.reg.E, cpu.reg.D)

	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu1B() int { // DEC DE
	cpu.dec16(&cpu.reg.E, &cpu.reg.D)

	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu1C() int { // INC E
	cpu.inc8(&cpu.reg.E, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu1D() int { // DEC E
	cpu.dec8(&cpu.reg.E, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu1E() int { // LD E, u8
	cpu.ld8(&cpu.reg.E, Read(cpu.reg.PC+1))

	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpu1F() int { // RRA
	cpu.rr8(&cpu.reg.A, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu20() int { // JR NZ, i8
	return cpu.jr(!cpu.flg.Z)
}

func (cpu *CPU) cpu21() int { // LD HL, u16
	cpu.ld16(&cpu.reg.L, &cpu.reg.H,
		Read(cpu.reg.PC+1), Read(cpu.reg.PC+2))
	cpu.reg.PC += 3
	return 3
}

func (cpu *CPU) cpu22() int { // LD (HL+), A
	cpu.ldToAddress(cpu.reg.L, cpu.reg.H, cpu.reg.A)
	cpu.inc16(&cpu.reg.L, &cpu.reg.H)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu23() int { // INC HL
	cpu.inc16(&cpu.reg.L, &cpu.reg.H)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu24() int { // INC H
	cpu.inc8(&cpu.reg.H, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu25() int { // DEC H
	cpu.dec8(&cpu.reg.H, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu26() int { // LD H, u8
	cpu.ld8(&cpu.reg.H, Read(cpu.reg.PC+1))
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpu27() int { // TODO: DAA
	cpu.flg.H = false
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu28() int { // JR Z, i8
	return cpu.jr(cpu.flg.Z)
}

func (cpu *CPU) cpu29() int { // ADD HL, HL
	cpu.add16(&cpu.reg.L, &cpu.reg.H, cpu.reg.L, cpu.reg.H)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu2A() int { // LD A, (HL+)
	cpu.ldFromAddress(&cpu.reg.A, cpu.reg.L, cpu.reg.H)
	cpu.inc16(&cpu.reg.L, &cpu.reg.H)

	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu2B() int { // DEC HL
	cpu.dec16(&cpu.reg.L, &cpu.reg.H)
	cpu.reg.PC++

	return 2
}

func (cpu *CPU) cpu2C() int { // INC L
	cpu.inc8(&cpu.reg.L, true)

	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu2D() int { // DEC L
	cpu.dec8(&cpu.reg.L, true)

	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu2E() int { // LD L, u8
	cpu.ld8(&cpu.reg.L, Read(cpu.reg.PC+1))
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpu2F() int { // CPL
	cpu.reg.A ^= 0xFF
	cpu.flg.N = true
	cpu.flg.H = true
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu30() int { // JR NC, i8
	return cpu.jr(!cpu.flg.C)
}

func (cpu *CPU) cpu31() int { // LD SP,u16
	cpu.ld16reg(&cpu.reg.SP, Read(cpu.reg.PC+1), Read(cpu.reg.PC+2))

	cpu.reg.PC += 3
	return 3
}

func (cpu *CPU) cpu32() int { // LD (HL-), A
	cpu.ldToAddress(cpu.reg.L, cpu.reg.H, cpu.reg.A)
	cpu.dec16(&cpu.reg.L, &cpu.reg.H)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu33() int { // INC SP
	cpu.reg.SP++
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu34() int { // INC (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	Write(address, Read(address)+1)

	cpu.flg.Z = Read(address) == 0
	cpu.flg.N = false
	cpu.flg.H = Read(address)&0xF == 0x0
	cpu.reg.PC++
	return 3
}

func (cpu *CPU) cpu35() int { // DEC (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	Write(address, Read(address)-1)

	cpu.flg.Z = Read(address) == 0
	cpu.flg.N = true
	cpu.flg.H = Read(address)&0xF == 0xF
	cpu.reg.PC++
	return 3
}

func (cpu *CPU) cpu36() int { // LD (HL),u8
	cpu.ldToAddress(cpu.reg.L, cpu.reg.H, Read(cpu.reg.PC+1))
	cpu.reg.PC += 2
	return 3
}

func (cpu *CPU) cpu37() int { // SCF
	cpu.flg.N = false
	cpu.flg.H = false
	cpu.flg.C = true
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu38() int { // JR C,i8
	return cpu.jr(cpu.flg.C)
}

func (cpu *CPU) cpu39() int { // ADD HL,SP
	cpu.add16(&cpu.reg.L, &cpu.reg.H, byte(cpu.reg.SP&0xFF), byte(cpu.reg.SP>>8))
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu3A() int { // LD A, (HL-)
	cpu.ldFromAddress(&cpu.reg.A, cpu.reg.L, cpu.reg.H)
	cpu.dec16(&cpu.reg.L, &cpu.reg.H)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu3B() int { // DEC SP
	cpu.reg.SP--
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu3C() int { // INC A
	cpu.inc8(&cpu.reg.A, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu3D() int { // DEC A
	cpu.dec8(&cpu.reg.A, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu3E() int { // LD A,u8
	cpu.ld8(&cpu.reg.A, Read(cpu.reg.PC+1))
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpu3F() int { // CCF
	cpu.flg.C = !cpu.flg.C
	cpu.flg.N = false
	cpu.flg.H = false

	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu40() int { // LD B,B
	cpu.ld8(&cpu.reg.B, cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu41() int { // LD B,C
	cpu.ld8(&cpu.reg.B, cpu.reg.C)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu42() int { // LD B,D
	cpu.ld8(&cpu.reg.B, cpu.reg.D)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu43() int { // LD B,E
	cpu.ld8(&cpu.reg.B, cpu.reg.E)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu44() int { // LD B,H
	cpu.ld8(&cpu.reg.B, cpu.reg.H)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu45() int { // LD B,L
	cpu.ld8(&cpu.reg.B, cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu46() int { // LD B,(HL)
	cpu.ldFromAddress(&cpu.reg.B, cpu.reg.L, cpu.reg.H)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu47() int { // LD B,A
	cpu.ld8(&cpu.reg.B, cpu.reg.A)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu48() int { // LD C,B
	cpu.ld8(&cpu.reg.C, cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu49() int { // LD C,C
	cpu.ld8(&cpu.reg.C, cpu.reg.C)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu4A() int { // LD C,D
	cpu.ld8(&cpu.reg.C, cpu.reg.D)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu4B() int { // LD C,E
	cpu.ld8(&cpu.reg.C, cpu.reg.E)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu4C() int { // LD C,H
	cpu.ld8(&cpu.reg.C, cpu.reg.H)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu4D() int { // LD C,L
	cpu.ld8(&cpu.reg.C, cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu4E() int { // LD C,(HL)
	cpu.ldFromAddress(&cpu.reg.C, cpu.reg.L, cpu.reg.H)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu4F() int { // LD C,A
	cpu.ld8(&cpu.reg.C, cpu.reg.A)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu50() int { // LD D,B
	cpu.ld8(&cpu.reg.D, cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu51() int { // LD D,C
	cpu.ld8(&cpu.reg.D, cpu.reg.C)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu52() int { // LD D,D
	cpu.ld8(&cpu.reg.D, cpu.reg.D)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu53() int { // LD D,E
	cpu.ld8(&cpu.reg.D, cpu.reg.E)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu54() int { // LD D,H
	cpu.ld8(&cpu.reg.D, cpu.reg.H)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu55() int { // LD D,L
	cpu.ld8(&cpu.reg.D, cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu56() int { // LD D,(HL)
	cpu.ldFromAddress(&cpu.reg.D, cpu.reg.L, cpu.reg.H)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu57() int { // LD D,A
	cpu.ld8(&cpu.reg.D, cpu.reg.A)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu58() int { // LD E,B
	cpu.ld8(&cpu.reg.E, cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu59() int { // LD E,C
	cpu.ld8(&cpu.reg.E, cpu.reg.C)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu5A() int { // LD E,D
	cpu.ld8(&cpu.reg.E, cpu.reg.D)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu5B() int { // LD E,E
	cpu.ld8(&cpu.reg.E, cpu.reg.E)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu5C() int { // LD E,H
	cpu.ld8(&cpu.reg.E, cpu.reg.H)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu5D() int { // LD E,L
	cpu.ld8(&cpu.reg.E, cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu5E() int { // LD E,(HL)
	cpu.ldFromAddress(&cpu.reg.E, cpu.reg.L, cpu.reg.H)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu5F() int { // LD E,A
	cpu.ld8(&cpu.reg.E, cpu.reg.A)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu60() int { // LD H,B
	cpu.ld8(&cpu.reg.H, cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu61() int { // LD H,C
	cpu.ld8(&cpu.reg.H, cpu.reg.C)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu62() int { // LD H,D
	cpu.ld8(&cpu.reg.H, cpu.reg.D)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu63() int { // LD H,E
	cpu.ld8(&cpu.reg.H, cpu.reg.E)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu64() int { // LD H,H
	cpu.ld8(&cpu.reg.H, cpu.reg.H)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu65() int { // LD H,L
	cpu.ld8(&cpu.reg.H, cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu66() int { // LD H,(HL)
	cpu.ldFromAddress(&cpu.reg.H, cpu.reg.L, cpu.reg.H)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu67() int { // LD H,A
	cpu.ld8(&cpu.reg.H, cpu.reg.A)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu68() int { // LD L,B
	cpu.ld8(&cpu.reg.L, cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu69() int { // LD L,C
	cpu.ld8(&cpu.reg.L, cpu.reg.C)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu6A() int { // LD L,D
	cpu.ld8(&cpu.reg.L, cpu.reg.D)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu6B() int { // LD L,E
	cpu.ld8(&cpu.reg.L, cpu.reg.E)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu6C() int { // LD L,H
	cpu.ld8(&cpu.reg.L, cpu.reg.H)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu6D() int { // LD L,L
	cpu.ld8(&cpu.reg.L, cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu6E() int { // LD L,(HL)
	cpu.ldFromAddress(&cpu.reg.L, cpu.reg.L, cpu.reg.H)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu6F() int { // LD L,A
	cpu.ld8(&cpu.reg.L, cpu.reg.A)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu70() int { // LD (HL),B
	cpu.ldToAddress(cpu.reg.L, cpu.reg.H, cpu.reg.B)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu71() int { // LD (HL),C
	cpu.ldToAddress(cpu.reg.L, cpu.reg.H, cpu.reg.C)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu72() int { // LD (HL),D
	cpu.ldToAddress(cpu.reg.L, cpu.reg.H, cpu.reg.D)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu73() int { // LD (HL),E
	cpu.ldToAddress(cpu.reg.L, cpu.reg.H, cpu.reg.E)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu74() int { // LD (HL),H
	cpu.ldToAddress(cpu.reg.L, cpu.reg.H, cpu.reg.H)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu75() int { // LD (HL),L
	cpu.ldToAddress(cpu.reg.L, cpu.reg.H, cpu.reg.L)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu76() int { // TODO: HALT
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu77() int { // LD (HL),A
	cpu.ldToAddress(cpu.reg.L, cpu.reg.H, cpu.reg.A)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu78() int { // LD A,B
	cpu.ld8(&cpu.reg.A, cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu79() int { // LD A,C
	cpu.ld8(&cpu.reg.A, cpu.reg.C)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu7A() int { // LD A,D
	cpu.ld8(&cpu.reg.A, cpu.reg.D)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu7B() int { // LD A,E
	cpu.ld8(&cpu.reg.A, cpu.reg.E)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu7C() int { // LD A,H
	cpu.ld8(&cpu.reg.A, cpu.reg.H)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu7D() int { // LD A,L
	cpu.ld8(&cpu.reg.A, cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu7E() int { // LD A,(HL)
	cpu.ldFromAddress(&cpu.reg.A, cpu.reg.L, cpu.reg.H)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpu7F() int { // LD A,A
	cpu.ld8(&cpu.reg.A, cpu.reg.A)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu80() int { // ADD A,B
	cpu.addA(cpu.reg.B, false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu81() int { // ADD A,C
	cpu.addA(cpu.reg.C, false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu82() int { // ADD A,D
	cpu.addA(cpu.reg.D, false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu83() int { // ADD A,E
	cpu.addA(cpu.reg.E, false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu84() int { // ADD A,H
	cpu.addA(cpu.reg.H, false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu85() int { // ADD A,L
	cpu.addA(cpu.reg.L, false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu86() int { // ADD A,(HL)
	address := uint16(cpu.reg.H)<<8 + cpu.reg.L
	cpu.addA(Read(address), false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu87() int { // ADD A,A
	cpu.addA(cpu.reg.A, false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu88() int { // ADC A,B
	cpu.addA(cpu.reg.B, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu89() int { // ADC A,C
	cpu.addA(cpu.reg.C, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu8A() int { // ADC A,D
	cpu.addA(cpu.reg.D, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu8B() int { // ADC A,E
	cpu.addA(cpu.reg.E, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu8C() int { // ADC A,H
	cpu.addA(cpu.reg.H, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu8D() int { // ADC A,L
	cpu.addA(cpu.reg.L, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu8E() int { // ADC A,(HL)
	address := uint16(cpu.reg.H)<<8 + cpu.reg.L
	cpu.addA(Read(address), true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu8F() int { // ADC A,A
	cpu.addA(cpu.reg.A, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu90() int { // SUB A,B
	cpu.subA(cpu.reg.B, false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu91() int { // SUB A,C
	cpu.subA(cpu.reg.C, false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu92() int { // SUB A,D
	cpu.subA(cpu.reg.D, false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu93() int { // SUB A,E
	cpu.subA(cpu.reg.E, false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu94() int { // SUB A,H
	cpu.subA(cpu.reg.H, false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu95() int { // SUB A,L
	cpu.subA(cpu.reg.L, false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu96() int { // SUB A,(HL)
	address := uint16(cpu.reg.H)<<8 + cpu.reg.L
	cpu.subA(Read(address), false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu97() int { // SUB A,A
	cpu.subA(cpu.reg.A, false)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu98() int { // SBC A,B
	cpu.subA(cpu.reg.B, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu99() int { // SBC A,C
	cpu.subA(cpu.reg.C, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu9A() int { // SBC A,D
	cpu.subA(cpu.reg.D, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu9B() int { // SBC A,E
	cpu.subA(cpu.reg.E, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu9C() int { // SBC A,H
	cpu.subA(cpu.reg.H, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu9D() int { // SBC A,L
	cpu.subA(cpu.reg.L, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu9E() int { // SBC A,(HL)
	address := uint16(cpu.reg.H)<<8 + cpu.reg.L
	cpu.subA(Read(address), true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpu9F() int { // SBC A,A
	cpu.subA(cpu.reg.A, true)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuA0() int { // AND A,B
	cpu.andA(cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuA1() int { // AND A,C
	cpu.andA(cpu.reg.C)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuA2() int { // AND A,D
	cpu.andA(cpu.reg.D)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuA3() int { // AND A,E
	cpu.andA(cpu.reg.E)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuA4() int { // AND A,H
	cpu.andA(cpu.reg.H)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuA5() int { // AND A,L
	cpu.andA(cpu.reg.L)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuA6() int { // AND A,(HL)
	address := uint16(cpu.reg.H)<<8 + cpu.reg.L
	cpu.andA(Read(address))
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuA7() int { // AND A,A
	cpu.andA(cpu.reg.A)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuA8() int { // XOR A,B
	cpu.xorA(cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuA9() int { // XOR A,C
	cpu.xorA(cpu.reg.C)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuAA() int { // XOR A,D
	cpu.xorA(cpu.reg.D)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuAB() int { // XOR A,E
	cpu.xorA(cpu.reg.E)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuAC() int { // XOR A,H
	cpu.xorA(cpu.reg.H)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuAD() int { // XOR A,L
	cpu.xorA(cpu.reg.L)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuAE() int { // XOR A,(HL)
	address := uint16(cpu.reg.H)<<8 + cpu.reg.L
	cpu.xorA(Read(address))
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuAF() int { // XOR A,A
	cpu.xorA(cpu.reg.A)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuB0() int { // OR A,B
	cpu.orA(cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuB1() int { // OR A,C
	cpu.orA(cpu.reg.C)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuB2() int { // OR A,D
	cpu.orA(cpu.reg.D)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuB3() int { // OR A,E
	cpu.orA(cpu.reg.E)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuB4() int { // OR A,H
	cpu.orA(cpu.reg.H)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuB5() int { // OR A,L
	cpu.orA(cpu.reg.L)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuB6() int { // OR A,(HL)
	address := uint16(cpu.reg.H)<<8 + cpu.reg.L
	cpu.orA(Read(address))
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuB7() int { // OR A,A
	cpu.orA(cpu.reg.A)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuB8() int { // CP A,B
	cpu.cpA(cpu.reg.B)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuB9() int { // CP A,C
	cpu.cpA(cpu.reg.C)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuBA() int { // CP A,D
	cpu.cpA(cpu.reg.D)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuBB() int { // CP A,E
	cpu.cpA(cpu.reg.E)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuBC() int { // CP A,H
	cpu.cpA(cpu.reg.H)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuBD() int { // CP A,L
	cpu.cpA(cpu.reg.L)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuBE() int { // CP A,(HL)
	address := uint16(cpu.reg.H)<<8 + cpu.reg.L
	cpu.cpA(Read(address))
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuBF() int { // CP A,A
	cpu.cpA(cpu.reg.A)
	cpu.reg.PC++
	return 1
}
