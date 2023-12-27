package main

import (
	"fmt"
	"go-nes/nes"
	"go-nes/ui"
	"log"
	"log/slog"
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
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

	frame := ui.NewFrame()

	//ui.View(c.CharacterRom)

	runtime.LockOSThread()

	window := ui.InitGlfw()
	defer glfw.Terminate()

	program := ui.InitOpenGL()

	b := nes.NewBus(c, func(p *nes.PPU) {
		ui.Render(p, frame)
		//for !window.ShouldClose() {
		ui.Draw(frame.Pixels, window, program)
		//}
	})
	//b := nes.NewBus(c, nil)
	cpu := nes.NewCPU(b)
	cpu.Reset()
	cpu.Run()
}
