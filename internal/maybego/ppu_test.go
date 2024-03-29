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
		{1, 0x10, 0x10, 0x0},   // mode 0, rows 0-143, ie, stat set
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

func TestMode1STATInterrupt(t *testing.T) {
	var tests = []struct {
		ly           byte
		stat         byte
		expectedSTAT byte
		expectedIF   byte
	}{
		// stat set, int enabled
		{1, 0x10, 0x10, 0x0},   // mode 0, rows 0-143, ie, stat set
		{143, 0x10, 0x10, 0x0}, // mode 0, rows 0-143, ie, stat set
		{144, 0x10, 0x11, 0x2}, // mode 0, row 145-153, ie
		{145, 0x10, 0x11, 0x2}, // mode 0, row 144, ie
		// stat not set, int enabled
		{144, 0x00, 0x01, 0x0}, // mode 0, row 143, ie
		{145, 0x01, 0x01, 0x0}, // mode 1, rows 144, ie
	}

	cpu.flg.IME = true
	for _, test := range tests {
		Write(IF, 0x0)
		Write(LY, test.ly)
		Write(STAT, test.stat)
		ppu.RenderRow()

		actualSTAT := Read(STAT)
		actualIF := (Read(IF) & 0x2)

		if actualSTAT != test.expectedSTAT {
			t.Errorf("Wrong STAT. Got %.2X, expected %.2X", actualSTAT, test.expectedSTAT)
			t.Errorf("Test: {LY: %.2d, STAT: %.2X", test.ly, test.stat)
		}

		if actualIF != test.expectedIF {
			t.Errorf("Wrong IF. Got %.2X, expected %.2X", actualIF, test.expectedIF)
			t.Errorf("Test: {LY: %.2d, STAT: %.2X", test.ly, test.stat)
		}
	}
}

func TestMode2STATInterrupt(t *testing.T) {
	var tests = []struct {
		dots         uint16
		cycles       byte
		stat         byte
		expectedSTAT byte
		expectedIF   byte
	}{
		// stat set, int enabled
		{0, 1, 0x20, 0x22, 0x2},   // mode 0, rows 0-143, ie, stat set
		{0, 80, 0x20, 0x22, 0x2},  // mode 0, rows 0-143, ie, stat set
		{0, 81, 0x20, 0x23, 0x0},  // mode 0, rows 0-143, ie, stat set
		{80, 0, 0x20, 0x22, 0x2},  // mode 0, rows 0-143, ie, stat set
		{80, 1, 0x20, 0x23, 0x0},  // mode 0, rows 0-143, ie, stat set
		{128, 0, 0x20, 0x23, 0x0}, // mode 0, rows 0-143, ie, stat set
		{456, 0, 0x20, 0x20, 0x0}, // mode 0, rows 0-143, ie, stat set
		{456, 1, 0x20, 0x22, 0x2}, // mode 0, rows 0-143, ie, stat set
		// stat not enabled
		{0, 1, 0x00, 0x02, 0x0},   // mode 0, rows 0-143, ie, stat set
		{0, 80, 0x00, 0x02, 0x0},  // mode 0, rows 0-143, ie, stat set
		{0, 81, 0x00, 0x03, 0x0},  // mode 0, rows 0-143, ie, stat set
		{80, 0, 0x00, 0x02, 0x0},  // mode 0, rows 0-143, ie, stat set
		{80, 1, 0x00, 0x03, 0x0},  // mode 0, rows 0-143, ie, stat set
		{128, 0, 0x00, 0x03, 0x0}, // mode 0, rows 0-143, ie, stat set
		{456, 0, 0x00, 0x00, 0x0}, // mode 0, rows 0-143, ie, stat set
		{456, 1, 0x00, 0x02, 0x0}, // mode 0, rows 0-143, ie, stat set
	}

	cpu.flg.IME = true
	for _, test := range tests {
		Write(IF, 0x0)
		Write(STAT, test.stat)
		ppu.dots = test.dots
		ppu.Render(test.cycles)

		actualSTAT := Read(STAT)
		actualIF := (Read(IF) & 0x2)

		if actualSTAT != test.expectedSTAT {
			t.Errorf("Wrong STAT. Got %.2X, expected %.2X", actualSTAT, test.expectedSTAT)
			t.Errorf("Test: {dots: %.2d, cycles: %.2d, STAT: %.2X", test.dots, test.cycles, test.stat)
		}

		if actualIF != test.expectedIF {
			t.Errorf("Wrong IF. Got %.2X, expected %.2X", actualIF, test.expectedIF)
			t.Errorf("Test: {dots: %.2d, cycles: %.2d, STAT: %.2X", test.dots, test.cycles, test.stat)
		}
	}
}

func TestMode0STATInterrupt(t *testing.T) {
	var tests = []struct {
		dots         uint16
		cycles       byte
		stat         byte
		expectedSTAT byte
		expectedIF   byte
	}{
		// stat set, int enabled
		{369, 0, 0x8, 0xB, 0x0},  // mode 0, rows 0-143, ie, stat set
		{369, 1, 0x8, 0x8, 0x2},  // mode 0, rows 0-143, ie, stat set
		{370, 0, 0x8, 0x8, 0x2},  // mode 0, rows 0-143, ie, stat set
		{370, 87, 0x8, 0xA, 0x0}, // mode 0, rows 0-143, ie, stat set
		{456, 0, 0x8, 0x8, 0x2},  // mode 0, rows 0-143, ie, stat set
		{456, 1, 0x8, 0xA, 0x0},  // mode 0, rows 0-143, ie, stat set
		// stat not enabled
		{369, 0, 0x3, 0x3, 0x0},  // mode 0, rows 0-143, ie, stat set
		{369, 1, 0x3, 0x0, 0x0},  // mode 0, rows 0-143, ie, stat set
		{370, 0, 0x0, 0x0, 0x0},  // mode 0, rows 0-143, ie, stat set
		{370, 87, 0x0, 0x2, 0x0}, // mode 0, rows 0-143, ie, stat set
		{456, 0, 0x0, 0x0, 0x0},  // mode 0, rows 0-143, ie, stat set
		{456, 1, 0x0, 0x2, 0x0},  // mode 0, rows 0-143, ie, stat set
	}

	cpu.flg.IME = true
	for _, test := range tests {
		Write(IF, 0x0)
		Write(STAT, test.stat)
		ppu.dots = test.dots
		ppu.Render(test.cycles)

		actualSTAT := Read(STAT)
		actualIF := (Read(IF) & 0x2)

		if actualSTAT != test.expectedSTAT {
			t.Errorf("Wrong STAT. Got %.2X, expected %.2X", actualSTAT, test.expectedSTAT)
			t.Errorf("Test: {dots: %.2d, cycles: %.2d, STAT: %.2X", test.dots, test.cycles, test.stat)
		}

		if actualIF != test.expectedIF {
			t.Errorf("Wrong IF. Got %.2X, expected %.2X", actualIF, test.expectedIF)
			t.Errorf("Test: {dots: %.2d, cycles: %.2d, STAT: %.2X", test.dots, test.cycles, test.stat)
		}
	}
}

func TestLYCInterrupt(t *testing.T) {
	var tests = []struct {
		ly           byte
		lyc          byte
		stat         byte
		expectedIF   byte
		expectedSTAT byte
	}{
		// stat set, int enabled
		{128, 127, 0x40, 0x0, 0x40}, // mode 0, rows 0-143, ie, stat set
		{128, 128, 0x40, 0x2, 0x44}, // mode 0, rows 0-143, ie, stat set
		{128, 129, 0x40, 0x0, 0x40}, // mode 0, rows 0-143, ie, stat set
		// stat disabled
		{128, 127, 0x00, 0x0, 0x00}, // mode 0, rows 0-143, ie, stat set
		{128, 128, 0x00, 0x0, 0x04}, // mode 0, rows 0-143, ie, stat set
		{128, 129, 0x00, 0x0, 0x00}, // mode 0, rows 0-143, ie, stat set
	}

	cpu.flg.IME = true
	for _, test := range tests {
		Write(IF, 0x0)
		Write(LY, test.ly)
		Write(LYC, test.lyc)
		Write(STAT, test.stat)
		ppu.RenderRow()

		actualSTAT := Read(STAT)
		actualIF := (Read(IF) & 0x2)

		if actualSTAT != test.expectedSTAT {
			t.Errorf("Wrong STAT. Got %.2X, expected %.2X", actualSTAT, test.expectedSTAT)
			t.Errorf("Test: {dots: %.2d, cycles: %.2d, STAT: %.2X", test.ly, test.lyc, test.stat)
		}

		if actualIF != test.expectedIF {
			t.Errorf("Wrong IF. Got %.2X, expected %.2X", actualIF, test.expectedIF)
			t.Errorf("Test: {dots: %.2d, cycles: %.2d, STAT: %.2X", test.ly, test.lyc, test.stat)
		}
	}
}
