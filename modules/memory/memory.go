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

	for {
		mem, per := h.SysInfo().Memory()
		block.FullText = fmt.Sprintf(" %v %v (%d)",
			mem[0], mem[1], per)

		h.Tick()

		<-ticker
	}
}
