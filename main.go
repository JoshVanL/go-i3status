package main

import (
	"github.com/joshvanl/go-i3status/blocks"
	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func main() {
	//blocks := []*protocol.Block{
	//	&protocol.Block{
	//		Name:       "echo",
	//		Instance:   "echoIn",
	//		FullText:   "hello",
	//		Align:      protocol.Right,
	//		Color:      "#dddddd",
	//		Background: "#2c2c2ccc",
	//		Border:     "#2c2c2ccc",
	//	},
	//}

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
