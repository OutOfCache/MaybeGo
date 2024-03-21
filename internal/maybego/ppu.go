package maybego

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

const (
	LCDC uint16 = 0xFF40
	STAT uint16 = 0xFF41
	LY   uint16 = 0xFF44
	LYC  uint16 = 0xFF45
)

type PPU struct {
	tilemap  uint16
	tiledata uint16
	dots     uint16
	logger   *Logger
}

var framebufferRGBA [160 * 144]uint32
var bgMapRGBA [256 * 256]uint32

const defaultColor = uint32(0xFF8080FF)

var palette = []uint32{0x000000FF, 0x080808FF, 0x808080FF, 0xFFFFFFFF}

var winWidth, winHeight int32 = 160, 144
var err error

var gWindow *sdl.Window
var gRenderer *sdl.Renderer

// var gTexture *sdl.Texture
var gTextureA *sdl.Texture

func NewPPU(logger *Logger) *PPU {
	ppu := &PPU{logger: logger, dots: 0}

	return ppu
}

func (ppu *PPU) StartSDL() {
	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Printf("SDL could not initialize! Error: %s\n", err)
		os.Exit(4)
	}

	gWindow, err = sdl.CreateWindow("MaybeGo", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)

	if err != nil {
		fmt.Printf("Could not create Window! Error: %s\n", err)
		os.Exit(4)
	}

	gRenderer, err = sdl.CreateRenderer(gWindow, -1, sdl.RENDERER_ACCELERATED)

	if err != nil {
		fmt.Printf("Could not create Renderer. Error: %s\n", err)
		os.Exit(4)
	}

	gTextureA, err = gRenderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, 160, 144)
	if err != nil {
		fmt.Printf("Could not create Texture. Error: %s\n", err)
		os.Exit(4)
	}
}

func (ppu *PPU) RenderBG(row byte) {
	y := int(row)
	for j := 0; j < 256; j += 1 {
		x := j
		tileID := Read(ppu.tilemap + uint16((y/8)*32) + uint16(x/8))
		tileX := uint16(x / 8)

		var tileY uint16

		if ppu.tiledata == 0x8800 {
			tileY = uint16(0x800 + uint16(int8(tileID*0x10)))
		} else {
			tileY = uint16(tileID * 0x10)
		}
		address := ppu.tiledata + tileY + tileX

		pixelcolor := (Read(address) >> (7 - (x % 8)) & 0x1) +
			(Read(address+1)>>(7-(x%8))&0x1)*2
		bgMapRGBA[y*256+x] = palette[pixelcolor]
	}

}

func (ppu *PPU) RenderRow() {
	cur_row := Read(LY)
	cur_lyc := Read(LYC)

	if cur_row == cur_lyc {
		cur_stat := Read(STAT)
		Write(STAT, (cur_stat | 0x4))
		if cur_stat&0x40 == 0x40 {
			RequestInterrupt(1)
		}
	}
	if cur_row < 144 {
		ppu.RenderBG(cur_row)
	} else {
		cur_stat := Read(STAT)
		Write(STAT, (cur_stat&0xFC)|0x1)
		if cur_stat&0x10 != 0 {
			RequestInterrupt(1)
		}

		if cur_row == 144 {
			RequestInterrupt(0)
		}
	}
	Write(LY, (cur_row+1)%153)
}

func (ppu *PPU) Render(cycles byte) {
	gRenderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)

	render := (ppu.dots + uint16(cycles)) > 456
	ppu.dots = (ppu.dots + uint16(cycles)) % 457

	cur_lcdc := Read(LCDC)
	cur_stat := Read(STAT)

	if ppu.dots <= 80 {
		Write(STAT, (cur_stat&0xFE)|0x2)
		if cur_stat&0x20 != 0 {
			RequestInterrupt(1)
		}
	} else if ppu.dots <= (80 + 289) {
		Write(STAT, (cur_stat&0xFC)|0x3)
	} else if ppu.dots <= 456 {
		Write(STAT, (cur_stat & 0xFC))
		if cur_stat&0x8 != 0 {
			RequestInterrupt(1)
		}
	}

	if !render {
		return
	}

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

	RequestInterrupt(0)
	Write(LCDC, (cur_lcdc&(0xFC))|0x1)
	if cur_stat&0x10 != 0 {
		RequestInterrupt(1)
	}
}

func (ppu *PPU) EndSDL() {
	gTextureA.Destroy()
	gRenderer.Destroy()
	gWindow.Destroy()

	sdl.Quit()
}
