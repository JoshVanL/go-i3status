package wifi

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/modules/utils"
	"github.com/joshvanl/go-i3status/protocol"
)

const (
	ifaceName = "wlan0"
)

func Wifi(block *protocol.Block, h *handler.Handler) {
	block.Name = "wifi"
	block.Separator = false
	block.SeparatorBlockWidth = 10

	update := func() {
		n := wirelessInt(h)

		if n == -1 {
			block.FullText = ""
			return
		}

		//if n >= 70 {
		//	block.Color = "#0050b8"
		//} else if n >= 50 {
		if n >= 50 {
			block.Color = "#aaddaa"
		} else if n >= 30 {
			block.Color = "#a65f3d"
			//block.Color = "#ffae88"
		} else {
			block.Color = "#ff0000"
		}

		block.FullText = fmt.Sprintf(" %d%%", n)
	}

	h.Scheduler().Register(time.Minute, update)

	ch := h.WatchSignal(protocol.RealTimeSignals["RTMIN+1"])

	go func() {
		for {
			update()
			h.Tick()

			<-ch
			time.Sleep(time.Second * 3)
		}
	}()
}

func wirelessInt(h *handler.Handler) int {
	b, err := utils.ReadFile(filepath.Join(
		"/sys/class/net", ifaceName, "operstate"))
	h.Must(err)

	if string(b) == "down" {
		return -1
	}

	b, err = utils.ReadFile("/proc/net/wireless")
	h.Must(err)

	for _, line := range strings.Split(string(b), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 || fields[0] != ifaceName+":" {
			continue
		}

		n, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return -1
		}

		return int(n * 100 / 70)
	}

	return -1
}
