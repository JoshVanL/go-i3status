package cpu

import (
	"fmt"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

type cpu struct {
	block *protocol.Block
	h     *handler.Handler
}

func CPU(block *protocol.Block, h *handler.Handler) {
	block.Name = "cpu"

	c := &cpu{
		block: block,
		h:     h,
	}
	//update := func() {
	//}

	//go func() {
	//	for {
	//		c.update()
	//		time.Sleep(time.Second * 4)
	//		h.Tick()
	//	}
	//}()

	h.Scheduler().Register(time.Second*2, c.update)
}

func (c *cpu) update() {
	load := c.h.SysInfo().CPULoad()

	if load == -1 {
		c.block.FullText = "cpu 0.00%"
	} else {
		c.block.FullText = fmt.Sprintf("cpu %.2f%%", load)
	}
}
