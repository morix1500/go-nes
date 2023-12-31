package ui

import (
	"go-nes/nes"
	"image/color"
)

func Render(ppu *nes.PPU, frame *Frame) {
	bank := ppu.ReadCTRLBackGroundTableAddress()
	index := uint(0)

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
				//var r, g, b uint32
				var rgb color.RGBA
				switch value {
				case 0:
					//r, g, b, _ = nes.Palletes[0x01].RGBA()
					rgb = nes.Palletes[0x01]
				case 1:
					//r, g, b, _ = nes.Palletes[0x23].RGBA()
					rgb = nes.Palletes[0x23]
				case 2:
					//r, g, b, _ = nes.Palletes[0x27].RGBA()
					rgb = nes.Palletes[0x27]
				case 3:
					//r, g, b, _ = nes.Palletes[0x30].RGBA()
					rgb = nes.Palletes[0x30]
				default:
					panic("unknown value")
				}
				tmpx := 7 - x + (tileX * 8)
				tmpy := tileY*8 + y
				//frame.SetPixel(index, uint(tmpx), uint(tmpy), uint8(r), uint8(g), uint8(b))
				frame.renderPixel(tmpx, tmpy, rgb)
			}
		}
		index++
	}
}
