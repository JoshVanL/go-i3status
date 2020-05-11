package mic

import (
	"bytes"
	"os/exec"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func Mic(block *protocol.Block, h *handler.Handler) {
	block.Name = "mic"
	block.Separator = true
	block.SeparatorBlockWidth = 10

	ch := h.WatchSignal(protocol.RealTimeSignals["RTMIN+4"])

	go func() {
		for {
			block.FullText = update(h)
			h.Tick()

			<-ch
		}
	}()
}

func update(h *handler.Handler) string {
	cmd := exec.Command("amixer",
		"sget", "Capture")

	bs, err := cmd.Output()
	h.Must(err)

	if bytes.Contains(bs, []byte("[off]")) {
		return " "
	}

	return " "
}
