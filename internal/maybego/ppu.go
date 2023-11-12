package maybego

import (
	"github.com/veandco/go-sdl2/sdl"
)

type PPU struct {
	logger *Logger
}

var err error

func NewPPU(logger *Logger) *PPU {
	ppu := &PPU{logger: logger}

	return ppu
}

func (ppu *PPU) StartSDL() {
	err = sdl.Init(sdl.INIT_EVERYTHING)
}

func (ppu *PPU) EndSDL() {
	sdl.Quit()
}
