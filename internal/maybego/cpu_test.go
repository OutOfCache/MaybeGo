package maybego

import (
	//	"fmt"
	"testing"
)

var cpu *CPU = NewCPU()

func TestFlagsToBytes(t *testing.T) {
	var tests = []struct {
		z            bool
		n            bool
		h            bool
		c            bool
		expectedByte byte
	}{
		{false, false, false, false, 0b00000000},
		{false, true, false, true, 0b01010000},
		{true, false, true, false, 0b10100000},
		{true, true, true, true, 0b11110000},
	}

	for _, test := range tests {
		cpu.flg.Z = test.z
		cpu.flg.N = test.n
		cpu.flg.H = test.h
		cpu.flg.C = test.c

		got := cpu.FlagsToBytes()
		if got != test.expectedByte {
			t.Errorf("Current byte %b; expected: %b", got, test.expectedByte)
		}
	}
}

func TestBytesToFlags(t *testing.T) {
	var tests = []struct {
		Byte      byte
		expectedZ bool
		expectedN bool
		expectedH bool
		expectedC bool
	}{
		{0b00000000, false, false, false, false},
		{0b01010000, false, true, false, true},
		{0b10100000, true, false, true, false},
		{0b11110000, true, true, true, true},
	}

	for _, test := range tests {
		cpu.BytesToFlags(test.Byte)

		actualZ := cpu.flg.Z
		actualN := cpu.flg.N
		actualH := cpu.flg.H
		actualC := cpu.flg.C

		if actualZ != test.expectedZ {
			t.Errorf("Current Z %t; expected: %t", actualZ, test.expectedZ)
		}
		if actualN != test.expectedN {
			t.Errorf("Current N %t; expected: %t", actualN, test.expectedN)
		}
		if actualH != test.expectedH {
			t.Errorf("Current H %t; expected: %t", actualH, test.expectedH)
		}
		if actualC != test.expectedC {
			t.Errorf("Current C %t; expected: %t", actualC, test.expectedC)
		}
	}
}

var registers8 = []struct {
	name string
	reg  *byte
}{
	{"B", &cpu.reg.B},
	{"C", &cpu.reg.C},
	{"D", &cpu.reg.D},
	{"E", &cpu.reg.E},
	{"H", &cpu.reg.H},
	{"L", &cpu.reg.L},
	{"A", &cpu.reg.A},
}

var registers16 = []struct {
	name string
	hi   *byte
	lo   *byte
}{
	{"BC", &cpu.reg.B, &cpu.reg.C},
	{"DE", &cpu.reg.D, &cpu.reg.E},
	{"HL", &cpu.reg.H, &cpu.reg.L},
}

func TestLD8(t *testing.T) { // {{{
	var tests = []struct {
		dest byte
		src  byte
	}{
		{0x73, 0x67},
		{0x38, 0xA9},
		{0xBA, 0xF0},
		{0xFF, 0x3F},
	}

	var commands = []struct {
		name  string
		instr [8]func() byte
	}{
		{"LD B, u8/r8", [8]func() byte{
			cpu.cpu06, cpu.cpu40, cpu.cpu41, cpu.cpu42,
			cpu.cpu43, cpu.cpu44, cpu.cpu45, cpu.cpu47},
		},
		{"LD C, u8/r8", [8]func() byte{
			cpu.cpu0E, cpu.cpu48, cpu.cpu49, cpu.cpu4A,
			cpu.cpu4B, cpu.cpu4C, cpu.cpu4D, cpu.cpu4F},
		},
		{"LD D, u8/r8", [8]func() byte{
			cpu.cpu16, cpu.cpu50, cpu.cpu51, cpu.cpu52,
			cpu.cpu53, cpu.cpu54, cpu.cpu55, cpu.cpu57},
		},
		{"LD E, u8/r8", [8]func() byte{
			cpu.cpu1E, cpu.cpu58, cpu.cpu59, cpu.cpu5A,
			cpu.cpu5B, cpu.cpu5C, cpu.cpu5D, cpu.cpu5F},
		},
		{"LD H, u8/r8", [8]func() byte{
			cpu.cpu26, cpu.cpu60, cpu.cpu61, cpu.cpu62,
			cpu.cpu63, cpu.cpu64, cpu.cpu65, cpu.cpu67},
		},
		{"LD L, u8/r8", [8]func() byte{
			cpu.cpu2E, cpu.cpu68, cpu.cpu69, cpu.cpu6A,
			cpu.cpu6B, cpu.cpu6C, cpu.cpu6D, cpu.cpu6F},
		},
		{"LD A, u8/r8", [8]func() byte{
			cpu.cpu3E, cpu.cpu78, cpu.cpu79, cpu.cpu7A,
			cpu.cpu7B, cpu.cpu7C, cpu.cpu7D, cpu.cpu7F},
		},
	}

	cpu.reg.PC = 0x66

	for dest_idx, destination := range registers8 {
		t.Run("LD "+destination.name+", u8", func(t *testing.T) {
			for _, test := range tests {
				*destination.reg = test.dest
				Write(cpu.reg.PC+1, test.src)
				expected_PC := cpu.reg.PC + 2

				cycles := commands[dest_idx].instr[0]()

				if *destination.reg != test.src {
					t.Errorf("Current %s: %x, expected: %x", destination.name, destination.reg, test.src)
				}
				if cycles != 2 {
					t.Errorf("Got %d cycles, expected 2", cycles)
				}
				if cpu.reg.PC != expected_PC {
					t.Errorf("Got PC: %x, expected %x", cpu.reg.PC, expected_PC)
				}
			}
		})

		for src_idx, source := range registers8 {
			t.Run("LD "+destination.name+", "+source.name, func(t *testing.T) {
				for _, test := range tests {
					*destination.reg = test.dest
					*source.reg = test.src
					expected_PC := cpu.reg.PC + 1

					cycles := commands[dest_idx].instr[src_idx+1]()

					if *destination.reg != test.src {
						t.Errorf("Current %s: %x, expected: %x", destination.name, destination.reg, test.src)
					}
					if cycles != 1 {
						t.Errorf("Got %d cycles, expected 1", cycles)
					}
					if cpu.reg.PC != expected_PC {
						t.Errorf("Got PC: %x, expected %x", cpu.reg.PC, expected_PC)
					}

				}
			})
		}
	}
} // }}}
func TestLD16(t *testing.T) { // {{{
	var tests = []struct {
		dest_hi byte
		dest_lo byte
		src_hi  byte
		src_lo  byte
	}{
		{0xDE, 0xAD, 0xBE, 0xEF},
		{0xCA, 0x55, 0xE7, 0x7E},
	}

	var commands = []struct {
		name   string
		instr  func() byte
		cycles byte
	}{
		{"LD BC, u16", cpu.cpu01, 3},
		{"LD DE, u16", cpu.cpu11, 3},
		{"LD HL, u16", cpu.cpu21, 3},
	}
	cpu.reg.PC = 0x66

	for index, r16 := range registers16 {
		t.Run(commands[index].name, func(t *testing.T) {
			for _, test := range tests {
				*r16.hi = test.dest_hi
				*r16.lo = test.dest_lo

				Write(cpu.reg.PC+1, test.src_lo)
				Write(cpu.reg.PC+2, test.src_hi)

				cycles := commands[index].instr()

				if *r16.hi != test.src_hi {
					t.Errorf("Current %s: %x, expected: %x", r16.name, r16.hi, test.src_hi)
				}
				if *r16.lo != test.src_lo {
					t.Errorf("Current %s: %x, expected: %x", r16.name, r16.lo, test.src_lo)
				}
				if cycles != commands[index].cycles {
					t.Errorf("Got %d cycles, expected %d", cycles, commands[index].cycles)
				}
			}
		})
	}
} // }}}
func TestLDToAdr(t *testing.T) { // {{{
	var tests = []struct {
		dest_hi byte
		dest_lo byte
		src     byte
	}{
		{0xDE, 0xAD, 0xBE},
		{0xCA, 0x55, 0xE7},
	}

	cpu.reg.PC = 0x66

	t.Run("To R16 from A", func(t *testing.T) { // {{{
		var commands = []struct {
			name   string
			instr  func() byte
			cycles byte
		}{
			{"LD (BC), A", cpu.cpu02, 2},
			{"LD (DE), A", cpu.cpu12, 2},
		}
		for index, command := range commands {
			t.Run(commands[index].name, func(t *testing.T) {
				r16 := registers16[index]
				for _, test := range tests {
					adress := (uint16(test.dest_hi) << 8) + uint16(test.dest_lo)
					Write(adress, 0x00)

					*r16.hi = test.dest_hi
					*r16.lo = test.dest_lo
					cpu.reg.A = test.src

					cycles := commands[index].instr()
					actual_byte := Read(adress)

					if actual_byte != test.src {
						t.Errorf("Current byte at %x: %x, expected: %x", adress, actual_byte, test.src)
					}
					if cycles != command.cycles {
						t.Errorf("Got %d cycles, expected %d", cycles, commands[index].cycles)
					}
				}
			})
		}
	}) // }}}

	t.Run("To (HL)", func(t *testing.T) {
		var commands = [8]struct {
			name   string
			instr  func() byte
			cycles byte
		}{
			{"LD (HL), u8", cpu.cpu36, 3},
			{"LD (HL), B", cpu.cpu70, 2},
			{"LD (HL), C", cpu.cpu71, 2},
			{"LD (HL), D", cpu.cpu72, 2},
			{"LD (HL), E", cpu.cpu73, 2},
			{"LD (HL), H", cpu.cpu74, 2},
			{"LD (HL), L", cpu.cpu75, 2},
			{"LD (HL), A", cpu.cpu77, 2},
		}

		for _, test := range tests {

			cpu.reg.H = test.dest_hi
			cpu.reg.L = test.dest_lo
			adress := (uint16(test.dest_hi) << 8) + uint16(test.dest_lo)

			t.Run(commands[0].name, func(t *testing.T) {
				Write(cpu.reg.PC+1, test.src)
				expected_cycles := commands[0].cycles

				actual_cycles := commands[0].instr()
				actual_byte := Read(adress)

				if actual_byte != test.src {
					t.Errorf("Got %x at adress %x, expected %x", actual_byte, adress, test.src)
				}
				if actual_cycles != expected_cycles {
					t.Errorf("Got %d cycles, expected %d", actual_cycles, expected_cycles)
				}
			})

			for src_idx, source_register := range registers8 {
				command := commands[src_idx+1]
				t.Run(command.name, func(t *testing.T) {
					registers8[src_idx] = test.src
					expected_cycles := command.cycles

					actual_cycles := command.instr()
					actual_byte := Read(adress)

					if actual_byte != test.src {
						t.Errorf("Got %x at adress %x, expected %x", actual_byte, adress, test.src)
					}
					if actual_cycles != expected_cycles {
						t.Errorf("Got %d cycles, expected %d", actual_cycles, expected_cycles)
					}
				})
			}
		}
	})
} // }}}

func TestCpu00(t *testing.T) {
	var tests = []struct {
		pc       uint16
		expected uint16
	}{
		{0x0, 0x1},
		{0x2F, 0x30},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.cpu00()
		if cpu.reg.PC != test.expected {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expected)
		}
	}
}

func TestCpu01(t *testing.T) {
	var tests = []struct {
		pc         uint16
		expectedB  byte
		expectedC  byte
		expectedPC uint16
	}{
		{0x9432, 0x13, 0x7F, 0x9435},
		{0x2F3C, 0x30, 0x49, 0x2F3F},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		Write(cpu.reg.PC+1, test.expectedC)
		Write(cpu.reg.PC+2, test.expectedB)
		cpu.cpu01()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if cpu.reg.B != test.expectedB {
			t.Errorf("Current B: %x; expected: %x", cpu.reg.B, test.expectedB)
		}
		if cpu.reg.C != test.expectedC {
			t.Errorf("Current C: %x; expected: %x", cpu.reg.C, test.expectedC)
		}
	}
}
func TestCpu02(t *testing.T) {
	var tests = []struct {
		pc         uint16
		B          byte
		C          byte
		address    uint16
		expectedPC uint16
		A          byte
	}{
		{0x9432, 0x13, 0x7F, 0x137F, 0x9433, 0x35},
		{0x2F3C, 0x30, 0x49, 0x3049, 0x2F3D, 0xF3},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.A = test.A
		cpu.reg.B = test.B
		cpu.reg.C = test.C
		cpu.cpu02()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if Read(test.address) != test.A {
			t.Errorf("Current [BC]: %x; expected: %x", Read(test.address), test.A)
		}
	}
}

func TestCpu07(t *testing.T) {
	var tests = []struct {
		pc         uint16
		a          byte
		carry      bool
		expectedCF bool
		expectedA  byte
		expectedPC uint16
	}{
		{0x1234, 0x80, false, true, 0x01, 0x1235},
		{0x1234, 0x80, true, true, 0x01, 0x1235},
		{0x63F8, 0x35, false, false, 0x6A, 0x63F9},
		{0x63F8, 0x35, true, false, 0x6A, 0x63F9},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.A = test.a
		cpu.flg.C = test.carry
		cpu.cpu07()
		if cpu.flg.C != test.expectedCF {
			t.Errorf("Carry: %t, expected: %t", cpu.flg.C, test.expectedCF)
		}
		if cpu.reg.A != test.expectedA {
			t.Errorf("Current A: %x, expected %x", cpu.reg.A, test.expectedA)
		}
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpu08(t *testing.T) {
	var tests = []struct {
		pc         uint16
		sp         uint16
		splo       byte
		sphi       byte
		lo         byte
		hi         byte
		address    uint16
		expectedPC uint16
	}{
		{0x1234, 0x385E, 0x5E, 0x38, 0x7D, 0x89, 0x897D, 0x1237},
		{0x63F8, 0x3582, 0x82, 0x35, 0x6A, 0x12, 0x126A, 0x63FB},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.SP = test.sp

		Write(cpu.reg.PC+1, test.lo)
		Write(cpu.reg.PC+2, test.hi)
		cpu.cpu08()
		if Read(test.address) != test.splo {
			t.Errorf("At Address: %x, expected: %x", Read(test.address), test.lo)
		}
		if Read(test.address+1) != test.sphi {
			t.Errorf("At Address + 1: %x, expected %x", Read(test.address+1), test.hi)
		}
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpu09(t *testing.T) {
	var tests = []struct {
		pc         uint16
		h          byte
		l          byte
		b          byte
		c          byte
		expectedH  byte
		expectedL  byte
		expectedHF bool
		expectedC  bool
		expectedPC uint16
	}{
		{0x1234, 0x5E, 0x38, 0x7D, 0x89, 0xDB, 0xC1, true, false, 0x1235},
		{0x63F8, 0x82, 0x35, 0x6A, 0x12, 0xEC, 0x47, false, false, 0x63F9},
		{0x63F8, 0x82, 0x35, 0x8A, 0x12, 0x0C, 0x47, false, true, 0x63F9},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.H = test.h
		cpu.reg.L = test.l
		cpu.reg.B = test.b
		cpu.reg.C = test.c

		cpu.cpu09()
		if cpu.reg.H != test.expectedH {
			t.Errorf("H: %x, expected: %x", cpu.reg.H, test.expectedH)
		}
		if cpu.reg.L != test.expectedL {
			t.Errorf("L: %x, expected: %x", cpu.reg.L, test.expectedL)
		}
		if cpu.flg.H != test.expectedHF {
			t.Errorf("H Flag: %t, expected %t", cpu.flg.H, test.expectedHF)
		}
		if cpu.flg.C != test.expectedC {
			t.Errorf("C Flag: %t, expected %t", cpu.flg.C, test.expectedC)
		}
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpu0A(t *testing.T) {
	var tests = []struct {
		pc         uint16
		B          byte
		C          byte
		address    uint16
		val        byte
		expectedPC uint16
		A          byte
	}{
		{0x9432, 0x13, 0x7F, 0x137F, 0x8E, 0x9433, 0x35},
		{0x2F3C, 0x30, 0x49, 0x3049, 0x3C, 0x2F3D, 0xF3},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.B = test.B
		cpu.reg.C = test.C
		Write(test.address, test.val)
		cpu.cpu0A()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if cpu.reg.A != test.val {
			t.Errorf("Current A: %x; expected: %x", cpu.reg.A, test.val)
		}
	}
}

func TestCpu0F(t *testing.T) {
	var tests = []struct {
		pc         uint16
		a          byte
		carry      bool
		expectedCF bool
		expectedA  byte
		expectedPC uint16
	}{
		{0x1234, 0x01, false, true, 0x80, 0x1235},
		{0x1234, 0x01, true, true, 0x80, 0x1235},
		{0x63F8, 0x32, false, false, 0x19, 0x63F9},
		{0x63F8, 0x32, true, false, 0x19, 0x63F9},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.A = test.a
		cpu.flg.C = test.carry
		cpu.cpu0F()
		if cpu.flg.C != test.expectedCF {
			t.Errorf("Carry: %t, expected: %t", cpu.flg.C, test.expectedCF)
		}
		if cpu.reg.A != test.expectedA {
			t.Errorf("Current A: %x, expected %x", cpu.reg.A, test.expectedA)
		}
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpu11(t *testing.T) {
	var tests = []struct {
		pc         uint16
		expectedD  byte
		expectedE  byte
		expectedPC uint16
	}{
		{0x9432, 0x13, 0x7F, 0x9435},
		{0x2F3C, 0x30, 0x49, 0x2F3F},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		Write(cpu.reg.PC+1, test.expectedE)
		Write(cpu.reg.PC+2, test.expectedD)
		cpu.cpu11()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if cpu.reg.D != test.expectedD {
			t.Errorf("Current B: %x; expected: %x", cpu.reg.D, test.expectedD)
		}
		if cpu.reg.E != test.expectedE {
			t.Errorf("Current C: %x; expected: %x", cpu.reg.E, test.expectedE)
		}
	}
}
func TestCpu12(t *testing.T) {
	var tests = []struct {
		pc         uint16
		D          byte
		E          byte
		address    uint16
		expectedPC uint16
		A          byte
	}{
		{0x9432, 0x13, 0x7F, 0x137F, 0x9433, 0x35},
		{0x2F3C, 0x30, 0x49, 0x3049, 0x2F3D, 0xF3},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.A = test.A
		cpu.reg.D = test.D
		cpu.reg.E = test.E
		cpu.cpu12()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if Read(test.address) != test.A {
			t.Errorf("Current [BC]: %x; expected: %x", Read(test.address), test.A)
		}
	}
}

func TestCpu17(t *testing.T) {
	var tests = []struct {
		pc         uint16
		a          byte
		carry      bool
		expectedCF bool
		expectedA  byte
		expectedPC uint16
	}{
		{0x1234, 0x80, false, true, 0x00, 0x1235},
		{0x1234, 0x80, true, true, 0x01, 0x1235},
		{0x63F8, 0x35, false, false, 0x6A, 0x63F9},
		{0x63F8, 0x35, true, false, 0x6B, 0x63F9},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.A = test.a
		cpu.flg.C = test.carry
		cpu.cpu17()
		if cpu.flg.C != test.expectedCF {
			t.Errorf("Carry: %t, expected: %t", cpu.flg.C, test.expectedCF)
		}
		if cpu.reg.A != test.expectedA {
			t.Errorf("Current A: %x, expected %x", cpu.reg.A, test.expectedA)
		}
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpu18(t *testing.T) {
	var tests = []struct {
		pc         uint16
		i8         int8
		expectedPC uint16
	}{
		{0x1234, 0x00, 0x1236},
		{0x1234, 0x08, 0x123E},
		{0x1234, -8, 0x122E},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		Write(cpu.reg.PC+1, byte(test.i8))

		cpu.cpu18()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpu19(t *testing.T) {
	var tests = []struct {
		pc         uint16
		h          byte
		l          byte
		d          byte
		e          byte
		expectedH  byte
		expectedL  byte
		expectedHF bool
		expectedC  bool
		expectedPC uint16
	}{
		{0x1234, 0x5E, 0x38, 0x7D, 0x89, 0xDB, 0xC1, true, false, 0x1235},
		{0x63F8, 0x82, 0x35, 0x6A, 0x12, 0xEC, 0x47, false, false, 0x63F9},
		{0x63F8, 0x82, 0x35, 0x8A, 0x12, 0x0C, 0x47, false, true, 0x63F9},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.H = test.h
		cpu.reg.L = test.l
		cpu.reg.D = test.d
		cpu.reg.E = test.e

		cpu.cpu19()
		if cpu.reg.H != test.expectedH {
			t.Errorf("H: %x, expected: %x", cpu.reg.H, test.expectedH)
		}
		if cpu.reg.L != test.expectedL {
			t.Errorf("L: %x, expected: %x", cpu.reg.L, test.expectedL)
		}
		if cpu.flg.H != test.expectedHF {
			t.Errorf("H Flag: %t, expected %t", cpu.flg.H, test.expectedHF)
		}
		if cpu.flg.C != test.expectedC {
			t.Errorf("C Flag: %t, expected %t", cpu.flg.C, test.expectedC)
		}
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpu1A(t *testing.T) {
	var tests = []struct {
		pc         uint16
		D          byte
		E          byte
		address    uint16
		val        byte
		expectedPC uint16
		A          byte
	}{
		{0x9432, 0x13, 0x7F, 0x137F, 0x8E, 0x9433, 0x35},
		{0x2F3C, 0x30, 0x49, 0x3049, 0x3C, 0x2F3D, 0xF3},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.D = test.D
		cpu.reg.E = test.E
		Write(test.address, test.val)
		cpu.cpu1A()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if cpu.reg.A != test.val {
			t.Errorf("Current A: %x; expected: %x", cpu.reg.A, test.val)
		}
	}
}

func TestCpu1F(t *testing.T) {
	var tests = []struct {
		pc         uint16
		a          byte
		carry      bool
		expectedCF bool
		expectedA  byte
		expectedPC uint16
	}{
		{0x1234, 0x80, false, false, 0x40, 0x1235},
		{0x1234, 0x80, true, false, 0xC0, 0x1235},
		{0x63F8, 0x35, false, true, 0x1A, 0x63F9},
		{0x63F8, 0x35, true, true, 0x9A, 0x63F9},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.A = test.a
		cpu.flg.C = test.carry
		cpu.cpu1F()
		if cpu.flg.C != test.expectedCF {
			t.Errorf("Carry: %t, expected: %t", cpu.flg.C, test.expectedCF)
		}
		if cpu.reg.A != test.expectedA {
			t.Errorf("Current A: %x, expected %x", cpu.reg.A, test.expectedA)
		}
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpu20(t *testing.T) {
	var tests = []struct {
		pc         uint16
		Z          bool
		i8         int8
		expectedPC uint16
	}{
		{0x1234, false, 0x00, 0x1236},
		{0x1234, false, 0x08, 0x123E},
		{0x1234, false, -8, 0x122E},
		{0x1234, true, 0x00, 0x1236},
		{0x1234, true, 0x08, 0x1236},
		{0x1234, true, -8, 0x1236},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		Write(cpu.reg.PC+1, byte(test.i8))
		cpu.flg.Z = test.Z

		cpu.cpu20()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpu21(t *testing.T) {
	var tests = []struct {
		pc         uint16
		expectedH  byte
		expectedL  byte
		expectedPC uint16
	}{
		{0x9432, 0x13, 0x7F, 0x9435},
		{0x2F3C, 0x30, 0x49, 0x2F3F},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		Write(cpu.reg.PC+1, test.expectedL)
		Write(cpu.reg.PC+2, test.expectedH)
		cpu.cpu21()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if cpu.reg.H != test.expectedH {
			t.Errorf("Current H: %x; expected: %x", cpu.reg.H, test.expectedH)
		}
		if cpu.reg.L != test.expectedL {
			t.Errorf("Current L: %x; expected: %x", cpu.reg.L, test.expectedL)
		}
	}
}

func TestCpu22(t *testing.T) {
	var tests = []struct {
		pc         uint16
		H          byte
		L          byte
		address    uint16
		expectedPC uint16
		A          byte
	}{
		{0x9432, 0x13, 0x7F, 0x137F, 0x9433, 0x35},
		{0x2F3C, 0x30, 0x49, 0x3049, 0x2F3D, 0xF3},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.A = test.A
		cpu.reg.H = test.H
		cpu.reg.L = test.L
		cpu.cpu22()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if Read(test.address) != test.A {
			t.Errorf("Current [BC]: %x; expected: %x", Read(test.address), test.A)
		}
		hl := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
		if hl != test.address+1 {
			t.Errorf("Current HL: %x, expected: %x", hl, test.address+1)
		}
	}
}

func TestCpu27(t *testing.T) {
	var tests = []struct {
		pc         uint16
		A          byte
		negative   bool
		halfcarry  bool
		carry      bool
		expectedPC uint16
		expectedA  byte
		expectedCF bool
	}{
		{0x9432, 0x6A, false, false, false, 0x9433, 0x70, false},
		{0x9432, 0x9A, false, false, false, 0x9433, 0x00, true},
		{0x2F3C, 0xFA, false, false, false, 0x2F3D, 0x60, true},
		{0x9432, 0x6A, false, false, true, 0x9433, 0xd0, true},
		{0x9432, 0x9A, false, false, true, 0x9433, 0x00, true},
		{0x2F3C, 0xFA, false, false, true, 0x2F3D, 0x60, true},
		{0x9432, 0x6A, false, true, false, 0x9433, 0x70, false},
		{0x9432, 0x9A, false, true, false, 0x9433, 0x00, true},
		{0x2F3C, 0xFA, false, true, false, 0x2F3D, 0x60, true},
		{0x9432, 0x29, false, true, true, 0x9433, 0x8F, true},
		{0x9432, 0x9A, false, true, true, 0x9433, 0x00, true},
		{0x2F3C, 0xFA, false, true, true, 0x2F3D, 0x60, true},
		{0x9432, 0x6A, true, false, false, 0x9433, 0x6A, false},
		{0x9432, 0x9A, true, false, false, 0x9433, 0x9A, false},
		{0x2F3C, 0xFA, true, false, false, 0x2F3D, 0xFA, false},
		{0x9432, 0x6A, true, false, true, 0x9433, 0x0A, true},
		{0x9432, 0x9A, true, false, true, 0x9433, 0x3A, true},
		{0x2F3C, 0xFA, true, false, true, 0x2F3D, 0x9A, true},
		{0x9432, 0x6A, true, true, false, 0x9433, 0x64, false},
		{0x9432, 0x9A, true, true, false, 0x9433, 0x94, false},
		{0x2F3C, 0xFA, true, true, false, 0x2F3D, 0xF4, false},
		{0x9432, 0x29, true, true, true, 0x9433, 0xC3, true},
		{0x9432, 0x9A, true, true, true, 0x9433, 0x34, true},
		{0x2F3C, 0xFA, true, true, true, 0x2F3D, 0x94, true},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.A = test.A
		cpu.flg.N = test.negative
		cpu.flg.H = test.halfcarry
		cpu.flg.C = test.carry
		cpu.cpu27()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if cpu.reg.A != test.expectedA {
			t.Errorf("Current A: %x; expected: %x", cpu.reg.A, test.expectedA)
		}
		if cpu.flg.C != test.expectedCF {
			t.Errorf("Current CF: %t, expected: %t", cpu.flg.C, test.expectedCF)
		}
	}
}

func TestCpu28(t *testing.T) {
	var tests = []struct {
		pc         uint16
		Z          bool
		i8         int8
		expectedPC uint16
	}{
		{0x1234, true, 0x00, 0x1236},
		{0x1234, true, 0x08, 0x123E},
		{0x1234, true, -8, 0x122E},
		{0x1234, false, 0x00, 0x1236},
		{0x1234, false, 0x08, 0x1236},
		{0x1234, false, -8, 0x1236},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		Write(cpu.reg.PC+1, byte(test.i8))
		cpu.flg.Z = test.Z

		cpu.cpu28()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}
}

func TestCpu29(t *testing.T) {
	var tests = []struct {
		pc         uint16
		h          byte
		l          byte
		expectedH  byte
		expectedL  byte
		expectedHF bool
		expectedC  bool
		expectedPC uint16
	}{
		{0x1234, 0x5E, 0x38, 0xBC, 0x70, true, false, 0x1235},
		{0x63F8, 0x82, 0x35, 0x04, 0x6A, false, true, 0x63F9},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.H = test.h
		cpu.reg.L = test.l

		cpu.cpu29()
		if cpu.reg.H != test.expectedH {
			t.Errorf("H: %x, expected: %x", cpu.reg.H, test.expectedH)
		}
		if cpu.reg.L != test.expectedL {
			t.Errorf("L: %x, expected: %x", cpu.reg.L, test.expectedL)
		}
		if cpu.flg.H != test.expectedHF {
			t.Errorf("H Flag: %t, expected %t", cpu.flg.H, test.expectedHF)
		}
		if cpu.flg.C != test.expectedC {
			t.Errorf("C Flag: %t, expected %t", cpu.flg.C, test.expectedC)
		}
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpu2A(t *testing.T) {
	var tests = []struct {
		pc         uint16
		H          byte
		L          byte
		address    uint16
		val        byte
		expectedPC uint16
		A          byte
	}{
		{0x9432, 0x13, 0x7F, 0x137F, 0x8E, 0x9433, 0x35},
		{0x2F3C, 0x30, 0x49, 0x3049, 0x3C, 0x2F3D, 0xF3},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.H = test.H
		cpu.reg.L = test.L
		Write(test.address, test.val)
		cpu.cpu2A()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if cpu.reg.A != test.val {
			t.Errorf("Current A: %x; expected: %x", cpu.reg.A, test.val)
		}
		hl := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
		if hl != test.address+1 {
			t.Errorf("Current HL: %x, expected: %x", hl, test.address+1)
		}
	}
}

func TestCpu2F(t *testing.T) {
	var tests = []struct {
		pc uint16
		a  byte
		//carry      bool
		// expectedCF bool
		expectedA  byte
		expectedPC uint16
	}{
		{0x1234, 0b10010110, 0b01101001, 0x1235},
		{0x1234, 0b10000000, 0b01111111, 0x1235},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.A = test.a
		cpu.cpu2F()
		if cpu.reg.A != test.expectedA {
			t.Errorf("Current A: %x, expected %x", cpu.reg.A, test.expectedA)
		}
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpu30(t *testing.T) {
	var tests = []struct {
		pc         uint16
		C          bool
		i8         int8
		expectedPC uint16
	}{
		{0x1234, false, 0x00, 0x1236},
		{0x1234, false, 0x08, 0x123E},
		{0x1234, false, -8, 0x122E},
		{0x1234, true, 0x00, 0x1236},
		{0x1234, true, 0x08, 0x1236},
		{0x1234, true, -8, 0x1236},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		Write(cpu.reg.PC+1, byte(test.i8))
		cpu.flg.C = test.C

		cpu.cpu30()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpu31(t *testing.T) {
	var tests = []struct {
		pc         uint16
		expectedHi byte
		expectedLo byte
		expectedSP uint16
		expectedPC uint16
	}{
		{0x9432, 0x13, 0x7F, 0x137F, 0x9435},
		{0x2F3C, 0x30, 0x49, 0x3049, 0x2F3F},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		Write(cpu.reg.PC+1, test.expectedLo)
		Write(cpu.reg.PC+2, test.expectedHi)
		cpu.cpu31()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if cpu.reg.SP != test.expectedSP {
			t.Errorf("Current SP: %x; expected: %x", cpu.reg.SP, test.expectedSP)
		}
	}
}
func TestCpu32(t *testing.T) {
	var tests = []struct {
		pc         uint16
		H          byte
		L          byte
		address    uint16
		expectedPC uint16
		A          byte
	}{
		{0x9432, 0x13, 0x7F, 0x137F, 0x9433, 0x35},
		{0x2F3C, 0x30, 0x49, 0x3049, 0x2F3D, 0xF3},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.A = test.A
		cpu.reg.H = test.H
		cpu.reg.L = test.L
		cpu.cpu32()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if Read(test.address) != test.A {
			t.Errorf("Current [BC]: %x; expected: %x", Read(test.address), test.A)
		}
		hl := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
		if hl != test.address-1 {
			t.Errorf("Current HL: %x, expected: %x", hl, test.address-1)
		}
	}
}

func TestCpu37(t *testing.T) {
	var tests = []struct {
		pc         uint16
		carry      bool
		expectedCF bool
		expectedPC uint16
	}{
		{0x1234, false, true, 0x1235},
		{0x1234, true, true, 0x1235},
		{0x63F8, false, true, 0x63F9},
		{0x63F8, true, true, 0x63F9},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.flg.C = test.carry
		cpu.cpu37()
		if cpu.flg.C != test.expectedCF {
			t.Errorf("Carry: %t, expected: %t", cpu.flg.C, test.expectedCF)
		}
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpu38(t *testing.T) {
	var tests = []struct {
		pc         uint16
		C          bool
		i8         int8
		expectedPC uint16
	}{
		{0x1234, true, 0x00, 0x1236},
		{0x1234, true, 0x08, 0x123E},
		{0x1234, true, -8, 0x122E},
		{0x1234, false, 0x00, 0x1236},
		{0x1234, false, 0x08, 0x1236},
		{0x1234, false, -8, 0x1236},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		Write(cpu.reg.PC+1, byte(test.i8))
		cpu.flg.C = test.C

		cpu.cpu38()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}
}

func TestCpu39(t *testing.T) {
	var tests = []struct {
		pc         uint16
		h          byte
		l          byte
		sp         uint16
		expectedH  byte
		expectedL  byte
		expectedHF bool
		expectedC  bool
		expectedPC uint16
	}{
		{0x1234, 0x5E, 0x38, 0x6F02, 0xCD, 0x3A, true, false, 0x1235},
		{0x63F8, 0xA2, 0x35, 0x6302, 0x05, 0x37, false, true, 0x63F9},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.H = test.h
		cpu.reg.L = test.l
		cpu.reg.SP = test.sp

		cpu.cpu39()
		if cpu.reg.H != test.expectedH {
			t.Errorf("H: %x, expected: %x", cpu.reg.H, test.expectedH)
		}
		if cpu.reg.L != test.expectedL {
			t.Errorf("L: %x, expected: %x", cpu.reg.L, test.expectedL)
		}
		if cpu.flg.H != test.expectedHF {
			t.Errorf("H Flag: %t, expected %t", cpu.flg.H, test.expectedHF)
		}
		if cpu.flg.C != test.expectedC {
			t.Errorf("C Flag: %t, expected %t", cpu.flg.C, test.expectedC)
		}
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpu3A(t *testing.T) {
	var tests = []struct {
		pc         uint16
		H          byte
		L          byte
		address    uint16
		val        byte
		expectedPC uint16
		A          byte
	}{
		{0x9432, 0x13, 0x7F, 0x137F, 0x8E, 0x9433, 0x35},
		{0x2F3C, 0x30, 0x49, 0x3049, 0x3C, 0x2F3D, 0xF3},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.H = test.H
		cpu.reg.L = test.L
		Write(test.address, test.val)
		cpu.cpu3A()
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if cpu.reg.A != test.val {
			t.Errorf("Current A: %x; expected: %x", cpu.reg.A, test.val)
		}
		hl := uint16(cpu.reg.H)<<8 + uint16(cpu.reg.L)
		if hl != test.address-1 {
			t.Errorf("Current HL: %x, expected: %x", hl, test.address-1)
		}
	}
}

func TestCpu3F(t *testing.T) {
	var tests = []struct {
		pc         uint16
		carry      bool
		expectedCF bool
		expectedPC uint16
	}{
		{0x1234, true, false, 0x1235},
		{0x1234, false, true, 0x1235},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.flg.C = test.carry
		cpu.cpu3F()
		if cpu.flg.C != test.expectedCF {
			t.Errorf("Current CF: %t, expected %t", cpu.flg.C, test.expectedCF)
		}
		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC: %x, expected: %x", cpu.reg.PC, test.expectedPC)
		}
	}

}

func TestCpuC3(t *testing.T) {
	var tests = []struct {
		pc       uint16
		lo       byte
		hi       byte
		expected uint16
	}{
		{0x0, 0x50, 0x01, 0x0150},
		{0x2F, 0x32, 0x7F, 0x7F32},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		Write(cpu.reg.PC+1, test.lo)
		Write(cpu.reg.PC+2, test.hi)
		cpu.cpuC3()
		if cpu.reg.PC != test.expected {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expected)
		}
	}
}

func TestCpuC6(t *testing.T) {
	var tests = []struct {
		pc         uint16
		a          byte
		u8         byte
		expectedPC uint16
		expectedA  byte
		expectedZF bool
		expectedNF bool
		expectedHF bool
		expectedCF bool
	}{
		{0xC8E9, 0xF8, 0x08, 0xC8EB, 0x00, true, false, true, true},
		{0xC8E9, 0x08, 0x08, 0xC8EB, 0x10, false, false, true, false},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.A = test.a
		Write(cpu.reg.PC+1, test.u8)
		cpu.cpuC6()

		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if cpu.reg.A != test.expectedA {
			t.Errorf("Current A %x; expected: %x", cpu.reg.A, test.expectedA)
		}
		if cpu.flg.Z != test.expectedZF {
			t.Errorf("Current Z %t; expected: %t", cpu.flg.Z, test.expectedZF)
		}
		if cpu.flg.N != test.expectedNF {
			t.Errorf("Current N %t; expected: %t", cpu.flg.N, test.expectedNF)
		}
		if cpu.flg.H != test.expectedHF {
			t.Errorf("Current H %t; expected: %t", cpu.flg.H, test.expectedHF)
		}
		if cpu.flg.C != test.expectedCF {
			t.Errorf("Current C %t; expected: %t", cpu.flg.C, test.expectedCF)
		}
	}
}

func TestCpuCE(t *testing.T) {
	var tests = []struct {
		pc         uint16
		a          byte
		u8         byte
		carry      bool
		expectedPC uint16
		expectedA  byte
		expectedZF bool
		expectedNF bool
		expectedHF bool
		expectedCF bool
	}{
		{0xC8E9, 0xF8, 0x08, false, 0xC8EB, 0x00, true, false, true, true},
		{0xC8E9, 0xF8, 0x08, true, 0xC8EB, 0x01, false, false, true, true},
		{0xC8E9, 0x08, 0x08, false, 0xC8EB, 0x10, false, false, true, false},
		{0xC8E9, 0x08, 0x08, true, 0xC8EB, 0x11, false, false, true, false},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.A = test.a
		cpu.flg.C = test.carry
		Write(cpu.reg.PC+1, test.u8)
		cpu.cpuCE()

		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if cpu.reg.A != test.expectedA {
			t.Errorf("Current A %x; expected: %x", cpu.reg.A, test.expectedA)
		}
		if cpu.flg.Z != test.expectedZF {
			t.Errorf("Current Z %t; expected: %t", cpu.flg.Z, test.expectedZF)
		}
		if cpu.flg.N != test.expectedNF {
			t.Errorf("Current N %t; expected: %t", cpu.flg.N, test.expectedNF)
		}
		if cpu.flg.H != test.expectedHF {
			t.Errorf("Current H %t; expected: %t", cpu.flg.H, test.expectedHF)
		}
		if cpu.flg.C != test.expectedCF {
			t.Errorf("Current C %t; expected: %t", cpu.flg.C, test.expectedCF)
		}
	}
}

func TestCpuDE(t *testing.T) {
	var tests = []struct {
		pc         uint16
		a          byte
		u8         byte
		carry      bool
		expectedPC uint16
		expectedA  byte
		expectedZF bool
		expectedNF bool
		expectedHF bool
		expectedCF bool
	}{
		{0xC8E9, 0xF8, 0x08, false, 0xC8EB, 0xF0, false, true, false, false},
		{0xC8E9, 0xF8, 0x08, true, 0xC8EB, 0xEF, false, true, true, false},
		{0xC8E9, 0x08, 0x08, false, 0xC8EB, 0x00, true, true, false, false},
		{0xC8E9, 0x08, 0x08, true, 0xC8EB, 0xFF, false, true, true, true},
		{0xC8E9, 0x00, 0xFF, true, 0xC8EB, 0x0, true, true, true, true},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.A = test.a
		cpu.flg.C = test.carry
		Write(cpu.reg.PC+1, test.u8)
		cpu.cpuDE()

		if cpu.reg.PC != test.expectedPC {
			t.Errorf("Current PC %x; expected: %x", cpu.reg.PC, test.expectedPC)
		}
		if cpu.reg.A != test.expectedA {
			t.Errorf("Current A %x; expected: %x", cpu.reg.A, test.expectedA)
		}
		if cpu.flg.Z != test.expectedZF {
			t.Errorf("Current Z %t; expected: %t", cpu.flg.Z, test.expectedZF)
		}
		if cpu.flg.N != test.expectedNF {
			t.Errorf("Current N %t; expected: %t", cpu.flg.N, test.expectedNF)
		}
		if cpu.flg.H != test.expectedHF {
			t.Errorf("Current H %t; expected: %t", cpu.flg.H, test.expectedHF)
		}
		if cpu.flg.C != test.expectedCF {
			t.Errorf("Current C %t; expected: %t", cpu.flg.C, test.expectedCF)
		}
	}
}

func TestSubA(t *testing.T) {
	var tests = []struct {
		a          byte
		reg        byte
		carry      bool
		cf         bool
		expectedA  byte
		expectedZF bool
		expectedNF bool
		expectedHF bool
		expectedCF bool
	}{
		{0x35, 0x37, false, false, 0xFE, false, true, true, true},
		{0x35, 0x37, false, true, 0xFE, false, true, true, true},
		{0xDE, 0x13, false, false, 0xCB, false, true, false, false},
		{0xDE, 0x13, false, true, 0xCB, false, true, false, false},
		{0x35, 0x37, true, false, 0xFE, false, true, true, true},
		{0x35, 0x37, true, true, 0xFD, false, true, true, true},
		{0xDE, 0x13, true, false, 0xCB, false, true, false, false},
		{0xDE, 0x13, true, true, 0xCA, false, true, false, false},
	}

	for _, test := range tests {
		cpu.reg.A = test.a
		cpu.flg.C = test.cf
		cpu.subA(test.reg, test.carry)

		if cpu.reg.A != test.expectedA {
			t.Errorf("Current A %x; expected: %x", cpu.reg.A, test.expectedA)
		}
		if cpu.flg.Z != test.expectedZF {
			t.Errorf("Current Z %t; expected: %t", cpu.flg.Z, test.expectedZF)
		}
		if cpu.flg.N != test.expectedNF {
			t.Errorf("Current N %t; expected: %t", cpu.flg.N, test.expectedNF)
		}
		if cpu.flg.H != test.expectedHF {
			t.Errorf("Current H %t; expected: %t", cpu.flg.H, test.expectedHF)
		}
		if cpu.flg.C != test.expectedCF {
			t.Errorf("Current C %t; expected: %t", cpu.flg.C, test.expectedCF)
		}
	}
}

func TestCpu98(t *testing.T) {
	var tests = []struct {
		a          byte
		reg        byte
		cf         bool
		expectedA  byte
		expectedZF bool
		expectedNF bool
		expectedHF bool
		expectedCF bool
	}{
		{0x35, 0x37, false, 0xFE, false, true, true, true},
		{0x35, 0x37, true, 0xFD, false, true, true, true},
		{0xDE, 0x13, false, 0xCB, false, true, false, false},
		{0xDE, 0x13, true, 0xCA, false, true, false, false},
	}

	for _, test := range tests {
		cpu.reg.A = test.a
		cpu.reg.B = test.reg
		cpu.flg.C = test.cf
		cpu.cpu98()

		if cpu.reg.A != test.expectedA {
			t.Errorf("Current A %x; expected: %x", cpu.reg.A, test.expectedA)
		}
		if cpu.flg.Z != test.expectedZF {
			t.Errorf("Current Z %t; expected: %t", cpu.flg.Z, test.expectedZF)
		}
		if cpu.flg.N != test.expectedNF {
			t.Errorf("Current N %t; expected: %t", cpu.flg.N, test.expectedNF)
		}
		if cpu.flg.H != test.expectedHF {
			t.Errorf("Current H %t; expected: %t", cpu.flg.H, test.expectedHF)
		}
		if cpu.flg.C != test.expectedCF {
			t.Errorf("Current C %t; expected: %t", cpu.flg.C, test.expectedCF)
		}
	}
}

func TestAddSP(t *testing.T) {
	var tests = []struct {
		sp         uint16
		i8         int8
		expectedSP uint16
	}{
		{0x0000, 1, 0x0001},
		{0x0001, 1, 0x0002},
		{0x000F, 1, 0x0010},
		{0x0010, 1, 0x0011},
		{0x001F, 1, 0x0020},
		{0x007F, 1, 0x0080},
		{0x0080, 1, 0x0081},
		{0x00FF, 1, 0x0100},
		{0x0F00, 1, 0x0F01},
		{0x1F00, 1, 0x1F01},
		{0x1000, 1, 0x1001},
		{0x7FFF, 1, 0x8000},
		{0x8000, 1, 0x8001},
		{0xFFFF, 1, 0x0000},
		{0x0000, -1, 0xFFFF},
		{0x0001, -1, 0x0000},
		{0x000F, -1, 0x000E},
		{0x0010, -1, 0x000F},
		{0x001F, -1, 0x001E},
		{0x007F, -1, 0x007E},
		{0x0080, -1, 0x007F},
		{0x00FF, -1, 0x00FE},
		{0x0F00, -1, 0x0EFF},
		{0x1F00, -1, 0x1EFF},
		{0x1000, -1, 0x0FFF},
		{0x7FFF, -1, 0x7FFE},
		{0x8000, -1, 0x7FFF},
		{0xFFFF, -1, 0xFFFE},
	}

	for _, test := range tests {
		cpu.reg.SP = test.sp
		var result uint16 = cpu.addSP(test.i8)

		if result != test.expectedSP {
			t.Errorf("Current SP %x; expected: %x", cpu.reg.SP, test.expectedSP)
		}
	}
}

func TestCpuE8(t *testing.T) {
	var tests = []struct {
		sp         uint16
		i8         int8
		expectedSP uint16
		expectedZF bool
		expectedNF bool
	}{
		{0x0000, 1, 0x0001, false, false},
		{0x0001, 1, 0x0002, false, false},
		{0x000F, 1, 0x0010, false, false},
		{0x0010, 1, 0x0011, false, false},
		{0x001F, 1, 0x0020, false, false},
		{0x007F, 1, 0x0080, false, false},
		{0x0080, 1, 0x0081, false, false},
		{0x00FF, 1, 0x0100, false, false},
		{0x0F00, 1, 0x0F01, false, false},
		{0x1F00, 1, 0x1F01, false, false},
		{0x1000, 1, 0x1001, false, false},
		{0x7FFF, 1, 0x8000, false, false},
		{0x8000, 1, 0x8001, false, false},
		{0xFFFF, 1, 0x0000, false, false},
		{0x0000, -1, 0xFFFF, false, false},
		{0x0001, -1, 0x0000, false, false},
		{0x000F, -1, 0x000E, false, false},
		{0x0010, -1, 0x000F, false, false},
		{0x001F, -1, 0x001E, false, false},
		{0x007F, -1, 0x007E, false, false},
		{0x0080, -1, 0x007F, false, false},
		{0x00FF, -1, 0x00FE, false, false},
		{0x0F00, -1, 0x0EFF, false, false},
		{0x1F00, -1, 0x1EFF, false, false},
		{0x1000, -1, 0x0FFF, false, false},
		{0x7FFF, -1, 0x7FFE, false, false},
		{0x8000, -1, 0x7FFF, false, false},
		{0xFFFF, -1, 0xFFFE, false, false},
	}

	for _, test := range tests {
		cpu.reg.PC = 0xF345
		cpu.reg.SP = test.sp
		Write(cpu.reg.PC+1, byte(test.i8))
		cpu.cpuE8()

		carries := test.sp ^ uint16(test.i8) ^ test.expectedSP
		expectedHF := carries&0x10 == 0x10
		expectedCF := carries&0x100 == 0x100

		if cpu.reg.SP != test.expectedSP {
			t.Errorf("Current SP %x; expected: %x", cpu.reg.SP, test.expectedSP)
		}
		if cpu.flg.Z != test.expectedZF {
			t.Errorf("Current Z %t; expected: %t", cpu.flg.Z, test.expectedZF)
		}
		if cpu.flg.N != test.expectedNF {
			t.Errorf("Current N %t; expected: %t", cpu.flg.N, test.expectedNF)
		}
		if cpu.flg.H != expectedHF {
			t.Errorf("Current H %t; expected: %t, i8: %d", cpu.flg.H, expectedHF, test.i8)
		}
		if cpu.flg.C != expectedCF {
			t.Errorf("Current C %t; expected: %t", cpu.flg.C, expectedCF)
		}
		if cpu.reg.PC != 0xF347 {
			t.Errorf("Current PC: %x; expected: %x", cpu.reg.PC, 0xF345)
		}
	}
}

func TestIncreaseDiv(t *testing.T) {
	var tests = []struct {
		div_clocksum  byte
		cycle         byte
		expected_div  byte
		expected_FF04 byte
	}{
		{0x00, 0x01, 0x01, 0x00},
		{0x01, 0x01, 0x02, 0x00},
		{0x0F, 0x01, 0x10, 0x00},
		{0x10, 0x01, 0x11, 0x00},
		{0x1F, 0x01, 0x20, 0x00},
		{0x7F, 0x01, 0x80, 0x00},
		{0x80, 0x01, 0x81, 0x00},
		{0xFF, 0x01, 0x00, 0x01},
		{0xF0, 0x11, 0x01, 0x01},
		{0x80, 0x81, 0x01, 0x01},
		{0xFF, 0x0F, 0x0E, 0x01},
		{0xFF, 0xF0, 0xEF, 0x01},
		{0xFF, 0xFF, 0xFE, 0x01},
	}

	for _, test := range tests {
		Write(0xFF04, 0x00)
		cpu.clk.div_clocksum = test.div_clocksum
		cpu.increase_div(test.cycle)

		actual_div := cpu.clk.div_clocksum
		actual_FF04 := Read(0xFF04)

		if actual_div != test.expected_div {
			t.Errorf("Current div_clocksum: %x; expected: %x", actual_div, test.expected_div)
		}
		if actual_FF04 != test.expected_FF04 {
			t.Errorf("Current DIV Register 0xFF04: %x; expected: %x", actual_FF04, test.expected_FF04)
		}
	}
}
func TestGetTimerFrequency(t *testing.T) {
	var tests = []struct {
		ff07         byte
		expected_div uint
	}{
		{0b0000, 1024},
		{0b0001, 16},
		{0b0010, 64},
		{0b0011, 256},
		{0b0100, 1024},
		{0b0101, 16},
		{0b0110, 64},
		{0b0111, 256},
		{0b1000, 1024},
		{0b1001, 16},
		{0b1010, 64},
		{0b1011, 256},
	}

	for _, test := range tests {
		Write(0xFF07, test.ff07)
		expected_freq := cpu.clk.MASTER_CLK / test.expected_div

		actual_freq := cpu.get_timer_frequency()

		if actual_freq != expected_freq {
			t.Errorf("Current frequency: %d; expected: %d; FF07: %x, divider: %d", actual_freq, expected_freq, Read(0xFF07), 4194304*actual_freq)
		}
	}
}
func TestIncreaseRegister(t *testing.T) {
	var tests = []struct {
		ff05          byte
		increment     byte
		expected_ovfl bool
	}{
		{0x00, 0x00, false},
		{0x00, 0x0F, false},
		{0x00, 0xFF, false},
		{0x01, 0x00, false},
		{0x01, 0x0F, false},
		{0x01, 0xFF, true},
		{0x08, 0x00, false},
		{0x08, 0x0F, false},
		{0x08, 0xFF, true},
		{0x0F, 0x00, false},
		{0x0F, 0x0F, false},
		{0x0F, 0xFF, true},
		{0xF0, 0x00, false},
		{0xF0, 0x0F, false},
		{0xF0, 0xFF, true},
		{0xFF, 0x00, false},
		{0xFF, 0x0F, true},
		{0xFF, 0xFF, true},
	}

	for _, test := range tests {
		Write(0xFF05, test.ff05)
		expected_ff05 := test.ff05 + test.increment

		actual_ovfl := cpu.increase_register(0xff05, test.increment)
		actual_ff05 := Read(0xff05)

		if actual_ovfl != test.expected_ovfl {
			t.Errorf("Current overflow: %t; expected: %t; FF05 before: %x, increment: %x; FF05 after: %x", actual_ovfl, test.expected_ovfl, test.ff05, test.increment, actual_ff05)
		}

		if actual_ff05 != expected_ff05 {
			t.Errorf("Current FF05: %x; expected: %x; FF05 before: %x, increment: %x", actual_ff05, expected_ff05, test.ff05, test.increment)
		}

	}
}

func TestSetInterruptTimer(t *testing.T) {
	var tests = []struct {
		initial_ff0f  byte
		request_bit   byte
		expected_ff0f byte
	}{
		{0b00000000, 0b00101100, 0b00101100},
		{0b00000011, 0b00101100, 0b00101111},
		{0b10000111, 0b01001011, 0b11001111},
	}

	for _, test := range tests {
		Write(0xFF0F, test.initial_ff0f)

		cpu.set_interrupt_request(test.request_bit)

		actual_ff0f := Read(0xFF0F)

		if actual_ff0f != test.expected_ff0f {
			t.Errorf("Current FF0F: %b; expected: %b; request_bits: %b", actual_ff0f, test.expected_ff0f, test.request_bit)
		}
	}
}

func TestHandleTimer(t *testing.T) {
	var tests = []struct {
		initial_tma  byte
		initial_tac  byte
		initial_tima byte
		cycles       uint16
		initial_if   byte
		expected_if  byte
	}{
		// blargg's test rom setup
		{0b00000000, 0x05, 0x0, 0x00, 0, 0},
		{0b00000000, 0x05, 0x0, 500, 0, 0},
		{0b00000000, 0x05, 0x0, 1500, 0, 4},
		// {0b00000011, 0b00101100, 0b00101111},
		// {0b10000111, 0b01001011, 0b11001111},
	}

	for _, test := range tests {
		Write(0xFF06, test.initial_tma)
		Write(0xFF07, test.initial_tac)
		Write(0xFF05, test.initial_tima)
		Write(0xFF0F, test.initial_if)

		remaining_cycles := test.cycles

		for remaining_cycles/256 > 0 {
			cpu.Handle_timer(255)

			remaining_cycles -= 255
		}

		actual_if := Read(0xFF0F)

		if actual_if != test.expected_if {
			t.Errorf("Current FF0F: %b; expected: %b; cycles: %d; TIMA: %x", actual_if, test.expected_if, test.cycles, Read(0xFF05))
		}
	}
}

// func TestCpuFB(t *testing.T) {
//     // blargg's test roms
//     cpu.cpuFB() // ei
//     cpu.ld16(&cpu.reg.C, &cpu.reg.B, 0, 0) // ld bc, 0
//     cpu.cpuC5() // push bc
//     cpu.cpuC1() // pop bc
//     cpu.cpu04() // inc b
//
//     Write(cpu.IF, 0x04) // wreg IF,$04
//     cpu.cpu05() // dec b
//     // jp nz, test failed
//     if !cpu.flg.Z {
// 	t.Errorf("Zero flag was not set after dec b")
//     }
//
//     cpu.reg.PC = 0x100
//     Write(PC+1, -2)
//     cpu.cpuF8() // ld hl, sp-2
//     cpu.cpu2A() // ldi a, (hl)
//     cpu.cpu() // cp < interrupt_addr
//     cpu.cpu() // jp nz, test_failed
//     cpu.cpu() // ld a,(hl)
//     cpu.cpu() // cp >interrupt_addr
//     cpu.cpu() // jp nz, test_failed
//     cpu.cpu() // lda IF
//     cpu.cpu() // and $04
//     cpu.cpu() // jp nz,test_failed
//
// }
