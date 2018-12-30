package disk

import (
	"fmt"
	"syscall"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func Disk(block *protocol.Block, h *handler.Handler) {
	ticker := time.NewTicker(time.Minute * 60).C

	block.Name = "disk"

	for {
		var stat syscall.Statfs_t
		err := syscall.Statfs("/", &stat)
		h.Must(err)

		block.FullText = fmt.Sprintf(" %.2fG",
			float64(stat.Bavail*uint64(stat.Bsize))/(1024*1024*1024))
		h.Tick()

		<-ticker
	}
}
