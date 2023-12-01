package main

import (
	"fmt"
	"os"
)

const (
	NES_HEADER_LENGTH               = 16
	NES_HEADER_PROGRAM_ROM_OFFSET   = 4
	NES_HEADER_CHARACTER_ROM_OFFSET = 5
	NES_HEADER_MAGIC_NUMBER         = "NES\x1a"
	PROGRAM_ROM_PAGE_SIZE           = 16384
	CHARACTER_ROM_SIZE              = 8192
)

type Cassette struct {
	characterROM []byte
	programROM   []byte
}

func NewCassette() *Cassette {
	return &Cassette{}
}

func (c *Cassette) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("ファイルの読み込みエラー:", err)
		return err
	}
	if len(data) < NES_HEADER_LENGTH {
		return fmt.Errorf("ファイルサイズが不正です")
	}
	if string(data[:4]) != NES_HEADER_MAGIC_NUMBER {
		fmt.Println("NESファイルではありません")
		return fmt.Errorf("NESファイルではありません")
	}

	characterROMPages := data[NES_HEADER_CHARACTER_ROM_OFFSET]
	characterROMStart := NES_HEADER_LENGTH + int(data[NES_HEADER_PROGRAM_ROM_OFFSET])*PROGRAM_ROM_PAGE_SIZE
	characterROMEnd := characterROMStart + int(characterROMPages)*CHARACTER_ROM_SIZE

	c.programROM = data[NES_HEADER_LENGTH : characterROMStart-1]
	c.characterROM = data[characterROMStart : characterROMEnd-1]

	return nil
}
