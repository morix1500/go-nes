package nes

import (
	"fmt"
	"strings"
)

func trace(cpu *CPU) string {
	var opsInfo OpeCode
	code := cpu.readMemory(cpu.programCounter)
	opsInfo = CPU_OPS_CODES[code]

	begin := cpu.programCounter
	hexDump := []uint8{}
	hexDump = append(hexDump, code)

	var memoryAddr uint16
	var storedValue uint8

	switch opsInfo.Mode {
	case IMMEDIATE, IMPLIED, ACCUMULATOR, RELATIVE, INDIRECT:
		memoryAddr = 0
		storedValue = 0
	default:
		memoryAddr, _ = cpu.getAbsoluteAddress(opsInfo, begin+1)
		storedValue = cpu.readMemory(memoryAddr)
	}
	if begin == 0xE545 {
		fmt.Print()
	}

	var tmp string

	switch opsInfo.Length {
	case 1:
		switch code {
		case 0x0a, 0x4a, 0x2a, 0x6a:
			tmp = "A "
		}
	case 2:
		address := cpu.readMemory(begin + 1)
		hexDump = append(hexDump, address)

		switch opsInfo.Mode {
		case IMMEDIATE:
			tmp = fmt.Sprintf("#$%02X", address)
		case ZERO_PAGE:
			tmp = fmt.Sprintf("$%02X = %02X", memoryAddr, storedValue)
		case ZERO_PAGE_X:
			tmp = fmt.Sprintf("$%02X,X @ %02X = %02X", address, memoryAddr, storedValue)
		case ZERO_PAGE_Y:
			tmp = fmt.Sprintf("$%02X,Y @ %02X = %02X", address, memoryAddr, storedValue)
		case INDIRECT_X:
			tmp = fmt.Sprintf("($%02X,X) @ %02X = %04X = %02X", address, (address + cpu.registerX), memoryAddr, storedValue)
		case INDIRECT_Y:
			tmp = fmt.Sprintf("($%02X),Y = %04X @ %04X = %02X", address, (memoryAddr - uint16(cpu.registerY)), memoryAddr, storedValue)
		case RELATIVE:
			address := uint16(cpu.readMemory(begin + 1))
			if address > 0x7f {
				address = uint16(address) - uint16(0x100)
			}
			address = begin + address + 2
			tmp = fmt.Sprintf("$%04X", address)
		default:
			var add uint
			add = uint(begin) + 2 + uint(address)
			tmp = fmt.Sprintf("$%04X", add)
		}
	case 3:
		addressLo := cpu.readMemory(begin + 1)
		addressHi := cpu.readMemory(begin + 2)
		hexDump = append(hexDump, addressLo)
		hexDump = append(hexDump, addressHi)

		address := cpu.readMemory16(begin + 1)

		switch opsInfo.Mode {
		case IMPLIED, ACCUMULATOR, RELATIVE, INDIRECT:
			if code == 0x6c {
				// jmp indirect
				jmpAddr := cpu.readMemory16(address)
				if address&0xff == 0xff {
					lo := cpu.readMemory(address)
					hi := cpu.readMemory(address & 0xff00)
					jmpAddr = uint16(hi)<<8 | uint16(lo)
				}
				tmp = fmt.Sprintf("($%04X) = %04X", address, jmpAddr)

			} else {
				tmp = fmt.Sprintf("$%04X", address)
			}
		case ABSOLUTE:
			if code == 0x4c || code == 0x20 {
				// jmp/jsr absolute
				tmp = fmt.Sprintf("$%04X", address)
			} else {
				tmp = fmt.Sprintf("$%04X = %02X", address, storedValue)
			}
		case ABSOLUTE_X:
			tmp = fmt.Sprintf("$%04X,X @ %04X = %02X", address, memoryAddr, storedValue)
		case ABSOLUTE_Y:
			tmp = fmt.Sprintf("$%04X,Y @ %04X = %02X", address, memoryAddr, storedValue)
		}
	}

	var hexStrs []string
	for _, b := range hexDump {
		hexStrs = append(hexStrs, fmt.Sprintf("%02X", b))
	}
	hexStr := strings.Join(hexStrs, " ")

	asmStr := fmt.Sprintf("%04X %-8s %4s %s", begin, hexStr, opsInfo.Mnemonic, tmp)

	return fmt.Sprintf("%-47s A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU: 0,0 CYC:%d", asmStr, cpu.registerA, cpu.registerX, cpu.registerY, cpu.status, cpu.stackPointer, cpu.bus.Cycles)
}
