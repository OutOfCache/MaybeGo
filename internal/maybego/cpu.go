package maybego

import (
	//	"fmt"
	"log"
	"os"
)

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
	Z    bool // zero flag
	C    bool // carry flag
	N    bool // sub flag
	H    bool // half carry
	HALT bool // HALT flag
	IME  bool // Interrupts Master Enable
	//	ICnt byte // Countdown to activate IME (since delay for one instruction)
}

type CPU struct {
	reg           *Registers
	flg           *Flags
	currentOpcode byte
	opcodes       [256]func() int
	cbOps         [256]func() int
	interrupts    [5]byte
}

// dummy "constructor"
func NewCPU() *CPU {
	cpu := &CPU{reg: new(Registers), flg: new(Flags)}
	cpu.reg.PC = 0x100  // to bypass boot rom for now
	cpu.reg.SP = 0xFFFE // bypassing boot rom
	cpu.reg.A = 0x1
	cpu.flg.Z = true
	cpu.flg.N = true
	cpu.flg.C = true
	cpu.reg.C = 0x13
	cpu.reg.E = 0xD8
	cpu.reg.H = 0x01
	cpu.reg.L = 0x4D

	cpu.opcodes = [256]func() int{
		cpu.cpu00, cpu.cpu01, cpu.cpu02, cpu.cpu03,
		cpu.cpu04, cpu.cpu05, cpu.cpu06, cpu.cpu07,
		cpu.cpu08, cpu.cpu09, cpu.cpu0A, cpu.cpu0B,
		cpu.cpu0C, cpu.cpu0D, cpu.cpu0E, cpu.cpu0F,
		cpu.cpu10, cpu.cpu11, cpu.cpu12, cpu.cpu13,
		cpu.cpu14, cpu.cpu15, cpu.cpu16, cpu.cpu17,
		cpu.cpu18, cpu.cpu19, cpu.cpu1A, cpu.cpu1B,
		cpu.cpu1C, cpu.cpu1D, cpu.cpu1E, cpu.cpu1F,
		cpu.cpu20, cpu.cpu21, cpu.cpu22, cpu.cpu23,
		cpu.cpu24, cpu.cpu25, cpu.cpu26, cpu.cpu27,
		cpu.cpu28, cpu.cpu29, cpu.cpu2A, cpu.cpu2B,
		cpu.cpu2C, cpu.cpu2D, cpu.cpu2E, cpu.cpu2F,
		cpu.cpu30, cpu.cpu31, cpu.cpu32, cpu.cpu33,
		cpu.cpu34, cpu.cpu35, cpu.cpu36, cpu.cpu37,
		cpu.cpu38, cpu.cpu39, cpu.cpu3A, cpu.cpu3B,
		cpu.cpu3C, cpu.cpu3D, cpu.cpu3E, cpu.cpu3F,
		cpu.cpu40, cpu.cpu41, cpu.cpu42, cpu.cpu43,
		cpu.cpu44, cpu.cpu45, cpu.cpu46, cpu.cpu47,
		cpu.cpu48, cpu.cpu49, cpu.cpu4A, cpu.cpu4B,
		cpu.cpu4C, cpu.cpu4D, cpu.cpu4E, cpu.cpu4F,
		cpu.cpu50, cpu.cpu51, cpu.cpu52, cpu.cpu53,
		cpu.cpu54, cpu.cpu55, cpu.cpu56, cpu.cpu57,
		cpu.cpu58, cpu.cpu59, cpu.cpu5A, cpu.cpu5B,
		cpu.cpu5C, cpu.cpu5D, cpu.cpu5E, cpu.cpu5F,
		cpu.cpu60, cpu.cpu61, cpu.cpu62, cpu.cpu63,
		cpu.cpu64, cpu.cpu65, cpu.cpu66, cpu.cpu67,
		cpu.cpu68, cpu.cpu69, cpu.cpu6A, cpu.cpu6B,
		cpu.cpu6C, cpu.cpu6D, cpu.cpu6E, cpu.cpu6F,
		cpu.cpu70, cpu.cpu71, cpu.cpu72, cpu.cpu73,
		cpu.cpu74, cpu.cpu75, cpu.cpu76, cpu.cpu77,
		cpu.cpu78, cpu.cpu79, cpu.cpu7A, cpu.cpu7B,
		cpu.cpu7C, cpu.cpu7D, cpu.cpu7E, cpu.cpu7F,
		cpu.cpu80, cpu.cpu81, cpu.cpu82, cpu.cpu83,
		cpu.cpu84, cpu.cpu85, cpu.cpu86, cpu.cpu87,
		cpu.cpu88, cpu.cpu89, cpu.cpu8A, cpu.cpu8B,
		cpu.cpu8C, cpu.cpu8D, cpu.cpu8E, cpu.cpu8F,
		cpu.cpu90, cpu.cpu91, cpu.cpu92, cpu.cpu93,
		cpu.cpu94, cpu.cpu95, cpu.cpu96, cpu.cpu97,
		cpu.cpu98, cpu.cpu99, cpu.cpu9A, cpu.cpu9B,
		cpu.cpu9C, cpu.cpu9D, cpu.cpu9E, cpu.cpu9F,
		cpu.cpuA0, cpu.cpuA1, cpu.cpuA2, cpu.cpuA3,
		cpu.cpuA4, cpu.cpuA5, cpu.cpuA6, cpu.cpuA7,
		cpu.cpuA8, cpu.cpuA9, cpu.cpuAA, cpu.cpuAB,
		cpu.cpuAC, cpu.cpuAD, cpu.cpuAE, cpu.cpuAF,
		cpu.cpuB0, cpu.cpuB1, cpu.cpuB2, cpu.cpuB3,
		cpu.cpuB4, cpu.cpuB5, cpu.cpuB6, cpu.cpuB7,
		cpu.cpuB8, cpu.cpuB9, cpu.cpuBA, cpu.cpuBB,
		cpu.cpuBC, cpu.cpuBD, cpu.cpuBE, cpu.cpuBF,
		cpu.cpuC0, cpu.cpuC1, cpu.cpuC2, cpu.cpuC3,
		cpu.cpuC4, cpu.cpuC5, cpu.cpuC6, cpu.cpuC7,
		cpu.cpuC8, cpu.cpuC9, cpu.cpuCA, cpu.cpuCB,
		cpu.cpuCC, cpu.cpuCD, cpu.cpuCE, cpu.cpuCF,
		cpu.cpuD0, cpu.cpuD1, cpu.cpuD2, cpu.cpuD3,
		cpu.cpuD4, cpu.cpuD5, cpu.cpuD6, cpu.cpuD7,
		cpu.cpuD8, cpu.cpuD9, cpu.cpuDA, cpu.cpuDB,
		cpu.cpuDC, cpu.cpuDD, cpu.cpuDE, cpu.cpuDF,
		cpu.cpuE0, cpu.cpuE1, cpu.cpuE2, cpu.cpuE3,
		cpu.cpuE4, cpu.cpuE5, cpu.cpuE6, cpu.cpuE7,
		cpu.cpuE8, cpu.cpuE9, cpu.cpuEA, cpu.cpuEB,
		cpu.cpuEC, cpu.cpuED, cpu.cpuEE, cpu.cpuEF,
		cpu.cpuF0, cpu.cpuF1, cpu.cpuF2, cpu.cpuF3,
		cpu.cpuF4, cpu.cpuF5, cpu.cpuF6, cpu.cpuF7,
		cpu.cpuF8, cpu.cpuF9, cpu.cpuFA, cpu.cpuFB,
		cpu.cpuFC, cpu.cpuFD, cpu.cpuFE, cpu.cpuFF,
	}

	cpu.cbOps = [256]func() int{
		cpu.cb00, cpu.cb01, cpu.cb02, cpu.cb03,
		cpu.cb04, cpu.cb05, cpu.cb06, cpu.cb07,
		cpu.cb08, cpu.cb09, cpu.cb0A, cpu.cb0B,
		cpu.cb0C, cpu.cb0D, cpu.cb0E, cpu.cb0F,
		cpu.cb10, cpu.cb11, cpu.cb12, cpu.cb13,
		cpu.cb14, cpu.cb15, cpu.cb16, cpu.cb17,
		cpu.cb18, cpu.cb19, cpu.cb1A, cpu.cb1B,
		cpu.cb1C, cpu.cb1D, cpu.cb1E, cpu.cb1F,
		cpu.cb20, cpu.cb21, cpu.cb22, cpu.cb23,
		cpu.cb24, cpu.cb25, cpu.cb26, cpu.cb27,
		cpu.cb28, cpu.cb29, cpu.cb2A, cpu.cb2B,
		cpu.cb2C, cpu.cb2D, cpu.cb2E, cpu.cb2F,
		cpu.cb30, cpu.cb31, cpu.cb32, cpu.cb33,
		cpu.cb34, cpu.cb35, cpu.cb36, cpu.cb37,
		cpu.cb38, cpu.cb39, cpu.cb3A, cpu.cb3B,
		cpu.cb3C, cpu.cb3D, cpu.cb3E, cpu.cb3F,
		cpu.cb40, cpu.cb41, cpu.cb42, cpu.cb43,
		cpu.cb44, cpu.cb45, cpu.cb46, cpu.cb47,
		cpu.cb48, cpu.cb49, cpu.cb4A, cpu.cb4B,
		cpu.cb4C, cpu.cb4D, cpu.cb4E, cpu.cb4F,
		cpu.cb50, cpu.cb51, cpu.cb52, cpu.cb53,
		cpu.cb54, cpu.cb55, cpu.cb56, cpu.cb57,
		cpu.cb58, cpu.cb59, cpu.cb5A, cpu.cb5B,
		cpu.cb5C, cpu.cb5D, cpu.cb5E, cpu.cb5F,
		cpu.cb60, cpu.cb61, cpu.cb62, cpu.cb63,
		cpu.cb64, cpu.cb65, cpu.cb66, cpu.cb67,
		cpu.cb68, cpu.cb69, cpu.cb6A, cpu.cb6B,
		cpu.cb6C, cpu.cb6D, cpu.cb6E, cpu.cb6F,
		cpu.cb70, cpu.cb71, cpu.cb72, cpu.cb73,
		cpu.cb74, cpu.cb75, cpu.cb76, cpu.cb77,
		cpu.cb78, cpu.cb79, cpu.cb7A, cpu.cb7B,
		cpu.cb7C, cpu.cb7D, cpu.cb7E, cpu.cb7F,
		cpu.cb80, cpu.cb81, cpu.cb82, cpu.cb83,
		cpu.cb84, cpu.cb85, cpu.cb86, cpu.cb87,
		cpu.cb88, cpu.cb89, cpu.cb8A, cpu.cb8B,
		cpu.cb8C, cpu.cb8D, cpu.cb8E, cpu.cb8F,
		cpu.cb90, cpu.cb91, cpu.cb92, cpu.cb93,
		cpu.cb94, cpu.cb95, cpu.cb96, cpu.cb97,
		cpu.cb98, cpu.cb99, cpu.cb9A, cpu.cb9B,
		cpu.cb9C, cpu.cb9D, cpu.cb9E, cpu.cb9F,
		cpu.cbA0, cpu.cbA1, cpu.cbA2, cpu.cbA3,
		cpu.cbA4, cpu.cbA5, cpu.cbA6, cpu.cbA7,
		cpu.cbA8, cpu.cbA9, cpu.cbAA, cpu.cbAB,
		cpu.cbAC, cpu.cbAD, cpu.cbAE, cpu.cbAF,
		cpu.cbB0, cpu.cbB1, cpu.cbB2, cpu.cbB3,
		cpu.cbB4, cpu.cbB5, cpu.cbB6, cpu.cbB7,
		cpu.cbB8, cpu.cbB9, cpu.cbBA, cpu.cbBB,
		cpu.cbBC, cpu.cbBD, cpu.cbBE, cpu.cbBF,
		cpu.cbC0, cpu.cbC1, cpu.cbC2, cpu.cbC3,
		cpu.cbC4, cpu.cbC5, cpu.cbC6, cpu.cbC7,
		cpu.cbC8, cpu.cbC9, cpu.cbCA, cpu.cbCB,
		cpu.cbCC, cpu.cbCD, cpu.cbCE, cpu.cbCF,
		cpu.cbD0, cpu.cbD1, cpu.cbD2, cpu.cbD3,
		cpu.cbD4, cpu.cbD5, cpu.cbD6, cpu.cbD7,
		cpu.cbD8, cpu.cbD9, cpu.cbDA, cpu.cbDB,
		cpu.cbDC, cpu.cbDD, cpu.cbDE, cpu.cbDF,
		cpu.cbE0, cpu.cbE1, cpu.cbE2, cpu.cbE3,
		cpu.cbE4, cpu.cbE5, cpu.cbE6, cpu.cbE7,
		cpu.cbE8, cpu.cbE9, cpu.cbEA, cpu.cbEB,
		cpu.cbEC, cpu.cbED, cpu.cbEE, cpu.cbEF,
		cpu.cbF0, cpu.cbF1, cpu.cbF2, cpu.cbF3,
		cpu.cbF4, cpu.cbF5, cpu.cbF6, cpu.cbF7,
		cpu.cbF8, cpu.cbF9, cpu.cbFA, cpu.cbFB,
		cpu.cbFC, cpu.cbFD, cpu.cbFE, cpu.cbFF,
	}

	cpu.interrupts = [5]byte{0x40, 0x48, 0x50, 0x58, 0x60}

	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	return cpu
}

func (cpu *CPU) Fetch() {
	if cpu.flg.IME {
		cpu.interrupt()
	}

	if cpu.flg.HALT {
		return
	}
	cpu.currentOpcode = Read(cpu.reg.PC)
	log.Printf("PC: %x, Opcode: %x, PC+1:%x, PC+2: %x\tZ: %t, N: %t, H: %t, C: %t\nA: %x, B: %x, C: %x, D: %x, E: %x, H: %x, L: %x",
		cpu.reg.PC, cpu.currentOpcode, Read(cpu.reg.PC+1), Read(cpu.reg.PC+2), cpu.flg.Z, cpu.flg.N, cpu.flg.H, cpu.flg.C,
		cpu.reg.A, cpu.reg.B, cpu.reg.C, cpu.reg.D, cpu.reg.E, cpu.reg.H, cpu.reg.L)

}

func (cpu *CPU) Decode() {
	cpu.opcodes[cpu.currentOpcode]()
}

func FlagToBit(flag bool) byte {
	if flag {
		return 1
	}
	return 0
}

func (cpu *CPU) FlagsToBytes() byte {
	z := FlagToBit(cpu.flg.Z)
	n := FlagToBit(cpu.flg.N)
	h := FlagToBit(cpu.flg.H)
	c := FlagToBit(cpu.flg.C)

	return (z << 7) + (n << 6) + (h << 5) + (c << 4)
}

func (cpu *CPU) BytesToFlags(flags byte) {
	cpu.flg.Z = flags&0x80 == 0x80
	cpu.flg.N = flags&0x40 == 0x40
	cpu.flg.H = flags&0x20 == 0x20
	cpu.flg.C = flags&0x10 == 0x10
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
		cpu.flg.H = *reg&0xF == 0x0
	}
}

func (cpu *CPU) inc16(destLo *byte, destHi *byte) {
	cpu.inc8(destLo, false)
	if *destLo == 0 { // increase if overflow in low byte
		cpu.inc8(destHi, false)
	}
}

func (cpu *CPU) dec8(reg *byte, flags bool) {
	*reg -= 1
	if flags {
		cpu.flg.N = true
		cpu.flg.Z = *reg == 0
		cpu.flg.H = (*reg)&0xF == 0xF
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
	cpu.flg.H = ((cpu.reg.A&0xF)+(reg&0xF)+FlagToBit(carry))&0x10 == 0x10

	sum := uint16(cpu.reg.A) + uint16(reg) + uint16(FlagToBit(carry))
	cpu.reg.A = byte(sum)

	cpu.flg.Z = cpu.reg.A == 0
	cpu.flg.C = sum&0x100 == 0x100
	cpu.flg.N = false
}

func (cpu *CPU) subA(reg byte, carry bool) {
	if carry {
		carry = cpu.flg.C
	}
	data := uint16(reg) + uint16(FlagToBit(carry))
	cpu.flg.H = (cpu.reg.A & 0xF) < (reg&0xF)+FlagToBit(carry)
	cpu.flg.C = uint16(cpu.reg.A) < data

	cpu.reg.A = byte((uint16(cpu.reg.A) - data))

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
	cpu.flg.H = (cpu.reg.A & 0xF) < (reg & 0xF)
	cpu.flg.C = cpu.reg.A < reg

	result := cpu.reg.A - reg

	cpu.flg.Z = result == 0
	cpu.flg.N = true
}

func (cpu *CPU) add16(destLo *byte, destHi *byte, srcLo byte, srcHi byte) {
	cpu.flg.N = false

	sum := uint16(*destLo) + uint16(srcLo)
	*destLo = byte(sum & 0xFF)
	cpu.flg.H = (byte(sum>>8)+(*destHi&0xf)+(srcHi&0xf))&0x10 == 0x10
	sum = (sum >> 8) + uint16(*destHi) + uint16(srcHi)
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

	*reg = (*reg >> 1) + (msb << 7)

	cpu.flg.N = false
	cpu.flg.H = false
	cpu.flg.Z = *reg == 0
}

func (cpu *CPU) sl8(reg *byte) {
	cpu.flg.C = *reg&0x80 == 0x80
	*reg <<= 1
	cpu.flg.Z = *reg == 0
	cpu.flg.N = false
	cpu.flg.H = false
}

func (cpu *CPU) sr8(reg *byte) {
	// arithmetic shift
	msb := *reg & 0x80
	cpu.flg.C = *reg&0x01 == 0x01
	*reg >>= 1
	*reg += msb
	cpu.flg.Z = *reg == 0
	cpu.flg.N = false
	cpu.flg.H = false
}

func (cpu *CPU) srl8(reg *byte) {
	cpu.flg.C = *reg&0x01 == 0x01
	*reg >>= 1
	cpu.flg.Z = *reg == 0
	cpu.flg.N = false
	cpu.flg.H = false
}

func (cpu *CPU) jr(flag bool) int {
	log.Printf("JR")
	if flag {
		cpu.reg.PC += uint16(2 + int8(Read(cpu.reg.PC+1)))
		return 3
	}
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) jp(flag bool) int {
	log.Printf("JP")
	if flag {
		cpu.reg.PC = uint16(Read(cpu.reg.PC+1)) + (uint16(Read(cpu.reg.PC+2)) << 8)
		return 4
	}
	cpu.reg.PC += 3
	return 3
}

func (cpu *CPU) call(flag bool) int {
	log.Printf("CALL")
	if flag {
		lo := byte(cpu.reg.PC + 3)
		hi := byte((cpu.reg.PC + 3) >> 8)
		cpu.push16(lo, hi)
		cpu.reg.PC = uint16(Read(cpu.reg.PC+1)) + (uint16(Read(cpu.reg.PC+2)) << 8)
		return 6
	}
	cpu.reg.PC += 3
	return 3
}

func (cpu *CPU) ret(flag bool) int {
	log.Printf("RET")
	// return if flag is true, otherwise continue to next instruction
	if flag {
		cpu.pop16reg(&cpu.reg.PC)
		return 5
	}
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) rst(vec byte) int {
	log.Printf("RST")
	lo := byte(cpu.reg.PC + 1)
	hi := byte((cpu.reg.PC + 1) >> 8)
	cpu.push16(lo, hi)
	cpu.reg.PC = uint16(vec)
	return 4
}

func (cpu *CPU) push16(lo byte, hi byte) {
	cpu.reg.SP -= 1
	Write(cpu.reg.SP, hi)
	cpu.reg.SP -= 1
	Write(cpu.reg.SP, lo)
}

func (cpu *CPU) pop16(destLo *byte, destHi *byte) {
	cpu.ldFromAddress(destLo, byte(cpu.reg.SP), byte(cpu.reg.SP>>8))
	cpu.reg.SP += 1
	cpu.ldFromAddress(destHi, byte(cpu.reg.SP), byte(cpu.reg.SP>>8))
	cpu.reg.SP += 1
}

func (cpu *CPU) pop16reg(dest *uint16) {
	var lo byte
	var hi byte
	cpu.ldFromAddress(&lo, byte(cpu.reg.SP), byte(cpu.reg.SP>>8))
	cpu.reg.SP += 1
	cpu.ldFromAddress(&hi, byte(cpu.reg.SP), byte(cpu.reg.SP>>8))
	cpu.reg.SP += 1

	*dest = (uint16(hi) << 8) + uint16(lo)
}

func (cpu *CPU) swap(dest *byte) {
	lo := *dest & 0x0F
	hi := *dest & 0xF0 >> 4
	*dest = lo<<4 + hi

	cpu.flg.Z = *dest == 0
	cpu.flg.N = false
	cpu.flg.H = false
	cpu.flg.C = false
}

func (cpu *CPU) bit(reg *byte, bit byte) {
	cpu.flg.N = false
	cpu.flg.H = true
	cpu.flg.Z = *reg&byte(0x01<<bit) == 0
}

func (cpu *CPU) res(reg *byte, bit byte) {
	// set bit to 0 by using AND
	*reg = *reg & ((0x01 << bit) ^ 0xFF) // shift the one to the right place and invert for a mask
}

func (cpu *CPU) set(reg *byte, bit byte) {
	// sest bit to 1 by using OR
	*reg = *reg | (0x01 << bit)
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
	cpu.flg.Z = false
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

func (cpu *CPU) cpu27() int { // DAA
	// from nesdev
	if !cpu.flg.N {
		if cpu.flg.C || cpu.reg.A > 0x99 {
			cpu.reg.A += 0x60
			cpu.flg.C = true
		}
		if cpu.flg.H || cpu.reg.A&0x0F > 0x09 {
			cpu.reg.A += 0x06
		}
	} else {
		if cpu.flg.C {
			cpu.reg.A -= 0x60
		}
		if cpu.flg.H {
			cpu.reg.A -= 0x06
		}
	}

	cpu.flg.Z = cpu.reg.A == 0
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
	cpu.ld8(&cpu.reg.B, cpu.reg.L)
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
	cpu.ld8(&cpu.reg.C, cpu.reg.L)
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
	cpu.ld8(&cpu.reg.D, cpu.reg.L)
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
	cpu.ld8(&cpu.reg.E, cpu.reg.L)
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
	cpu.ld8(&cpu.reg.H, cpu.reg.L)
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
	cpu.ld8(&cpu.reg.L, cpu.reg.L)
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

func (cpu *CPU) cpu76() int { // HALT
	cpu.flg.HALT = true
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
	cpu.ld8(&cpu.reg.A, cpu.reg.L)
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
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
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
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
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
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
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
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
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
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
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
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
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
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
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
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	cpu.cpA(Read(address))
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuBF() int { // CP A,A
	cpu.cpA(cpu.reg.A)
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuC0() int { // RET NZ
	return cpu.ret(!cpu.flg.Z)
}

func (cpu *CPU) cpuC1() int { // POP BC
	cpu.pop16(&cpu.reg.C, &cpu.reg.B)
	cpu.reg.PC++
	return 3
}

func (cpu *CPU) cpuC2() int { // JP NZ, u16
	return cpu.jp(!cpu.flg.Z)
}

func (cpu *CPU) cpuC3() int { // JP u16
	return cpu.jp(true)
}

func (cpu *CPU) cpuC4() int { // CALL NZ, u16
	return cpu.call(!cpu.flg.Z)
}

func (cpu *CPU) cpuC5() int { // PUSH BC
	cpu.push16(cpu.reg.C, cpu.reg.B)
	cpu.reg.PC++
	return 4
}

func (cpu *CPU) cpuC6() int { // ADD A, u8
	cpu.addA(Read(cpu.reg.PC+1), false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpuC7() int { // RST 0x00
	return cpu.rst(0x00)
}

func (cpu *CPU) cpuC8() int { // RET Z
	return cpu.ret(cpu.flg.Z)
}

func (cpu *CPU) cpuC9() int { // RET
	return cpu.ret(true)
}

func (cpu *CPU) cpuCA() int { // JP Z,u16
	return cpu.jp(cpu.flg.Z)
}

func (cpu *CPU) cpuCB() int { // Prefix 0xCB
	return cpu.cbOps[Read(cpu.reg.PC+1)]()
}

func (cpu *CPU) cpuCC() int { // CALL Z,u16
	return cpu.call(cpu.flg.Z)
}

func (cpu *CPU) cpuCD() int { // CALL u16
	return cpu.call(true)
}

func (cpu *CPU) cpuCE() int { // ADC A,u8
	cpu.addA(Read(cpu.reg.PC+1), true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpuCF() int { // RST 08
	return cpu.rst(0x08)
}

func (cpu *CPU) cpuD0() int { // RET NC
	return cpu.ret(!cpu.flg.C)
}

func (cpu *CPU) cpuD1() int { // POP DE
	cpu.pop16(&cpu.reg.E, &cpu.reg.D)
	cpu.reg.PC++
	return 3
}

func (cpu *CPU) cpuD2() int { // JP NC, u16
	return cpu.jp(!cpu.flg.C)
}

func (cpu *CPU) cpuD3() int { // invalid
	return cpu.nop()
}

func (cpu *CPU) cpuD4() int { // CALL NC, u16
	return cpu.call(!cpu.flg.C)
}

func (cpu *CPU) cpuD5() int { // PUSH DE
	cpu.push16(cpu.reg.E, cpu.reg.D)
	cpu.reg.PC++
	return 4
}

func (cpu *CPU) cpuD6() int { // SUB A, u8
	cpu.subA(Read(cpu.reg.PC+1), false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpuD7() int { // RST 0x10
	return cpu.rst(0x10)
}

func (cpu *CPU) cpuD8() int { // RET C
	return cpu.ret(cpu.flg.C)
}

func (cpu *CPU) cpuD9() int { // RETI
	cpu.flg.IME = true
	cpu.ret(true)
	return 4
}

func (cpu *CPU) cpuDA() int { // JP C,u16
	return cpu.jp(cpu.flg.C)
}

func (cpu *CPU) cpuDB() int { // invalid
	return cpu.nop()
}

func (cpu *CPU) cpuDC() int { // CALL C,u16
	return cpu.call(cpu.flg.C)
}

func (cpu *CPU) cpuDD() int { // invalid
	return cpu.nop()
}

func (cpu *CPU) cpuDE() int { // SBC A,u8
	cpu.subA(Read(cpu.reg.PC+1), true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpuDF() int { // RST 18
	return cpu.rst(0x18)
}

func (cpu *CPU) nop() int { // invalid
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuE0() int { // LD (FF00+u8),A
	cpu.ldToAddress(Read(cpu.reg.PC+1), 0xFF, cpu.reg.A)
	cpu.reg.PC += 2
	return 3
}

func (cpu *CPU) cpuE1() int { // POP HL
	cpu.pop16(&cpu.reg.L, &cpu.reg.H)
	cpu.reg.PC++
	return 3
}

func (cpu *CPU) cpuE2() int { // LD (FF00+C),A
	cpu.ldToAddress(cpu.reg.C, 0xFF, cpu.reg.A)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpuE3() int { // invalid
	return cpu.nop()
}

func (cpu *CPU) cpuE4() int { // invalid
	return cpu.nop()
}

func (cpu *CPU) cpuE5() int { // PUSH HL
	cpu.push16(cpu.reg.L, cpu.reg.H)
	cpu.reg.PC++
	return 4
}

func (cpu *CPU) cpuE6() int { // AND A,u8
	cpu.andA(Read(cpu.reg.PC + 1))
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpuE7() int { // RST 20
	return cpu.rst(0x20)
}

func (cpu *CPU) cpuE8() int { // ADD SP,i8
	cpu.flg.H = cpu.reg.SP&0x0F+uint16(Read(cpu.reg.PC+1)&0x0F)&0x010 == 0x10
	cpu.flg.C = cpu.reg.SP&0xFF+uint16(Read(cpu.reg.PC+1)&0xFF)&0x100 == 0x100
	cpu.reg.SP = uint16(int16(cpu.reg.SP) + int16(Read(cpu.reg.PC+1)))

	cpu.flg.Z = false
	cpu.flg.N = false
	cpu.reg.PC += 2
	return 4
}

func (cpu *CPU) cpuE9() int { // JP HL
	cpu.reg.PC = uint16((cpu.reg.H))<<8 + uint16(cpu.reg.L)
	return 1
}

func (cpu *CPU) cpuEA() int { // LD (u16),A
	cpu.ldToAddress(Read(cpu.reg.PC+1), Read(cpu.reg.PC+2), cpu.reg.A)
	cpu.reg.PC += 3
	return 4
}

func (cpu *CPU) cpuEB() int { // invalid
	return cpu.nop()
}

func (cpu *CPU) cpuEC() int { // invalid
	return cpu.nop()
}

func (cpu *CPU) cpuED() int { // invalid
	return cpu.nop()
}

func (cpu *CPU) cpuEE() int { // XOR A,u8
	cpu.xorA(Read(cpu.reg.PC + 1))
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpuEF() int { // RST 0x28
	return cpu.rst(0x28)
}

func (cpu *CPU) cpuF0() int { // LD A,(FF00+u8)
	cpu.ldFromAddress(&cpu.reg.A, Read(cpu.reg.PC+1), 0xFF)
	cpu.reg.PC += 2
	return 3
}

func (cpu *CPU) cpuF1() int { // POP AF
	flags := cpu.FlagsToBytes()
	cpu.pop16(&flags, &cpu.reg.A)
	cpu.BytesToFlags(flags)
	cpu.reg.PC++
	return 3
}

func (cpu *CPU) cpuF2() int { // LD A,(FF00+C)
	cpu.ldFromAddress(&cpu.reg.A, cpu.reg.C, 0xFF)
	cpu.reg.PC++
	return 2
}

func (cpu *CPU) cpuF3() int { // DI
	cpu.flg.IME = false
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuF4() int { // invalid
	return cpu.nop()
}

func (cpu *CPU) cpuF5() int { // PUSH AF
	cpu.push16(cpu.FlagsToBytes(), cpu.reg.A)
	cpu.reg.PC++
	return 4
}

func (cpu *CPU) cpuF6() int { // OR A,u8
	cpu.orA(Read(cpu.reg.PC + 1))
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpuF7() int { // RST 30
	return cpu.rst(0x30)
}

func (cpu *CPU) cpuF8() int { // LD HL,SP+i8
	cpu.flg.H = cpu.reg.SP&0xF+uint16(Read(cpu.reg.PC+1)&0xF)&0x10 == 0x10
	cpu.flg.C = cpu.reg.SP&0xFF+uint16(Read(cpu.reg.PC+1)&0xFF)&0x100 == 0x100
	hl := uint16(int16(cpu.reg.SP) + int16(Read(cpu.reg.PC+1)))

	cpu.reg.L = byte(hl)
	cpu.reg.H = byte(hl >> 8)

	cpu.flg.Z = false
	cpu.flg.N = false
	cpu.reg.PC += 2
	return 4
}

func (cpu *CPU) cpuF9() int { // LD SP,HL
	cpu.ld16reg(&cpu.reg.SP, cpu.reg.L, cpu.reg.H)
	return 2
}

func (cpu *CPU) cpuFA() int { // LD A,(u16)
	cpu.ldFromAddress(&cpu.reg.A, Read(cpu.reg.PC+1), Read(cpu.reg.PC+2))
	cpu.reg.PC += 3
	return 4
}

func (cpu *CPU) cpuFB() int { // EI
	cpu.flg.IME = true
	cpu.reg.PC++
	return 1
}

func (cpu *CPU) cpuFC() int { // invalid
	return cpu.nop()
}

func (cpu *CPU) cpuFD() int { // invalid
	return cpu.nop()
}

func (cpu *CPU) cpuFE() int { // CP A,u8
	cpu.cpA(Read(cpu.reg.PC + 1))
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cpuFF() int { // RST 0x38
	return cpu.rst(0x38)
}

// Prefix CB functions

func (cpu *CPU) cb00() int { // RLC B
	cpu.rl8(&cpu.reg.B, false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb01() int { // RLC C
	cpu.rl8(&cpu.reg.C, false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb02() int { // RLC D
	cpu.rl8(&cpu.reg.D, false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb03() int { // RLC E
	cpu.rl8(&cpu.reg.E, false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb04() int { // RLC H
	cpu.rl8(&cpu.reg.H, false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb05() int { // RLC L
	cpu.rl8(&cpu.reg.L, false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb06() int { // RLC (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.rl8(&val, false)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb07() int { // RLC A
	cpu.rl8(&cpu.reg.A, false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb08() int { // RRC B
	cpu.rr8(&cpu.reg.B, false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb09() int { // RRC C
	cpu.rr8(&cpu.reg.C, false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb0A() int { // RRC D
	cpu.rr8(&cpu.reg.D, false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb0B() int { // RRC E
	cpu.rr8(&cpu.reg.E, false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb0C() int { // RRC H
	cpu.rr8(&cpu.reg.H, false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb0D() int { // RRC L
	cpu.rr8(&cpu.reg.L, false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb0E() int { // RRC (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.rr8(&val, false)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb0F() int { // RRC A
	cpu.rr8(&cpu.reg.A, false)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb10() int { // RL B
	cpu.rl8(&cpu.reg.B, true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb11() int { // RL C
	cpu.rl8(&cpu.reg.C, true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb12() int { // RL D
	cpu.rl8(&cpu.reg.D, true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb13() int { // RL E
	cpu.rl8(&cpu.reg.E, true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb14() int { // RL H
	cpu.rl8(&cpu.reg.H, true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb15() int { // RL L
	cpu.rl8(&cpu.reg.L, true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb16() int { // RL (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.rl8(&val, true)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb17() int { // RL A
	cpu.rl8(&cpu.reg.A, true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb18() int { // RR B
	cpu.rr8(&cpu.reg.B, true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb19() int { // RR C
	cpu.rr8(&cpu.reg.C, true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb1A() int { // RR D
	cpu.rr8(&cpu.reg.D, true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb1B() int { // RR E
	cpu.rr8(&cpu.reg.E, true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb1C() int { // RR H
	cpu.rr8(&cpu.reg.H, true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb1D() int { // RR L
	cpu.rr8(&cpu.reg.L, true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb1E() int { // RR (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.rr8(&val, true)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb1F() int { // RR A
	cpu.rr8(&cpu.reg.A, true)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb20() int { // SLA B
	cpu.sl8(&cpu.reg.B)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb21() int { // SLA C
	cpu.sl8(&cpu.reg.C)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb22() int { // SLA D
	cpu.sl8(&cpu.reg.D)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb23() int { // SLA E
	cpu.sl8(&cpu.reg.E)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb24() int { // SLA H
	cpu.sl8(&cpu.reg.H)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb25() int { // SLA L
	cpu.sl8(&cpu.reg.L)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb26() int { // SLA (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.sl8(&val)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb27() int { // SLA A
	cpu.sl8(&cpu.reg.A)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb28() int { // SRA B
	cpu.sr8(&cpu.reg.B)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb29() int { // SRA C
	cpu.sr8(&cpu.reg.C)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb2A() int { // SRA D
	cpu.sr8(&cpu.reg.D)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb2B() int { // SRA E
	cpu.sr8(&cpu.reg.E)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb2C() int { // SRA H
	cpu.sr8(&cpu.reg.H)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb2D() int { // SRA L
	cpu.sr8(&cpu.reg.L)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb2E() int { // SRA (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.sr8(&val)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb2F() int { // SRA A
	cpu.sr8(&cpu.reg.A)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb30() int { // SWAP B
	cpu.swap(&cpu.reg.B)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb31() int { // SWAP C
	cpu.swap(&cpu.reg.C)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb32() int { // SWAP D
	cpu.swap(&cpu.reg.D)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb33() int { // SWAP E
	cpu.swap(&cpu.reg.E)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb34() int { // SWAP H
	cpu.swap(&cpu.reg.H)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb35() int { // SWAP L
	cpu.swap(&cpu.reg.L)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb36() int { // SWAP (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.swap(&val)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb37() int { // SWAP A
	cpu.swap(&cpu.reg.A)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb38() int { // SRL B
	cpu.srl8(&cpu.reg.B)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb39() int { // SRL C
	cpu.srl8(&cpu.reg.C)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb3A() int { // SRL D
	cpu.srl8(&cpu.reg.D)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb3B() int { // SRL E
	cpu.srl8(&cpu.reg.E)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb3C() int { // SRL H
	cpu.srl8(&cpu.reg.H)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb3D() int { // SRL L
	cpu.srl8(&cpu.reg.L)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb3E() int { // SRL (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.srl8(&val)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb3F() int { // SRL A
	cpu.srl8(&cpu.reg.A)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb40() int { // BIT 0, B
	cpu.bit(&cpu.reg.B, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb41() int { // BIT 0, C
	cpu.bit(&cpu.reg.C, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb42() int { // BIT 0, D
	cpu.bit(&cpu.reg.D, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb43() int { // BIT 0, E
	cpu.bit(&cpu.reg.E, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb44() int { // BIT 0, H
	cpu.bit(&cpu.reg.H, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb45() int { // BIT 0, L
	cpu.bit(&cpu.reg.L, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb46() int { // BIT 0, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.bit(&val, 0)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb47() int { // BIT 0, A
	cpu.bit(&cpu.reg.A, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb48() int { // BIT 1, B
	cpu.bit(&cpu.reg.B, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb49() int { // BIT 1, C
	cpu.bit(&cpu.reg.C, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb4A() int { // BIT 1, D
	cpu.bit(&cpu.reg.D, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb4B() int { // BIT 1, E
	cpu.bit(&cpu.reg.E, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb4C() int { // BIT 1, H
	cpu.bit(&cpu.reg.H, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb4D() int { // BIT 1, L
	cpu.bit(&cpu.reg.L, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb4E() int { // BIT 1, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.bit(&val, 1)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb4F() int { // BIT 1, A
	cpu.bit(&cpu.reg.A, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb50() int { // BIT 2, B
	cpu.bit(&cpu.reg.B, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb51() int { // BIT 2, C
	cpu.bit(&cpu.reg.C, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb52() int { // BIT 2, D
	cpu.bit(&cpu.reg.D, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb53() int { // BIT 2, E
	cpu.bit(&cpu.reg.E, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb54() int { // BIT 2, H
	cpu.bit(&cpu.reg.H, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb55() int { // BIT 2, L
	cpu.bit(&cpu.reg.L, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb56() int { // BIT 2, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.bit(&val, 2)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb57() int { // BIT 2, A
	cpu.bit(&cpu.reg.A, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb58() int { // BIT 3, B
	cpu.bit(&cpu.reg.B, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb59() int { // BIT 3, C
	cpu.bit(&cpu.reg.C, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb5A() int { // BIT 3, D
	cpu.bit(&cpu.reg.D, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb5B() int { // BIT 3, E
	cpu.bit(&cpu.reg.E, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb5C() int { // BIT 3, H
	cpu.bit(&cpu.reg.H, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb5D() int { // BIT 3, L
	cpu.bit(&cpu.reg.L, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb5E() int { // BIT 3, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.bit(&val, 3)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb5F() int { // BIT 3, A
	cpu.bit(&cpu.reg.A, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb60() int { // BIT 4, B
	cpu.bit(&cpu.reg.B, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb61() int { // BIT 4, C
	cpu.bit(&cpu.reg.C, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb62() int { // BIT 4, D
	cpu.bit(&cpu.reg.D, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb63() int { // BIT 4, E
	cpu.bit(&cpu.reg.E, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb64() int { // BIT 4, H
	cpu.bit(&cpu.reg.H, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb65() int { // BIT 4, L
	cpu.bit(&cpu.reg.L, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb66() int { // BIT 4, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.bit(&val, 4)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb67() int { // BIT 4, A
	cpu.bit(&cpu.reg.A, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb68() int { // BIT 5, B
	cpu.bit(&cpu.reg.B, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb69() int { // BIT 5, C
	cpu.bit(&cpu.reg.C, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb6A() int { // BIT 5, D
	cpu.bit(&cpu.reg.D, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb6B() int { // BIT 5, E
	cpu.bit(&cpu.reg.E, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb6C() int { // BIT 5, H
	cpu.bit(&cpu.reg.H, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb6D() int { // BIT 5, L
	cpu.bit(&cpu.reg.L, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb6E() int { // BIT 5, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.bit(&val, 5)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb6F() int { // BIT 5, A
	cpu.bit(&cpu.reg.A, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb70() int { // BIT 6, B
	cpu.bit(&cpu.reg.B, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb71() int { // BIT 6, C
	cpu.bit(&cpu.reg.C, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb72() int { // BIT 6, D
	cpu.bit(&cpu.reg.D, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb73() int { // BIT 6, E
	cpu.bit(&cpu.reg.E, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb74() int { // BIT 6, H
	cpu.bit(&cpu.reg.H, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb75() int { // BIT 6, L
	cpu.bit(&cpu.reg.L, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb76() int { // BIT 6, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.bit(&val, 6)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb77() int { // BIT 6, A
	cpu.bit(&cpu.reg.A, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb78() int { // BIT 7, B
	cpu.bit(&cpu.reg.B, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb79() int { // BIT 7, C
	cpu.bit(&cpu.reg.C, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb7A() int { // BIT 7, D
	cpu.bit(&cpu.reg.D, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb7B() int { // BIT 7, E
	cpu.bit(&cpu.reg.E, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb7C() int { // BIT 7, H
	cpu.bit(&cpu.reg.H, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb7D() int { // BIT 7, L
	cpu.bit(&cpu.reg.L, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb7E() int { // BIT 7, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.bit(&val, 7)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb7F() int { // BIT 7, A
	cpu.bit(&cpu.reg.A, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb80() int { // RES 0, B
	cpu.res(&cpu.reg.B, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb81() int { // RES 0, C
	cpu.res(&cpu.reg.C, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb82() int { // RES 0, D
	cpu.res(&cpu.reg.D, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb83() int { // RES 0, E
	cpu.res(&cpu.reg.E, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb84() int { // RES 0, H
	cpu.res(&cpu.reg.H, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb85() int { // RES 0, L
	cpu.res(&cpu.reg.L, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb86() int { // RES 0, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.res(&val, 0)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb87() int { // RES 0, A
	cpu.res(&cpu.reg.A, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb88() int { // RES 1, B
	cpu.res(&cpu.reg.B, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb89() int { // RES 1, C
	cpu.res(&cpu.reg.C, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb8A() int { // RES 1, D
	cpu.res(&cpu.reg.D, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb8B() int { // RES 1, E
	cpu.res(&cpu.reg.E, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb8C() int { // RES 1, H
	cpu.res(&cpu.reg.H, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb8D() int { // RES 1, L
	cpu.res(&cpu.reg.L, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb8E() int { // RES 1, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.res(&val, 1)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb8F() int { // RES 1, A
	cpu.res(&cpu.reg.A, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb90() int { // RES 2, B
	cpu.res(&cpu.reg.B, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb91() int { // RES 2, C
	cpu.res(&cpu.reg.C, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb92() int { // RES 2, D
	cpu.res(&cpu.reg.D, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb93() int { // RES 2, E
	cpu.res(&cpu.reg.E, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb94() int { // RES 2, H
	cpu.res(&cpu.reg.H, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb95() int { // RES 2, L
	cpu.res(&cpu.reg.L, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb96() int { // RES 2, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.res(&val, 2)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb97() int { // RES 2, A
	cpu.res(&cpu.reg.A, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb98() int { // RES 3, B
	cpu.res(&cpu.reg.B, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb99() int { // RES 3, C
	cpu.res(&cpu.reg.C, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb9A() int { // RES 3, D
	cpu.res(&cpu.reg.D, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb9B() int { // RES 3, E
	cpu.res(&cpu.reg.E, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb9C() int { // RES 3, H
	cpu.res(&cpu.reg.H, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb9D() int { // RES 3, L
	cpu.res(&cpu.reg.L, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb9E() int { // RES 3, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.res(&val, 3)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cb9F() int { // RES 3, A
	cpu.res(&cpu.reg.A, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbA0() int { // RES 4, B
	cpu.res(&cpu.reg.B, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbA1() int { // RES 4, C
	cpu.res(&cpu.reg.C, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbA2() int { // RES 4, D
	cpu.res(&cpu.reg.D, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbA3() int { // RES 4, E
	cpu.res(&cpu.reg.E, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbA4() int { // RES 4, H
	cpu.res(&cpu.reg.H, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbA5() int { // RES 4, L
	cpu.res(&cpu.reg.L, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbA6() int { // RES 4, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.res(&val, 4)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbA7() int { // RES 4, A
	cpu.res(&cpu.reg.A, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbA8() int { // RES 5, B
	cpu.res(&cpu.reg.B, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbA9() int { // RES 5, C
	cpu.res(&cpu.reg.C, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbAA() int { // RES 5, D
	cpu.res(&cpu.reg.D, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbAB() int { // RES 5, E
	cpu.res(&cpu.reg.E, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbAC() int { // RES 5, H
	cpu.res(&cpu.reg.H, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbAD() int { // RES 5, L
	cpu.res(&cpu.reg.L, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbAE() int { // RES 5, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.res(&val, 5)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbAF() int { // RES 5, A
	cpu.res(&cpu.reg.A, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbB0() int { // RES 6, B
	cpu.res(&cpu.reg.B, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbB1() int { // RES 6, C
	cpu.res(&cpu.reg.C, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbB2() int { // RES 6, D
	cpu.res(&cpu.reg.D, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbB3() int { // RES 6, E
	cpu.res(&cpu.reg.E, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbB4() int { // RES 6, H
	cpu.res(&cpu.reg.H, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbB5() int { // RES 6, L
	cpu.res(&cpu.reg.L, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbB6() int { // RES 6, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.res(&val, 6)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbB7() int { // RES 6, A
	cpu.res(&cpu.reg.A, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbB8() int { // RES 7, B
	cpu.res(&cpu.reg.B, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbB9() int { // RES 7, C
	cpu.res(&cpu.reg.C, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbBA() int { // RES 7, D
	cpu.res(&cpu.reg.D, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbBB() int { // RES 7, E
	cpu.res(&cpu.reg.E, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbBC() int { // RES 7, H
	cpu.res(&cpu.reg.H, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbBD() int { // RES 7, L
	cpu.res(&cpu.reg.L, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbBE() int { // RES 7, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.res(&val, 7)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbBF() int { // RES 7, A
	cpu.res(&cpu.reg.A, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbC0() int { // SET 0, B
	cpu.set(&cpu.reg.B, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbC1() int { // SET 0, C
	cpu.set(&cpu.reg.C, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbC2() int { // SET 0, D
	cpu.set(&cpu.reg.D, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbC3() int { // SET 0, E
	cpu.set(&cpu.reg.E, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbC4() int { // SET 0, H
	cpu.set(&cpu.reg.H, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbC5() int { // SET 0, L
	cpu.set(&cpu.reg.L, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbC6() int { // SET 0, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.set(&val, 0)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbC7() int { // SET 0, A
	cpu.set(&cpu.reg.A, 0)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbC8() int { // SET 1, B
	cpu.set(&cpu.reg.B, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbC9() int { // SET 1, C
	cpu.set(&cpu.reg.C, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbCA() int { // SET 1, D
	cpu.set(&cpu.reg.D, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbCB() int { // SET 1, E
	cpu.set(&cpu.reg.E, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbCC() int { // SET 1, H
	cpu.set(&cpu.reg.H, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbCD() int { // SET 1, L
	cpu.set(&cpu.reg.L, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbCE() int { // SET 1, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.set(&val, 1)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbCF() int { // SET 1, A
	cpu.set(&cpu.reg.A, 1)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbD0() int { // SET 2, B
	cpu.set(&cpu.reg.B, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbD1() int { // SET 2, C
	cpu.set(&cpu.reg.C, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbD2() int { // SET 2, D
	cpu.set(&cpu.reg.D, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbD3() int { // SET 2, E
	cpu.set(&cpu.reg.E, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbD4() int { // SET 2, H
	cpu.set(&cpu.reg.H, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbD5() int { // SET 2, L
	cpu.set(&cpu.reg.L, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbD6() int { // SET 2, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.set(&val, 2)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbD7() int { // SET 2, A
	cpu.set(&cpu.reg.A, 2)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbD8() int { // SET 3, B
	cpu.set(&cpu.reg.B, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbD9() int { // SET 3, C
	cpu.set(&cpu.reg.C, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbDA() int { // SET 3, D
	cpu.set(&cpu.reg.D, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbDB() int { // SET 3, E
	cpu.set(&cpu.reg.E, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbDC() int { // SET 3, H
	cpu.set(&cpu.reg.H, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbDD() int { // SET 3, L
	cpu.set(&cpu.reg.L, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbDE() int { // SET 3, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.set(&val, 3)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbDF() int { // SET 3, A
	cpu.set(&cpu.reg.A, 3)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbE0() int { // SET 4, B
	cpu.set(&cpu.reg.B, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbE1() int { // SET 4, C
	cpu.set(&cpu.reg.C, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbE2() int { // SET 4, D
	cpu.set(&cpu.reg.D, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbE3() int { // SET 4, E
	cpu.set(&cpu.reg.E, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbE4() int { // SET 4, H
	cpu.set(&cpu.reg.H, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbE5() int { // SET 4, L
	cpu.set(&cpu.reg.L, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbE6() int { // SET 4, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.set(&val, 4)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbE7() int { // SET 4, A
	cpu.set(&cpu.reg.A, 4)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbE8() int { // SET 5, B
	cpu.set(&cpu.reg.B, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbE9() int { // SET 5, C
	cpu.set(&cpu.reg.C, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbEA() int { // SET 5, D
	cpu.set(&cpu.reg.D, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbEB() int { // SET 5, E
	cpu.set(&cpu.reg.E, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbEC() int { // SET 5, H
	cpu.set(&cpu.reg.H, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbED() int { // SET 5, L
	cpu.set(&cpu.reg.L, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbEE() int { // SET 5, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.set(&val, 5)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbEF() int { // SET 5, A
	cpu.set(&cpu.reg.A, 5)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbF0() int { // SET 6, B
	cpu.set(&cpu.reg.B, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbF1() int { // SET 6, C
	cpu.set(&cpu.reg.C, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbF2() int { // SET 6, D
	cpu.set(&cpu.reg.D, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbF3() int { // SET 6, E
	cpu.set(&cpu.reg.E, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbF4() int { // SET 6, H
	cpu.set(&cpu.reg.H, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbF5() int { // SET 6, L
	cpu.set(&cpu.reg.L, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbF6() int { // SET 6, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.set(&val, 6)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbF7() int { // SET 6, A
	cpu.set(&cpu.reg.A, 6)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbF8() int { // SET 7, B
	cpu.set(&cpu.reg.B, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbF9() int { // SET 7, C
	cpu.set(&cpu.reg.C, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbFA() int { // SET 7, D
	cpu.set(&cpu.reg.D, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbFB() int { // SET 7, E
	cpu.set(&cpu.reg.E, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbFC() int { // SET 7, H
	cpu.set(&cpu.reg.H, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbFD() int { // SET 7, L
	cpu.set(&cpu.reg.L, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbFE() int { // SET 7, (HL)
	address := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
	val := Read(address)
	cpu.set(&val, 7)
	Write(address, val)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) cbFF() int { // SET 7, A
	cpu.set(&cpu.reg.A, 7)
	cpu.reg.PC += 2
	return 2
}

func (cpu *CPU) interrupt() int { // handle interrupts
	// check if interrupt occurred
	// loop through every bit in 0xFF0F until we find one
	for i := 0; i < 5; i++ {
		if Read(0xFF0F)&(0x01<<i) > 0 {
			// check if found int is enabled in 0xFFFF
			if Read(0xFFFF)&(0x01<<i) > 0 {
				// reset IME and the interrupt bit
				cpu.flg.IME = false
				Write(0xFFFF, Read(0xFFFF)&((0x01<<i)^0xFF))
				// originally, rst(byte) was just for the RST instruction
				// however, it allows easy calling of a specific address
				// and pushing the current PC to stack already
				// so I won't write the same code here
				cpu.rst(cpu.interrupts[i])
				// cpu.ret(true)
			}
		}
	}
	return 20 + int(FlagToBit(cpu.flg.HALT)*4) // according to "The Cycle-Accurate GB" doc, "It takes 20 clocks to dispatch an interrupt. If CPU is in HALT mode, another extra 4 clocks are needed"
}
