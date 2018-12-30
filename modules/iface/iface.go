package iface

import (
	//"net"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

const (
	iface = "wlp58s0"
)

func IFace(block *protocol.Block, h *handler.Handler) {
	block.Name = "iface"

	ch, err := h.WatchSocket("iface")
	h.Must(err)

	for {
		//ifaces, err := net.Interfaces()
		//if err != nil {
		//}
		<-ch
	}
}
