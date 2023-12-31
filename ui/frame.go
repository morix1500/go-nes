package ui

import (
	"go-nes/nes"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	WIDTH  = 256
	HEIGHT = 240
	SCALE  = 3
)

type Frame struct {
	Front *image.RGBA
}

func NewFrame() *Frame {
	return &Frame{
		Front: image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT)),
	}
}

func Init() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(WIDTH*SCALE, HEIGHT*SCALE, "Morix NES", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetKeyCallback(onKey)

	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
	gl.Enable(gl.TEXTURE_2D)

	return window
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyEscape:
			w.SetShouldClose(true)
			os.Exit(0)
		}
	}
}

func (f *Frame) renderPixel(x, y uint, c color.RGBA) {
	f.Front.SetRGBA(int(x), int(y), c)
}

func (f *Frame) Render(ppu *nes.PPU) {
	bank := ppu.ReadCTRLBackGroundTableAddress()

	for i := 0; i <= 0x03c0; i++ {
		tileAddr := uint16(ppu.VRAM[i])
		tileX := uint(i % 32)
		tileY := uint(i / 32)
		tile := ppu.CharacterRom[(bank + tileAddr*16):(bank + tileAddr*16 + 15 + 1)]

		for y := uint(0); y < 8; y++ {
			upper := tile[y]
			lower := tile[y+8]

			for x := uint(0); x < 8; x++ {
				value := (1&lower)<<1 | (1 & upper)
				upper = upper >> 1
				lower = lower >> 1
				var rgb color.RGBA
				switch value {
				case 0:
					rgb = nes.Palletes[0x01]
				case 1:
					rgb = nes.Palletes[0x23]
				case 2:
					rgb = nes.Palletes[0x27]
				case 3:
					rgb = nes.Palletes[0x30]
				default:
					panic("unknown value")
				}
				tmpx := 7 - x + (tileX * 8)
				tmpy := tileY*8 + y
				f.renderPixel(tmpx, tmpy, rgb)
			}
		}
	}
}
