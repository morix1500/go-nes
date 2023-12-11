package main

import (
	"go-nes/nes"
	"log"
	"os"
)

func main() {
	filepath := "test.nes"
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	c, err := nes.NewCartridge(data)
	if err != nil {
		log.Fatal(err)
	}

	b := nes.NewBus(c)
	cpu := nes.NewCPU(b)
	cpu.Reset()
	cpu.Run()
}
