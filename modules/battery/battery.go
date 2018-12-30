package battery

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/modules/utils"
	"github.com/joshvanl/go-i3status/protocol"
)

const (
	path        = "/sys/class/power_supply"
	batteryName = "BAT0"
)

func Battery(block *protocol.Block, h *handler.Handler) {
	block.Name = "battery"

	capPath := filepath.Join(path, batteryName, "capacity")
	statPath := filepath.Join(path, batteryName, "status")

	status := utils.ReadFile(statPath)
	capacity := utils.ReadFile(capPath)
	setBatteryString(block, status, capacity)
	h.Tick()

	ch, err := h.WatchSocket("battery")
	if err != nil {
		h.Kill(err)
	}

	ticker := time.NewTicker(time.Minute * 3).C

	for {
		tickNow := false

		select {
		case <-ch:
			tickNow = true
		case <-ticker:
		}

		status := utils.ReadFile(statPath)
		capacity := utils.ReadFile(capPath)
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
