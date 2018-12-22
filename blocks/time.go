package blocks

import (
	"fmt"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func Time(block *protocol.Block, h *handler.Handler) {
	block.Name = "time"
	now := time.Now()
	m := now.Minute()
	block.FullText = getString(now)
	h.Tick()

	ticker := time.NewTicker(time.Second)
	for t := range ticker.C {
		if m != t.Minute() {
			block.FullText = getString(t)
			h.Tick()
			m = t.Minute()
		}
	}
}

func getString(t time.Time) string {
	return fmt.Sprintf("%s %d %s %d  |  %s ",
		t.Format("Mon"), t.Day(), t.Format("Dec"), t.Year(), t.Format("15:04"))
}
