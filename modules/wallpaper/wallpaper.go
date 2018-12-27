package wallpaper

import (
	"fmt"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/modules/utils"
	"github.com/joshvanl/go-i3status/protocol"
)

const (
	wallPath = "/home/josh/scripts/currWall"
)

func Wallpaper(block *protocol.Block, h *handler.Handler) {
	block.Name = "wallpaper"

	ch, err := h.WatchFile(wallPath)
	if err != nil {
		panic(err)
	}

	for {
		block.FullText = fmt.Sprintf(" %s",
			utils.ReadFile(wallPath))
		h.Tick()

		<-ch
	}
}
