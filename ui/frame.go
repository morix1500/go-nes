package ui

import (
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

var (
	keyMap = map[glfw.Key]uint8{
		glfw.KeyA:     nes.JOYPAD_BUTTON_A,
		glfw.KeyS:     nes.JOYPAD_BUTTON_B,
		glfw.KeyEnter: nes.JOYPAD_START,
		glfw.KeySpace: nes.JOYPAD_SELECT,
		glfw.KeyUp:    nes.JOYPAD_UP,
		glfw.KeyDown:  nes.JOYPAD_DOWN,
		glfw.KeyLeft:  nes.JOYPAD_LEFT,
		glfw.KeyRight: nes.JOYPAD_RIGHT,
	}
)

type Frame struct {
	Front  *image.RGBA
	Joypad *nes.Joypad
}

func NewFrame(joypad *nes.Joypad) *Frame {
	return &Frame{
		Front:  image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT)),
		Joypad: joypad,
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
	if action == glfw.Press {
		if _, exist := keyMap[key]; exist {
			f.Joypad.Press(keyMap[key])
		} else if key == glfw.KeyEscape {
			w.SetShouldClose(true)
		}
	}
	if action == glfw.Release {
		if _, exist := keyMap[key]; exist {
			f.Joypad.Release(keyMap[key])
		}
	}
}

func (f *Frame) renderPixel(x, y uint, c color.RGBA) {
	f.Front.SetRGBA(int(x), int(y), c)
}

func (f *Frame) Render(ppu *nes.PPU) {
	bank := ppu.ReadCTRLBackGroundTableAddress()

	// Render background
	for i := 0; i <= 0x03c0; i++ {
		tileAddr := uint16(ppu.VRAM[i])
		tileColumn := uint(i % 32)
		tileRow := uint(i / 32)
		tile := ppu.CharacterRom[(bank + tileAddr*16):(bank + tileAddr*16 + 15 + 1)]
		palette := backgroundPallette(ppu, tileColumn, tileRow)

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
				tmpx := (tileColumn * 8) + uint(x)
				tmpy := tileRow*8 + y
				f.renderPixel(tmpx, tmpy, rgb)
			}
		}
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

func backgroundPallette(ppu *nes.PPU, tileColumn uint, tileRow uint) []uint8 {
	attrTableIdx := (tileRow/4)*8 + tileColumn/4
	attrByte := ppu.VRAM[0x03c0+attrTableIdx]

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
