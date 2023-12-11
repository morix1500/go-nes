package nes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestCartridgeForCPUTest(program []uint8) *Cartridge {
	programRom := make([]uint8, 2*PROGRAM_ROM_PAGE_SIZE)
	copy(programRom, program)
	programRom[0x7ffc] = 0x00
	programRom[0x7ffd] = 0x80

	dummyCharacterRom := createDummyRom(2, 1*CHARACTER_ROM_PAGE_SIZE)
	testRom := createTestCartridge(TestCartridge{
		header: []uint8{
			0x4E, 0x45, 0x53, 0x1A, 0x02, 0x01, 0x31, 00, 00, 00, 00, 00, 00, 00, 00, 00,
		},
		trainer:      nil,
		programRom:   programRom,
		characterRom: dummyCharacterRom,
	})

	cartridge, _ := NewCartridge(testRom)
	return cartridge
}

func TestCPULDA(t *testing.T) {
	cases := []struct {
		name            string
		memory          map[uint16]uint8
		program         []uint8
		expectRegisterA uint8
	}{
		{
			name:            "LDA Immediate",
			program:         []uint8{0xa9, 0x05, 0x00},
			expectRegisterA: uint8(0x05),
		},
		{
			name:            "LDA ZeroPage",
			memory:          map[uint16]uint8{0x10: 0x05},
			program:         []uint8{0xa5, 0x10, 0x00},
			expectRegisterA: uint8(0x05),
		},
		{
			name:            "LDA ZeroPageX",
			memory:          map[uint16]uint8{0x10: 0x04, 0x11: 0x05},
			program:         []uint8{0xa9, 0x01, 0xaa, 0xb5, 0x10, 0x00},
			expectRegisterA: uint8(0x05),
		},
		// LDA Absolute
		{
			name:            "LDA Absolute",
			memory:          map[uint16]uint8{0x0010: 0x05},
			program:         []uint8{0xad, 0x10, 0x00, 0x00},
			expectRegisterA: uint8(0x05),
		},
		// LDA AbsoluteX
		{
			name:            "LDA AbsoluteX",
			memory:          map[uint16]uint8{0x0011: 0x05},
			program:         []uint8{0xa9, 0x01, 0xaa, 0xbd, 0x10, 0x00, 0x00},
			expectRegisterA: uint8(0x05),
		},
		// LDA IndirectX
		{
			name:            "LDA IndirectX",
			memory:          map[uint16]uint8{0x11: 0x05, 0x12: 0x06, 0x0605: 0x07},
			program:         []uint8{0xa9, 0x01, 0xaa, 0xa1, 0x10, 0x00},
			expectRegisterA: uint8(0x07),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectRegisterA, cpu.registerA)
		})
	}
}

func TestCPULDX(t *testing.T) {
	cases := []struct {
		name            string
		program         []uint8
		expectRegisterX uint8
	}{
		{
			name:            "LDX Immediate",
			program:         []uint8{0xa2, 0x05, 0x00},
			expectRegisterX: uint8(0x05),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectRegisterX, cpu.registerX)
		})
	}
}

func TestCPULDY(t *testing.T) {
	cases := []struct {
		name            string
		program         []uint8
		expectRegisterY uint8
	}{
		{
			name:            "LDY Immediate",
			program:         []uint8{0xa0, 0x05, 0x00},
			expectRegisterY: uint8(0x05),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectRegisterY, cpu.registerY)
		})
	}
}

func TestCPUSTA(t *testing.T) {
	cases := []struct {
		name         string
		memory       map[uint16]uint8
		program      []uint8
		expectMemory map[uint16]uint8
	}{
		// STA ZeroPage
		{
			name:         "STA ZeroPage",
			program:      []uint8{0xa9, 0x05, 0x85, 0x10, 0x00},
			expectMemory: map[uint16]uint8{0x10: 0x05},
		},
		// STA ZeroPageX
		{
			name:         "STA ZeroPageX",
			program:      []uint8{0xa9, 0x01, 0xaa, 0xa9, 0x05, 0x95, 0x10, 0x00},
			expectMemory: map[uint16]uint8{0x11: 0x05},
		},
		// STA Absolute
		{
			name:         "STA Absolute",
			program:      []uint8{0xa9, 0x05, 0x8d, 0x10, 0x00, 0x00},
			expectMemory: map[uint16]uint8{0x0010: 0x05},
		},
		// STA AbsoluteX
		{
			name:         "STA AbsoluteX",
			program:      []uint8{0xa9, 0x01, 0xaa, 0xa9, 0x05, 0x9d, 0x10, 0x00, 0x00},
			expectMemory: map[uint16]uint8{0x0011: 0x05},
		},
		// STA IndirectX
		{
			name:         "STA IndirectX",
			memory:       map[uint16]uint8{0x11: 0x05, 0x12: 0x06},
			program:      []uint8{0xa9, 0x01, 0xaa, 0xa9, 0x05, 0x81, 0x10, 0x00},
			expectMemory: map[uint16]uint8{0x0605: 0x05},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			for addr, value := range tt.expectMemory {
				assert.Equal(t, value, cpu.readMemory(addr))
			}
		})
	}
}

func TestCPUSTX(t *testing.T) {
	cases := []struct {
		name         string
		memory       map[uint16]uint8
		program      []uint8
		expectMemory map[uint16]uint8
	}{
		{
			name:         "STX ZeroPage",
			program:      []uint8{0xa2, 0x05, 0x86, 0x10, 0x00},
			expectMemory: map[uint16]uint8{0x10: 0x05},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			for addr, value := range tt.expectMemory {
				assert.Equal(t, value, cpu.readMemory(addr))
			}
		})
	}
}

func TestCPUSTY(t *testing.T) {
	cases := []struct {
		name         string
		memory       map[uint16]uint8
		program      []uint8
		expectMemory map[uint16]uint8
	}{
		{
			name:         "STY ZeroPage",
			program:      []uint8{0xa0, 0x05, 0x84, 0x10, 0x00},
			expectMemory: map[uint16]uint8{0x10: 0x05},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			for addr, value := range tt.expectMemory {
				assert.Equal(t, value, cpu.readMemory(addr))
			}
		})
	}
}

func TestCPUInterpretLDAImmediateLoad(t *testing.T) {
	bus := NewBus(createTestCartridgeForCPUTest([]uint8{0xa9, 0x05, 0x00}))
	cpu := NewCPU(bus)
	cpu.Reset()
	cpu.Run()
	assert.Equal(t, uint8(0x05), cpu.registerA)
	assert.True(t, cpu.status&0b0000_0010 == 0b00)
	assert.True(t, cpu.status&0b1000_0000 == 0)
}

func TestCPUInterpretLDAZeroFlag(t *testing.T) {
	bus := NewBus(createTestCartridgeForCPUTest([]uint8{0xa9, 0x00, 0x00}))
	cpu := NewCPU(bus)
	cpu.Reset()
	cpu.Run()
	assert.True(t, cpu.status&0b0000_0010 == 0b10)
}

func TestCPUInterpretaTaxMoveAToX(t *testing.T) {
	bus := NewBus(createTestCartridgeForCPUTest([]uint8{0xa9, 0x0a, 0xaa, 0x00}))
	cpu := NewCPU(bus)
	cpu.Reset()
	cpu.Run()
	assert.Equal(t, uint8(10), cpu.registerX)
}

func TestCPUInterpretaTaxMoveAToY(t *testing.T) {
	bus := NewBus(createTestCartridgeForCPUTest([]uint8{0xa9, 0x0a, 0xa8, 0x00}))
	cpu := NewCPU(bus)
	cpu.Reset()
	cpu.Run()
	assert.Equal(t, uint8(10), cpu.registerY)
}

func TestCPUINC(t *testing.T) {
	bus := NewBus(createTestCartridgeForCPUTest([]uint8{0xe6, 0x10, 0x00}))
	cpu := NewCPU(bus)
	cpu.writeMemory(0x10, 0x05)
	cpu.Reset()
	cpu.Run()
	assert.Equal(t, uint8(0x06), cpu.readMemory(0x10))
}

func TestCPUInterpretInx(t *testing.T) {
	bus := NewBus(createTestCartridgeForCPUTest([]uint8{0xa9, 0x02, 0xaa, 0xe8, 0x00}))
	cpu := NewCPU(bus)
	cpu.Reset()
	cpu.Run()
	assert.Equal(t, uint8(3), cpu.registerX)
}

func TestCPUInterpretInxOverflow(t *testing.T) {
	bus := NewBus(createTestCartridgeForCPUTest([]uint8{0xa9, 0xff, 0xaa, 0xe8, 0xe8, 0x00}))
	cpu := NewCPU(bus)
	cpu.Reset()
	cpu.Run()
	assert.Equal(t, uint8(1), cpu.registerX)
}

func TestCPUInterpretIny(t *testing.T) {
	bus := NewBus(createTestCartridgeForCPUTest([]uint8{0xa9, 0x02, 0xa8, 0xc8, 0x00}))
	cpu := NewCPU(bus)
	cpu.Reset()
	cpu.Run()
	assert.Equal(t, uint8(3), cpu.registerY)
}

func TestCPUInterpretInyOverflow(t *testing.T) {
	bus := NewBus(createTestCartridgeForCPUTest([]uint8{0xa9, 0xff, 0xa8, 0xc8, 0xc8, 0x00}))
	cpu := NewCPU(bus)
	cpu.Reset()
	cpu.Run()
	assert.Equal(t, uint8(1), cpu.registerY)
}

func TestCPUInterpret5OpsWorkingTogether(t *testing.T) {
	bus := NewBus(createTestCartridgeForCPUTest([]uint8{0xa9, 0xc0, 0xaa, 0xe8, 0x00}))
	cpu := NewCPU(bus)
	cpu.Reset()
	cpu.Run()
	assert.Equal(t, uint8(0xc1), cpu.registerX)
}

func TestCPUInterpretLDAFromMemory(t *testing.T) {
	bus := NewBus(createTestCartridgeForCPUTest([]uint8{0xa5, 0x10, 0x00}))
	cpu := NewCPU(bus)
	cpu.writeMemory(0x10, 0x55)
	cpu.Reset()
	cpu.Run()
	assert.Equal(t, uint8(0x55), cpu.registerA)
}

func TestCPUADC(t *testing.T) {
	cases := []struct {
		name            string
		program         []uint8
		expectRegisterA uint8
		expectStatus    uint8
	}{
		{
			name:            "ADC Immediate",
			program:         []uint8{0xa9, 0x05, 0x69, 0x05, 0x00},
			expectRegisterA: uint8(0x0a),
			expectStatus:    0b0010_0100,
		},
		{
			name:            "ADC Immediate with carry",
			program:         []uint8{0xa9, 0xff, 0x69, 0x02, 0x00},
			expectRegisterA: uint8(0x01),
			expectStatus:    0b0010_0101,
		},
		{
			name:            "ADC Immediate with overflow",
			program:         []uint8{0xa9, 0x7f, 0x69, 0x01, 0x00},
			expectRegisterA: uint8(0x80),
			expectStatus:    0b1110_0100,
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectRegisterA, cpu.registerA)
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

// ANDのテストコードを生成する
func TestCPUAND(t *testing.T) {
	cases := []struct {
		name            string
		program         []uint8
		expectRegisterA uint8
		expectStatus    uint8
	}{
		{
			name:            "AND Immediate",
			program:         []uint8{0xa9, 0x05, 0x29, 0x04, 0x00},
			expectRegisterA: uint8(0x04),
			expectStatus:    0b0010_0100,
		},
		{
			name:            "AND Immediate with negative",
			program:         []uint8{0xa9, 0xff, 0x29, 0x80, 0x00},
			expectRegisterA: uint8(0x80),
			expectStatus:    0b1010_0100,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectRegisterA, cpu.registerA)
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUASLAccumulator(t *testing.T) {
	cases := []struct {
		name            string
		program         []uint8
		expectRegisterA uint8
		expectStatus    uint8
	}{
		{
			name:            "ASL Accumulator",
			program:         []uint8{0xa9, 0x05, 0x0a, 0x00},
			expectRegisterA: uint8(0x0a),
			expectStatus:    0b0010_0100,
		},
		{
			name:            "ASL Accumulator with carry",
			program:         []uint8{0xa9, 0xff, 0x0a, 0x00},
			expectRegisterA: uint8(0xfe),
			expectStatus:    0b1010_0101,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectRegisterA, cpu.registerA)
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUASL(t *testing.T) {
	cases := []struct {
		name         string
		memory       map[uint16]uint8
		program      []uint8
		expectMemory map[uint16]uint8
		expectStatus uint8
	}{
		{
			name:         "ASL ZeroPage",
			memory:       map[uint16]uint8{0x10: 0x05},
			program:      []uint8{0x06, 0x10, 0x00},
			expectMemory: map[uint16]uint8{0x10: 0x0a},
			expectStatus: 0b0010_0100,
		},
		{
			name: "ASL ZeroPage with carry",
			memory: map[uint16]uint8{
				0x10: 0xff,
			},
			program:      []uint8{0x06, 0x10, 0x00},
			expectMemory: map[uint16]uint8{0x10: 0xfe},
			expectStatus: 0b1010_0101,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			for addr, value := range tt.expectMemory {
				assert.Equal(t, value, cpu.readMemory(addr))
			}
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUBCC(t *testing.T) {
	cases := []struct {
		name     string
		program  []uint8
		expectPC uint16
	}{
		{
			name:     "BCC Branch",
			program:  []uint8{0x90, 0x02, 0x00},
			expectPC: uint16(0x8005),
		},
		{
			name:     "BCC Branch with negative",
			program:  []uint8{0x90, 0x80, 0x00},
			expectPC: uint16(0x7f83),
		},
		{
			name:     "BCC No Branch",
			program:  []uint8{0xa9, 0xff, 0x0a, 0x90, 0x02, 0x00},
			expectPC: uint16(0x8006),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectPC, cpu.programCounter)
		})
	}
}

func TestCPUBCS(t *testing.T) {
	cases := []struct {
		name     string
		program  []uint8
		expectPC uint16
	}{
		{
			name:     "BCS Branch",
			program:  []uint8{0xa9, 0xff, 0x0a, 0xb0, 0x02, 0x00},
			expectPC: uint16(0x8008),
		},
		{
			name:     "BCS Branch with negative",
			program:  []uint8{0xa9, 0xff, 0x0a, 0xb0, 0x80, 0x00},
			expectPC: uint16(0x7f86),
		},
		{
			name:     "BCS No Branch",
			program:  []uint8{0xb0, 0x02, 0x00},
			expectPC: uint16(0x8003),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectPC, cpu.programCounter)
		})
	}
}

func TestCPUBEQ(t *testing.T) {
	cases := []struct {
		name     string
		program  []uint8
		expectPC uint16
	}{
		{
			name:     "BEQ Branch",
			program:  []uint8{0xa9, 0x00, 0xf0, 0x02, 0x00},
			expectPC: uint16(0x8007),
		},
		{
			name:     "BEQ No Branch",
			program:  []uint8{0xa9, 0x01, 0xf0, 0x02, 0x00},
			expectPC: uint16(0x8005),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectPC, cpu.programCounter)
		})
	}
}

func TestCPUBIT(t *testing.T) {
	cases := []struct {
		name         string
		memory       map[uint16]uint8
		program      []uint8
		expectStatus uint8
	}{
		{
			name:         "BIT ZeroPage",
			memory:       map[uint16]uint8{0x10: 0x05},
			program:      []uint8{0xa9, 0x05, 0x24, 0x10, 0x00},
			expectStatus: 0b0010_0100,
		},
		{
			name:         "BIT ZeroPage with zero",
			memory:       map[uint16]uint8{0x10: 0x05},
			program:      []uint8{0xa9, 0x0a, 0x24, 0x10, 0x00},
			expectStatus: 0b0010_0110,
		},
		{
			name:         "BIT ZeroPage with negative and overflow",
			memory:       map[uint16]uint8{0x10: 0xff},
			program:      []uint8{0x24, 0x10, 0x00},
			expectStatus: 0b1110_0110,
		},
		{
			name:         "BIT Absolute",
			memory:       map[uint16]uint8{0x0010: 0xff},
			program:      []uint8{0xa9, 0x05, 0x2c, 0x10, 0x00, 0x00},
			expectStatus: 0b1110_0100,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUBMI(t *testing.T) {
	cases := []struct {
		name     string
		program  []uint8
		expectPC uint16
	}{
		{
			name:     "BMI Branch",
			program:  []uint8{0xa9, 0xf1, 0x30, 0x02, 0x00},
			expectPC: uint16(0x8007),
		},
		{
			name:     "BMI No Branch",
			program:  []uint8{0x30, 0x02, 0x00},
			expectPC: uint16(0x8003),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectPC, cpu.programCounter)
		})
	}
}

func TestCPUBNE(t *testing.T) {
	cases := []struct {
		name     string
		program  []uint8
		expectPC uint16
	}{
		{
			name:     "BNE Branch",
			program:  []uint8{0xa9, 0x01, 0xd0, 0x02, 0x00},
			expectPC: uint16(0x8007),
		},
		{
			name:     "BNE No Branch",
			program:  []uint8{0xa9, 0x00, 0xd0, 0x02, 0x00},
			expectPC: uint16(0x8005),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectPC, cpu.programCounter)
		})
	}
}

func TestCPUBPL(t *testing.T) {
	cases := []struct {
		name     string
		program  []uint8
		expectPC uint16
	}{
		{
			name:     "BPL Branch",
			program:  []uint8{0x10, 0x02, 0x00},
			expectPC: uint16(0x8005),
		},
		{
			name:     "BPL No Branch",
			program:  []uint8{0x10, 0x00, 0x00},
			expectPC: uint16(0x8003),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectPC, cpu.programCounter)
		})
	}
}

func TestCPUBVC(t *testing.T) {
	cases := []struct {
		name     string
		program  []uint8
		expectPC uint16
	}{
		{
			name:     "BVC Branch",
			program:  []uint8{0x50, 0x02, 0x00},
			expectPC: uint16(0x8005),
		},
		{
			name:     "BVC No Branch",
			program:  []uint8{0xa9, 0x7f, 0xaa, 0x69, 0x01, 0x50, 0x02, 0x00},
			expectPC: uint16(0x8008),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectPC, cpu.programCounter)
		})
	}
}

func TestCPUBVS(t *testing.T) {
	cases := []struct {
		name     string
		program  []uint8
		expectPC uint16
	}{
		{
			name:     "BVS Branch",
			program:  []uint8{0xa9, 0x7f, 0xaa, 0x69, 0x01, 0x70, 0x02, 0x00},
			expectPC: uint16(0x800a),
		},
		{
			name:     "BVS No Branch",
			program:  []uint8{0x70, 0x02, 0x00},
			expectPC: uint16(0x8003),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectPC, cpu.programCounter)
		})
	}
}

func TestCPUCLC(t *testing.T) {
	cases := []struct {
		name         string
		program      []uint8
		expectStatus uint8
	}{
		{
			name:         "CLC",
			program:      []uint8{0xa9, 0xff, 0x69, 0x02, 0x18, 0x00},
			expectStatus: 0b0010_0100,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUCLD(t *testing.T) {
	cases := []struct {
		name         string
		program      []uint8
		expectStatus uint8
	}{
		{
			name:         "CLD",
			program:      []uint8{0xf8, 0xd8, 0x00},
			expectStatus: 0b0010_0100,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUSEC(t *testing.T) {
	cases := []struct {
		name         string
		program      []uint8
		expectStatus uint8
	}{
		{
			name:         "SEC",
			program:      []uint8{0x38, 0x00},
			expectStatus: 0b0010_0101,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUSED(t *testing.T) {
	cases := []struct {
		name         string
		program      []uint8
		expectStatus uint8
	}{
		{
			name:         "SED",
			program:      []uint8{0xf8, 0x00},
			expectStatus: 0b0010_1100,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUCLI(t *testing.T) {
	cases := []struct {
		name         string
		program      []uint8
		expectStatus uint8
	}{
		{
			name:         "CLI",
			program:      []uint8{0x58, 0x00},
			expectStatus: 0b0010_0000,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUSEI(t *testing.T) {
	cases := []struct {
		name         string
		program      []uint8
		expectStatus uint8
	}{
		{
			name:         "SEI",
			program:      []uint8{0x58, 0x78, 0x00},
			expectStatus: 0b0010_0100,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUCLV(t *testing.T) {
	cases := []struct {
		name         string
		program      []uint8
		expectStatus uint8
	}{
		{
			name:         "CLV",
			program:      []uint8{0xa9, 0x7f, 0x69, 0x01, 0xb8, 0x00},
			expectStatus: 0b1010_0100,
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUCMP(t *testing.T) {
	cases := []struct {
		name         string
		memory       map[uint16]uint8
		program      []uint8
		expectStatus uint8
	}{
		{
			name:         "CMP Immediate with zero",
			memory:       map[uint16]uint8{0x10: 0x05},
			program:      []uint8{0xa9, 0x05, 0xc9, 0x05, 0x00},
			expectStatus: 0b0010_0111,
		},
		{
			name:         "CMP Immediate with carry",
			memory:       map[uint16]uint8{0x10: 0x05},
			program:      []uint8{0xa9, 0x05, 0xc9, 0x04, 0x00},
			expectStatus: 0b0010_0101,
		},
		{
			name:         "CMP Immediate with negative",
			memory:       map[uint16]uint8{0x10: 0x05},
			program:      []uint8{0xa9, 0x05, 0xc9, 0x06, 0x00},
			expectStatus: 0b1010_0100,
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUCPX(t *testing.T) {
	cases := []struct {
		name         string
		memory       map[uint16]uint8
		program      []uint8
		expectStatus uint8
	}{
		{
			name:         "CPX Immediate with zero",
			memory:       map[uint16]uint8{0x10: 0x05},
			program:      []uint8{0xa9, 0x05, 0xaa, 0xc9, 0x05, 0x00},
			expectStatus: 0b0010_0111,
		},
		{
			name:         "CPX Immediate with carry",
			memory:       map[uint16]uint8{0x10: 0x05},
			program:      []uint8{0xa9, 0x05, 0xaa, 0xc9, 0x04, 0x00},
			expectStatus: 0b0010_0101,
		},
		{
			name:         "CPX Immediate with negative",
			memory:       map[uint16]uint8{0x10: 0x05},
			program:      []uint8{0xa9, 0x05, 0xaa, 0xc9, 0x06, 0x00},
			expectStatus: 0b1010_0100,
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUCPY(t *testing.T) {
	cases := []struct {
		name         string
		memory       map[uint16]uint8
		program      []uint8
		expectStatus uint8
	}{
		{
			name:         "CPY Immediate with zero",
			memory:       map[uint16]uint8{0x10: 0x05},
			program:      []uint8{0xa9, 0x05, 0xa8, 0xc0, 0x05, 0x00},
			expectStatus: 0b0010_0111,
		},
		{
			name:         "CPY Immediate with carry",
			memory:       map[uint16]uint8{0x10: 0x05},
			program:      []uint8{0xa9, 0x05, 0xa8, 0xc0, 0x04, 0x00},
			expectStatus: 0b0010_0101,
		},
		{
			name:         "CPY Immediate with negative",
			memory:       map[uint16]uint8{0x10: 0x05},
			program:      []uint8{0xa9, 0x05, 0xa8, 0xc0, 0x06, 0x00},
			expectStatus: 0b1010_0100,
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUDEC(t *testing.T) {
	bus := NewBus(createTestCartridgeForCPUTest([]uint8{0xc6, 0x10, 0x00}))
	cpu := NewCPU(bus)
	cpu.writeMemory(0x10, 0x05)
	cpu.Reset()
	cpu.Run()
	assert.Equal(t, uint8(0x04), cpu.readMemory(0x10))
}

func TestCPUDEX(t *testing.T) {
	bus := NewBus(createTestCartridgeForCPUTest([]uint8{0xa2, 0x10, 0xca, 0x00}))
	cpu := NewCPU(bus)
	cpu.Reset()
	cpu.Run()
	assert.Equal(t, uint8(0x0f), cpu.registerX)
}

func TestCPUDEY(t *testing.T) {
	bus := NewBus(createTestCartridgeForCPUTest([]uint8{0xa0, 0x10, 0x88, 0x00}))
	cpu := NewCPU(bus)
	cpu.Reset()
	cpu.Run()
	assert.Equal(t, uint8(0x0f), cpu.registerY)
}

func TestCPUEOR(t *testing.T) {
	cases := []struct {
		name            string
		program         []uint8
		expectRegisterA uint8
		expectStatus    uint8
	}{
		{
			name:            "EOR Immediate",
			program:         []uint8{0xa9, 0x05, 0x49, 0x04, 0x00},
			expectRegisterA: uint8(0x01),
			expectStatus:    0b0010_0100,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectRegisterA, cpu.registerA)
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUORA(t *testing.T) {
	cases := []struct {
		name            string
		program         []uint8
		expectRegisterA uint8
		expectStatus    uint8
	}{
		{
			name:            "ORA Immediate",
			program:         []uint8{0xa9, 0x05, 0x09, 0x0a, 0x00},
			expectRegisterA: uint8(0x0f),
			expectStatus:    0b0010_0100,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectRegisterA, cpu.registerA)
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUJMP(t *testing.T) {
	cases := []struct {
		name     string
		memory   map[uint16]uint8
		program  []uint8
		expectPC uint16
	}{
		{
			name:     "JMP Absolute",
			program:  []uint8{0x4c, 0x05, 0x10, 0x00},
			expectPC: uint16(0x1006),
		},
		{
			name:     "JMP Indirect",
			program:  []uint8{0x6c, 0x05, 0x10, 0x00},
			expectPC: uint16(0x1006),
		},
		{
			name:     "JMP Indirect with page boundary",
			memory:   map[uint16]uint8{0x1100: 0x10, 0x11ff: 0x70},
			program:  []uint8{0x6c, 0xff, 0x11, 0x00},
			expectPC: uint16(0x1071),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectPC, cpu.programCounter)
		})
	}
}

func TestCPUJSR(t *testing.T) {
	cases := []struct {
		name        string
		memory      map[uint16]uint8
		program     []uint8
		expectPC    uint16
		expectStack map[uint16]uint8
	}{
		{
			name:        "JSR",
			program:     []uint8{0x20, 0x05, 0x80, 0x00},
			expectPC:    uint16(0x8006),
			expectStack: map[uint16]uint8{0x01fd: 0x80, 0x01fc: 0x02},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectPC, cpu.programCounter)
			for addr, value := range tt.expectStack {
				assert.Equal(t, value, cpu.readMemory(addr))
			}
		})
	}
}

func TestCPURTS(t *testing.T) {
	cases := []struct {
		name     string
		memory   map[uint16]uint8
		program  []uint8
		expectPC uint16
	}{
		{
			name:     "RTS",
			memory:   map[uint16]uint8{0x1005: 0x60},
			program:  []uint8{0x20, 0x05, 0x10, 0x00},
			expectPC: uint16(0x8004),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectPC, cpu.programCounter)
		})
	}
}

func TestCPULSRAccumulator(t *testing.T) {
	cases := []struct {
		name            string
		memory          map[uint16]uint8
		program         []uint8
		expectRegisterA uint8
		expectStatus    uint8
	}{
		{
			name:            "LSR Accumulator",
			program:         []uint8{0xa9, 0x04, 0x4a, 0x00},
			expectRegisterA: uint8(0x02),
			expectStatus:    0b0010_0100,
		},
		{
			name:            "LSR Accumulator with carry",
			program:         []uint8{0xa9, 0x05, 0x4a, 0x00},
			expectRegisterA: uint8(0x02),
			expectStatus:    0b0010_0101,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPULSR(t *testing.T) {
	cases := []struct {
		name         string
		memory       map[uint16]uint8
		program      []uint8
		expectMemory map[uint16]uint8
		expectStatus uint8
	}{
		{
			name:         "LSR ZeroPage",
			memory:       map[uint16]uint8{},
			program:      []uint8{0xa9, 0x04, 0x85, 0x0a, 0x46, 0x0a, 0x00},
			expectMemory: map[uint16]uint8{0x0a: 0x02},
			expectStatus: 0b0010_0100,
		},
		{
			name:         "LSR ZeroPage with carry",
			memory:       map[uint16]uint8{0x10: 0xff},
			program:      []uint8{0xa9, 0x05, 0x85, 0x0a, 0x46, 0x0a, 0x00},
			expectMemory: map[uint16]uint8{0x0a: 0x02},
			expectStatus: 0b0010_0101,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			for addr, value := range tt.expectMemory {
				assert.Equal(t, value, cpu.readMemory(addr))
			}
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUPHA(t *testing.T) {
	cases := []struct {
		name        string
		memory      map[uint16]uint8
		program     []uint8
		expectStack map[uint16]uint8
	}{
		{
			name:        "PHA",
			memory:      map[uint16]uint8{0x10: 0x05},
			program:     []uint8{0xa9, 0x05, 0x48, 0x00},
			expectStack: map[uint16]uint8{0x01fd: 0x05},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)

			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()

			for addr, value := range tt.expectStack {
				assert.Equal(t, value, cpu.readMemory(addr))
			}
		})
	}
}

func TestCPUPLA(t *testing.T) {
	cases := []struct {
		name            string
		memory          map[uint16]uint8
		program         []uint8
		expectRegisterA uint8
	}{
		{
			name:            "PHL",
			memory:          map[uint16]uint8{0x10: 0x05},
			program:         []uint8{0xa9, 0x05, 0x48, 0xa9, 0x1a, 0x68, 0x00},
			expectRegisterA: 0x05,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)

			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectRegisterA, cpu.registerA)

		})
	}
}

func TestCPUPHP(t *testing.T) {
	cases := []struct {
		name        string
		program     []uint8
		expectStack map[uint16]uint8
	}{
		{
			name:        "PHP",
			program:     []uint8{0xa9, 0x00, 0x08, 0x00},
			expectStack: map[uint16]uint8{0x01fd: 0b0011_0110},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			for addr, value := range tt.expectStack {
				assert.Equal(t, value, cpu.readMemory(addr))
			}
		})
	}
}

func TestCPUPLP(t *testing.T) {
	cases := []struct {
		name         string
		program      []uint8
		expectStatus uint8
	}{
		{
			name:         "PLP",
			program:      []uint8{0xa9, 0x00, 0x08, 0xa9, 0x01, 0x28, 0x00},
			expectStatus: uint8(0b0010_0110),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUROLAccumulator(t *testing.T) {
	cases := []struct {
		name            string
		program         []uint8
		expectRegisterA uint8
		expectStatus    uint8
	}{
		{
			name:            "ROL Accumulator",
			program:         []uint8{0xa9, 0x04, 0x2a, 0x00},
			expectRegisterA: uint8(0x08),
			expectStatus:    0b0010_0100,
		},
		{
			name:            "ROL Accumulator with carry",
			program:         []uint8{0xa9, 0x85, 0x38, 0x2a, 0x00},
			expectRegisterA: uint8(0b0000_1011),
			expectStatus:    0b0010_0101,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectRegisterA, cpu.registerA)
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUROL(t *testing.T) {
	cases := []struct {
		name         string
		memory       map[uint16]uint8
		program      []uint8
		expectMemory map[uint16]uint8
		expectStatus uint8
	}{
		{
			name:         "ROL ZeroPage",
			program:      []uint8{0xa9, 0x04, 0x85, 0xb1, 0x26, 0xb1, 0x00},
			expectMemory: map[uint16]uint8{0xb1: 0x08},
			expectStatus: 0b0010_0100,
		},
		{
			name:         "ROL ZeroPage with carry",
			program:      []uint8{0xa9, 0x85, 0x85, 0xb1, 0x38, 0x26, 0xb1, 0x00},
			expectMemory: map[uint16]uint8{0xb1: 0b0000_1011},
			expectStatus: 0b0010_0101,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}

			cpu.Reset()
			cpu.Run()
			for addr, value := range tt.expectMemory {
				assert.Equal(t, value, cpu.readMemory(addr))
			}
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPURORAccumulator(t *testing.T) {
	cases := []struct {
		name            string
		program         []uint8
		expectRegisterA uint8
		expectStatus    uint8
	}{
		{
			name:            "ROR Accumulator",
			program:         []uint8{0xa9, 0x04, 0x6a, 0x00},
			expectRegisterA: uint8(0x02),
			expectStatus:    0b0010_0100,
		},
		{
			name:            "ROR Accumulator with carry",
			program:         []uint8{0xa9, 0x85, 0x38, 0x6a, 0x00},
			expectRegisterA: uint8(0b1100_0010),
			expectStatus:    0b1010_0101,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectRegisterA, cpu.registerA)
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUROR(t *testing.T) {
	cases := []struct {
		name         string
		memory       map[uint16]uint8
		program      []uint8
		expectMemory map[uint16]uint8
		expectStatus uint8
	}{
		{
			name:         "ROR ZeroPage",
			program:      []uint8{0xa9, 0x04, 0x85, 0xb1, 0x66, 0xb1, 0x00},
			expectMemory: map[uint16]uint8{0xb1: 0x02},
			expectStatus: 0b0010_0100,
		},
		{
			name:         "ROR ZeroPage with carry",
			program:      []uint8{0xa9, 0x85, 0x85, 0xb1, 0x38, 0x66, 0xb1, 0x00},
			expectMemory: map[uint16]uint8{0xb1: 0b1100_0010},
			expectStatus: 0b1010_0101,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			for addr, value := range tt.expectMemory {
				assert.Equal(t, value, cpu.readMemory(addr))
			}
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPURTI(t *testing.T) {
	program := []uint8{0x40, 0x00}
	bus := NewBus(createTestCartridgeForCPUTest(program))
	cpu := NewCPU(bus)
	cpu.writeMemory(0x01fd, 0x80)
	cpu.writeMemory(0x01fc, 0x11)
	cpu.writeMemory(0x01fb, 0b1000_0011)
	cpu.programCounter = cpu.readMemory16(0xFFFC)
	cpu.stackPointer = 0xfa
	cpu.Run()

	assert.Equal(t, uint8(0b1010_0011), cpu.status)
	assert.Equal(t, uint16(0x8012), cpu.programCounter)
}

func TestCPUSBC(t *testing.T) {
	cases := []struct {
		name            string
		memory          map[uint16]uint8
		program         []uint8
		expectRegisterA uint8
		expectStatus    uint8
	}{
		{
			name:            "SBC Immediate",
			program:         []uint8{0xa9, 0x05, 0xe9, 0x03, 0x00},
			expectRegisterA: uint8(0x01),
			expectStatus:    0b0010_0101,
		},
		{
			name:            "SBC Immediate with carry",
			program:         []uint8{0xa9, 0x05, 0x38, 0xe9, 0x05, 0x00},
			expectRegisterA: uint8(0x00),
			expectStatus:    0b0010_0111,
		},
		{
			name:            "SBC Immediate with minux",
			program:         []uint8{0xa9, 0x05, 0xe9, 0x05, 0x00},
			expectRegisterA: uint8(0xff),
			expectStatus:    0b1010_0100,
		},

		{
			name:            "SBC Immediate with overflow",
			program:         []uint8{0xa9, 0xb0, 0x38, 0xe9, 0x7f, 0x00},
			expectRegisterA: uint8(0x31),
			expectStatus:    0b0110_0101,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bus := NewBus(createTestCartridgeForCPUTest(tt.program))
			cpu := NewCPU(bus)
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.Reset()
			cpu.Run()
			assert.Equal(t, tt.expectRegisterA, cpu.registerA)
			assert.Equal(t, tt.expectStatus, cpu.status)
		})
	}
}

func TestCPUTSX(t *testing.T) {
	program := []uint8{0xba, 0x00}
	bus := NewBus(createTestCartridgeForCPUTest(program))
	cpu := NewCPU(bus)
	cpu.stackPointer = 0x05
	cpu.programCounter = 0x8000
	cpu.Run()
	assert.Equal(t, uint8(0x05), cpu.registerX)
}

func TestCPUTXA(t *testing.T) {
	program := []uint8{0x8a, 0x00}
	bus := NewBus(createTestCartridgeForCPUTest(program))
	cpu := NewCPU(bus)
	cpu.registerX = 0x05
	cpu.programCounter = 0x8000
	cpu.Run()
	assert.Equal(t, uint8(0x05), cpu.registerA)
}

func TestCPUTXS(t *testing.T) {
	program := []uint8{0x9a, 0x00}
	bus := NewBus(createTestCartridgeForCPUTest(program))
	cpu := NewCPU(bus)
	cpu.registerX = 0x05
	cpu.programCounter = 0x8000
	cpu.Run()
	assert.Equal(t, uint8(0x05), cpu.stackPointer)
}

func TestCPUTYA(t *testing.T) {
	program := []uint8{0x98, 0x00}
	bus := NewBus(createTestCartridgeForCPUTest(program))
	cpu := NewCPU(bus)
	cpu.registerY = 0x05
	cpu.programCounter = 0x8000
	cpu.Run()
	assert.Equal(t, uint8(0x05), cpu.registerA)
}
