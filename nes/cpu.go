package nes

import (
	"fmt"
	"os"
)

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
	bus            *Bus
	cycle          uint8
}

func NewCPU(bus *Bus) *CPU {
	return &CPU{
		registerA:      0,
		registerX:      0,
		registerY:      0,
		status:         0,
		programCounter: 0,
		stackPointer:   0xfd,
		bus:            bus,
		cycle:          0,
	}
}

func (c *CPU) lda(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
	value := c.readMemory(addr)
	c.setRegisterA(value)
}

func (c *CPU) ldx(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
	value := c.readMemory(addr)
	c.registerX = value
	c.updateZeroAndNegativeFlags(c.registerX)
}

func (c *CPU) ldy(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
	value := c.readMemory(addr)
	c.registerY = value
	c.updateZeroAndNegativeFlags(c.registerY)
}

func (c *CPU) sta(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
	c.writeMemory(addr, c.registerA)
}

func (c *CPU) stx(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
	c.writeMemory(addr, c.registerX)
}

func (c *CPU) sty(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
	c.writeMemory(addr, c.registerY)
}

func (cpu *CPU) adc(opsInfo OpeCode) {
	addr := cpu.getOperandAddress(opsInfo)
	a := cpu.registerA
	b := cpu.readMemory(addr)
	c := cpu.status & CPU_FLAG_CARRY

	cpu.setRegisterA(a + b + c)

	if int(a)+int(b)+int(c) > 0xff {
		cpu.status |= CPU_FLAG_CARRY
	} else {
		cpu.status &= ^CPU_FLAG_CARRY
	}

	// 同じ符号同士の演算かつ計算結果の符号が元の符号と異なる場合オーバーフローとなる
	if (a^b)&0x80 == 0 && (a^cpu.registerA)&0x80 != 0 {
		cpu.status |= CPU_FLAG_OVERFLOW
	} else {
		cpu.status &= ^CPU_FLAG_OVERFLOW
	}
}

func (c *CPU) and(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
	c.setRegisterA(c.registerA & c.readMemory(addr))
}

func (c *CPU) asl(opsInfo OpeCode) {
	if opsInfo.Mode == ACCUMULATOR {
		if c.registerA&0x80 != 0 {
			c.status |= CPU_FLAG_CARRY
		} else {
			c.status &= ^CPU_FLAG_CARRY
		}
		c.setRegisterA(c.registerA << 1)
		return
	}

	addr := c.getOperandAddress(opsInfo)
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
		addr := c.getOperandAddress(CPU_OPS_CODES[0x90])
		c.addBranchCycles(addr)
		c.programCounter += addr + 1
	}
}

func (c *CPU) bcs() {
	if c.status&CPU_FLAG_CARRY != 0 {
		addr := c.getOperandAddress(CPU_OPS_CODES[0xb0])
		c.addBranchCycles(addr)
		c.programCounter += addr + 1
	}
}

func (c *CPU) beq() {
	if c.status&CPU_FLAG_ZERO != 0 {
		addr := c.getOperandAddress(CPU_OPS_CODES[0xf0])
		c.addBranchCycles(addr)
		c.programCounter += addr + 1
	}
}

func (c *CPU) bit(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
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
		addr := c.getOperandAddress(CPU_OPS_CODES[0x30])
		c.addBranchCycles(addr)
		c.programCounter += addr + 1
	}
}

func (c *CPU) bne() {
	if c.status&CPU_FLAG_ZERO == 0 {
		addr := c.getOperandAddress(CPU_OPS_CODES[0xd0])
		c.addBranchCycles(addr)
		c.programCounter += addr + 1
	}
}

func (c *CPU) bpl() {
	if c.status&CPU_FLAG_NEGATIVE == 0 {
		addr := c.getOperandAddress(CPU_OPS_CODES[0x10])
		c.addBranchCycles(addr)
		c.programCounter += addr + 1
	}
}

func (c *CPU) bvc() {
	if c.status&CPU_FLAG_OVERFLOW == 0 {
		addr := c.getOperandAddress(CPU_OPS_CODES[0x50])
		c.addBranchCycles(addr)
		c.programCounter += addr + 1
	}
}

func (c *CPU) bvs() {
	if c.status&CPU_FLAG_OVERFLOW != 0 {
		addr := c.getOperandAddress(CPU_OPS_CODES[0x70])
		c.addBranchCycles(addr)
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

func (c *CPU) cmp(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
	value := c.readMemory(addr)
	if c.registerA >= value {
		c.status |= CPU_FLAG_CARRY
	} else {
		c.status &= ^CPU_FLAG_CARRY
	}

	c.updateZeroAndNegativeFlags(c.registerA - value)
}

func (c *CPU) cpx(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
	value := c.readMemory(addr)
	if c.registerX >= value {
		c.status |= CPU_FLAG_CARRY
	} else {
		c.status &= ^CPU_FLAG_CARRY
	}

	c.updateZeroAndNegativeFlags(c.registerX - value)
}

func (c *CPU) cpy(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
	value := c.readMemory(addr)
	if c.registerY >= value {
		c.status |= CPU_FLAG_CARRY
	} else {
		c.status &= ^CPU_FLAG_CARRY
	}

	c.updateZeroAndNegativeFlags(c.registerY - value)
}

func (c *CPU) dec(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
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

func (c *CPU) eor(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
	c.setRegisterA(c.registerA ^ c.readMemory(addr))
}

func (c *CPU) inc(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
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

func (c *CPU) jmp(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)

	if opsInfo.Mode == ABSOLUTE {
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
		indirectAddr = c.readMemory16(addr)
	}

	c.programCounter = indirectAddr
}

func (c *CPU) jsr() {
	addr := c.getOperandAddress(CPU_OPS_CODES[0x20])
	// 2バイト加算している理由は、JSR命令の次の命令を実行するため
	c.stackPush16(c.programCounter + 2 - 1)
	c.programCounter = addr
}

func (c *CPU) lsr(opsInfo OpeCode) {
	if opsInfo.Mode == ACCUMULATOR {
		if c.registerA&CPU_FLAG_CARRY != 0 {
			c.status |= CPU_FLAG_CARRY
		} else {
			c.status &= ^CPU_FLAG_CARRY
		}
		c.setRegisterA(c.registerA >> 1)
		return
	}

	addr := c.getOperandAddress(opsInfo)
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

func (c *CPU) ora(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
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

func (c *CPU) rol(opsInfo OpeCode) {
	oldCarry := c.status & CPU_FLAG_CARRY

	if opsInfo.Mode == ACCUMULATOR {
		if c.registerA&0x80 != 0 {
			c.status |= CPU_FLAG_CARRY
		} else {
			c.status &= ^CPU_FLAG_CARRY
		}
		c.setRegisterA(c.registerA<<1 | oldCarry)
		return
	}

	addr := c.getOperandAddress(opsInfo)
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

func (c *CPU) ror(opsInfo OpeCode) {
	oldCarry := c.status & CPU_FLAG_CARRY
	oldCarry <<= 7

	if opsInfo.Mode == ACCUMULATOR {
		if c.registerA&CPU_FLAG_CARRY != 0 {
			c.status |= CPU_FLAG_CARRY
		} else {
			c.status &= ^CPU_FLAG_CARRY
		}

		c.setRegisterA(c.registerA>>1 | oldCarry)
		return
	}

	addr := c.getOperandAddress(opsInfo)
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

func (c *CPU) rti() {
	c.status = c.stackPop()&^CPU_FLAG_BREAK | CPU_FLAG_BREAK2
	c.programCounter = c.stackPop16()
}

func (c *CPU) rts() {
	c.programCounter = c.stackPop16() + 1
}

func (cpu *CPU) sbc(opsInfo OpeCode) {
	addr := cpu.getOperandAddress(opsInfo)
	a := cpu.registerA
	b := cpu.readMemory(addr)
	c := cpu.status & CPU_FLAG_CARRY

	cpu.setRegisterA(a - b - (1 - c))

	// 引き算の結果が0未満の場合、キャリーフラグをクリアします。それ以外の場合はキャリーフラグをセットします。
	if int(a)-int(b)-int(1-c) >= 0 {
		cpu.status |= CPU_FLAG_CARRY
	} else {
		cpu.status &= ^CPU_FLAG_CARRY
	}

	// 異なる符号同士の演算かつ計算結果の符号が 1 つ目の整数の符号と異なる場合オーバーフローとなる
	if (a^b)&0x80 != 0 && (a^cpu.registerA)&0x80 != 0 {
		cpu.status |= CPU_FLAG_OVERFLOW
	} else {
		cpu.status &= ^CPU_FLAG_OVERFLOW
	}
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

func (c *CPU) tsx() {
	c.registerX = c.stackPointer
	c.updateZeroAndNegativeFlags(c.registerX)
}

func (c *CPU) txa() {
	c.setRegisterA(c.registerX)
}

func (c *CPU) txs() {
	c.stackPointer = c.registerX
}

func (c *CPU) tya() {
	c.setRegisterA(c.registerY)
}

func (c *CPU) brk() {
	c.stackPush16(c.programCounter)
	c.php()
	c.sei()
	c.programCounter = c.readMemory16(0xfffe)
}

func (c *CPU) lax(opsInfo OpeCode) {
	c.lda(opsInfo)
	c.tax()
}

func (c *CPU) sax(opsInfo OpeCode) {
	addr := c.getOperandAddress(opsInfo)
	c.writeMemory(addr, c.registerA&c.registerX)
}

func (c *CPU) dcp(opsInfo OpeCode) {
	c.dec(opsInfo)
	c.cmp(opsInfo)
}

func (c *CPU) isb(opsInfo OpeCode) {
	c.inc(opsInfo)
	c.sbc(opsInfo)
}

func (c *CPU) slo(opsInfo OpeCode) {
	c.asl(opsInfo)
	c.ora(opsInfo)
}

func (c *CPU) rla(opsInfo OpeCode) {
	c.rol(opsInfo)
	c.and(opsInfo)
}

func (c *CPU) sre(opsInfo OpeCode) {
	c.lsr(opsInfo)
	c.eor(opsInfo)
}

func (c *CPU) rra(opsInfo OpeCode) {
	c.ror(opsInfo)
	c.adc(opsInfo)
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
	return c.bus.ReadMemory(address)
}

func (c *CPU) readMemory16(address uint16) uint16 {
	lo := uint16(c.readMemory(address))
	hi := uint16(c.readMemory(address + 1))
	return (hi << 8) | lo
}

func (c *CPU) writeMemory(address uint16, value uint8) {
	c.bus.WriteMemory(address, value)
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

func (c *CPU) Reset() {
	c.registerA = 0
	c.registerX = 0
	c.registerY = 0
	c.status = 0b00100100
	c.programCounter = c.readMemory16(0xFFFC)
	if os.Getenv("CPU_TEST") == "true" {
		c.programCounter = 0xc000
	}
	c.stackPointer = 0xfd
}

func (c *CPU) Run() {
	var opsInfo OpeCode
	var ok bool
	for {
		if c.bus.PollNMIStatus() {
			c.InterruptNMI()
		}

		//fmt.Println(trace(c))
		code := c.readMemory(c.programCounter)
		c.programCounter++
		programCounterState := c.programCounter

		if opsInfo, ok = CPU_OPS_CODES[code]; !ok {
			panic(fmt.Sprintf("unknown code: %d", code))
		}

		c.cycle = opsInfo.Cycles

		switch opsInfo.Mnemonic {
		case "BRK":
			//c.brk()
			return
		case "ADC":
			c.adc(opsInfo)
		case "AND":
			c.and(opsInfo)
		case "ASL":
			c.asl(opsInfo)
		case "BCC":
			c.bcc()
		case "BCS":
			c.bcs()
		case "BEQ":
			c.beq()
		case "BIT":
			c.bit(opsInfo)
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
			c.cmp(opsInfo)
		case "CPX":
			c.cpx(opsInfo)
		case "CPY":
			c.cpy(opsInfo)
		case "DEC":
			c.dec(opsInfo)
		case "DEX":
			c.dex()
		case "DEY":
			c.dey()
		case "EOR":
			c.eor(opsInfo)
		case "INC":
			c.inc(opsInfo)
		case "INX":
			c.inx()
		case "INY":
			c.iny()
		case "JMP":
			c.jmp(opsInfo)
		case "JSR":
			c.jsr()
		case "LDA":
			c.lda(opsInfo)
		case "LDX":
			c.ldx(opsInfo)
		case "LDY":
			c.ldy(opsInfo)
		case "LSR":
			c.lsr(opsInfo)
		case "NOP", "*NOP":
			// 何もしない
		case "ORA":
			c.ora(opsInfo)
		case "PHA":
			c.pha()
		case "PHP":
			c.php()
		case "PLA":
			c.pla()
		case "PLP":
			c.plp()
		case "ROL":
			c.rol(opsInfo)
		case "ROR":
			c.ror(opsInfo)
		case "RTI":
			c.rti()
		case "RTS":
			c.rts()
		case "SBC", "*SBC":
			c.sbc(opsInfo)
		case "SEC":
			c.sec()
		case "SED":
			c.sed()
		case "SEI":
			c.sei()
		case "STA":
			c.sta(opsInfo)
		case "STX":
			c.stx(opsInfo)
		case "STY":
			c.sty(opsInfo)
		case "TAX":
			c.tax()
		case "TAY":
			c.tay()
		case "TSX":
			c.tsx()
		case "TXA":
			c.txa()
		case "TXS":
			c.txs()
		case "TYA":
			c.tya()
		case "*LAX":
			c.lax(opsInfo)
		case "*SAX":
			c.sax(opsInfo)
		case "*DCP":
			c.dcp(opsInfo)
		case "*ISB":
			c.isb(opsInfo)
		case "*SLO":
			c.slo(opsInfo)
		case "*RLA":
			c.rla(opsInfo)
		case "*RRA":
			c.rra(opsInfo)
		case "*SRE":
			c.sre(opsInfo)
		}

		c.bus.Tick(c.cycle)

		if programCounterState == c.programCounter {
			c.programCounter += uint16(opsInfo.Length - 1)
		}
	}
}

func (c *CPU) getOperandAddress(opsInfo OpeCode) uint16 {
	return c.getAbsoluteAddress(opsInfo, c.programCounter)
}

func (c *CPU) getAbsoluteAddress(opsInfo OpeCode, addr uint16) uint16 {
	var result uint16
	pageCrossed := false

	switch opsInfo.Mode {
	case IMMEDIATE:
		result = addr
	case ZERO_PAGE:
		result = uint16(c.readMemory(addr))
	case ZERO_PAGE_X:
		result = uint16(c.readMemory(addr) + c.registerX)
	case ZERO_PAGE_Y:
		result = uint16(c.readMemory(addr) + c.registerY)
	case ABSOLUTE:
		result = c.readMemory16(addr)
	case ABSOLUTE_X:
		address := c.readMemory16(addr)

		pageCrossed = PageDiffer(address-uint16(c.registerX), address)
		result = c.readMemory16(addr) + uint16(c.registerX)
	case ABSOLUTE_Y:
		address := c.readMemory16(addr)

		pageCrossed = PageDiffer(address-uint16(c.registerY), address)
		result = c.readMemory16(addr) + uint16(c.registerY)
	case INDIRECT:
		result = c.readMemory16(addr)
	case INDIRECT_X:
		base := c.readMemory16(addr)
		ptr := uint8(base + uint16(c.registerX))
		lo := c.readMemory(uint16(ptr))
		hi := c.readMemory(uint16(ptr + 1))

		result = uint16(hi)<<8 | uint16(lo)
	case INDIRECT_Y:
		base := c.readMemory(addr)
		lo := c.readMemory(uint16(base))
		hi := c.readMemory(uint16(base + 1))
		derefBase := uint16(hi)<<8 | uint16(lo)
		deref := derefBase + uint16(c.registerY)

		pageCrossed = PageDiffer(deref-uint16(c.registerY), deref)
		result = deref
	case ACCUMULATOR:
		result = 0
	case RELATIVE:
		// オペランドは符号付き8ビットのオフセットとして解釈される
		address := uint16(c.readMemory(addr))
		if address > 0x7f {
			address = uint16(address) - uint16(0x100)
		}
		result = address
	case IMPLIED:
		result = 0
	default:
		panic(fmt.Sprintf("unknown addressing mode: %d", opsInfo.Mode))
	}

	if pageCrossed && opsInfo.AddCycleIfPageCrossed {
		c.cycle++
	}

	return result
}

func PageDiffer(a, b uint16) bool {
	return a&0xff00 != b&0xff00
}

func (c *CPU) addBranchCycles(address uint16) {
	c.cycle++
	if PageDiffer(c.programCounter, address) {
		c.cycle++
	}
}

func (c *CPU) InterruptNMI() {
	c.stackPush16(c.programCounter)
	status := c.status & ^CPU_FLAG_BREAK | CPU_FLAG_BREAK2
	c.stackPush(status)
	c.status |= CPU_FLAG_INTERRUPT_DISABLE
	c.bus.Tick(2)
	c.programCounter = c.readMemory16(0xfffa)
}
