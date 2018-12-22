package blocks

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

const (
	path        = "/sys/class/power_supply"
	batteryName = "BAT0"
)

func Battery(block *protocol.Block, h *handler.Handler) {
	block.Name = "battery"

	statPath := filepath.Join(path, batteryName, "status")
	capPath := filepath.Join(path, batteryName, "capacity")

	status := readFile(statPath)
	capacity := readFile(capPath)

	block.FullText = fmt.Sprintf("%s %s%%", status, capacity)
	h.Tick()

	statCh, err := h.WatchFile(statPath)
	if err != nil {
		return
	}

	capCh, err := h.WatchFile(capPath)
	if err != nil {
		return
	}

	for {
		select {
		case <-statCh:
			status = readFile(statPath)

		case <-capCh:
			capacity = readFile(capPath)
		}

		block.FullText = fmt.Sprintf("%s %s%%", status, capacity)
		h.Tick()
	}
}

func readFile(path string) []byte {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	return bytes.TrimSpace(f)
}
