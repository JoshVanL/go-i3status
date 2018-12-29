package memory

import (
	"fmt"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func Memory(block *protocol.Block, h *handler.Handler) {
	ticker := time.NewTicker(time.Second * 20).C

	block.Name = "memory"
	block.Separator = false
	block.SeparatorBlockWidth = 15

	for {
		mem := h.SysInfo().Memory()
		block.FullText = fmt.Sprintf(" %.2f/%.2f (%.1f%%)",
			mem[0], mem[1], mem[2])

		h.Tick()

		<-ticker
	}
}
