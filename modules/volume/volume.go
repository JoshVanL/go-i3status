package volume

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func Volume(block *protocol.Block, h *handler.Handler) {
	block.Name = "volume"
	//block.Separator = false
	//block.SeparatorBlockWidth = 10

	ch := h.WatchSignal(protocol.RealTimeSignals["RTMIN+3"])

	go func() {
		for {
			block.FullText = update(h)
			h.Tick()

			<-ch
		}
	}()
}

func update(h *handler.Handler) string {
	cmd := exec.Command("amixer",
		"sget", "Master")

	bs, err := cmd.Output()
	h.Must(err)

	split := bytes.Split(bs, []byte{'\n'})
	if len(split) < 5 {
		h.Must(fmt.Errorf("%s", bs))
	}

	if bytes.Contains(bs, []byte{'o', 'f', 'f'}) {
		return "♪ "
	}

	var out string
	var index int
	for i, b := range split[4] {
		if b == '[' {
			index = i
			break
		}
	}

	if index == 0 {
		h.Must(fmt.Errorf("filed to find volume: %s", split[4]))
	}

	for _, b := range split[4][index+1:] {
		if b == ']' {
			break
		}

		out = out + string(b)
	}

	return fmt.Sprintf("♪ %s", out)
}
