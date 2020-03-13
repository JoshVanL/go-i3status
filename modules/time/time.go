package time

import (
	"fmt"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func Time(block *protocol.Block, h *handler.Handler) {
	block.Name = "time"
	now := time.Now()
	block.FullText = getTimeString(now)
	h.Tick()

	ticker := time.NewTicker(time.Second / 2)

	go func() {
		<-ticker.C

		block.FullText = getTimeString(time.Now())
		h.Tick()

		ticker := time.NewTicker(time.Minute)

		for t := range ticker.C {
			block.FullText = getTimeString(t)
			h.Tick()
		}
	}()
}

func getTimeString(t time.Time) string {
	return fmt.Sprintf("%s ", t.Format("15:04"))
}
