package wallpaper

import (
	"fmt"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/modules/utils"
	"github.com/joshvanl/go-i3status/protocol"
)

const (
	wallPath = "/home/josh/.config/wallpaper"
)

func Wallpaper(block *protocol.Block, h *handler.Handler) {
	block.Name = "wallpaper"

	ch, err := h.WatchFile(wallPath)
	h.Must(err)

	go func() {
		for {
			f, err := utils.ReadFile(wallPath)
			h.Must(err)

			block.FullText = fmt.Sprintf(" %s", f)
			h.Tick()

			<-ch
		}
	}()
}
