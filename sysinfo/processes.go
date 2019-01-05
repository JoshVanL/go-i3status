package sysinfo

import (
	"io/ioutil"
	"os"
	"sync"
)

type Processes struct {
	names map[string]struct{}
	mu    sync.Mutex
}

func newProcesses() *Processes {
	p := new(Processes)
	p.Refresh()

	return p
}

func (p *Processes) Refresh() {
	names := make(map[string]struct{})

	proc, err := os.OpenFile("/proc", os.O_RDONLY, 0333)
	if err != nil {
		return
	}
	defer proc.Close()

	fs, err := proc.Readdirnames(-1)
	if err != nil {
		return
	}
	proc.Close()

	for _, f := range fs {

		if f[0] < '0' || f[0] > '9' {
			continue
		}

		b, err := ioutil.ReadFile("/proc/" + f + "/cmdline")
		if err != nil {
			continue
		}

		if len(b) > 0 {
			names[string(b[:len(b)-1])] = struct{}{}
		}
	}

	p.mu.Lock()
	p.names = names
	p.mu.Unlock()
}

func (p *Processes) Running(processCmd string) bool {
	p.mu.Lock()
	_, ok := p.names[processCmd]
	p.mu.Unlock()

	if !ok {
		return false
	}

	return true
}
