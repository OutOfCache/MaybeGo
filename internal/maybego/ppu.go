package maybego

const (
	LCDC      uint16 = 0xFF40
	STAT      uint16 = 0xFF41
	LY        uint16 = 0xFF44
	LYC       uint16 = 0xFF45
	MODE2_END uint16 = 80
	MODE3_END uint16 = 80 + 289
	MODE0_END uint16 = 456
)

type PPU struct {
	tilemap  uint16
	tiledata uint16
	dots     uint16
	scanline byte
	logger   *Logger
}

var framebufferPalette [160 * 144]byte
var BGMapPalette [256 * 256]byte

var winWidth, winHeight int32 = 160, 144
var err error

func NewPPU(logger *Logger) *PPU {
	ppu := &PPU{logger: logger, dots: 0, scanline: 0}
	ppu.Reset()

	return ppu
}

func (ppu *PPU) GetCurrentFrame() *[160 * 144]byte {
	return &framebufferPalette
}

func (ppu *PPU) RenderBG(row byte) {
	y := int(row)
	// FIXME: tileID only changes every 8 pixels
	// for tileID := 0; tileID < 8; tileID += 1 {
	// 	var tileStart uint16
	// 	if ppu.tiledata == 0x8800 {
	// 		tileStart = uint16(0x800 + uint16(int8(tileID*0x10)))
	// 	} else {
	// 		tileStart = uint16(tileID * 0x10)
	// 	}
	// 	address := ppu.tiledata + tileStart
	// 	// fmt.Printf("TileID: %d @ %x\n\n", tileID, address)
	// 	for y := uint16(0); y < 8; y += 1 {

	// 		data1 := int64(Read(address + y*2))
	// 		if data1 != 0 {
	// 			fmt.Printf("data1 @ address %x:\t\t%s\n", address+y*2, strconv.FormatInt(data1, 2))
	// 		}
	// 		data2 := int64(Read(address + 1 + y*2))
	// 		if data2 != 0 {
	// 			fmt.Printf("data2 @ address %x:\t\t%s\n", address+y*2+1, strconv.FormatInt(data2, 2))
	// 		}
	// 	}
	// }
	for j := 0; j < /*(SCX + */ 256; j += 1 {
		x := j // + SCX
		tileX := uint16(x / 8)
		tileID := Read(ppu.tilemap + uint16((y/8)*32) + tileX)

		var tileY uint16

		if ppu.tiledata == 0x8800 {
			tileY = uint16(0x800 + uint16(int8(uint16(tileID)*0x10)))
		} else {
			tileY = uint16(tileID) * uint16(0x10)
		}
		address := ppu.tiledata + tileY + uint16((y%8)*2)

		pixelcolor := (Read(address) >> (7 - (x % 8)) & 0x1) +
			(Read(address+1)>>(7-(x%8))&0x1)*2
		// pixelcolor := address
		// fmt.Printf("")
		// pixelcolor := (Read(address) >> (7 - (x % 8)) & 0x1) +
		// 	(Read(address+1)>>(7-(x%8))&0x1)*2
		// if (x >= (2 * 8) && x < (3 * 8) && y < 8) {
		// 	pixelcolor := Read(uint16(0x82d0 + y * 2))
		// 	fmt.Printf("(%d, %d): tileID: %x, tileY: %x, tileID * 0x10: %x, address: %x, tiledata: %x, color: %x\n", x, y, tileID, tileY, uint16(tileID) * uint16(0x10), address, ppu.tiledata, pixelcolor)
		// }
		// fmt.Printf("(%d, %d): tileID: %x, address: %x, tiledata: %x, color: %d\n", x, y, tileID, address, ppu.tiledata, pixelcolor)
		// if pixelcolor != 0 {
		// 	fmt.Printf("Color @ (%d, %d): %d\n", x, y, pixelcolor)
		// }
		BGMapPalette[y*256+x] = pixelcolor
		if x /* - SCX */ < 160 && y < 144 {
			framebufferPalette[(int(ppu.scanline)*160)+x] = pixelcolor
		}
	}
}

func (ppu *PPU) Render(cycles byte) bool {
	// lcd_on := cur_lcdc&0x1 != 0

	// if !lcd_on {
	// 	ppu.Reset()
	// 	return false
	// }

	new_dots := uint16(cycles * 4)
	row_done := (ppu.dots + new_dots) > 455
	ppu.dots = (ppu.dots + new_dots) % 456
	// ppu.logger.LogValue("dots", ppu.dots)

	cur_lcdc := Read(LCDC)
	cur_stat := Read(STAT)
	cur_mode := cur_stat & 0x3

	// if LCD is turned off
	// if cur_lcdc & 0x80 == 0 {
	// 	return
	// }

	if ppu.dots <= MODE2_END {
		if cur_mode != 2 {
			if cur_stat&0x20 != 0 {
				Write(STAT, (cur_stat&0xFE)|0x2)
				RequestInterrupt(1)
			}
		}
	} else if ppu.dots <= MODE3_END {
		if cur_mode != 3 {
			Write(STAT, (cur_stat&0xFC)|0x3)
		}
	} else if ppu.dots <= MODE0_END {
		if cur_mode != 0 {
			Write(STAT, (cur_stat & 0xFC))
			if cur_stat&0x8 != 0 {
				RequestInterrupt(1)
			}
		}
	}
	if !row_done {
		// fmt.Println("not render")
		return false
	}

	Write(LCDC, (cur_lcdc&(0xFC))|0x1)
	if cur_stat&0x10 != 0 {
		RequestInterrupt(1)
	}

	cur_row := Read(LY)
	Write(LY, (cur_row+1)%154)

	if cur_row >= 144 {
		if cur_mode != 1 {
			RequestInterrupt(0)
		}
		Write(STAT, (cur_stat&0xFC | 0x01))
	}

	// cur_lcdc := Read(LCDC)
	if cur_lcdc&0x8 == 0 {
		ppu.tilemap = 0x9800
	} else {
		ppu.tilemap = 0x9C00
	}

	if cur_lcdc&0x10 == 0 {
		ppu.tiledata = 0x8800
	} else {
		ppu.tiledata = 0x8000
	}

	// ppu.logger.LogValue("LY", uint16(cur_row))
	ppu.RenderBG(cur_row)
	if cur_row < 144 {
		ppu.scanline = (ppu.scanline + byte(1)) % 144
	}
	if cur_row == 144 {
		RequestInterrupt(0)
		Write(STAT, (cur_stat&0xFC | 0x01))
		return true
	}

	return false
}

func (ppu *PPU) Reset() {
	ppu.dots = 0
	ppu.scanline = 0
	ppu.tiledata = 0
	ppu.tilemap = 0

	framebufferPalette = [160 * 144]byte{}
	BGMapPalette = [256 * 256]byte{}
}
