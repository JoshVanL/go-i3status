package blocks

import (
	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func All() []func(*protocol.Block, *handler.Handler) {
	return []func(*protocol.Block, *handler.Handler){
		Time,
	}
}
