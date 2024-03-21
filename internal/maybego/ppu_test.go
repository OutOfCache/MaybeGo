package maybego

import (
	"testing"
)

var logger *Logger = NewLogger(false, "")
var ppu *PPU = NewPPU(logger)
var cpu *CPU = NewCPU(logger)

func TestRowTransition(t *testing.T) {
	var tests = []struct {
		ly         byte
		expectedLY byte
	}{
		{0, 1},
		{143, 144},
		{152, 0},
	}

	for _, test := range tests {
		Write(LY, test.ly)
		ppu.RenderRow()
		if Read(LY) != test.expectedLY {
			t.Errorf("Current LY: %3d; expected: %3d", Read(LY), test.expectedLY)
		}
	}
}

func TestVBlankInterrupt(t *testing.T) {
	var tests = []struct {
		ly           byte
		stat         byte
		expectedSTAT byte
		expectedIF   byte
	}{
		// stat set, int enabled
		{0, 0x10, 0x10, 0x0},   // mode 0, rows 0-143, ie, stat set
		{143, 0x10, 0x10, 0x0}, // mode 0, rows 0-143, ie, stat set
		{144, 0x10, 0x11, 0x1}, // mode 0, row 145-153, ie
		{145, 0x10, 0x11, 0x0}, // mode 0, row 144, ie
		// stat not set, int enabled
		{144, 0x00, 0x01, 0x1}, // mode 0, row 143, ie
		{145, 0x01, 0x01, 0x0}, // mode 1, rows 144, ie
	}

	cpu.flg.IME = true
	for _, test := range tests {
		Write(IF, 0x0)
		Write(LY, test.ly)
		Write(STAT, test.stat)
		ppu.RenderRow()

		actualSTAT := Read(STAT)
		actualIF := (Read(IF) & 0x1)

		if actualSTAT != test.expectedSTAT {
			t.Errorf("Wrong STAT. Got %.2X, expected %.2X", actualSTAT, test.expectedSTAT)
			t.Errorf("Test: {LY: %3d, STAT: %.2X", test.ly, test.stat)
		}

		if actualIF != test.expectedIF {
			t.Errorf("Wrong IF. Got %.2X, expected %.2X", actualIF, test.expectedIF)
			t.Errorf("Test: {LY: %3d, STAT: %.2X", test.ly, test.stat)
		}
	}
}
