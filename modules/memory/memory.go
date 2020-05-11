package memory

import (
	"fmt"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

type memory struct {
	block *protocol.Block
	h     *handler.Handler
}

func Memory(block *protocol.Block, h *handler.Handler) {
	block.Name = "memory"
	block.Separator = false
	block.SeparatorBlockWidth = 15

	update := func() {
		mem := h.SysInfo().Memory()
		//block.FullText = fmt.Sprintf(" %.2f/%.2f (%.1f%%)",
		//	mem[0], mem[1], mem[2])
		block.FullText = fmt.Sprintf(" %.1f/%.1f",
			//mem[1]-mem[0], mem[1])
			mem[0], mem[1])
	}

	//h.Scheduler().Register(time.Second*20, update)
	h.Scheduler().Register(time.Second/2, update)
}
