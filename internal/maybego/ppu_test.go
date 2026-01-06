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
		{153, 0},
	}

	for _, test := range tests {
		Write(LY, test.ly)
		ppu.dots = 456
		ppu.Render(0)
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
		{144, 0x10, 0x11, 0x1}, // mode 0, row 144, ie
		{145, 0x10, 0x11, 0x1}, // mode 0, row 145-153, ie
		{145, 0x11, 0x11, 0x0}, // mode 1, row 145-153, ie
		// stat not set, int enabled
		{144, 0x00, 0x01, 0x1}, // mode 0, row 143, ie
		{145, 0x01, 0x01, 0x0}, // mode 1, rows 144, ie
	}

	cpu.flg.IME = true
	for _, test := range tests {
		Write(IF, 0x0)
		Write(LY, test.ly)
		Write(STAT, test.stat)
		Write(LCDC, 0x1) // LCD enable
		ppu.dots = MODE0_END + 1
		ppu.Render(0)

		actualSTAT := Read(STAT)
		actualIF := (Read(IF) & 0x1)

		if actualSTAT != test.expectedSTAT {
			t.Errorf("Wrong STAT. Got %.2X, expected %.2X", actualSTAT, test.expectedSTAT)
			t.Errorf("Test: {LY: %3d, STAT: %.2X}", test.ly, test.stat)
		}

		if actualIF != test.expectedIF {
			t.Errorf("Wrong IF. Got %.2X, expected %.2X", actualIF, test.expectedIF)
			t.Errorf("Test: {LY: %3d, STAT: %.2X}", test.ly, test.stat)
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
		{0, 20, 0x20, 0x22, 0x2},  // mode 0, rows 0-143, ie, stat set
		{0, 21, 0x20, 0x23, 0x0},  // mode 0, rows 0-143, ie, stat set
		{80, 0, 0x20, 0x22, 0x2},  // mode 0, rows 0-143, ie, stat set
		{80, 1, 0x20, 0x23, 0x0},  // mode 0, rows 0-143, ie, stat set
		{128, 0, 0x20, 0x23, 0x0}, // mode 0, rows 0-143, ie, stat set
		{455, 0, 0x20, 0x20, 0x0}, // mode 0, rows 0-143, ie, stat set
		{455, 1, 0x20, 0x22, 0x2}, // mode 0, rows 0-143, ie, stat set
		// stat not enabled
		{0, 1, 0x00, 0x02, 0x0},   // mode 0, rows 0-143, ie, stat set
		{0, 20, 0x00, 0x02, 0x0},  // mode 0, rows 0-143, ie, stat set
		{0, 21, 0x00, 0x03, 0x0},  // mode 0, rows 0-143, ie, stat set
		{80, 0, 0x00, 0x02, 0x0},  // mode 0, rows 0-143, ie, stat set
		{80, 1, 0x00, 0x03, 0x0},  // mode 0, rows 0-143, ie, stat set
		{128, 0, 0x00, 0x03, 0x0}, // mode 0, rows 0-143, ie, stat set
		{455, 0, 0x00, 0x00, 0x0}, // mode 0, rows 0-143, ie, stat set
		{455, 1, 0x00, 0x02, 0x0}, // mode 0, rows 0-143, ie, stat set
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
		{455, 0, 0x8, 0x8, 0x2},  // mode 0, rows 0-143, ie, stat set
		{455, 1, 0x8, 0xA, 0x0},  // mode 0, rows 0-143, ie, stat set
		// stat not enabled
		{369, 0, 0x3, 0x3, 0x0},  // mode 0, rows 0-143, ie, stat set
		{369, 1, 0x3, 0x0, 0x0},  // mode 0, rows 0-143, ie, stat set
		{370, 0, 0x0, 0x0, 0x0},  // mode 0, rows 0-143, ie, stat set
		{370, 87, 0x0, 0x2, 0x0}, // mode 0, rows 0-143, ie, stat set
		{455, 0, 0x0, 0x0, 0x0},  // mode 0, rows 0-143, ie, stat set
		{455, 1, 0x0, 0x2, 0x0},  // mode 0, rows 0-143, ie, stat set
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

// tile data from https://gbdev.io/pandocs/Tile_Data.html
var tile = []byte{
	0x3C, 0x7E,
	0x42, 0x42,
	0x42, 0x42,
	0x42, 0x42,
	0x7E, 0x5E,
	0x7E, 0x0A,
	0x7C, 0x56,
	0x38, 0x7C,
}

var tileColors = []uint32{
	0, 2, 3, 3, 3, 3, 2, 0,
	0, 3, 0, 0, 0, 0, 3, 0,
	0, 3, 0, 0, 0, 0, 3, 0,
	0, 3, 0, 0, 0, 0, 3, 0,
	0, 3, 1, 3, 3, 3, 3, 0,
	0, 1, 1, 1, 3, 1, 3, 0,
	0, 3, 1, 3, 1, 3, 2, 0,
	0, 2, 3, 3, 3, 2, 0, 0,
}

func TestUnsignedTileData(t *testing.T) {
	var tests = []struct {
		tileNr      int
	}{
		{0}, {31}, {32}, {255}, {1023},
	}

	// LCD & PPU enable	// BG data area: 8000-8FFF, unsigned
	ppu.tiledata = 0x8000
	// Setup tile data for tileID 1
	for i := 0; i < 16; i+= 1 {
		Write(uint16(0x8010 + i), tile[i]);
	}

	cpu.flg.IME = true
	for _, test := range tests {
		// Set the tested tiles to tileID 1.
		Write(ppu.tilemap + uint16(test.tileNr), 0x1)
		startRow := 8 * (test.tileNr / 32);
		for i := 0; i < 8; i++ {
		    ppu.RenderBG(byte(startRow + i));
		}

		y := startRow;
		x := (test.tileNr % 32) * 8;
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				actualColor := BGMapRGBA[(((y + i) * 256)) + (x + j)];
				expectedColor := Palette[tileColors[i * 8 + j]];
				if actualColor != expectedColor {
					t.Errorf("Wrong color in tile %d. Got %.6X @ (%d,%d), expected %.6X", test.tileNr, actualColor, x + j, y + i, expectedColor);
				}
			}
		}
	}
}

func TestSignedTileData(t *testing.T) {
	var tests = []struct {
		tileNr      int
		tileID      int
	}{
		{0, 1},   {31, 1},   {32, 1},   {255, 1},   {1023, 1},   // block 2 (0x9000-0x9FFF)
		{0, 127}, {31, 127}, {32, 127}, {255, 127}, {1023, 127}, // block 2 (0x9000-0x9FFF)
		{0, 128}, {31, 128}, {32, 128}, {255, 128}, {1023, 128}, // block 1 (0x8800-0x8FFF)
		{0, 255}, {31, 255}, {32, 255}, {255, 255}, {1023, 255}, // block 1 (0x8800-0x8FFF)
	}

	// LCD & PPU enable	// BG data area: 0x8800-0x97FF, signed
	ppu.tiledata = 0x8800

	cpu.flg.IME = true
	for _, test := range tests {
		// Setup tile data for tileID 1
		for i := 0; i < 16; i+= 1 {
			Write(uint16(0x9000 + uint16(int8(test.tileID * 0x10)) + uint16(i)), tile[i]);
		}
		// Set the tested tiles to the tested tileID.
		Write(ppu.tilemap + uint16(test.tileNr), byte(test.tileID))
		startRow := 8 * (test.tileNr / 32);
		for i := 0; i < 8; i++ {
		    ppu.RenderBG(byte(startRow + i));
		}

		y := startRow;
		x := (test.tileNr % 32) * 8;
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				actualColor := BGMapRGBA[(((y + i) * 256)) + (x + j)];
				expectedColor := Palette[tileColors[i * 8 + j]];
				if actualColor != expectedColor {
					t.Errorf("Wrong color in tile %d, ID %d. Got %.6X @ (%d,%d), expected %.6X", test.tileNr, test.tileID, actualColor, x + j, y + i, expectedColor);
				}
			}
		}
	}
}

func TestLCDCSettings(t *testing.T) {
	var tests = []struct {
		lcdc    byte
		expectedBGTiledata uint16
		expectedBGTilemap uint16
	}{
	    {0x00, 0x8800, 0x9800},
	    {0x08, 0x8800, 0x9C00},
	    {0x10, 0x8000, 0x9800},
	    {0x18, 0x8000, 0x9C00},
	}

	for _, test := range tests {
		Write(LCDC, test.lcdc)
		ppu.dots = 456
		ppu.Render(1)

		if ppu.tiledata != test.expectedBGTiledata {
			t.Errorf("Tiledata is %.4X; expected %.4X for LCDC %.2X", ppu.tiledata, test.expectedBGTiledata, test.lcdc)
		}
		if ppu.tilemap != test.expectedBGTilemap {
			t.Errorf("Tilemap is %.4X; expected %.4X for LCDC %.2X", ppu.tilemap, test.expectedBGTilemap, test.lcdc)
		}
	}
}
