package bandwidth

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
	statPath  = "/proc/net/dev"
	secs      = 3
)

type netStat struct {
	block *protocol.Block
	h     *handler.Handler

	down        bool
	received    uint64
	transmitted uint64
}

func Bandwidth(block *protocol.Block, h *handler.Handler) {
	ticker := time.NewTicker(time.Second * secs).C
	ch := h.WatchSignal(protocol.RealTimeSignals["RTMIN+1"])

	block.Name = "bandwidth"
	block.Separator = false
	block.SeparatorBlockWidth = 10

	n := netStat{
		block: block,
		h:     h,
		down:  ifaceDown(h),
	}

	for {
		n.setString()
		h.Tick()

		select {
		case <-ticker:
		case <-ch:
			time.Sleep(time.Second * 3)
			n.down = ifaceDown(n.h)
		}
	}
}

func ifaceDown(h *handler.Handler) bool {
	b, err := utils.ReadFile(filepath.Join(
		"/sys/class/net", ifaceName, "operstate"))
	h.Must(err)

	return string(b) == "down"
}

func (n *netStat) setString() {
	if n.down {
		n.block.FullText = ""
		return
	}

	b, err := utils.ReadFile(statPath)
	n.h.Must(err)

	for _, lines := range strings.Split(string(b), "\n") {
		fields := strings.Fields(lines)

		if len(fields) < 4 || fields[0] != ifaceName+":" {
			continue
		}

		rec, err := strconv.ParseUint(fields[1], 10, 64)
		n.h.Must(err)

		tran, err := strconv.ParseUint(fields[9], 10, 64)
		n.h.Must(err)

		recf := float64(rec-n.received) / (secs * 1024 * 1024)
		tranf := float64(tran-n.transmitted) / (secs * 1024 * 1024)

		if n.received != 0 {
			n.block.FullText = fmt.Sprintf(" %.1f / %.1f Mb/s",
				recf, tranf)
		} else {
			n.block.FullText = " 0.0 / 0.0 Mb/s"
		}

		n.received = rec
		n.transmitted = tran

		return
	}

	n.block.FullText = ""
}
