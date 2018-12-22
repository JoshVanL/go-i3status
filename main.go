package main

import (
	"fmt"
	"os"

	"github.com/joshvanl/go-i3status/blocks"
	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

var (
	enabledBlocks = []func(*protocol.Block, *handler.Handler){
		blocks.Battery,
		blocks.Time,
	}
)

func main() {
	h, err := handler.New()
	if err != nil {
		fmt.Printf("error creating handler: %s", err)
		os.Exit(1)
	}

	for _, f := range enabledBlocks {
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
