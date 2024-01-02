package nes

import "github.com/go-gl/glfw/v3.2/glfw"

const (
	JOYPAD_A = iota
	JOYPAD_B
	JOYPAD_SELECT
	JOYPAD_START
	JOYPAD_UP
	JOYPAD_DOWN
	JOYPAD_LEFT
	JOYPAD_RIGHT
)

type Joypad struct {
	Strobe       bool
	ButtonIndex  uint8
	ButtonStatus [8]bool
}

func NewJoypad() *Joypad {
	return &Joypad{
		Strobe:       false,
		ButtonIndex:  0,
		ButtonStatus: [8]bool{false, false, false, false, false, false, false, false},
	}
}

func (j *Joypad) Read() uint8 {
	v := uint8(0)
	if j.ButtonIndex <= 7 && j.ButtonStatus[j.ButtonIndex] {
		v = 1
	}
	j.ButtonIndex++
	if j.Strobe {
		j.ButtonIndex = 0
	}
	return v
}

func (j *Joypad) Write(data uint8) {
	j.Strobe = data&1 == 1
	if j.Strobe {
		j.ButtonIndex = 0
	}
}

func (j *Joypad) ReadKeys(window *glfw.Window) {
	j.ButtonStatus[JOYPAD_A] = window.GetKey(glfw.KeyA) == glfw.Press
	j.ButtonStatus[JOYPAD_B] = window.GetKey(glfw.KeyS) == glfw.Press
	j.ButtonStatus[JOYPAD_SELECT] = window.GetKey(glfw.KeySpace) == glfw.Press
	j.ButtonStatus[JOYPAD_START] = window.GetKey(glfw.KeyEnter) == glfw.Press
	j.ButtonStatus[JOYPAD_UP] = window.GetKey(glfw.KeyUp) == glfw.Press
	j.ButtonStatus[JOYPAD_DOWN] = window.GetKey(glfw.KeyDown) == glfw.Press
	j.ButtonStatus[JOYPAD_LEFT] = window.GetKey(glfw.KeyLeft) == glfw.Press
	j.ButtonStatus[JOYPAD_RIGHT] = window.GetKey(glfw.KeyRight) == glfw.Press
}
