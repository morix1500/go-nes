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
	INDIRECT
	INDIRECT_X
	INDIRECT_Y
	ACCUMULATOR
	RELATIVE
	IMPLIED
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
	stackPointer   uint8
	memory         [0xFFFF]uint8
}

func NewCPU() *CPU {
	return &CPU{
		registerA:      0,
		registerX:      0,
		registerY:      0,
		status:         0,
		programCounter: 0,
		stackPointer:   0xfd,
	}
}

func (c *CPU) lda(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.readMemory(addr)
	c.setRegisterA(value)
}

func (c *CPU) ldx(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.readMemory(addr)
	c.registerX = value
	c.updateZeroAndNegativeFlags(c.registerX)
}

func (c *CPU) ldy(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.readMemory(addr)
	c.registerY = value
	c.updateZeroAndNegativeFlags(c.registerY)
}

func (c *CPU) sta(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	c.writeMemory(addr, c.registerA)
}

func (c *CPU) stx(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	c.writeMemory(addr, c.registerX)
}

func (c *CPU) sty(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	c.writeMemory(addr, c.registerY)
}

func (c *CPU) adc(mode AddressingMode) {
	addr := c.getOperandAddress(mode)

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
	c.setRegisterA(uint8(tmpValue & 0xff))

	// 結果がオーバーフロー（符号付き加算において結果が-128から127の範囲を超える）した場合、オーバーフローフラグをセットします。それ以外の場合はオーバーフローフラグをクリアします。
	// 符号付き整数の加算は存在しないため、マイナスの場合は考慮していない
	if tmpValue > 0x7f {
		c.status |= CPU_FLAG_OVERFLOW
	} else {
		c.status &= ^CPU_FLAG_OVERFLOW
	}
}

func (c *CPU) and(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	c.setRegisterA(c.registerA & c.readMemory(addr))
}

func (c *CPU) asl(mode AddressingMode) {
	if mode == ACCUMULATOR {
		if c.registerA&0x80 != 0 {
			c.status |= CPU_FLAG_CARRY
		} else {
			c.status &= ^CPU_FLAG_CARRY
		}
		c.setRegisterA(c.registerA << 1)
		return
	}

	addr := c.getOperandAddress(mode)
	value := c.readMemory(addr)
	if value&0x80 != 0 {
		c.status |= CPU_FLAG_CARRY
	} else {
		c.status &= ^CPU_FLAG_CARRY
	}
	value <<= 1
	c.writeMemory(addr, value)
	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) bcc() {
	if c.status&CPU_FLAG_CARRY == 0 {
		addr := c.getOperandAddress(RELATIVE)
		c.programCounter += addr + 1
	}
}

func (c *CPU) bcs() {
	if c.status&CPU_FLAG_CARRY != 0 {
		addr := c.getOperandAddress(RELATIVE)
		c.programCounter += addr + 1
	}
}

func (c *CPU) beq() {
	if c.status&CPU_FLAG_ZERO != 0 {
		addr := c.getOperandAddress(RELATIVE)
		c.programCounter += addr + 1
	}
}

func (c *CPU) bit(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.readMemory(addr)
	if value&c.registerA == 0 {
		c.status |= CPU_FLAG_ZERO
	} else {
		c.status &= ^CPU_FLAG_ZERO
	}

	if value&CPU_FLAG_NEGATIVE != 0 {
		c.status |= CPU_FLAG_NEGATIVE
	} else {
		c.status &= ^CPU_FLAG_NEGATIVE
	}

	if value&CPU_FLAG_OVERFLOW != 0 {
		c.status |= CPU_FLAG_OVERFLOW
	} else {
		c.status &= ^CPU_FLAG_OVERFLOW
	}
}

func (c *CPU) bmi() {
	if c.status&CPU_FLAG_NEGATIVE != 0 {
		addr := c.getOperandAddress(RELATIVE)
		c.programCounter += addr + 1
	}
}

func (c *CPU) bne() {
	if c.status&CPU_FLAG_ZERO == 0 {
		addr := c.getOperandAddress(RELATIVE)
		c.programCounter += addr + 1
	}
}

func (c *CPU) bpl() {
	if c.status&CPU_FLAG_NEGATIVE == 0 {
		addr := c.getOperandAddress(RELATIVE)
		c.programCounter += addr + 1
	}
}

func (c *CPU) bvc() {
	if c.status&CPU_FLAG_OVERFLOW == 0 {
		addr := c.getOperandAddress(RELATIVE)
		c.programCounter += addr + 1
	}
}

func (c *CPU) bvs() {
	if c.status&CPU_FLAG_OVERFLOW != 0 {
		addr := c.getOperandAddress(RELATIVE)
		c.programCounter += addr + 1
	}
}

func (c *CPU) clc() {
	c.status &= ^CPU_FLAG_CARRY
}

func (c *CPU) cld() {
	c.status &= ^CPU_FLAG_DECIMAL_MODE
}

func (c *CPU) cli() {
	c.status &= ^CPU_FLAG_INTERRUPT_DISABLE
}

func (c *CPU) clv() {
	c.status &= ^CPU_FLAG_OVERFLOW
}

func (c *CPU) cmp(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.readMemory(addr)
	if c.registerA >= value {
		c.status |= CPU_FLAG_CARRY
	} else {
		c.status &= ^CPU_FLAG_CARRY
	}

	c.updateZeroAndNegativeFlags(c.registerA - value)
}

func (c *CPU) cpx(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.readMemory(addr)
	if c.registerX >= value {
		c.status |= CPU_FLAG_CARRY
	} else {
		c.status &= ^CPU_FLAG_CARRY
	}

	c.updateZeroAndNegativeFlags(c.registerX - value)
}

func (c *CPU) cpy(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.readMemory(addr)
	if c.registerY >= value {
		c.status |= CPU_FLAG_CARRY
	} else {
		c.status &= ^CPU_FLAG_CARRY
	}

	c.updateZeroAndNegativeFlags(c.registerY - value)
}

func (c *CPU) dec(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.readMemory(addr)
	value--
	c.writeMemory(addr, value)
	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) dex() {
	c.registerX--
	c.updateZeroAndNegativeFlags(c.registerX)
}

func (c *CPU) dey() {
	c.registerY--
	c.updateZeroAndNegativeFlags(c.registerY)
}

func (c *CPU) eor(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	c.setRegisterA(c.registerA ^ c.readMemory(addr))
}

func (c *CPU) inc(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	value := c.readMemory(addr)
	value++
	c.writeMemory(addr, value)
	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) inx() {
	c.registerX++
	c.updateZeroAndNegativeFlags(c.registerX)
}

func (c *CPU) iny() {
	c.registerY++
	c.updateZeroAndNegativeFlags(c.registerY)
}

func (c *CPU) jmp(mode AddressingMode) {
	addr := c.getOperandAddress(mode)

	if mode == ABSOLUTE {
		c.programCounter = addr
		return
	}

	// An original 6502 has does not correctly fetch the target address
	// if the indirect vector falls on a page boundary (e.g. $xxFF where xx is any value from $00 to $FF).
	// In this case fetches the LSB from $xxFF as expected but takes the MSB from $xx00.
	// This is fixed in some later chips like the 65SC02 so for compatibility always ensure
	// the indirect vector is not at the end of the page.
	var indirectAddr uint16
	if addr&0x00ff == 0x00ff {
		lo := c.readMemory(addr)
		hi := c.readMemory(addr & 0xff00)
		indirectAddr = uint16(hi)<<8 | uint16(lo)
	} else {
		indirectAddr = addr
	}

	c.programCounter = indirectAddr
}

func (c *CPU) jsr() {
	addr := c.getOperandAddress(ABSOLUTE)
	// 2バイト加算している理由は、JSR命令の次の命令を実行するため
	c.stackPush16(c.programCounter + 2 - 1)
	c.programCounter = addr
}

func (c *CPU) lsr(mode AddressingMode) {
	if mode == ACCUMULATOR {
		if c.registerA&CPU_FLAG_CARRY != 0 {
			c.status |= CPU_FLAG_CARRY
		} else {
			c.status &= ^CPU_FLAG_CARRY
		}
		c.setRegisterA(c.registerA >> 1)
		return
	}

	addr := c.getOperandAddress(mode)
	value := c.readMemory(addr)
	if value&CPU_FLAG_CARRY != 0 {
		c.status |= CPU_FLAG_CARRY
	} else {
		c.status &= ^CPU_FLAG_CARRY
	}
	value >>= 1
	c.writeMemory(addr, value)
	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) ora(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	c.setRegisterA(c.registerA | c.readMemory(addr))
}

func (c *CPU) pha() {
	c.stackPush(c.registerA)
}

func (c *CPU) php() {
	// https://www.nesdev.org/wiki/Status_flags#The_B_flag
	c.stackPush(c.status | CPU_FLAG_BREAK | CPU_FLAG_BREAK2)
}

func (c *CPU) pla() {
	c.setRegisterA(c.stackPop())
}

func (c *CPU) plp() {
	c.status = c.stackPop()&^CPU_FLAG_BREAK | CPU_FLAG_BREAK2
}

func (c *CPU) rol(mode AddressingMode) {
	oldCarry := c.status & CPU_FLAG_CARRY

	if mode == ACCUMULATOR {
		if c.registerA&0x80 != 0 {
			c.status |= CPU_FLAG_CARRY
		} else {
			c.status &= ^CPU_FLAG_CARRY
		}
		c.setRegisterA(c.registerA<<1 | oldCarry)
		return
	}

	addr := c.getOperandAddress(mode)
	value := c.readMemory(addr)
	if value&0x80 != 0 {
		c.status |= CPU_FLAG_CARRY
	} else {
		c.status &= ^CPU_FLAG_CARRY
	}
	value <<= 1
	value |= oldCarry
	c.writeMemory(addr, value)
	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) ror(mode AddressingMode) {
	oldCarry := c.status & CPU_FLAG_CARRY
	oldCarry <<= 7

	if mode == ACCUMULATOR {
		if c.registerA&CPU_FLAG_CARRY != 0 {
			c.status |= CPU_FLAG_CARRY
		} else {
			c.status &= ^CPU_FLAG_CARRY
		}

		c.setRegisterA(c.registerA>>1 | oldCarry)
		return
	}

	addr := c.getOperandAddress(mode)
	value := c.readMemory(addr)
	if value&CPU_FLAG_CARRY != 0 {
		c.status |= CPU_FLAG_CARRY
	} else {
		c.status &= ^CPU_FLAG_CARRY
	}
	value >>= 1
	value |= oldCarry
	c.writeMemory(addr, value)
	c.updateZeroAndNegativeFlags(value)
}

func (c *CPU) rts() {
	c.programCounter = c.stackPop16() + 1
}

func (c *CPU) sec() {
	c.status |= CPU_FLAG_CARRY
}

func (c *CPU) sed() {
	c.status |= CPU_FLAG_DECIMAL_MODE
}

func (c *CPU) sei() {
	c.status |= CPU_FLAG_INTERRUPT_DISABLE
}

func (c *CPU) tax() {
	c.registerX = c.registerA
	c.updateZeroAndNegativeFlags(c.registerX)
}

func (c *CPU) tay() {
	c.registerY = c.registerA
	c.updateZeroAndNegativeFlags(c.registerY)
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

func (c *CPU) setRegisterA(value uint8) {
	c.registerA = value
	c.updateZeroAndNegativeFlags(c.registerA)
}

func (c *CPU) stackPush(value uint8) {
	c.writeMemory(0x0100+uint16(c.stackPointer), value)
	c.stackPointer--
}

func (c *CPU) stackPop() uint8 {
	c.stackPointer++
	return c.readMemory(0x0100 + uint16(c.stackPointer))
}

func (c *CPU) stackPush16(value uint16) {
	hi := uint8(value >> 8)
	lo := uint8(value & 0xff)
	c.stackPush(hi)
	c.stackPush(lo)
}

func (c *CPU) stackPop16() uint16 {
	lo := uint16(c.stackPop())
	hi := uint16(c.stackPop())
	return (hi << 8) | lo
}

func (c *CPU) LoadAndRun(program []uint8) {
	c.Load(program)
	c.Reset()
	c.Run()
}

func (c *CPU) Load(program []uint8) {
	copy(c.memory[0x8000:(0x8000+len(program))], program)
	c.writeMemory16(0xFFFC, 0x8000)
}

func (c *CPU) Reset() {
	c.registerA = 0
	c.registerX = 0
	c.registerY = 0
	c.status = 0
	c.programCounter = c.readMemory16(0xFFFC)
	c.stackPointer = 0xfd
}

func (c *CPU) Run() {
	var opsInfo OpeCode
	var ok bool
	for {
		code := c.readMemory(c.programCounter)
		c.programCounter++
		programCounterState := c.programCounter

		if opsInfo, ok = CPU_OPS_CODES[code]; !ok {
			panic(fmt.Sprintf("unknown code: %d", code))
		}

		switch opsInfo.Mnemonic {
		case "BRK":
			return
		case "ADC":
			c.adc(opsInfo.Mode)
		case "AND":
			c.and(opsInfo.Mode)
		case "ASL":
			c.asl(opsInfo.Mode)
		case "BCC":
			c.bcc()
		case "BCS":
			c.bcs()
		case "BEQ":
			c.beq()
		case "BIT":
			c.bit(opsInfo.Mode)
		case "BMI":
			c.bmi()
		case "BNE":
			c.bne()
		case "BPL":
			c.bpl()
		case "BVC":
			c.bvc()
		case "BVS":
			c.bvs()
		case "CLC":
			c.clc()
		case "CLD":
			c.cld()
		case "CLI":
			c.cli()
		case "CLV":
			c.clv()
		case "CMP":
			c.cmp(opsInfo.Mode)
		case "CPX":
			c.cpx(opsInfo.Mode)
		case "CPY":
			c.cpy(opsInfo.Mode)
		case "DEC":
			c.dec(opsInfo.Mode)
		case "DEX":
			c.dex()
		case "DEY":
			c.dey()
		case "EOR":
			c.eor(opsInfo.Mode)
		case "INC":
			c.inc(opsInfo.Mode)
		case "INX":
			c.inx()
		case "INY":
			c.iny()
		case "JMP":
			c.jmp(opsInfo.Mode)
		case "JSR":
			c.jsr()
		case "LDA":
			c.lda(opsInfo.Mode)
		case "LDX":
			c.ldx(opsInfo.Mode)
		case "LDY":
			c.ldy(opsInfo.Mode)
		case "LSR":
			c.lsr(opsInfo.Mode)
		case "NOP":
			// 何もしない
		case "ORA":
			c.ora(opsInfo.Mode)
		case "PHA":
			c.pha()
		case "PHP":
			c.php()
		case "PLA":
			c.pla()
		case "PLP":
			c.plp()
		case "ROL":
			c.rol(opsInfo.Mode)
		case "ROR":
			c.ror(opsInfo.Mode)
		case "RTS":
			c.rts()
		case "SEC":
			c.sec()
		case "SED":
			c.sed()
		case "SEI":
			c.sei()
		case "STA":
			c.sta(opsInfo.Mode)
		case "STX":
			c.stx(opsInfo.Mode)
		case "STY":
			c.sty(opsInfo.Mode)
		case "TAX":
			c.tax()
		case "TAY":
			c.tay()
		}
		if programCounterState == c.programCounter {
			c.programCounter += uint16(opsInfo.Length - 1)
		}
	}
}

func (c *CPU) getOperandAddress(mode AddressingMode) uint16 {
	switch mode {
	case IMMEDIATE:
		return c.programCounter
	case ZERO_PAGE:
		return uint16(c.readMemory(c.programCounter))
	case ZERO_PAGE_X:
		return uint16(c.readMemory(c.programCounter) + c.registerX)
	case ZERO_PAGE_Y:
		return uint16(c.readMemory(c.programCounter) + c.registerY)
	case ABSOLUTE:
		return c.readMemory16(c.programCounter)
	case ABSOLUTE_X:
		return c.readMemory16(c.programCounter) + uint16(c.registerX)
	case ABSOLUTE_Y:
		return c.readMemory16(c.programCounter) + uint16(c.registerY)
	case INDIRECT:
		return c.readMemory16(c.programCounter)
	case INDIRECT_X:
		base := c.readMemory16(c.programCounter)
		ptr := uint8(base + uint16(c.registerX))
		lo := c.readMemory(uint16(ptr))
		hi := c.readMemory(uint16(ptr + 1))

		return uint16(hi)<<8 | uint16(lo)
	case INDIRECT_Y:
		base := c.readMemory(c.programCounter)
		lo := c.readMemory(uint16(base))
		hi := c.readMemory(uint16(base + 1))
		derefBase := uint16(hi)<<8 | uint16(lo)
		deref := derefBase + uint16(c.registerY)

		return deref
	case ACCUMULATOR:
		return 0
	case RELATIVE:
		// オペランドは符号付き8ビットのオフセットとして解釈される
		address := uint16(c.readMemory(c.programCounter))
		if address > 0x7f {
			address = uint16(address) - uint16(0x100)
		}
		return address
	case IMPLIED:
		return 0
	default:
		panic(fmt.Sprintf("unknown addressing mode: %d", mode))
	}
}
