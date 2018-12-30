package temp

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/modules/utils"
	"github.com/joshvanl/go-i3status/protocol"
)

const (
	thermalDir  = "/sys/class/thermal"
	thermalType = "x86_pkg_temp"
)

func Temp(block *protocol.Block, h *handler.Handler) {
	ticker := time.NewTicker(time.Second * 5).C

	block.Name = "temp"

	files, err := ioutil.ReadDir(thermalDir)
	h.Must(err)

	var path string
	for _, f := range files {
		b, err := utils.ReadFile(filepath.Join(thermalDir, f.Name(), "type"))
		h.Must(err)

		if string(b) == thermalType {
			path = filepath.Join(thermalDir, f.Name(), "temp")
			break
		}
	}

	if path == "" {
		h.Must(fmt.Errorf("failed to find thermal with type %s",
			thermalType))
	}

	normalColor := block.Color

	for {
		b, err := utils.ReadFile(path)
		h.Must(err)

		temp, err := strconv.ParseFloat(string(b), 64)
		h.Must(err)

		temp = temp / 1000

		if temp > 90 {
			block.Color = "#ff3333"
		} else if temp > 70 {
			block.Color = "#ffaa33"
		} else {
			block.Color = normalColor
		}
		block.FullText = fmt.Sprintf("%1.f°C", temp)

		h.Tick()

		<-ticker
	}
}
