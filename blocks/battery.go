package blocks

import (
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func Battery(block *protocol.Block, h *handler.Handler) {
	block.Name = "battery"
	now := time.Now()
	m := now.Minute()
	block.FullText = now.Format("15:04") + " "
	h.Tick()

	ticker := time.NewTicker(time.Second)
	for t := range ticker.C {
		if m != t.Minute() {
			m = t.Minute()
			block.FullText = t.Format("15:04") + " "
			h.Tick()
		}
	}
}
