package nes

import "fmt"

type Mirroring uint8

const (
	MIRROR_VERTICAL Mirroring = iota
	MIRROR_HORIZONTAL
	MIRROR_FOUR_SCREEN
)

const (
	NES_TAG                 = "NES\x1a"
	PROGRAM_ROM_PAGE_SIZE   = 0x4000 // 16KB
	CHARACTER_ROM_PAGE_SIZE = 0x2000 // 8KB
)

type Cartridge struct {
	ProgramRom      []uint8
	CharacterRom    []uint8
	Mapper          uint8
	ScreenMirroring Mirroring
}

func NewCartridge(raw []uint8) (*Cartridge, error) {
	if len(raw) < 16 {
		return nil, fmt.Errorf("invalid ROM file")
	}
	if string(raw[:4]) != NES_TAG {
		return nil, fmt.Errorf("invalid ROM file")
	}
	header := raw[:16]

	mapper := (header[7] & 0xf0) | (header[6] >> 4)
	inesVersion := header[7] >> 2 & 0b11
	if inesVersion != 0 {
		return nil, fmt.Errorf("NES2.0 format is not supported")
	}

	fourScreen := header[6]&0b1000 != 0
	verticalMirroring := header[6]&0b1 != 0
	var screenMirroring Mirroring

	if fourScreen {
		screenMirroring = MIRROR_FOUR_SCREEN
	} else if verticalMirroring {
		screenMirroring = MIRROR_VERTICAL
	} else {
		screenMirroring = MIRROR_HORIZONTAL
	}

	prgSize := int(header[4]) * PROGRAM_ROM_PAGE_SIZE
	charSize := int(header[5]) * CHARACTER_ROM_PAGE_SIZE
	skipTrainer := header[6]&0b100 != 0

	var prgRomStart uint16
	if skipTrainer {
		prgRomStart = 16 + 512
	} else {
		prgRomStart = 16
	}
	charRomStart := prgRomStart + uint16(prgSize)

	return &Cartridge{
		ProgramRom:      raw[prgRomStart : prgRomStart+uint16(prgSize)],
		CharacterRom:    raw[charRomStart : charRomStart+uint16(charSize)],
		Mapper:          mapper,
		ScreenMirroring: screenMirroring,
	}, nil
}
