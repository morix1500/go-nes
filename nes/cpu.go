package nes

import "fmt"

type AddressingMode uint8

const (
	Immediate AddressingMode = iota
	ZeroPage
	ZeroPageX
	ZeroPageY
	Absolute
	AbsoluteX
	AbsoluteY
	IndirectX
	IndirectY
	NoneAddressing
)

type CPU struct {
	registerA      uint8
	registerX      uint8
	registerY      uint8
	status         uint8
	programCounter uint16
	memory         [0xFFFF]uint8
}

func NewCPU() *CPU {
	return &CPU{
		registerA:      0,
		registerX:      0,
		registerY:      0,
		status:         0,
		programCounter: 0,
	}
}

func (c *CPU) lda(mode AddressingMode) error {
	addr, err := c.getOperandAddress(mode)
	if err != nil {
		return err
	}
	value := c.readMemory(addr)
	c.registerA = value
	c.updateZeroAndNegativeFlags(c.registerA)

	return nil
}

func (c *CPU) sta(mode AddressingMode) error {
	addr, err := c.getOperandAddress(mode)
	if err != nil {
		return err
	}
	c.writeMemory(addr, c.registerA)

	return nil
}

func (c *CPU) tax() {
	c.registerX = c.registerA
	c.updateZeroAndNegativeFlags(c.registerX)
}

func (c *CPU) inx() {
	c.registerX++
	c.updateZeroAndNegativeFlags(c.registerX)
}

func (c *CPU) updateZeroAndNegativeFlags(result uint8) {
	if result == 0 {
		c.status |= 0b0000_0010
	} else {
		c.status &= 0b1111_1101
	}

	if result&0b1000_0000 != 0 {
		c.status |= 0b1000_0000
	} else {
		c.status &= 0b0111_1111
	}
}

func (c *CPU) readMemory(address uint16) uint8 {
	return c.memory[address]
}

func (c *CPU) readMemory16(address uint16) uint16 {
	lo := uint16(c.readMemory(address))
	hi := uint16(c.readMemory(address + 1))
	return (hi << 8) | lo
}

func (c *CPU) writeMemory(address uint16, value uint8) {
	c.memory[address] = value
}

func (c *CPU) writeMemory16(address uint16, value uint16) {
	hi := uint8(value >> 8)
	lo := uint8(value & 0xff)
	c.writeMemory(address, lo)
	c.writeMemory(address+1, hi)
}

func (c *CPU) LoadAndRun(program []uint8) error {
	c.Load(program)
	c.Reset()
	return c.Run()
}

func (c *CPU) Load(program []uint8) {
	copy(c.memory[0x8000:(0x8000+len(program))], program)
	c.writeMemory16(0xFFFC, 0x8000)
}

func (c *CPU) Reset() {
	c.registerA = 0
	c.registerX = 0
	c.status = 0
	c.programCounter = c.readMemory16(0xFFFC)
}

func (c *CPU) Run() error {
	var opsInfo OpeCode
	var ok bool
	for {
		code := c.readMemory(c.programCounter)
		c.programCounter++

		if opsInfo, ok = CPU_OPS_CODES[code]; !ok {
			return fmt.Errorf("unknown code: %d", code)
		}

		switch opsInfo.Mnemonic {
		case "BRK":
			return nil
		case "TAX":
			c.tax()
		case "INX":
			c.inx()
		case "LDA":
			if err := c.lda(opsInfo.Mode); err != nil {
				return err
			}
		case "STA":
			if err := c.sta(opsInfo.Mode); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown code: %d", code)
		}
		c.programCounter += uint16(opsInfo.Length - 1)
	}
}

func (c *CPU) getOperandAddress(mode AddressingMode) (uint16, error) {
	switch mode {
	case Immediate:
		return c.programCounter, nil
	case ZeroPage:
		return uint16(c.readMemory(c.programCounter)), nil
	case ZeroPageX:
		return uint16(c.readMemory(c.programCounter) + c.registerX), nil
	case ZeroPageY:
		return uint16(c.readMemory(c.programCounter) + c.registerY), nil
	case Absolute:
		return c.readMemory16(c.programCounter), nil
	case AbsoluteX:
		return c.readMemory16(c.programCounter) + uint16(c.registerX), nil
	case AbsoluteY:
		return c.readMemory16(c.programCounter) + uint16(c.registerY), nil
	case IndirectX:
		base := c.readMemory16(c.programCounter)
		ptr := uint8(base + uint16(c.registerX))
		lo := c.readMemory(uint16(ptr))
		hi := c.readMemory(uint16(ptr + 1))

		return uint16(hi)<<8 | uint16(lo), nil
	case IndirectY:
		base := c.readMemory(c.programCounter)
		lo := c.readMemory(uint16(base))
		hi := c.readMemory(uint16(base + 1))
		derefBase := uint16(hi)<<8 | uint16(lo)
		deref := derefBase + uint16(c.registerY)

		return deref, nil
	default:
		return 0, fmt.Errorf("unknown addressing mode: %d", mode)
	}
}
