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

	hour := now.Hour()
	min := now.Minute()

	min = (min + 1) % 60
	if min == 0 {
		hour = (hour + 1) % 24
	}

	till := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
	time.Sleep(time.Until(till))

	block.FullText = getTimeString(time.Now())
	h.Tick()

	ticker := time.NewTicker(time.Minute)
	for t := range ticker.C {
		block.FullText = getTimeString(t)
		h.Tick()
	}
}

func getTimeString(t time.Time) string {
	return fmt.Sprintf("%s ", t.Format("15:04"))
}
