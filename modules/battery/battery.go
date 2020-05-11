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

var (
	capPath  = filepath.Join(path, batteryName, "capacity")
	statPath = filepath.Join(path, batteryName, "status")
)

type battery struct {
	block *protocol.Block
	h     *handler.Handler
}

func Battery(block *protocol.Block, h *handler.Handler) {
	block.Name = "battery"

	ticker := time.NewTicker(time.Second * 10)
	ch := h.WatchSignal(protocol.RealTimeSignals["RTMIN+2"])

	b := &battery{
		block: block,
		h:     h,
	}

	h.Scheduler().Register(time.Minute*3, b.setBatteryString)

	go func() {
		for {
			b.setBatteryString()
			b.h.Tick()

			select {
			case <-ticker.C:
				break
			case <-ch:
				break
			}
		}
	}()
}

func (b *battery) setBatteryString() {
	status, capacity := b.getFiles()
	i, err := strconv.Atoi(string(capacity))
	if err != nil {
		b.block.FullText = err.Error()
		return
	}

	bat := getIcon(b.block, i)
	var charging string
	if string(status) == "Charging" {
		charging = " "
	}

	b.block.FullText = fmt.Sprintf("%s%s %s%%", bat, charging, capacity)
}

func (b *battery) getFiles() (status, capacity []byte) {
	status, err := utils.ReadFile(statPath)
	b.h.Must(err)

	capacity, err = utils.ReadFile(capPath)
	b.h.Must(err)

	if string(capacity) == "100" {
		status = []byte("full")
	}

	return status, capacity
}

func getIcon(b *protocol.Block, capacity int) string {
	if capacity == 100 {
		//b.Color = "#FFFFFF"
		b.Color = "#000000"
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

	if capacity >= 30 {
		b.Color = "#ffdd77"
		return ""
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
