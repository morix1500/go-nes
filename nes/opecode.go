package nes

type OpeCode struct {
	Code     uint8
	Mnemonic string
	Length   uint8
	Cycles   uint8
	Mode     AddressingMode
}

var CPU_OPS_CODES map[uint8]OpeCode = map[uint8]OpeCode{
	0x69: {Mnemonic: "ADC", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0x65: {Mnemonic: "ADC", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0x75: {Mnemonic: "ADC", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0x6d: {Mnemonic: "ADC", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0x7d: {Mnemonic: "ADC", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},
	0x79: {Mnemonic: "ADC", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_Y},
	0x61: {Mnemonic: "ADC", Length: 2, Cycles: 6, Mode: INDIRECT_X},
	0x71: {Mnemonic: "ADC", Length: 2, Cycles: 5 /* +1 if page crossed*/, Mode: INDIRECT_Y},

	0x29: {Mnemonic: "AND", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0x25: {Mnemonic: "AND", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0x35: {Mnemonic: "AND", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0x2d: {Mnemonic: "AND", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0x3d: {Mnemonic: "AND", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},
	0x39: {Mnemonic: "AND", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_Y},
	0x21: {Mnemonic: "AND", Length: 2, Cycles: 6, Mode: INDIRECT_X},
	0x31: {Mnemonic: "AND", Length: 2, Cycles: 5 /* +1 if page crossed*/, Mode: INDIRECT_Y},

	0x0a: {Mnemonic: "ASL", Length: 1, Cycles: 2, Mode: ACCUMULATOR},
	0x06: {Mnemonic: "ASL", Length: 2, Cycles: 5, Mode: ZERO_PAGE},
	0x16: {Mnemonic: "ASL", Length: 2, Cycles: 6, Mode: ZERO_PAGE_X},
	0x0e: {Mnemonic: "ASL", Length: 3, Cycles: 6, Mode: ABSOLUTE},
	0x1e: {Mnemonic: "ASL", Length: 3, Cycles: 7, Mode: ABSOLUTE_X},

	0x90: {Mnemonic: "BCC", Length: 2, Cycles: 2 /* +1 if branch succeeds, +2 if to a new page*/, Mode: RELATIVE},
	0xb0: {Mnemonic: "BCS", Length: 2, Cycles: 2 /* +1 if branch succeeds, +2 if to a new page*/, Mode: RELATIVE},
	0xf0: {Mnemonic: "BEQ", Length: 2, Cycles: 2 /* +1 if branch succeeds, +2 if to a new page*/, Mode: RELATIVE},

	0x24: {Mnemonic: "BIT", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0x2c: {Mnemonic: "BIT", Length: 3, Cycles: 4, Mode: ABSOLUTE},

	0x30: {Mnemonic: "BMI", Length: 2, Cycles: 2 /* +1 if branch succeeds, +2 if to a new page*/, Mode: RELATIVE},
	0xd0: {Mnemonic: "BNE", Length: 2, Cycles: 2 /* +1 if branch succeeds, +2 if to a new page*/, Mode: RELATIVE},
	0x10: {Mnemonic: "BPL", Length: 2, Cycles: 2 /* +1 if branch succeeds, +2 if to a new page*/, Mode: RELATIVE},
	0x00: {Mnemonic: "BRK", Length: 1, Cycles: 7, Mode: IMPLIED},

	0x50: {Mnemonic: "BVC", Length: 2, Cycles: 2 /* +1 if branch succeeds, +2 if to a new page*/, Mode: RELATIVE},
	0x70: {Mnemonic: "BVS", Length: 2, Cycles: 2 /* +1 if branch succeeds, +2 if to a new page*/, Mode: RELATIVE},
	0x18: {Mnemonic: "CLC", Length: 1, Cycles: 2, Mode: IMPLIED},
	0xd8: {Mnemonic: "CLD", Length: 1, Cycles: 2, Mode: IMPLIED},
	0x58: {Mnemonic: "CLI", Length: 1, Cycles: 2, Mode: IMPLIED},
	0xb8: {Mnemonic: "CLV", Length: 1, Cycles: 2, Mode: IMPLIED},

	0xc9: {Mnemonic: "CMP", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0xc5: {Mnemonic: "CMP", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0xd5: {Mnemonic: "CMP", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0xcd: {Mnemonic: "CMP", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0xdd: {Mnemonic: "CMP", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},
	0xd9: {Mnemonic: "CMP", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_Y},
	0xc1: {Mnemonic: "CMP", Length: 2, Cycles: 6, Mode: INDIRECT_X},
	0xd1: {Mnemonic: "CMP", Length: 2, Cycles: 5 /* +1 if page crossed*/, Mode: INDIRECT_Y},
	0xe0: {Mnemonic: "CPX", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0xe4: {Mnemonic: "CPX", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0xec: {Mnemonic: "CPX", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0xc0: {Mnemonic: "CPY", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0xc4: {Mnemonic: "CPY", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0xcc: {Mnemonic: "CPY", Length: 3, Cycles: 4, Mode: ABSOLUTE},

	0xc6: {Mnemonic: "DEC", Length: 2, Cycles: 5, Mode: ZERO_PAGE},
	0xd6: {Mnemonic: "DEC", Length: 2, Cycles: 6, Mode: ZERO_PAGE_X},
	0xce: {Mnemonic: "DEC", Length: 3, Cycles: 6, Mode: ABSOLUTE},
	0xde: {Mnemonic: "DEC", Length: 3, Cycles: 7, Mode: ABSOLUTE_X},
	0xca: {Mnemonic: "DEX", Length: 1, Cycles: 2, Mode: IMPLIED},
	0x88: {Mnemonic: "DEY", Length: 1, Cycles: 2, Mode: IMPLIED},

	0x49: {Mnemonic: "EOR", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0x45: {Mnemonic: "EOR", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0x55: {Mnemonic: "EOR", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0x4d: {Mnemonic: "EOR", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0x5d: {Mnemonic: "EOR", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},
	0x59: {Mnemonic: "EOR", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_Y},
	0x41: {Mnemonic: "EOR", Length: 2, Cycles: 6, Mode: INDIRECT_X},
	0x51: {Mnemonic: "EOR", Length: 2, Cycles: 5 /* +1 if page crossed*/, Mode: INDIRECT_Y},

	0xe6: {Mnemonic: "INC", Length: 2, Cycles: 5, Mode: ZERO_PAGE},
	0xf6: {Mnemonic: "INC", Length: 2, Cycles: 6, Mode: ZERO_PAGE_X},
	0xee: {Mnemonic: "INC", Length: 3, Cycles: 6, Mode: ABSOLUTE},
	0xfe: {Mnemonic: "INC", Length: 3, Cycles: 7, Mode: ABSOLUTE_X},
	0xe8: {Mnemonic: "INX", Length: 1, Cycles: 2, Mode: IMPLIED},
	0xc8: {Mnemonic: "INY", Length: 1, Cycles: 2, Mode: IMPLIED},

	0x4c: {Mnemonic: "JMP", Length: 3, Cycles: 3, Mode: ABSOLUTE},
	0x6c: {Mnemonic: "JMP", Length: 3, Cycles: 5, Mode: INDIRECT},

	0x20: {Mnemonic: "JSR", Length: 3, Cycles: 6, Mode: ABSOLUTE},

	0xa9: {Mnemonic: "LDA", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0xa5: {Mnemonic: "LDA", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0xb5: {Mnemonic: "LDA", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0xad: {Mnemonic: "LDA", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0xbd: {Mnemonic: "LDA", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},
	0xb9: {Mnemonic: "LDA", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_Y},
	0xa1: {Mnemonic: "LDA", Length: 2, Cycles: 6, Mode: INDIRECT_X},
	0xb1: {Mnemonic: "LDA", Length: 2, Cycles: 5 /* +1 if page crossed*/, Mode: INDIRECT_Y},

	0xa2: {Mnemonic: "LDX", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0xa6: {Mnemonic: "LDX", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0xb6: {Mnemonic: "LDX", Length: 2, Cycles: 4, Mode: ZERO_PAGE_Y},
	0xae: {Mnemonic: "LDX", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0xbe: {Mnemonic: "LDX", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_Y},

	0xa0: {Mnemonic: "LDY", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0xa4: {Mnemonic: "LDY", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0xb4: {Mnemonic: "LDY", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0xac: {Mnemonic: "LDY", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0xbc: {Mnemonic: "LDY", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},

	0x4a: {Mnemonic: "LSR", Length: 1, Cycles: 2, Mode: ACCUMULATOR},
	0x46: {Mnemonic: "LSR", Length: 2, Cycles: 5, Mode: ZERO_PAGE},
	0x56: {Mnemonic: "LSR", Length: 2, Cycles: 6, Mode: ZERO_PAGE_X},
	0x4e: {Mnemonic: "LSR", Length: 3, Cycles: 6, Mode: ABSOLUTE},
	0x5e: {Mnemonic: "LSR", Length: 3, Cycles: 7, Mode: ABSOLUTE_X},

	0xea: {Mnemonic: "NOP", Length: 1, Cycles: 2, Mode: IMPLIED},

	0x09: {Mnemonic: "ORA", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0x05: {Mnemonic: "ORA", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0x15: {Mnemonic: "ORA", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0x0d: {Mnemonic: "ORA", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0x1d: {Mnemonic: "ORA", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},
	0x19: {Mnemonic: "ORA", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_Y},
	0x01: {Mnemonic: "ORA", Length: 2, Cycles: 6, Mode: INDIRECT_X},
	0x11: {Mnemonic: "ORA", Length: 2, Cycles: 5 /* +1 if page crossed*/, Mode: INDIRECT_Y},

	0x48: {Mnemonic: "PHA", Length: 1, Cycles: 3, Mode: IMPLIED},
	0x08: {Mnemonic: "PHP", Length: 1, Cycles: 3, Mode: IMPLIED},
	0x68: {Mnemonic: "PLA", Length: 1, Cycles: 4, Mode: IMPLIED},
	0x28: {Mnemonic: "PLP", Length: 1, Cycles: 4, Mode: IMPLIED},

	0x2a: {Mnemonic: "ROL", Length: 1, Cycles: 2, Mode: ACCUMULATOR},
	0x26: {Mnemonic: "ROL", Length: 2, Cycles: 5, Mode: ZERO_PAGE},
	0x36: {Mnemonic: "ROL", Length: 2, Cycles: 6, Mode: ZERO_PAGE_X},
	0x2e: {Mnemonic: "ROL", Length: 3, Cycles: 6, Mode: ABSOLUTE},
	0x3e: {Mnemonic: "ROL", Length: 3, Cycles: 7, Mode: ABSOLUTE_X},

	0x6a: {Mnemonic: "ROR", Length: 1, Cycles: 2, Mode: ACCUMULATOR},
	0x66: {Mnemonic: "ROR", Length: 2, Cycles: 5, Mode: ZERO_PAGE},
	0x76: {Mnemonic: "ROR", Length: 2, Cycles: 6, Mode: ZERO_PAGE_X},
	0x6e: {Mnemonic: "ROR", Length: 3, Cycles: 6, Mode: ABSOLUTE},
	0x7e: {Mnemonic: "ROR", Length: 3, Cycles: 7, Mode: ABSOLUTE_X},

	0x40: {Mnemonic: "RTI", Length: 1, Cycles: 6, Mode: IMPLIED},
	0x60: {Mnemonic: "RTS", Length: 1, Cycles: 6, Mode: IMPLIED},

	0xe9: {Mnemonic: "SBC", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0xe5: {Mnemonic: "SBC", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0xf5: {Mnemonic: "SBC", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0xed: {Mnemonic: "SBC", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0xfd: {Mnemonic: "SBC", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},
	0xf9: {Mnemonic: "SBC", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_Y},
	0xe1: {Mnemonic: "SBC", Length: 2, Cycles: 6, Mode: INDIRECT_X},
	0xf1: {Mnemonic: "SBC", Length: 2, Cycles: 5 /* +1 if page crossed*/, Mode: INDIRECT_Y},

	0x38: {Mnemonic: "SEC", Length: 1, Cycles: 2, Mode: IMPLIED},
	0xf8: {Mnemonic: "SED", Length: 1, Cycles: 2, Mode: IMPLIED},
	0x78: {Mnemonic: "SEI", Length: 1, Cycles: 2, Mode: IMPLIED},

	0x85: {Mnemonic: "STA", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0x95: {Mnemonic: "STA", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0x8d: {Mnemonic: "STA", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0x9d: {Mnemonic: "STA", Length: 3, Cycles: 5, Mode: ABSOLUTE_X},
	0x99: {Mnemonic: "STA", Length: 3, Cycles: 5, Mode: ABSOLUTE_Y},
	0x81: {Mnemonic: "STA", Length: 2, Cycles: 6, Mode: INDIRECT_X},
	0x91: {Mnemonic: "STA", Length: 2, Cycles: 6, Mode: INDIRECT_Y},
	0x86: {Mnemonic: "STX", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0x96: {Mnemonic: "STX", Length: 2, Cycles: 4, Mode: ZERO_PAGE_Y},
	0x8e: {Mnemonic: "STX", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0x84: {Mnemonic: "STY", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0x94: {Mnemonic: "STY", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0x8c: {Mnemonic: "STY", Length: 3, Cycles: 4, Mode: ABSOLUTE},

	0xaa: {Mnemonic: "TAX", Length: 1, Cycles: 2, Mode: IMPLIED},
	0xa8: {Mnemonic: "TAY", Length: 1, Cycles: 2, Mode: IMPLIED},
	0xba: {Mnemonic: "TSX", Length: 1, Cycles: 2, Mode: IMPLIED},
	0x8a: {Mnemonic: "TXA", Length: 1, Cycles: 2, Mode: IMPLIED},
	0x9a: {Mnemonic: "TXS", Length: 1, Cycles: 2, Mode: IMPLIED},
	0x98: {Mnemonic: "TYA", Length: 1, Cycles: 2, Mode: IMPLIED},

	// -------------------------
	// Unofficial Opecodes
	// -------------------------

	// Combined Operations
	0xa3: {Mnemonic: "*LAX", Length: 2, Cycles: 6, Mode: INDIRECT_X},
	0xa7: {Mnemonic: "*LAX", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0xaf: {Mnemonic: "*LAX", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0xb3: {Mnemonic: "*LAX", Length: 2, Cycles: 5 /* +1 if page crossed*/, Mode: INDIRECT_Y},
	0xb7: {Mnemonic: "*LAX", Length: 2, Cycles: 4, Mode: ZERO_PAGE_Y},
	0xbf: {Mnemonic: "*LAX", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_Y},

	0x83: {Mnemonic: "*SAX", Length: 2, Cycles: 6, Mode: INDIRECT_X},
	0x87: {Mnemonic: "*SAX", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0x8f: {Mnemonic: "*SAX", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0x97: {Mnemonic: "*SAX", Length: 2, Cycles: 4, Mode: ZERO_PAGE_Y},

	// RMW instructions
	0xc3: {Mnemonic: "*DCP", Length: 2, Cycles: 8, Mode: INDIRECT_X},
	0xc7: {Mnemonic: "*DCP", Length: 2, Cycles: 5, Mode: ZERO_PAGE},
	0xcf: {Mnemonic: "*DCP", Length: 3, Cycles: 6, Mode: ABSOLUTE},
	0xd3: {Mnemonic: "*DCP", Length: 2, Cycles: 8, Mode: INDIRECT_Y},
	0xd7: {Mnemonic: "*DCP", Length: 2, Cycles: 6, Mode: ZERO_PAGE_X},
	0xdb: {Mnemonic: "*DCP", Length: 3, Cycles: 7, Mode: ABSOLUTE_Y},
	0xdf: {Mnemonic: "*DCP", Length: 3, Cycles: 7, Mode: ABSOLUTE_X},

	// Duplicated instructions
	0xeb: {Mnemonic: "*SBC", Length: 2, Cycles: 2, Mode: IMMEDIATE},

	// NOPs
	0x1a: {Mnemonic: "*NOP", Length: 1, Cycles: 2, Mode: IMPLIED},
	0x3a: {Mnemonic: "*NOP", Length: 1, Cycles: 2, Mode: IMPLIED},
	0x5a: {Mnemonic: "*NOP", Length: 1, Cycles: 2, Mode: IMPLIED},
	0x7a: {Mnemonic: "*NOP", Length: 1, Cycles: 2, Mode: IMPLIED},
	0xda: {Mnemonic: "*NOP", Length: 1, Cycles: 2, Mode: IMPLIED},
	0xfa: {Mnemonic: "*NOP", Length: 1, Cycles: 2, Mode: IMPLIED},
	0x80: {Mnemonic: "*NOP", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0x82: {Mnemonic: "*NOP", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0x89: {Mnemonic: "*NOP", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0xc2: {Mnemonic: "*NOP", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0xe2: {Mnemonic: "*NOP", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0x0c: {Mnemonic: "*NOP", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0x1c: {Mnemonic: "*NOP", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},
	0x3c: {Mnemonic: "*NOP", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},
	0x5c: {Mnemonic: "*NOP", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},
	0x7c: {Mnemonic: "*NOP", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},
	0xdc: {Mnemonic: "*NOP", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},
	0xfc: {Mnemonic: "*NOP", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},
	0x04: {Mnemonic: "*NOP", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0x44: {Mnemonic: "*NOP", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0x64: {Mnemonic: "*NOP", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0x14: {Mnemonic: "*NOP", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0x34: {Mnemonic: "*NOP", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0x54: {Mnemonic: "*NOP", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0x74: {Mnemonic: "*NOP", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0xd4: {Mnemonic: "*NOP", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0xf4: {Mnemonic: "*NOP", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
}
