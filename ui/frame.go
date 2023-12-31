package ui

import (
	"fmt"
	"go-nes/nes"
	"image"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	WIDTH  = 256
	HEIGHT = 240
	SCALE  = 3

	VERTEX_SHADER_SOURCE = `
	#version 410
	in vec3 vp;
	in vec4 color;
	out vec4 vertex_color;
	void main() {
	  vertex_color = color;
	  gl_Position = vec4(vp, 1.0);
	}
  ` + "\x00"
	FRAGMENT_SHADER_SOURCE = `
	#version 410
	in vec4 vertex_color;
	out vec4 frag_colour;
	void main() {
	  frag_colour = vertex_color;
	}
  ` + "\x00"
)

var (
	squareArrays = [][]float32{
		{-0.5, 0.5, 0},
		{-0.5, -0.5, 0},
		{0.5, -0.5, 0},
		{-0.5, 0.5, 0},
		{0.5, 0.5, 0},
		{0.5, -0.5, 0},
	}
)

type pixel struct {
	drawable uint32
}

func (p *pixel) draw() {
	gl.BindVertexArray(p.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(squareArrays)))
}

type Frame struct {
	Pixels  [][]*pixel
	Front   *image.RGBA
	Texture uint32
}

func NewFrame() *Frame {
	return &Frame{
		Pixels:  make([][]*pixel, HEIGHT*8),
		Front:   image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT)),
		Texture: createTexture(),
	}
}

func (f *Frame) SetPixel(ind, x, y uint, r, g, b uint8) {
	points := make([]float32, 36)
	index := 0

	for i := 0; i < len(squareArrays); i++ {
		for j := 0; j < len(squareArrays[i]); j++ {
			var position float32
			var size float32
			switch j {
			case 0:
				size = 1.0 / float32(WIDTH)
				position = float32(x) * size
			case 1:
				size = 1.0 / float32(HEIGHT)
				position = float32(HEIGHT-1-y) * size
			default:
				index++
				continue
			}

			if squareArrays[i][j] < 0 {
				points[index] = (position * 2) - 1
			} else {
				points[index] = ((position + size) * 2) - 1
			}
			index++
		}
		col := []float32{float32(r) / 255, float32(g) / 255, float32(b) / 255}
		copy(points[index:], col)
		index += 3
	}

	pixel := &pixel{
		drawable: makeVao(points),
	}
	f.Pixels[ind] = append(f.Pixels[ind], pixel)
}

func ShowTile(ind uint, characterRom []uint8, bank uint, tileN uint) *Frame {
	frame := NewFrame()
	bank = (bank * 0x1000)

	tile := characterRom[bank+tileN*16 : bank+tileN*16+15+1]

	for y := uint(0); y < 8; y++ {
		upper := tile[y]
		lower := tile[y+8]

		for x := uint(0); x < 8; x++ {
			value := (1&upper)<<1 | (1 & lower)
			upper = upper >> 1
			lower = lower >> 1
			var r, g, b uint32
			switch value {
			case 0:
				r, g, b, _ = nes.Palletes[0x01].RGBA()
			case 1:
				r, g, b, _ = nes.Palletes[0x23].RGBA()
			case 2:
				r, g, b, _ = nes.Palletes[0x27].RGBA()
			case 3:
				r, g, b, _ = nes.Palletes[0x30].RGBA()
			default:
				panic("unknown value")
			}
			tmpx := 7 - x
			frame.SetPixel(ind, uint(tmpx), uint(y), uint8(r), uint8(g), uint8(b))
		}
	}

	return frame
}

func ShowTileBank(characterRom []uint8, bank uint) *Frame {
	frame := NewFrame()
	bank = (bank * 0x1000)

	tileX := uint(0)
	tileY := uint(0)
	ind := uint(0)

	for tileN := uint(0); tileN < 255; tileN++ {
		if tileN%20 == 0 {
			tileX = 0
			tileY += 10
		}
		tile := characterRom[bank+tileN*16 : bank+tileN*16+15+1]

		for y := uint(0); y < 8; y++ {
			upper := tile[y]
			lower := tile[y+8]

			for x := uint(0); x < 8; x++ {
				value := (1&upper)<<1 | (1 & lower)
				upper = upper >> 1
				lower = lower >> 1
				var r, g, b uint32
				switch value {
				case 0:
					r, g, b, _ = nes.Palletes[0x01].RGBA()
				case 1:
					r, g, b, _ = nes.Palletes[0x23].RGBA()
				case 2:
					r, g, b, _ = nes.Palletes[0x27].RGBA()
				case 3:
					r, g, b, _ = nes.Palletes[0x30].RGBA()
				default:
					panic("unknown value")
				}
				tmpx := 7 - x + (tileX * 8) + (tileX * 2)
				tmpy := y + tileY
				frame.SetPixel(ind, uint(tmpx), uint(tmpy), uint8(r), uint8(g), uint8(b))
			}
		}
		ind++
		tileX++
	}

	return frame
}

//func View(characterRom []uint8) {
//	runtime.LockOSThread()
//
//	window := InitGlfw()
//	defer glfw.Terminate()
//
//	program := InitOpenGL()
//	//frame := ShowTile(0, characterRom, 1, 0)
//	frame := ShowTileBank(characterRom, 1)
//
//	for !window.ShouldClose() {
//		Draw(frame.Pixels, window, program)
//	}
//}

func Draw(cells [][]*pixel, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	for x := range cells {
		for _, c := range cells[x] {
			c.draw()
		}
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

func Init() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(WIDTH*SCALE, HEIGHT*SCALE, "Morix NES", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetKeyCallback(onKey)

	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
	gl.Enable(gl.TEXTURE_2D)

	return window
}

func InitGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(WIDTH*SCALE, HEIGHT*SCALE, "Morix NES", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetKeyCallback(onKey)

	return window
}

func InitOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(VERTEX_SHADER_SOURCE, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(FRAGMENT_SHADER_SOURCE, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))

	return vao
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyEscape:
			w.SetShouldClose(true)
			os.Exit(0)
		}
	}
}

func createTexture() uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}

func SetTexture(im *image.RGBA) {
	size := im.Rect.Size()
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(im.Pix))
}

func (f *Frame) renderPixel(x, y uint, c color.RGBA) {
	f.Front.SetRGBA(int(x), int(y), c)
}

func DrawBuffer(window *glfw.Window) {
	w, h := window.GetFramebufferSize()
	s1 := float32(w) / 256
	s2 := float32(h) / 240
	f := float32(1 - 0)
	var x, y float32
	if s1 >= s2 {
		x = f * s2 / s1
		y = f
	} else {
		x = f
		y = f * s1 / s2
	}
	vertices := []float32{
		-x, -y, 0, 1, // Position and texture coordinates for each vertex
		x, -y, 1, 1,
		x, y, 1, 0,
		-x, y, 0, 0,
	}

	// Create a VBO and bind it
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	// Create a VAO and bind it
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Specify the layout of the vertex data
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	// Draw the quad
	gl.DrawArrays(gl.QUADS, 0, 4)

	// Cleanup
	//gl.DeleteBuffers(1, &vbo)
	//gl.DeleteVertexArrays(1, &vao)
}

func CreateVAO(vertices []float32, indices []uint32) uint32 {

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
