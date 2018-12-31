package date

import (
	"fmt"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func Date(block *protocol.Block, h *handler.Handler) {
	block.Name = "date"
	now := time.Now()

	block.FullText = getDateString(now)
	h.Tick()

	till := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, now.Location())
	go func() {
		time.Sleep(time.Until(till))

		block.FullText = getDateString(time.Now())
		h.Tick()

		ticker := time.NewTicker(time.Hour * 24)

		for t := range ticker.C {
			block.FullText = getDateString(t)
			h.Tick()
		}
	}()

}

func getDateString(t time.Time) string {
	return fmt.Sprintf("%s %d %s %d",
		t.Format("Mon"), t.Day(), t.Format("Dec"), t.Year())
}
