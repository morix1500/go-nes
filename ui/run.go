package ui

import (
	"go-nes/nes"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	vertices = []float32{
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

	indices = []uint32{
		// rectangle
		0, 1, 2, // top triangle
		0, 2, 3, // bottom triangle
	}
)

func Run(cpu *nes.CPU, bus *nes.Bus, window *glfw.Window, program *Program) error {
	frame := NewFrame(bus.JoyPad1, bus.JoyPad2)
	VAO := CreateVAO()

	window.SetKeyCallback(frame.OnKey)

	for !window.ShouldClose() {
		cpu.Step()
		if bus.RenderFlag {
			glfw.PollEvents()
			gl.Clear(gl.COLOR_BUFFER_BIT)
			frame.Render(bus.PPU)

			program.Use()

			tex, err := NewTexture(frame.Front, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)
			if err != nil {
				return err
			}

			tex.Bind(gl.TEXTURE0)
			//tex.SetUniform(shaderProgram.GetUniformLocation("ourTexture"))

			gl.BindVertexArray(VAO)
			gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, unsafe.Pointer(nil))
			gl.BindVertexArray(0)

			tex.UnBind()

			bus.RenderFlag = false
			window.SwapBuffers()
		}
	}

	return nil
}

func CreateVAO() uint32 {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	var EBO uint32
	gl.GenBuffers(1, &EBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// copy indices into element buffer
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// size of one whole vertex (sum of attrib sizes)
	var stride int32 = 3*4 + 3*4 + 2*4
	var offset int = 0

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(0)
	offset += 3 * 4

	// color
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(1)
	offset += 3 * 4

	// texture position
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(2)
	offset += 2 * 4

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	return VAO
}
