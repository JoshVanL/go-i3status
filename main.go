package main

import (
	"github.com/joshvanl/go-i3status/blocks"
	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func main() {
	h := handler.New()

	for _, f := range blocks.All() {
		b := &protocol.Block{
			Align:      protocol.Right,
			Color:      "#dddddd",
			Background: "#2c2c2ccc",
			Border:     "#2c2c2ccc",
		}
		h.RegisterBlock(b)
		go f(b, h)
	}

	h.Tick()

	select {}
}
