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
	ifaceName = "wlp58s0"
)

func Wifi(block *protocol.Block, h *handler.Handler) {
	ticker := time.NewTicker(time.Second * 15).C

	block.Name = "wifi"
	block.Separator = false
	block.SeparatorBlockWidth = 15

	for {
		setString(block, h)
		h.Tick()

		<-ticker
	}
}

func setString(block *protocol.Block, h *handler.Handler) {
	n := wirelessInt(h)

	if n == -1 {
		block.FullText = ""
		return
	}

	if n >= 80 {
		block.Color = "#aaddaa"
	} else if n >= 60 {
		block.Color = "#fff600"
	} else if n >= 40 {
		block.Color = "#ffae88"
	} else {
		block.Color = "#ff0000"
	}

	block.FullText = fmt.Sprintf(" %d%%", n)
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

		//return fmt.Sprintf(" %.0f%%", (n * 100 / 70))
		return int(n * 100 / 70)
	}

	return -1
}
