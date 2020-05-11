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

	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Second / 2)

	go func() {
		for {
			block.FullText = getTimeString(time.Now().In(loc))
			h.Tick()

			<-ticker.C
		}
	}()
}

func getTimeString(t time.Time) string {
	return fmt.Sprintf("%s ", t.Format("15:04"))
}
