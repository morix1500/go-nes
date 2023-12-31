package main

import (
	"fmt"
	"go-nes/nes"
	"go-nes/ui"
	"log"
	"log/slog"
	"os"
	"runtime"
	"unsafe"

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

	//window := ui.InitGlfw()
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

	frame := ui.NewFrame()
	//program := ui.InitOpenGL()

	vertices := []float32{
		// top left
		-1.0, 1.0, 0.0, // position
		1.0, 0.0, 0.0, // Color
		0.0, 0.0, // texture coordinates

		// top right
		1.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		1.0, 0.0,

		// bottom right
		1.0, -1.0, 0.0,
		0.0, 0.0, 1.0,
		1.0, 1.0,

		// bottom left
		-1.0, -1.0, 0.0,
		1.0, 1.0, 1.0,
		0.0, 1.0,
	}

	indices := []uint32{
		// rectangle
		0, 1, 2, // top triangle
		0, 2, 3, // bottom triangle
	}
	VAO := ui.CreateVAO(vertices, indices)

	b := nes.NewBus(c, nil)
	cpu := nes.NewCPU(b)
	cpu.Reset()

	for !window.ShouldClose() {
		cpu.Step()
		if b.RenderFlag {
			glfw.PollEvents()
			gl.Clear(gl.COLOR_BUFFER_BIT)
			ui.Render(b.PPU, frame)

			shaderProgram.Use()

			tex, err := ui.NewTexture(frame.Front, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)
			if err != nil {
				panic(err)
			}

			tex.Bind(gl.TEXTURE0)
			tex.SetUniform(shaderProgram.GetUniformLocation("ourTexture"))

			gl.BindVertexArray(VAO)
			gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, unsafe.Pointer(nil))
			gl.BindVertexArray(0)

			tex.UnBind()

			b.RenderFlag = false
			window.SwapBuffers()
		}
	}
}
