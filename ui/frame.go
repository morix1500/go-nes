package ui

import (
	"fmt"
	"go-nes/nes"
	"image"
	"image/color"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	WIDTH  = 256
	HEIGHT = 240
	SCALE  = 3
)

type Frame struct {
	Front   *image.RGBA
	Joypad1 *nes.Joypad
	Joypad2 *nes.Joypad
}

func NewFrame(joypad1, joypad2 *nes.Joypad) *Frame {
	return &Frame{
		Front:   image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT)),
		Joypad1: joypad1,
		Joypad2: joypad2,
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

	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
	gl.Enable(gl.TEXTURE_2D)

	return window
}

func (f *Frame) OnKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	f.Joypad1.ReadKeys(w)
	f.Joypad2.ReadKeys(w)
	if action == glfw.Press {
		if key == glfw.KeyEscape {
			w.SetShouldClose(true)
		}
	}
}

func (f *Frame) renderPixel(x, y uint, c color.RGBA) {
	f.Front.SetRGBA(int(x), int(y), c)
}

func (f *Frame) Render(ppu *nes.PPU) {
	scrollX := ppu.ReadPPUScrollX()
	scrollY := ppu.ReadPPUScrollY()

	var mainNameTable, secondNameTable []uint8
	if ppu.Mirroring == nes.MIRROR_VERTICAL &&
		(ppu.ReadCTRLNameTableAddress() == 0x2000 || ppu.ReadCTRLNameTableAddress() == 0x2800) {
		mainNameTable = ppu.VRAM[0:0x400]
		secondNameTable = ppu.VRAM[0x400:0x800]
	} else if ppu.Mirroring == nes.MIRROR_VERTICAL &&
		(ppu.ReadCTRLNameTableAddress() == 0x2400 || ppu.ReadCTRLNameTableAddress() == 0x2c00) {
		mainNameTable = ppu.VRAM[0x400:0x800]
		secondNameTable = ppu.VRAM[0:0x400]
	} else if ppu.Mirroring == nes.MIRROR_HORIZONTAL &&
		(ppu.ReadCTRLNameTableAddress() == 0x2000 || ppu.ReadCTRLNameTableAddress() == 0x2400) {
		mainNameTable = ppu.VRAM[0:0x400]
		secondNameTable = ppu.VRAM[0x400:0x800]
	} else if ppu.Mirroring == nes.MIRROR_HORIZONTAL &&
		(ppu.ReadCTRLNameTableAddress() == 0x2800 || ppu.ReadCTRLNameTableAddress() == 0x2c00) {
		mainNameTable = ppu.VRAM[0x400:0x800]
		secondNameTable = ppu.VRAM[0:0x400]
	} else {
		panic(fmt.Sprintf("not supported mirroring type: %d", ppu.Mirroring))
	}

	f.RenderNameTable(ppu, mainNameTable, *NewRect(uint(scrollX), uint(scrollY), WIDTH, HEIGHT), -int(scrollX), -int(scrollY))
	if scrollX > 0 {
		f.RenderNameTable(ppu, secondNameTable, *NewRect(0, 0, uint(scrollX), HEIGHT), WIDTH-int(scrollX), 0)
	} else if scrollY > 0 {
		f.RenderNameTable(ppu, secondNameTable, *NewRect(0, 0, WIDTH, uint(scrollY)), 0, HEIGHT-int(scrollY))
	}

	// Render Sprite
	for i := 0; i < len(ppu.OAMData); i += 4 {
		tileIndex := uint16(ppu.OAMData[i+1])
		tileX := uint(ppu.OAMData[i+3])
		tileY := uint(ppu.OAMData[i])

		var flipVertical, flipHorizontal bool
		if ppu.OAMData[i+2]>>7&1 == 1 {
			flipVertical = true
		} else {
			flipVertical = false
		}

		if ppu.OAMData[i+2]>>6&1 == 1 {
			flipHorizontal = true
		} else {
			flipHorizontal = false
		}
		palletteIndex := ppu.OAMData[i+2] & 0b11
		spritePallet := spritePallete(ppu, palletteIndex)

		bank := ppu.ReadCTRLSpriteTableAddress()

		tile := ppu.CharacterRom[(bank + tileIndex*16):(bank + tileIndex*16 + 15 + 1)]

		for y := uint(0); y < 8; y++ {
			upper := tile[y]
			lower := tile[y+8]
			for x := 7; x >= 0; x-- {
				value := (1&lower)<<1 | (1 & upper)
				upper = upper >> 1
				lower = lower >> 1
				var rgb color.RGBA
				switch value {
				case 0:
					continue
				case 1:
					rgb = nes.Palletes[spritePallet[1]]
				case 2:
					rgb = nes.Palletes[spritePallet[2]]
				case 3:
					rgb = nes.Palletes[spritePallet[3]]
				default:
					panic("unknown value")
				}

				if flipHorizontal && flipVertical {
					f.renderPixel(tileX+7-uint(x), tileY+7-y, rgb)
				} else if flipHorizontal && !flipVertical {
					f.renderPixel(tileX+7-uint(x), tileY+y, rgb)
				} else if !flipHorizontal && flipVertical {
					f.renderPixel(tileX+uint(x), tileY+7-y, rgb)
				} else {
					f.renderPixel(tileX+uint(x), tileY+y, rgb)
				}
			}
		}
	}
}

func backgroundPallette(ppu *nes.PPU, attributeTable []uint8, tileColumn uint, tileRow uint) []uint8 {
	attrTableIdx := (tileRow/4)*8 + tileColumn/4
	attrByte := attributeTable[attrTableIdx]

	var palletIndex uint8
	palletIdx1 := tileColumn % 4 / 2
	palletIdx2 := tileRow % 4 / 2

	if palletIdx1 == 0 && palletIdx2 == 0 {
		palletIndex = attrByte & 0b11
	} else if palletIdx1 == 1 && palletIdx2 == 0 {
		palletIndex = (attrByte >> 2) & 0b11
	} else if palletIdx1 == 0 && palletIdx2 == 1 {
		palletIndex = (attrByte >> 4) & 0b11
	} else {
		palletIndex = (attrByte >> 6) & 0b11
	}

	palletStart := 1 + (palletIndex * 4)

	return []uint8{
		ppu.PaletteTable[0],
		ppu.PaletteTable[palletStart],
		ppu.PaletteTable[palletStart+1],
		ppu.PaletteTable[palletStart+2],
	}
}

func spritePallete(ppu *nes.PPU, palleteIndex uint8) []uint8 {
	palletStart := 0x11 + (palleteIndex * 4)

	return []uint8{
		0,
		ppu.PaletteTable[palletStart],
		ppu.PaletteTable[palletStart+1],
		ppu.PaletteTable[palletStart+2],
	}
}

type Rect struct {
	x1 uint
	y1 uint
	x2 uint
	y2 uint
}

func NewRect(x1, y1, x2, y2 uint) *Rect {
	return &Rect{
		x1: x1,
		y1: y1,
		x2: x2,
		y2: y2,
	}
}

func (f *Frame) RenderNameTable(ppu *nes.PPU, nametable []uint8, viewPort Rect, shiftX int, shiftY int) {
	bank := ppu.ReadCTRLBackGroundTableAddress()
	attributeTable := nametable[0x3c0:0x400]

	for i := 0; i < 0x3c0; i++ {
		tileColumn := uint(i % 32)
		tileRow := uint(i / 32)
		tileIndex := uint16(nametable[i])
		tile := ppu.CharacterRom[(bank + tileIndex*16):(bank + tileIndex*16 + 15 + 1)]
		palette := backgroundPallette(ppu, attributeTable, tileColumn, tileRow)

		for y := uint(0); y < 8; y++ {
			upper := tile[y]
			lower := tile[y+8]

			for x := 7; x >= 0; x-- {
				value := (1&lower)<<1 | (1 & upper)
				upper = upper >> 1
				lower = lower >> 1
				var rgb color.RGBA
				switch value {
				case 0:
					rgb = nes.Palletes[ppu.PaletteTable[0]]
				case 1:
					rgb = nes.Palletes[palette[1]]
				case 2:
					rgb = nes.Palletes[palette[2]]
				case 3:
					rgb = nes.Palletes[palette[3]]
				default:
					panic("unknown value")
				}
				pixelX := tileColumn*8 + uint(x)
				pixelY := tileRow*8 + y

				if pixelX >= viewPort.x1 && pixelX < viewPort.x2 && pixelY >= viewPort.y1 && pixelY < viewPort.y2 {
					f.renderPixel(uint(shiftX+int(pixelX)), uint(shiftY+int(pixelY)), rgb)
				}
			}
		}
	}
}
