package cpu

import (
	"fmt"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func CPU(block *protocol.Block, h *handler.Handler) {
	ticker := time.NewTicker(time.Second * 3).C

	block.Name = "cpu"

	for {
		block.FullText = fmt.Sprintf("cpu %.2f%%",
			h.SysInfo().CPULoad())

		h.Tick()

		<-ticker
	}
}
