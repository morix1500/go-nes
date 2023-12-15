package main

import (
	"fmt"
	"go-nes/nes"
	"log"
	"log/slog"
	"os"
)

func main() {
	filepath := os.Args[1]
	if filepath == "" {
		log.Fatal("Please specify a file path")
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	c, err := nes.NewCartridge(data)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info(fmt.Sprintf("Program Rom Length: %d", len(c.ProgramRom)))
	slog.Info(fmt.Sprintf("Charactor Rom Length: %d", len(c.CharacterRom)))

	b := nes.NewBus(c)
	cpu := nes.NewCPU(b)
	cpu.Reset()
	cpu.Run()
}
