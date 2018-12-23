package battery

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

const (
	path        = "/sys/class/power_supply"
	batteryName = "BAT0"
	// TODO: do proper path
	watchPath = "/home/josh/go/src/github.com/joshvanl/go-i3status/battery_watch"
)

func Battery(block *protocol.Block, h *handler.Handler) {
	block.Name = "battery"

	capPath := filepath.Join(path, batteryName, "capacity")
	statPath := filepath.Join(path, batteryName, "status")

	status := readFile(statPath)
	capacity := readFile(capPath)
	setBatteryString(block, status, capacity)
	h.Tick()

	// TODO: use signals
	_, err := os.Create(watchPath)
	if err != nil {
		panic(err)
	}

	ch, err := h.WatchFile(watchPath)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Minute * 3).C

	for {
		tickNow := false

		select {
		case <-ch:
			tickNow = true
		case <-ticker:
		}

		status := readFile(statPath)
		capacity := readFile(capPath)
		setBatteryString(block, status, capacity)

		if tickNow {
			h.Tick()
		}
	}
}

func setBatteryString(b *protocol.Block, status, capacity []byte) {
	i, err := strconv.Atoi(string(capacity))
	if err != nil {
		b.FullText = err.Error()
		return
	}

	bat := getIcon(b, i)
	var charging string
	if string(status) == "Charging" {
		charging = " "
	}

	b.FullText = fmt.Sprintf("%s%s %s%%", bat, charging, capacity)
}

func getIcon(b *protocol.Block, capacity int) string {
	if capacity == 100 {
		b.Color = "#FFFFFF"
		return ""
	}

	if capacity >= 90 {
		b.Color = "#ccffcc"
	}

	if capacity >= 70 {
		b.Color = "#bbffbb"
	}

	if capacity >= 75 {
		return ""
	}

	if capacity >= 60 {
		b.Color = "#ddffaa"
	}

	if capacity >= 50 {
		b.Color = "#eeffaa"
		return ""
	}

	if capacity >= 40 {
		b.Color = "#ffdd77"
		return ""
	}

	if capacity >= 30 {
		b.Color = "#ffbb44"
	}

	if capacity < 30 && capacity >= 20 {
		b.Color = "#ffaaaa"
	}

	if capacity < 20 {
		b.Color = "#ff1c1c"
	}

	if capacity >= 25 {
		return ""
	}

	if capacity <= 10 {
		b.Color = "#FF0000"
		return ""
	}

	return ""
}

func readFile(path string) []byte {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	return bytes.TrimSpace(f)
}
