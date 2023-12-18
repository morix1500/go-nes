package nes

import (
	"fmt"
)

type PPU struct {
	CharacterRom       []uint8
	PaletteTable       [32]uint8
	VRAM               [2048]uint8
	InternalDataBuffer uint8
	Mirroring          Mirroring
	OAMAddress         uint8
	OAMData            [256]uint8

	// PPU internal registers
	v uint16 // current vram address(15bit)
	t uint16 // temporary vram address(15bit)
	x uint8  // fine x scroll(3bit)
	w uint8  // write toggle(1bit)
	f uint8  // even/odd frame flag(1bit)

	// $2000 PPUCTRL
	flagNameTable       uint8 // 0: $2000; 1: $2400; 2: $2800; 3: $2C00
	flagIncrement       uint8 // 0: add 1; 1: add 32
	flagSpriteTable     uint8 // 0: $0000; 1: $1000; ignored in 8x16 mode
	flagBackgroundTable uint8 // 0: $0000; 1: $1000
	flagSpriteSize      uint8 // 0: 8x8; 1: 8x16
	flagMasterSlave     uint8 // 0: read EXT; 1: write EXT
	flagNMI             bool  // 0: off; 1: on

	// $2001 PPUMASK
	flagGrayscale                 uint8 // 0: normal color; 1: grayscale
	flagShowBackgroundLeftMost8px uint8 // 0: hide; 1: show
	flagShowSpriteLeftMost8px     uint8 // 0: hide; 1: show
	flagShowBackground            uint8 // 0: hide; 1: show
	flagShowSprite                uint8 // 0: hide; 1: show
	flagEmphasizeRed              uint8 // 1: emphasize
	flagEmphasizeGreen            uint8 // 1: emphasize
	flagEmphasizeBlue             uint8 // 1: emphasize

	// $2002 PPUSTATUS
	flagSpriteOverflow uint8
	flagSpriteZeroHit  uint8
	flagVblankStarted  uint8
}

func NewPPU(characterRom []uint8, mirroring Mirroring) *PPU {
	return &PPU{
		CharacterRom: characterRom,
		PaletteTable: [32]uint8{},
		VRAM:         [2048]uint8{},
		OAMData:      [256]uint8{},
		Mirroring:    mirroring,
		w:            0,
	}
}

func (p *PPU) WriteToPPUAddr(value uint8) {
	// 15 14 13 12 11 10 9 8 7 6 5 4 3 2 1 0
	// -  0  h  h  h  h  h h l l l l l l l l

	if p.w == 0 {
		// high byte
		p.t = (p.t & 0b1000_0000_1111_1111) | ((uint16(value) & 0b0011_1111) << 8)
		p.w = 1
	} else {
		// low byte
		p.t = (p.t & 0b1111_1111_0000_0000) | (uint16(value))
		p.v = p.t
		p.w = 0
	}
}

// $2005 PPU scroll register
func (p *PPU) WriteToPPUScroll(value uint8) {
	// 15 14 13 12 11 10 9 8 7 6 5 4 3 2 1 0
	// -  y  y  y  -  -  y y y y y x x x x x

	if p.w == 0 {
		// X scroll
		// 1回目の書き込みでは、値の0-2ビット目がxレジスタに入り、3-8ビット目がtレジスタの0-4ビット目に入る
		p.x = value & 0b0000_0111
		p.t = (p.t & 0b1111_1111_1110_0000) | (uint16(value) >> 3)
		p.w = 1
	} else {
		// Y scroll
		// 2回目の書き込みでは、値の0-2ビット目がtレジスタの12-14ビット目に入り、3-8ビット目がtレジスタの5-9ビット目に入る
		p.t = (p.t & 0b1000_1111_1111_1111) | ((uint16(value) & 0b0000_0111) << 12)
		p.t = (p.t & 0b1111_1100_0001_1111) | ((uint16(value) & 0b1111_1000) << 2)
		p.w = 0
	}
}

func (p *PPU) WriteToPPUCTRL(value uint8) {
	p.flagNameTable = (value & 0b0000_0011)
	p.flagIncrement = (value & 0b0000_0100) >> 2
	p.flagSpriteTable = (value & 0b0000_1000) >> 3
	p.flagBackgroundTable = (value & 0b0001_0000) >> 4
	p.flagSpriteSize = (value & 0b0010_0000) >> 5
	p.flagMasterSlave = (value & 0b0100_0000) >> 6
	p.flagNMI = (value & 0b1000_0000) != 0

	p.t = (p.t & 0b1111_0011_1111_1111) | ((uint16(value) & 0b0000_0011) << 10)
}

func (p *PPU) WriteToPPUMask(value uint8) {
	p.flagGrayscale = (value & 0b0000_0001)
	p.flagShowBackgroundLeftMost8px = (value & 0b0000_0010) >> 1
	p.flagShowSpriteLeftMost8px = (value & 0b0000_0100) >> 2
	p.flagShowBackground = (value & 0b0000_1000) >> 3
	p.flagShowSprite = (value & 0b0001_0000) >> 4
	p.flagEmphasizeRed = (value & 0b0010_0000) >> 5
	p.flagEmphasizeGreen = (value & 0b0100_0000) >> 6
	p.flagEmphasizeBlue = (value & 0b1000_0000) >> 7
}

func (p *PPU) WriteToOAMAddr(value uint8) {
	p.OAMAddress = value
}

func (p *PPU) WriteToOAMData(value uint8) {
	p.OAMData[p.OAMAddress] = value
	p.OAMAddress++
}

func (p *PPU) WriteOAMDMA(data []uint8) {
	for i := 0; i < len(data); i++ {
		p.OAMData[p.OAMAddress] = data[i]
		p.OAMAddress++
	}
}

func (p *PPU) ReadOAMData() uint8 {
	return p.OAMData[p.OAMAddress]
}

func (p *PPU) VRAMAddrIncrement() uint8 {
	if p.flagIncrement == 0 {
		return 1
	} else {
		return 32
	}
}

func (p *PPU) ReadData() uint8 {
	addr := p.v
	p.v += uint16(p.VRAMAddrIncrement())

	var result uint8

	if addr <= 0x1fff {
		result = p.InternalDataBuffer
		p.InternalDataBuffer = p.CharacterRom[addr]
	} else if addr >= 0x2000 && addr <= 0x2fff {
		result = p.InternalDataBuffer
		p.InternalDataBuffer = p.VRAM[p.mirrorVRAMAddr(addr)]
	} else if addr >= 0x3000 && addr <= 0x3eff {
		panic(fmt.Sprintf("addr space 0x3000..0x3eff is not expected to be used, requested = %d", addr))
	} else if addr >= 0x3f00 && addr <= 0x3fff {
		result = p.PaletteTable[addr-0x3f00]
	} else {
		panic(fmt.Sprintf("unexpected access to mirrored space %d", addr))
	}
	return result
}

func (p *PPU) WriteData(value uint8) {
	addr := p.v

	if addr <= 0x1fff {
		panic(fmt.Sprintf("attempt to write to character rom space: %d", addr))
	} else if addr >= 0x2000 && addr <= 0x2fff {
		p.VRAM[p.mirrorVRAMAddr(addr)] = value
	} else if addr >= 0x3000 && addr <= 0x3eff {
		fmt.Println("not implemented")
	} else if addr == 0x3f10 || addr == 0x3f14 || addr == 0x3f18 || addr == 0x3f1c {
		addMirror := addr - 0x10
		p.PaletteTable[addMirror-0x3f00] = value
	} else if addr >= 0x3f00 && addr <= 0x3fff {
		p.PaletteTable[addr-0x3f00] = value
	} else {
		panic(fmt.Sprintf("unexpected access to mirrored space %d", addr))
	}

	p.v += uint16(p.VRAMAddrIncrement())
}

// Horizontal:
//
//	[ A ] [ a ]
//	[ B ] [ b ]
//
// Vertical:
//
//	[ A ] [ B ]
//	[ a ] [ b ]
func (p *PPU) mirrorVRAMAddr(addr uint16) uint16 {
	// mirror down 0x3000-0x3eff to 0x2000 - 0x2eff
	mirroredVram := addr & 0x2fff
	vramIndex := mirroredVram - 0x2000 // to vram vector
	nameTable := vramIndex / 0x400     // to the name table index

	var result uint16
	if p.Mirroring == MIRROR_VERTICAL && (nameTable == 2 || nameTable == 3) {
		result = vramIndex - 0x800
	} else if p.Mirroring == MIRROR_HORIZONTAL && nameTable == 2 {
		result = vramIndex - 0x400
	} else if p.Mirroring == MIRROR_HORIZONTAL && nameTable == 1 {
		result = vramIndex - 0x400
	} else if p.Mirroring == MIRROR_HORIZONTAL && nameTable == 3 {
		result = vramIndex - 0x800
	} else {
		result = vramIndex
	}
	return result
}

func (p *PPU) ReadStatus() uint8 {
	var result uint8
	result = result | (p.flagSpriteOverflow << 5)
	result = result | (p.flagSpriteZeroHit << 6)
	result = result | (p.flagVblankStarted << 7)
	p.flagVblankStarted = 0
	p.w = 0
	return result
}
