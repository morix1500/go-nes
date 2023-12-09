package nes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCartridge struct {
	header       []uint8
	trainer      []uint8
	programRom   []uint8
	characterRom []uint8
}

func createTestCartridge(cartridge TestCartridge) []uint8 {
	var raw []uint8
	raw = append(raw, cartridge.header...)
	raw = append(raw, cartridge.trainer...)
	raw = append(raw, cartridge.programRom...)
	raw = append(raw, cartridge.characterRom...)
	return raw
}

func createDummyRom(val uint8, size int) []uint8 {
	var rom []uint8
	for i := 0; i < size; i++ {
		rom = append(rom, val)
	}
	return rom
}

func TestCartridgeInitialize(t *testing.T) {
	dummyProgramRom := createDummyRom(1, 2*PROGRAM_ROM_PAGE_SIZE)
	dummyCharacterRom := createDummyRom(2, 1*CHARACTER_ROM_PAGE_SIZE)
	testRom := createTestCartridge(TestCartridge{
		header: []uint8{
			0x4E, 0x45, 0x53, 0x1A, 0x02, 0x01, 0x31, 00, 00, 00, 00, 00, 00, 00, 00, 00,
		},
		trainer:      nil,
		programRom:   dummyProgramRom,
		characterRom: dummyCharacterRom,
	})

	cartridge, err := NewCartridge(testRom)
	assert.NoError(t, err)
	assert.Equal(t, dummyProgramRom, cartridge.ProgramRom)
	assert.Equal(t, dummyCharacterRom, cartridge.CharacterRom)
	assert.Equal(t, uint8(3), cartridge.Mapper)
	assert.Equal(t, MIRROR_VERTICAL, cartridge.ScreenMirroring)
}

func TestCartridgeWithTrainer(t *testing.T) {
	dummyProgramRom := createDummyRom(1, 2*PROGRAM_ROM_PAGE_SIZE)
	dummyCharacterRom := createDummyRom(2, 1*CHARACTER_ROM_PAGE_SIZE)
	testRom := createTestCartridge(TestCartridge{
		header: []uint8{
			0x4E,
			0x45,
			0x53,
			0x1A,
			0x02,
			0x01,
			0x31 | 0b100,
			00,
			00,
			00,
			00,
			00,
			00,
			00,
			00,
			00,
		},
		trainer:      createDummyRom(0, 512),
		programRom:   dummyProgramRom,
		characterRom: dummyCharacterRom,
	})

	cartridge, err := NewCartridge(testRom)
	assert.NoError(t, err)
	assert.Equal(t, dummyProgramRom, cartridge.ProgramRom)
	assert.Equal(t, dummyCharacterRom, cartridge.CharacterRom)
	assert.Equal(t, uint8(3), cartridge.Mapper)
	assert.Equal(t, MIRROR_VERTICAL, cartridge.ScreenMirroring)
}

func TestCartridgeNES2NotSupported(t *testing.T) {
	testRom := createTestCartridge(TestCartridge{
		header:       []uint8{0x4E, 0x45, 0x53, 0x1A, 0x01, 0x01, 0x31, 0x8, 00, 00, 00, 00, 00, 00, 00, 00},
		trainer:      nil,
		programRom:   nil,
		characterRom: nil,
	})

	_, err := NewCartridge(testRom)
	assert.Error(t, err)
	assert.Equal(t, "NES2.0 format is not supported", err.Error())
}
