package nes

import "fmt"

type PPU struct {
	CharacterRom       []uint8
	PaletteTable       [32]uint8
	VRAM               [2048]uint8
	InternalDataBuffer uint8
	Mirroring          Mirroring
	Ctrl               ControlRegister
	Mask               MaskRegister
	Status             StatusRegister
	OAMAddress         uint8
	OAMData            [256]uint8

	// PPU internal registers
	v uint16 // current vram address(15bit)
	t uint16 // temporary vram address(15bit)
	x uint8  // fine x scroll(3bit)
	w uint8  // write toggle(1bit)
	f uint8  // even/odd frame flag(1bit)
}

func NewPPU(characterRom []uint8, mirroring Mirroring) *PPU {
	return &PPU{
		CharacterRom: characterRom,
		PaletteTable: [32]uint8{},
		VRAM:         [2048]uint8{},
		OAMData:      [256]uint8{},
		Mirroring:    mirroring,
		Ctrl:         *NewControlRegister(),
		Status:       *NewStatusRegister(),
		Mask:         *NewMaskRegister(),
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
	p.Ctrl.Update(value)
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

func (p *PPU) ReadData() uint8 {
	addr := p.v
	p.v += uint16(p.Ctrl.VRAMAddrIncrement())

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

	p.v += uint16(p.Ctrl.VRAMAddrIncrement())
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
	result := p.Status.Bits
	p.Status.SetVblankStatus(false)
	p.w = 0
	return result
}

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

const (
	// Status Register
	// 7  bit  0
	// ---- ----
	// VSO. ....
	// |||| ||||
	// |||+-++++- Least significant bits previously written into a PPU register
	// |||        (due to register not being updated for this address)
	// ||+------- Sprite overflow. The intent was for this flag to be set
	// ||         whenever more than eight sprites appear on a scanline, but a
	// ||         hardware bug causes the actual behavior to be more complicated
	// ||         and generate false positives as well as false negatives; see
	// ||         PPU sprite evaluation. This flag is set during sprite
	// ||         evaluation and cleared at dot 1 (the second dot) of the
	// ||         pre-render line.
	// |+-------- Sprite 0 Hit.  Set when a nonzero pixel of sprite 0 overlaps
	// |          a nonzero background pixel; cleared at dot 1 of the pre-render
	// |          line.  Used for raster timing.
	// +--------- Vertical blank has started (0: not in vblank; 1: in vblank).
	//            Set at dot 1 of line 241 (the line *after* the post-render
	//            line); cleared after reading $2002 and at dot 1 of the
	//            pre-render line.
	S_NOTUSED         uint8 = 0b0000_0001
	S_NOTUSED2        uint8 = 0b0000_0010
	S_NOTUSED3        uint8 = 0b0000_0100
	S_NOTUSED4        uint8 = 0b0000_1000
	S_NOTUSED5        uint8 = 0b0001_0000
	S_SPRITE_OVERFLOW uint8 = 0b0010_0000
	S_SPRITE_ZERO_HIT uint8 = 0b0100_0000
	S_VBLANK_STARTED  uint8 = 0b1000_0000
)

type StatusRegister struct {
	Bits uint8
}

func NewStatusRegister() *StatusRegister {
	return &StatusRegister{
		Bits: 0,
	}
}

func (s *StatusRegister) SetVblankStatus(status bool) {
	if status {
		s.Bits = s.Bits | S_VBLANK_STARTED
	} else {
		s.Bits = s.Bits &^ S_VBLANK_STARTED
	}
}

func (s *StatusRegister) SetSpriteZeroHitStatus(status bool) {
	if status {
		s.Bits = s.Bits | S_SPRITE_ZERO_HIT
	} else {
		s.Bits = s.Bits &^ S_SPRITE_ZERO_HIT
	}
}

func (s *StatusRegister) SetSpriteOverflowStatus(status bool) {
	if status {
		s.Bits = s.Bits | S_SPRITE_OVERFLOW
	} else {
		s.Bits = s.Bits &^ S_SPRITE_OVERFLOW
	}
}

const (
	// 7  bit  0
	// ---- ----
	// BGRs bMmG
	// |||| ||||
	// |||| |||+- Greyscale (0: normal color, 1: produce a greyscale display)
	// |||| ||+-- 1: Show background in leftmost 8 pixels of screen, 0: Hide
	// |||| |+--- 1: Show sprites in leftmost 8 pixels of screen, 0: Hide
	// |||| +---- 1: Show background
	// |||+------ 1: Show sprites
	// ||+------- Emphasize red (green on PAL/Dendy)
	// |+-------- Emphasize green (red on PAL/Dendy)
	// +--------- Emphasize blue
	MASK_GREY_SCALE               uint8 = 0b0000_0001
	MASK_LEFTMOST_8PXL_BACKGROUND uint8 = 0b0000_0010
	MASK_LEFTMOST_8PXL_SPRITE     uint8 = 0b0000_0100
	MASK_SHOW_BACKGROUND          uint8 = 0b0000_1000
	MASK_SHOW_SPRITES             uint8 = 0b0001_0000
	MASK_ENPHASIZE_RED            uint8 = 0b0010_0000
	MASK_ENPHASIZE_GREEN          uint8 = 0b0100_0000
	MASK_EMPHASIZE_BLUE           uint8 = 0b1000_0000
)

type MaskRegister struct {
	Bits uint8
}

func NewMaskRegister() *MaskRegister {
	return &MaskRegister{
		Bits: 0,
	}
}

func (m *MaskRegister) Update(data uint8) {
	m.Bits = data
}
