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
		expectedA  byte
		expectedPC uint16
	}{
		{0x1234, 0x80, true, 0x01, 0x1235},
		{0x63F8, 0x35, false, 0x6A, 0x63F9},
	}

	for _, test := range tests {
		cpu.reg.PC = test.pc
		cpu.reg.A = test.a
		cpu.cpu07()
		if cpu.flg.C != test.carry {
			t.Errorf("Carry: %t, expected: %t", cpu.flg.C, test.carry)
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
