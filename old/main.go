package main

import (
	"fmt"
	"log/slog"
)

func main() {
	path := "./sample1/sample1.nes"
	c := NewCassette()
	if err := c.Load(path); err != nil {
		slog.Error("Failed to load cassette.", "path", path, "err", err)
		return
	}

	fmt.Println(c.characterROM)
}
