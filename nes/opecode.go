package nes

type OpeCode struct {
	Code     uint8
	Mnemonic string
	Length   uint8
	Cycles   uint8
	Mode     AddressingMode
}

var CPU_OPS_CODES map[uint8]OpeCode = map[uint8]OpeCode{
	0x00: {Mnemonic: "BRK", Length: 1, Cycles: 7, Mode: NoneAddressing},
	0xaa: {Mnemonic: "TAX", Length: 1, Cycles: 2, Mode: NoneAddressing},
	0xe8: {Mnemonic: "INX", Length: 1, Cycles: 2, Mode: NoneAddressing},

	0xa9: {Mnemonic: "LDA", Length: 2, Cycles: 2, Mode: Immediate},
	0xa5: {Mnemonic: "LDA", Length: 2, Cycles: 3, Mode: ZeroPage},
	0xb5: {Mnemonic: "LDA", Length: 2, Cycles: 4, Mode: ZeroPageX},
	0xad: {Mnemonic: "LDA", Length: 3, Cycles: 4, Mode: Absolute},
	0xbd: {Mnemonic: "LDA", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: AbsoluteX},
	0xb9: {Mnemonic: "LDA", Length: 3, Cycles: 4 /* +1 if page crossed*/, Mode: AbsoluteY},
	0xa1: {Mnemonic: "LDA", Length: 2, Cycles: 6, Mode: IndirectX},
	0xb1: {Mnemonic: "LDA", Length: 2, Cycles: 5 /* +1 if page crossed*/, Mode: IndirectY},

	0x85: {Mnemonic: "STA", Length: 2, Cycles: 3, Mode: ZeroPage},
	0x95: {Mnemonic: "STA", Length: 2, Cycles: 4, Mode: ZeroPageX},
	0x8d: {Mnemonic: "STA", Length: 3, Cycles: 4, Mode: Absolute},
	0x9d: {Mnemonic: "STA", Length: 3, Cycles: 5, Mode: AbsoluteX},
	0x99: {Mnemonic: "STA", Length: 3, Cycles: 5, Mode: AbsoluteY},
	0x81: {Mnemonic: "STA", Length: 2, Cycles: 6, Mode: IndirectX},
	0x91: {Mnemonic: "STA", Length: 2, Cycles: 6, Mode: IndirectY},
}
