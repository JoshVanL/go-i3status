package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"syscall"

	"github.com/joshvanl/go-i3status/errors"
	"github.com/joshvanl/go-i3status/protocol"
	"github.com/joshvanl/go-i3status/sysinfo"
	"github.com/joshvanl/go-i3status/watcher"
)

type Handler struct {
	blocks []*protocol.Block
	stdin  *os.File
	stdout *os.File

	registeredEvents map[string][]func(*protocol.ClickEvent)
	watcher          *watcher.Watcher
	sysinfo          *sysinfo.SysInfo

	paused     bool
	blocksLock sync.Mutex
}

func New() (*Handler, error) {
	w, err := watcher.New()
	if err != nil {
		return nil, err
	}

	s, err := sysinfo.New()
	if err != nil {
		return nil, err
	}

	h := &Handler{
		stdin:            os.Stdin,
		stdout:           os.Stdout,
		paused:           false,
		registeredEvents: make(map[string][]func(*protocol.ClickEvent)),
		watcher:          w,
		sysinfo:          s,
	}

	go h.signalHandler()
	go h.clickEvents()

	b, err := json.Marshal(&protocol.Header{
		Version:        1,
		StopSignal:     int(syscall.SIGSTOP),
		ContinueSignal: int(syscall.SIGCONT),
		ClickEvents:    true,
	})
	if err != nil {
		return nil, err
	}

	h.stdout.Write(append(b, '['))

	return h, nil
}

func (h *Handler) RegisterClickEvent(name string, f func(*protocol.ClickEvent)) {
	h.blocksLock.Lock()
	h.registeredEvents[name] = append(h.registeredEvents[name], f)
	h.blocksLock.Unlock()
}

func (h *Handler) RegisterBlock(b *protocol.Block) {
	h.blocksLock.Lock()
	h.blocks = append(h.blocks, b)
	h.blocksLock.Unlock()
}

func (h *Handler) Tick() {
	h.blocksLock.Lock()
	defer h.blocksLock.Unlock()

	if h.paused {
		return
	}

	b, err := json.Marshal(h.blocks)
	if err != nil {
		return
	}

	h.stdout.Write(append(b, ','))
}

func (h *Handler) clickEvents() {
	if h.paused {
		return
	}

	scanner := bufio.NewScanner(h.stdin)

	for scanner.Scan() && !h.paused {
		clickEvent := new(protocol.ClickEvent)
		err := json.Unmarshal(scanner.Bytes(), clickEvent)
		if err != nil {
			continue
		}

		fs, ok := h.registeredEvents[clickEvent.Name]
		if !ok {
			continue
		}

		for _, f := range fs {
			go f(clickEvent)
		}
	}
}

func (h *Handler) WatchFile(path string) (<-chan struct{}, error) {
	ch := make(chan struct{})
	err := h.watcher.AddFile(path, ch)
	return ch, err
}

func (h *Handler) WatchSocket(module string) (<-chan struct{}, error) {
	ch := make(chan struct{})
	err := h.watcher.AddSocket(module, ch)
	return ch, err
}

func (h *Handler) SysInfo() *sysinfo.SysInfo {
	return h.sysinfo
}

func (h *Handler) Must(err error) {
	if err == nil {
		return
	}

	h.watcher.Kill()
	errors.Kill(fmt.Errorf("go-i3status was killed: %v\n", err))
}

func (h *Handler) signalHandler() {
	ch := errors.NewSignalHandler()

	for s := range ch {
		switch s {
		case 1:
			h.paused = false
			break

		case 0:
			h.paused = true
			break

		case -1:
			h.Must(fmt.Errorf("got signal interupt"))

		default:
			break
		}
	}
}
