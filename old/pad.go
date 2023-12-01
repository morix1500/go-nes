package main

type Button uint8

const (
	ButtonA      Button = 1
	ButtonB      Button = 2
	ButtonSELECT Button = 4
	ButtonSTART  Button = 8
	ButtonUP     Button = 16
	ButtonDOWN   Button = 32
	ButtonLEFT   Button = 64
	ButtonRIGHT  Button = 128
)

type Pad struct {
	status      uint8
	readIndex   uint8
	stobeEnable bool
}

func (p *Pad) ReadButton(b Button) {
	switch b {
	case ButtonA:
		p.status |= uint8(ButtonA)
	case ButtonB:
		p.status |= uint8(ButtonB)
	case ButtonSELECT:
		p.status |= uint8(ButtonSELECT)
	case ButtonSTART:
		p.status |= uint8(ButtonSTART)
	case ButtonUP:
		p.status |= uint8(ButtonUP)
	case ButtonDOWN:
		p.status |= uint8(ButtonDOWN)
	case ButtonLEFT:
		p.status |= uint8(ButtonLEFT)
	case ButtonRIGHT:
		p.status |= uint8(ButtonRIGHT)
	}
}

func (p *Pad) ReleaseButton(b Button) {
	switch b {
	case ButtonA:
		p.status &= ^uint8(ButtonA)
	case ButtonB:
		p.status &= ^uint8(ButtonB)
	case ButtonSELECT:
		p.status &= ^uint8(ButtonSELECT)
	case ButtonSTART:
		p.status &= ^uint8(ButtonSTART)
	case ButtonUP:
		p.status &= ^uint8(ButtonUP)
	case ButtonDOWN:
		p.status &= ^uint8(ButtonDOWN)
	case ButtonLEFT:
		p.status &= ^uint8(ButtonLEFT)
	case ButtonRIGHT:
		p.status &= ^uint8(ButtonRIGHT)
	}
}

func (p *Pad) ReadPad() uint8 {
	data := (p.status >> p.readIndex) & 1
	if !p.stobeEnable {
		p.readIndex++
		p.readIndex %= 8
	}
	return data
}
