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
	0x00: {Mnemonic: "BRK", Length: 1, Cycles: 7, Mode: NONE_ADDRESSING},

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

	0xe8: {Mnemonic: "INX", Length: 1, Cycles: 2, Mode: IMPLIED},
	0xc8: {Mnemonic: "INY", Length: 1, Cycles: 2, Mode: IMPLIED},

	0xa9: {Mnemonic: "LDA", Length: 2, Cycles: 2, Mode: IMMEDIATE},
	0xa5: {Mnemonic: "LDA", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0xb5: {Mnemonic: "LDA", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0xad: {Mnemonic: "LDA", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0xbd: {Mnemonic: "LDA", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_X},
	0xb9: {Mnemonic: "LDA", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: ABSOLUTE_Y},
	0xa1: {Mnemonic: "LDA", Length: 2, Cycles: 6, Mode: INDIRECT_X},
	0xb1: {Mnemonic: "LDA", Length: 2, Cycles: 5 /* +1 if page crossed*/, Mode: INDIRECT_Y},

	0xf8: {Mnemonic: "SED", Length: 1, Cycles: 2, Mode: IMPLIED},
	0x78: {Mnemonic: "SEI", Length: 1, Cycles: 2, Mode: IMPLIED},

	0x85: {Mnemonic: "STA", Length: 2, Cycles: 3, Mode: ZERO_PAGE},
	0x95: {Mnemonic: "STA", Length: 2, Cycles: 4, Mode: ZERO_PAGE_X},
	0x8d: {Mnemonic: "STA", Length: 3, Cycles: 4, Mode: ABSOLUTE},
	0x9d: {Mnemonic: "STA", Length: 3, Cycles: 5, Mode: ABSOLUTE_X},
	0x99: {Mnemonic: "STA", Length: 3, Cycles: 5, Mode: ABSOLUTE_Y},
	0x81: {Mnemonic: "STA", Length: 2, Cycles: 6, Mode: INDIRECT_X},
	0x91: {Mnemonic: "STA", Length: 2, Cycles: 6, Mode: INDIRECT_Y},

	0xaa: {Mnemonic: "TAX", Length: 1, Cycles: 2, Mode: IMPLIED},
	0xa8: {Mnemonic: "TAY", Length: 1, Cycles: 2, Mode: IMPLIED},
}
