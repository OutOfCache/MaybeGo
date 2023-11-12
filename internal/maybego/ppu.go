package maybego

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

type PPU struct {
	logger *Logger
}

var framebufferRGBA [160 * 144]uint32

const defaultColor = uint32(0xFF8080FF)

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

	gTextureA, err = gRenderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, 160, 144)
	if err != nil {
		fmt.Printf("Could not create Texture. Error: %s\n", err)
		os.Exit(4)
	}
}

func (ppu *PPU) Render() {
	gRenderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)

	for row := 0; row < 144; row += 1 {
		for col := 0; col < 160; col += 1 {
			framebufferRGBA[row*160+col] = defaultColor
		}
	}
	gTextureA.UpdateRGBA(nil, framebufferRGBA[:], 160)
	gRenderer.Copy(gTextureA, nil, nil)
	gRenderer.Present()
}

func (ppu *PPU) EndSDL() {
	gTextureA.Destroy()
	gRenderer.Destroy()
	gWindow.Destroy()

	sdl.Quit()
}
