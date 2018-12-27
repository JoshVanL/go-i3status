package cpu

import (
	"fmt"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func CPU(block *protocol.Block, h *handler.Handler) {
	ticker := time.NewTicker(time.Second / 2).C

	block.Name = "cpu"

	for {
		block.FullText = fmt.Sprintf("cpu %v %v",
			h.SysInfo().CPULoads()[0],
			h.SysInfo().CPULoads()[2])

		h.Tick()

		<-ticker
	}
}
