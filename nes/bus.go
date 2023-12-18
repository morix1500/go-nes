package nes

import "fmt"

//  _______________ $10000  _______________
// | PRG-ROM       |       |               |
// | Upper Bank    |       |               |
// |_ _ _ _ _ _ _ _| $C000 | PRG-ROM       |
// | PRG-ROM       |       |               |
// | Lower Bank    |       |               |
// |_______________| $8000 |_______________|
// | SRAM          |       | SRAM          |
// |_______________| $6000 |_______________|
// | Expansion ROM |       | Expansion ROM |
// |_______________| $4020 |_______________|
// | I/O Registers |       |               |
// |_ _ _ _ _ _ _ _| $4000 |               |
// | Mirrors       |       | I/O Registers |
// | $2000-$2007   |       |               |
// |_ _ _ _ _ _ _ _| $2008 |               |
// | I/O Registers |       |               |
// |_______________| $2000 |_______________|
// | Mirrors       |       |               |
// | $0000-$07FF   |       |               |
// |_ _ _ _ _ _ _ _| $0800 |               |
// | RAM           |       | RAM           |
// |_ _ _ _ _ _ _ _| $0200 |               |
// | Stack         |       |               |
// |_ _ _ _ _ _ _ _| $0100 |               |
// | Zero Page     |       |               |
// |_______________| $0000 |_______________|

type Bus struct {
	CpuVRAM   [2048]uint8
	Cartridge *Cartridge
	PPU       *PPU
}

const (
	RAM                       uint16 = 0x0000
	RAM_MIRRORS_END           uint16 = 0x1fff
	PPU_REGISTERS             uint16 = 0x2000
	PPU_REGISTERS_MIRRORS_END uint16 = 0x3fff
)

func NewBus(cartridge *Cartridge) *Bus {
	ppu := NewPPU(cartridge.CharacterRom, cartridge.ScreenMirroring)
	return &Bus{
		CpuVRAM:   [2048]uint8{},
		Cartridge: cartridge,
		PPU:       ppu,
	}
}

func (b *Bus) ReadMemory(addr uint16) uint8 {
	if addr >= RAM && addr <= RAM_MIRRORS_END {
		mirrorDownAddr := addr & 0b111_1111_1111
		return b.CpuVRAM[mirrorDownAddr]
	} else if addr == 0x2000 || addr == 0x2001 || addr == 0x2003 || addr == 0x2005 || addr == 0x2006 || addr == 0x4014 {
		panic(fmt.Sprintf("attempt to read from write-only PPU address: %d", addr))
	} else if addr == 0x2002 {
		return b.PPU.ReadStatus()
	} else if addr == 0x2004 {
		return b.PPU.ReadOAMData()
	} else if addr == 0x2007 {
		return b.PPU.ReadData()
	} else if addr >= 0x2008 && addr <= PPU_REGISTERS_MIRRORS_END {
		mirrorDownAddr := addr & 0b00100000_00000111
		return b.ReadMemory(mirrorDownAddr)
	} else if addr >= 0x8000 && addr <= 0xFFFF {
		return b.ReadProgramRom(addr)
	}
	return 0
}

func (b *Bus) WriteMemory(addr uint16, data uint8) {
	if addr >= RAM && addr <= RAM_MIRRORS_END {
		mirrorDownAddr := addr & 0b111_1111_1111
		b.CpuVRAM[mirrorDownAddr] = data
	} else if addr == 0x2000 {
		b.PPU.WriteToPPUCTRL(data)
	} else if addr == 0x2001 {
		b.PPU.WriteToPPUMask(data)
	} else if addr == 0x2003 {
		b.PPU.WriteToPPUOAMAddr(data)
	} else if addr == 0x2004 {
		b.PPU.WriteToPPUOAMData(data)
	} else if addr == 0x2005 {
		b.PPU.WriteToPPUScroll(data)
	} else if addr == 0x2006 {
		b.PPU.WriteToPPUAddr(data)
	} else if addr == 0x2007 {
		b.PPU.WriteData(data)
	} else if addr == 0x4014 {
		hi := uint16(data) << 8
		buf := make([]uint8, 256)
		for i := 0; i < 256; i++ {
			buf[i] = b.ReadMemory(hi + uint16(i))
		}
		b.PPU.WriteOAMDMA(buf)
	} else if addr >= 0x2008 && addr <= PPU_REGISTERS_MIRRORS_END {
		mirrorDownAddr := addr & 0b00100000_00000111
		b.CpuVRAM[mirrorDownAddr] = data
	} else if addr >= 0x8000 && addr <= 0xFFFF {
		panic("Attempt to write to Cartridge ROM space")
	}
}

func (b *Bus) ReadProgramRom(addr uint16) uint8 {
	addr -= 0x8000

	// プログラムROMは16kbまたは32kbのいずれか。なぜならマップアドレススペースが32kbのため、ROMが16kbの場合は上位16kbを下位16kbにミラーする必要がある
	if len(b.Cartridge.ProgramRom) == 0x4000 && addr >= 0x4000 {
		// mirror if needed
		addr = addr % 0x4000
	}
	return b.Cartridge.ProgramRom[addr]
}
