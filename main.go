package main

import (
	"fmt"
	"os"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/modules/battery"
	"github.com/joshvanl/go-i3status/modules/cpu"
	"github.com/joshvanl/go-i3status/modules/date"
	"github.com/joshvanl/go-i3status/modules/time"
	"github.com/joshvanl/go-i3status/protocol"
)

var (
	enabledBlocks = []func(*protocol.Block, *handler.Handler){
		cpu.CPU,
		battery.Battery,
		date.Date,
		time.Time,
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
			Align:               protocol.Right,
			Color:               "#dddddddd",
			Background:          "#2c2c2ccc",
			Border:              "#2c2c2ccc",
			SeperatorBlockWidth: 30,
		}
		h.RegisterBlock(b)
		go f(b, h)
	}

	h.Tick()

	select {}
}
