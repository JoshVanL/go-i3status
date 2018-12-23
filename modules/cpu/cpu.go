package cpu

//#include <stdlib.h>
import "C"
import (
	"fmt"
	"time"

	"github.com/joshvanl/go-i3status/handler"
	"github.com/joshvanl/go-i3status/protocol"
)

func CPU(block *protocol.Block, h *handler.Handler) {
	ticker := time.NewTicker(time.Second / 2).C

	block.Name = "cpu"

	for {
		loads := h.SysInfo().CPULoads()
		_, err := getloadavg(&loads, 3)
		if err != nil {
			panic(err)
		}

		block.FullText = fmt.Sprintf("cpu %v",
			loads[0])

		h.Tick()

		<-ticker
	}
}

// To allow tests to mock out getloadavg.
var getloadavg = func(out *[3]float64, count int) (int, error) {
	read, err := C.getloadavg((*C.double)(&out[0]), (C.int)(count))
	return int(read), err
}
