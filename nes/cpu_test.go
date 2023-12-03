package nes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			memory:          map[uint16]uint8{0x2010: 0x05},
			program:         []uint8{0xad, 0x10, 0x20, 0x00},
			expectRegisterA: uint8(0x05),
		},
		// LDA AbsoluteX
		{
			name:            "LDA AbsoluteX",
			memory:          map[uint16]uint8{0x2011: 0x05},
			program:         []uint8{0xa9, 0x01, 0xaa, 0xbd, 0x10, 0x20, 0x00},
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
			cpu := NewCPU()
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.LoadAndRun(tt.program)
			assert.Equal(t, tt.expectRegisterA, cpu.registerA)
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
			program:      []uint8{0xa9, 0x05, 0x8d, 0x10, 0x20, 0x00},
			expectMemory: map[uint16]uint8{0x2010: 0x05},
		},
		// STA AbsoluteX
		{
			name:         "STA AbsoluteX",
			program:      []uint8{0xa9, 0x01, 0xaa, 0xa9, 0x05, 0x9d, 0x10, 0x20, 0x00},
			expectMemory: map[uint16]uint8{0x2011: 0x05},
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
			cpu := NewCPU()
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.LoadAndRun(tt.program)
			for addr, value := range tt.expectMemory {
				assert.Equal(t, value, cpu.readMemory(addr))
			}
		})
	}
}

func TestCPUInterpretLDAImmediateLoad(t *testing.T) {
	cpu := NewCPU()
	cpu.LoadAndRun([]uint8{0xa9, 0x05, 0x00})
	assert.Equal(t, uint8(0x05), cpu.registerA)
	assert.True(t, cpu.status&0b0000_0010 == 0b00)
	assert.True(t, cpu.status&0b1000_0000 == 0)
}

func TestCPUInterpretLDAZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.LoadAndRun([]uint8{0xa9, 0x00, 0x00})
	assert.True(t, cpu.status&0b0000_0010 == 0b10)
}

func TestCPUInterpretaTaxMoveAToX(t *testing.T) {
	cpu := NewCPU()
	// LDA $0a
	// TAX
	cpu.LoadAndRun([]uint8{0xa9, 0x0a, 0xaa, 0x00})
	assert.Equal(t, uint8(10), cpu.registerX)
}

func TestCPUInterpretInx(t *testing.T) {
	cpu := NewCPU()
	// LDA $02
	// TAX
	// INX
	cpu.LoadAndRun([]uint8{0xa9, 0x02, 0xaa, 0xe8, 0x00})
	assert.Equal(t, uint8(3), cpu.registerX)
}

func TestCPUInterpretInxOverflow(t *testing.T) {
	cpu := NewCPU()
	// LDA $ff
	// TAX
	// INX
	// INX
	cpu.LoadAndRun([]uint8{0xa9, 0xff, 0xaa, 0xe8, 0xe8, 0x00})
	assert.Equal(t, uint8(1), cpu.registerX)
}

func TestCPUInterpret5OpsWorkingTogether(t *testing.T) {
	cpu := NewCPU()
	// LDA $c0
	// TAX
	// INX
	cpu.LoadAndRun([]uint8{0xa9, 0xc0, 0xaa, 0xe8, 0x00})
	assert.Equal(t, uint8(0xc1), cpu.registerX)
}

func TestCPUInterpretLDAFromMemory(t *testing.T) {
	cpu := NewCPU()
	cpu.writeMemory(0x10, 0x55)
	cpu.LoadAndRun([]uint8{0xa5, 0x10, 0x00})
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
			expectStatus:    0b0000_0000,
		},
		{
			name:            "ADC Immediate with carry",
			program:         []uint8{0xa9, 0xff, 0x69, 0x02, 0x00},
			expectRegisterA: uint8(0x01),
			expectStatus:    0b0000_0001,
		},
		{
			name:            "ADC Immediate with overflow",
			program:         []uint8{0xa9, 0x7f, 0x69, 0x01, 0x00},
			expectRegisterA: uint8(0x80),
			expectStatus:    0b1100_0000,
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cpu := NewCPU()
			cpu.LoadAndRun(tt.program)
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
			expectStatus:    0b0000_0000,
		},
		{
			name:            "AND Immediate with negative",
			program:         []uint8{0xa9, 0xff, 0x29, 0x80, 0x00},
			expectRegisterA: uint8(0x80),
			expectStatus:    0b1000_0000,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cpu := NewCPU()
			cpu.LoadAndRun(tt.program)
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
			expectStatus:    0b0000_0000,
		},
		{
			name:            "ASL Accumulator with carry",
			program:         []uint8{0xa9, 0xff, 0x0a, 0x00},
			expectRegisterA: uint8(0xfe),
			expectStatus:    0b1000_0001,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cpu := NewCPU()
			cpu.LoadAndRun(tt.program)
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
			expectStatus: 0b0000_0000,
		},
		{
			name: "ASL ZeroPage with carry",
			memory: map[uint16]uint8{
				0x10: 0xff,
			},
			program:      []uint8{0x06, 0x10, 0x00},
			expectMemory: map[uint16]uint8{0x10: 0xfe},
			expectStatus: 0b1000_0001,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cpu := NewCPU()
			for addr, value := range tt.memory {
				cpu.writeMemory(addr, value)
			}
			cpu.LoadAndRun(tt.program)
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
			cpu := NewCPU()
			cpu.LoadAndRun(tt.program)
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
			cpu := NewCPU()
			cpu.LoadAndRun(tt.program)
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
			cpu := NewCPU()
			cpu.LoadAndRun(tt.program)
			assert.Equal(t, tt.expectPC, cpu.programCounter)
		})
	}
}
