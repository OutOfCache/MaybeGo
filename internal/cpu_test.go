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
