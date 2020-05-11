package bluetooth

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func Bluetooth(block *protocol.Block, h *handler.Handler) {
	block.Name = "bluetooth"
	block.Separator = false
	block.SeparatorBlockWidth = 10

	ch := h.WatchSignal(protocol.RealTimeSignals["RTMIN+6"])
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			block.FullText = update(h)
			h.Tick()

			select {
			case <-ch:
			case <-ticker.C:
			}
		}
	}()
}

func update(h *handler.Handler) string {
	p := h.SysInfo().Processes()
	p.Refresh()

	if p.Running("/usr/lib/bluetooth/bluetoothd") {

		cmd := exec.Command("bluetoothctl", "info")

		out, err := cmd.Output()
		if err != nil {
			return ""
		}

		split := bytes.Split(out, []byte{'\n'})

		var outString [][]byte
		for _, ss := range split {
			if bytes.Contains(ss, []byte("Name: ")) {
				outString = append(outString, bytes.TrimSpace(bytes.Split(ss, []byte("Name: "))[1]))
			}
		}

		return fmt.Sprintf("%s ", bytes.Join(outString, []byte(", ")))
	}

	return ""
}
