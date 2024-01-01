package nes

const (
	JOYPAD_BUTTON_A uint8 = 0b0000_0001
	JOYPAD_BUTTON_B uint8 = 0b0000_0010
	JOYPAD_SELECT   uint8 = 0b0000_0100
	JOYPAD_START    uint8 = 0b0000_1000
	JOYPAD_UP       uint8 = 0b0001_0000
	JOYPAD_DOWN     uint8 = 0b0010_0000
	JOYPAD_LEFT     uint8 = 0b0100_0000
	JOYPAD_RIGHT    uint8 = 0b1000_0000
)

type Joypad struct {
	Strobe       bool
	ButtonIndex  uint8
	ButtonStatus uint8
}

func NewJoypad() *Joypad {
	return &Joypad{
		Strobe:       false,
		ButtonIndex:  0,
		ButtonStatus: 0,
	}
}

func (j *Joypad) Read() uint8 {
	if j.ButtonIndex > 7 {
		return 1
	}
	response := (j.ButtonStatus & (1 << j.ButtonIndex)) >> j.ButtonIndex
	if !j.Strobe && j.ButtonIndex <= 7 {
		j.ButtonIndex++
	}
	return response
}

func (j *Joypad) Write(data uint8) {
	j.Strobe = data&1 == 1
	if j.Strobe {
		j.ButtonIndex = 0
	}
}

func (j *Joypad) Press(button uint8) {
	j.ButtonStatus |= button
}

func (j *Joypad) Release(button uint8) {
	j.ButtonStatus &= ^button
}
