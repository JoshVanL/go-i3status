package iface

import (
	"net"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

const (
	ifaceName = "wlp58s0"
)

func IFace(block *protocol.Block, h *handler.Handler) {
	block.Name = "iface"

	ch, err := h.WatchSocket("iface")
	h.Must(err)

	for {
		iface, err := net.InterfaceByName(ifaceName)
		h.Must(err)

		addrs, err := iface.Addrs()
		h.Must(err)

		if len(addrs) == 0 {
			block.FullText = "down"
			block.Color = "#ee9999"
		}

		for _, addr := range addrs {
			v, ok := addr.(*net.IPNet)
			if !ok || v.IP.To4() == nil {
				continue
			}

			block.FullText = v.IP.String()
			block.Color = "#aaddaa"
			break
		}

		h.Tick()

		<-ch

		// netctl uses a pre down hook so have some buffer
		// so we don't catch the last up interfaces
		if block.FullText != "down" {
			time.Sleep(time.Second * 2)
		}
	}
}
