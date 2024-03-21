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
