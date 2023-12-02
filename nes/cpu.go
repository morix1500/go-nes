package nes

import "fmt"

type AddressingMode uint8

const (
	IMMEDIATE AddressingMode = iota
	ZERO_PAGE
	ZERO_PAGE_X
	ZERO_PAGE_Y
	ABSOLUTE
	ABSOLUTE_X
	ABSOLUTE_Y
	INDIRECT_X
	INDIRECT_Y
	NONE_ADDRESSING
)

const (
	CPU_FLAG_CARRY             uint8 = 0b0000_0001
	CPU_FLAG_ZERO              uint8 = 0b0000_0010
	CPU_FLAG_INTERRUPT_DISABLE uint8 = 0b0000_0100
	CPU_FLAG_DECIMAL_MODE      uint8 = 0b0000_1000
	CPU_FLAG_BREAK             uint8 = 0b0001_0000
	CPU_FLAG_BREAK2            uint8 = 0b0010_0000
	CPU_FLAG_OVERFLOW          uint8 = 0b0100_0000
	CPU_FLAG_NEGATIVE          uint8 = 0b1000_0000
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

func (c *CPU) adc(mode AddressingMode) error {
	addr, err := c.getOperandAddress(mode)
	if err != nil {
		return err
	}

	// Aレジスタの値とメモリの値、そしてキャリーフラグの値（0または1）を加えます。
	tmpValue := uint16(c.registerA) + uint16(c.readMemory(addr)) + uint16(c.status&CPU_FLAG_CARRY)

	// 加算の結果が255を超える場合、キャリーフラグをセットします。それ以外の場合はキャリーフラグをクリアします。
	if tmpValue&0x100 != 0 {
		tmpValue -= (0xff + 1)
		c.status |= CPU_FLAG_CARRY
	} else {
		c.status &= ^CPU_FLAG_CARRY
	}

	// 加算の結果をAレジスタに格納します。結果が256以上の場合、結果から256を引いた値がAレジスタに格納されます
	c.registerA = uint8(tmpValue & 0xff)
	c.updateZeroAndNegativeFlags(c.registerA)

	// 結果がオーバーフロー（符号付き加算において結果が-128から127の範囲を超える）した場合、オーバーフローフラグをセットします。それ以外の場合はオーバーフローフラグをクリアします。
	if tmpValue > 0x7f {
		c.status |= CPU_FLAG_OVERFLOW
	} else {
		c.status &= ^CPU_FLAG_OVERFLOW
	}

	return nil
}

func (c *CPU) and(mode AddressingMode) error {
	addr, err := c.getOperandAddress(mode)
	if err != nil {
		return err
	}

	c.registerA &= c.readMemory(addr)
	c.updateZeroAndNegativeFlags(c.registerA)

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
		c.status |= CPU_FLAG_ZERO
	} else {
		c.status &= ^CPU_FLAG_ZERO
	}

	if result&CPU_FLAG_NEGATIVE != 0 {
		c.status |= CPU_FLAG_NEGATIVE
	} else {
		c.status &= ^CPU_FLAG_NEGATIVE
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
		case "ADC":
			if err := c.adc(opsInfo.Mode); err != nil {
				return err
			}
		case "AND":
			if err := c.and(opsInfo.Mode); err != nil {
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
	case IMMEDIATE:
		return c.programCounter, nil
	case ZERO_PAGE:
		return uint16(c.readMemory(c.programCounter)), nil
	case ZERO_PAGE_X:
		return uint16(c.readMemory(c.programCounter) + c.registerX), nil
	case ZERO_PAGE_Y:
		return uint16(c.readMemory(c.programCounter) + c.registerY), nil
	case ABSOLUTE:
		return c.readMemory16(c.programCounter), nil
	case ABSOLUTE_X:
		return c.readMemory16(c.programCounter) + uint16(c.registerX), nil
	case ABSOLUTE_Y:
		return c.readMemory16(c.programCounter) + uint16(c.registerY), nil
	case INDIRECT_X:
		base := c.readMemory16(c.programCounter)
		ptr := uint8(base + uint16(c.registerX))
		lo := c.readMemory(uint16(ptr))
		hi := c.readMemory(uint16(ptr + 1))

		return uint16(hi)<<8 | uint16(lo), nil
	case INDIRECT_Y:
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
