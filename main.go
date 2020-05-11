package main

import (
	"fmt"

	"github.com/joshvanl/go-i3status/errors"
	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/modules/backlight"
	"github.com/joshvanl/go-i3status/modules/bandwidth"
	"github.com/joshvanl/go-i3status/modules/battery"
	"github.com/joshvanl/go-i3status/modules/bluetooth"
	"github.com/joshvanl/go-i3status/modules/cpu"
	"github.com/joshvanl/go-i3status/modules/date"
	"github.com/joshvanl/go-i3status/modules/disk"
	"github.com/joshvanl/go-i3status/modules/iface"
	"github.com/joshvanl/go-i3status/modules/memory"
	"github.com/joshvanl/go-i3status/modules/mic"
	"github.com/joshvanl/go-i3status/modules/temp"
	"github.com/joshvanl/go-i3status/modules/time"
	"github.com/joshvanl/go-i3status/modules/volume"

	//"github.com/joshvanl/go-i3status/modules/wallpaper"
	"github.com/joshvanl/go-i3status/modules/wifi"
	"github.com/joshvanl/go-i3status/protocol"
)

var (
	enabledBlocks = []func(*protocol.Block, *handler.Handler){
		bluetooth.Bluetooth,
		mic.Mic,
		volume.Volume,
		//wallpaper.Wallpaper,
		cpu.CPU,
		memory.Memory,
		disk.Disk,
		wifi.Wifi,
		bandwidth.Bandwidth,
		iface.IFace,
		temp.Temp,
		backlight.Backlight,
		battery.Battery,
		date.Date,
		time.Time,
	}
)

func main() {
	h, err := handler.New()
	if err != nil {
		errors.Kill(fmt.Errorf("error creating handler: %s\n", err))
	}

	for _, f := range enabledBlocks {
		b := &protocol.Block{
			Align:               protocol.Right,
			Color:               "#000000FF",
			Background:          "#aeaeb2BB",
			Border:              "#aeaeb2BB",
			SeparatorBlockWidth: 30,
			Separator:           true,
		}
		h.RegisterBlock(b)
		f(b, h)
	}

	h.Scheduler().Run()

	select {}
}
