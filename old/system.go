package main

import "log/slog"

const (
	WRAM_SIZE       = 2048
	APU_IO_REG_SIZE = 32
)

type System struct {
	Wram     []byte
	IOReg    []byte
	Cassette *Cassette
	Pads     [2]*Pad
}

func NewSystem(cassettePath string) (*System, error) {
	c := NewCassette()
	if err := c.Load(cassettePath); err != nil {
		slog.Error("Failed to load cassette.", "path", cassettePath, "err", err)
		return nil, err
	}

	return &System{
		Wram:     make([]byte, WRAM_SIZE),
		IOReg:    make([]byte, APU_IO_REG_SIZE),
		Cassette: c,
	}, nil
}
