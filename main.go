package main

import (
	"fmt"
	"go-nes/nes"
	"go-nes/ui"
	"log"
	"log/slog"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
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

	runtime.LockOSThread()

	window := ui.Init()
	defer glfw.Terminate()

	vertShader, err := ui.NewShaderFromFile("ui/shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragShader, err := ui.NewShaderFromFile("ui/shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	shaderProgram, err := ui.NewProgram(vertShader, fragShader)
	if err != nil {
		panic(err)
	}
	defer shaderProgram.Delete()

	b := nes.NewBus(c, nil)
	cpu := nes.NewCPU(b)
	cpu.Reset()

	if err := ui.Run(cpu, b, window, shaderProgram); err != nil {
		panic(err)
	}
}
