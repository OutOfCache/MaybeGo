package maybego

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

type PPU struct {
	logger *Logger
}

var winWidth, winHeight int32 = 160, 144
var err error

var gWindow *sdl.Window
var gRenderer *sdl.Renderer

// var gTexture *sdl.Texture
var gTextureA *sdl.Texture

func NewPPU(logger *Logger) *PPU {
	ppu := &PPU{logger: logger}

	return ppu
}

func (ppu *PPU) StartSDL() {
	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Println("SDL could not initialize! Error: %s\n", err)
		os.Exit(4)
	}

	gWindow, err = sdl.CreateWindow("MaybeGo", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)

	if err != nil {
		fmt.Println("Could not create Window! Error: %s\n", err)
		os.Exit(4)
	}

	gRenderer, err = sdl.CreateRenderer(gWindow, -1, sdl.RENDERER_ACCELERATED)

	if err != nil {
		fmt.Printf("Could not create Renderer. Error: %s\n", err)
		os.Exit(4)
	}
}

func (ppu *PPU) Render() {
	gRenderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
	gRenderer.Clear()
	gRenderer.Present()
}

func (ppu *PPU) EndSDL() {
	gRenderer.Destroy()
	gWindow.Destroy()

	sdl.Quit()
}
