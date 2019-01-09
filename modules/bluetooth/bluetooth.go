package bluetooth

import (
	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func Bluetooth(block *protocol.Block, h *handler.Handler) {
	block.Name = "bluetooth"
	block.Separator = false
	block.SeparatorBlockWidth = 10

	ch := h.WatchSignal(protocol.RealTimeSignals["RTMIN+6"])

	go func() {
		for {
			block.FullText = update(h)
			h.Tick()

			<-ch
		}
	}()
}

func update(h *handler.Handler) string {
	p := h.SysInfo().Processes()
	p.Refresh()

	if p.Running("/usr/lib/bluetooth/bluetoothd") {
		return ""
	}

	return ""
}
