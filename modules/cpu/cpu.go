package cpu

import (
	"fmt"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func CPU(block *protocol.Block, h *handler.Handler) {
	block.Name = "cpu"

	update := func() {
		load := h.SysInfo().CPULoad()

		if load == -1 {
			block.FullText = "cpu 0.00%"
		} else {
			block.FullText = fmt.Sprintf("cpu %.2f%%", load)
		}
	}

	h.Scheduler().Register(time.Second, update)
}
