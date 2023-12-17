package nes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPPUVramWrites(t *testing.T) {
	ppu := NewPPU(nil, MIRROR_HORIZONTAL)
	ppu.WriteToPPUAddr(0x23)
	ppu.WriteToPPUAddr(0x05)
	ppu.WriteData(0x66)

	assert.Equal(t, uint8(0x66), ppu.VRAM[0x0305])
}

func TestPPUVramReads(t *testing.T) {
	ppu := NewPPU(nil, MIRROR_HORIZONTAL)
	ppu.WriteToPPUCTRL(0x00)
	ppu.VRAM[0x0305] = 0x66

	ppu.WriteToPPUAddr(0x23)
	ppu.WriteToPPUAddr(0x05)

	ppu.ReadData()
	assert.Equal(t, uint16(0x2306), ppu.Addr.Get())
	assert.Equal(t, uint8(0x66), ppu.ReadData())
}

func TestPPUVramReadsCrossPage(t *testing.T) {
	ppu := NewPPU(nil, MIRROR_HORIZONTAL)
	ppu.WriteToPPUCTRL(0x00)
	ppu.VRAM[0x01ff] = 0x66
	ppu.VRAM[0x0200] = 0x77

	ppu.WriteToPPUAddr(0x21)
	ppu.WriteToPPUAddr(0xff)

	ppu.ReadData()
	assert.Equal(t, uint8(0x66), ppu.ReadData())
	assert.Equal(t, uint8(0x77), ppu.ReadData())
}

func TestPPUVramReadsStep32(t *testing.T) {
	ppu := NewPPU(nil, MIRROR_HORIZONTAL)
	ppu.WriteToPPUCTRL(CR_VRAM_ADD_INCREMENT)
	ppu.VRAM[0x01ff] = 0x66
	ppu.VRAM[0x01ff+32] = 0x77
	ppu.VRAM[0x01ff+64] = 0x88

	ppu.WriteToPPUAddr(0x21)
	ppu.WriteToPPUAddr(0xff)

	ppu.ReadData()
	assert.Equal(t, uint8(0x66), ppu.ReadData())
	assert.Equal(t, uint8(0x77), ppu.ReadData())
	assert.Equal(t, uint8(0x88), ppu.ReadData())
}

// Horizontal: https://wiki.nesdev.com/w/index.php/Mirroring
//
//	[0x2000 A ] [0x2400 a ]
//	[0x2800 B ] [0x2C00 b ]
func TestPPUVramHorizontalMirror(t *testing.T) {
	ppu := NewPPU(nil, MIRROR_HORIZONTAL)
	ppu.WriteToPPUAddr(0x24)
	ppu.WriteToPPUAddr(0x05)

	ppu.WriteData(0x66) // write to a

	ppu.WriteToPPUAddr(0x28)
	ppu.WriteToPPUAddr(0x05)

	ppu.WriteData(0x77) // write to B

	ppu.WriteToPPUAddr(0x20)
	ppu.WriteToPPUAddr(0x05)

	ppu.ReadData()                               // load into buffer
	assert.Equal(t, uint8(0x66), ppu.ReadData()) // read from A

	ppu.WriteToPPUAddr(0x2C)
	ppu.WriteToPPUAddr(0x05)

	ppu.ReadData()                               // load into buffer
	assert.Equal(t, uint8(0x77), ppu.ReadData()) // read from b

}

// Vertical: https://wiki.nesdev.com/w/index.php/Mirroring
//
//	[0x2000 A ] [0x2400 B ]
//	[0x2800 a ] [0x2C00 b ]
func TestPPUVramVerticalMirror(t *testing.T) {
	ppu := NewPPU(nil, MIRROR_VERTICAL)

	ppu.WriteToPPUAddr(0x20)
	ppu.WriteToPPUAddr(0x05)

	ppu.WriteData(0x66) // write to A

	ppu.WriteToPPUAddr(0x2C)
	ppu.WriteToPPUAddr(0x05)

	ppu.WriteData(0x77) // write to b

	ppu.WriteToPPUAddr(0x28)
	ppu.WriteToPPUAddr(0x05)

	ppu.ReadData()
	assert.Equal(t, uint8(0x66), ppu.ReadData()) // read from a

	ppu.WriteToPPUAddr(0x24)
	ppu.WriteToPPUAddr(0x05)

	ppu.ReadData()
	assert.Equal(t, uint8(0x77), ppu.ReadData()) // read from B
}

func TestReadStatusResetsLatch(t *testing.T) {
	ppu := NewPPU(make([]uint8, 2048), MIRROR_HORIZONTAL)
	ppu.VRAM[0x0305] = 0x66

	ppu.WriteToPPUAddr(0x21)
	ppu.WriteToPPUAddr(0x23)
	ppu.WriteToPPUAddr(0x05)

	ppu.ReadData()
	assert.NotEqual(t, uint8(0x66), ppu.ReadData())

	ppu.ReadStatus()

	ppu.WriteToPPUAddr(0x23)
	ppu.WriteToPPUAddr(0x05)

	ppu.ReadData()
	assert.Equal(t, uint8(0x66), ppu.ReadData())
}

func TestPPUVramMirroring(t *testing.T) {
	ppu := NewPPU(nil, MIRROR_HORIZONTAL)
	ppu.WriteToPPUCTRL(0)
	ppu.VRAM[0x0305] = 0x66

	ppu.WriteToPPUAddr(0x63) // 0x6305 -> 0x2305
	ppu.WriteToPPUAddr(0x05)

	ppu.ReadData()
	assert.Equal(t, uint8(0x66), ppu.ReadData())
}

func TestReadStatusResetsVblank(t *testing.T) {
	ppu := NewPPU(nil, MIRROR_HORIZONTAL)
	ppu.Status.SetVblankStatus(true)

	status := ppu.ReadStatus()

	assert.Equal(t, uint8(1), status>>7)
	assert.Equal(t, uint8(0), ppu.Status.Bits>>7)
}
