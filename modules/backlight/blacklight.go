package backlight

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func Backlight(block *protocol.Block, h *handler.Handler) {
	block.Name = "backlight"

	ch := h.WatchSignal(protocol.RealTimeSignals["RTMIN+5"])

	go func() {
		for {
			block.FullText = update(h)
			h.Tick()

			<-ch
		}
	}()
}

func update(h *handler.Handler) string {
	cmd := exec.Command("light")

	bs, err := cmd.Output()
	h.Must(err)

	f, err := strconv.ParseFloat(
		strings.TrimSpace(string(bs)), 64)
	if err != nil {
		return "-"
	}

	return fmt.Sprintf(" %.0f%%", f)
}
