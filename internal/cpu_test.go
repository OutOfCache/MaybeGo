package maybego

import (
	//	"fmt"
	"testing"
)

var cpu *CPU = NewCPU()

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
