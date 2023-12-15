package nes

import "fmt"

const (
	// Control Register
	// 7  bit  0
	// ---- ----
	// VPHB SINN
	// |||| ||||
	// |||| ||++- Base nametable address
	// |||| ||    (0 = $2000; 1 = $2400; 2 = $2800; 3 = $2C00)
	// |||| |+--- VRAM address increment per CPU read/write of PPUDATA
	// |||| |     (0: add 1, going across; 1: add 32, going down)
	// |||| +---- Sprite pattern table address for 8x8 sprites
	// ||||       (0: $0000; 1: $1000; ignored in 8x16 mode)
	// |||+------ Background pattern table address (0: $0000; 1: $1000)
	// ||+------- Sprite size (0: 8x8 pixels; 1: 8x16 pixels)
	// |+-------- PPU master/slave select
	// |          (0: read backdrop from EXT pins; 1: output color on EXT pins)
	// +--------- Generate an NMI at the start of the
	//            vertical blanking interval (0: off; 1: on)
	CR_NAMETABLE1              uint8 = 0b0000_0001
	CR_NAMETABLE2              uint8 = 0b0000_0010
	CR_VRAM_ADD_INCREMENT      uint8 = 0b0000_0100
	CR_SPRITE_PATTERN_ADDR     uint8 = 0b0000_1000
	CR_BACKGROUND_PATTERN_ADDR uint8 = 0b0001_0000
	CR_SPRITE_SIZE             uint8 = 0b0010_0000
	CR_MASTER_SLAVE_SELECT     uint8 = 0b0100_0000
	CR_GENERATE_NMI            uint8 = 0b1000_0000
)

type PPU struct {
	CharacterRom       []uint8
	PaletteTable       [32]uint8
	VRAM               [2048]uint8
	OAMData            [256]uint8
	InternalDataBuffer uint8
	Mirroring          Mirroring
	Addr               AddrRegister
	Ctrl               ControlRegister
}

func NewPPU(characterRom []uint8, mirroring Mirroring) *PPU {
	return &PPU{
		CharacterRom: characterRom,
		PaletteTable: [32]uint8{},
		VRAM:         [2048]uint8{},
		OAMData:      [256]uint8{},
		Mirroring:    mirroring,
		Addr:         *NewAddrRegister(),
		Ctrl:         *NewControlRegister(),
	}
}

func (p *PPU) WriteToPPUAddr(value uint8) {
	p.Addr.Update(value)
}

func (p *PPU) WriteToPPUCTRL(value uint8) {
	p.Ctrl.Update(value)
}

func (p *PPU) incrementVRAMAddr() {
	p.Addr.Increment(p.Ctrl.VRAMAddrIncrement())
}

func (p *PPU) ReadData() uint8 {
	addr := p.Addr.Get()
	p.incrementVRAMAddr()

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
	addr := p.Addr.Get()

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

	p.incrementVRAMAddr()
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
	if p.Mirroring == MIRROR_VERTICAL && nameTable == 2 {
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

type AddrRegister struct {
	High  uint8
	Low   uint8
	HiPrt bool
}

func NewAddrRegister() *AddrRegister {
	return &AddrRegister{
		High:  0,
		Low:   0,
		HiPrt: true,
	}
}

func (a *AddrRegister) set(data uint16) {
	a.High = uint8(data >> 8)
	a.Low = uint8(data & 0xff)
}

func (a *AddrRegister) Update(data uint8) {
	if a.HiPrt {
		a.High = data
	} else {
		a.Low = data
	}
	// mirror down addr above 0x3fff
	if a.Get() > 0x3fff {
		a.set(a.Get() & 0x3fff)
	}
	a.HiPrt = !a.HiPrt
}

func (a *AddrRegister) Increment(inc uint8) {
	lo := a.Low
	a.Low = a.Low + inc
	if lo > a.Low {
		a.High++
	}
	if a.Get() > 0x3fff {
		a.set(a.Get() & 0x3fff) // mirror down addr above 0x3fff
	}
}

func (a *AddrRegister) ResetLatch() {
	a.HiPrt = true
}

func (a *AddrRegister) Get() uint16 {
	return uint16(a.High)<<8 | uint16(a.Low)
}

type ControlRegister struct {
	Bits uint8
}

func NewControlRegister() *ControlRegister {
	return &ControlRegister{
		Bits: 0,
	}
}

func (c *ControlRegister) VRAMAddrIncrement() uint8 {
	if c.Bits&CR_VRAM_ADD_INCREMENT == 0 {
		return 1
	} else {
		return 32
	}
}

func (c *ControlRegister) Update(data uint8) {
	c.Bits = data
}
