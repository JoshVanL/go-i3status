package handler

import (
	"bufio"
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joshvanl/go-i3status/protocol"
)

type Handler struct {
	blocks []*protocol.Block
	stdin  *os.File
	stdout *os.File

	registeredEvents map[string][]func(*protocol.ClickEvent)

	stopCh     chan struct{}
	paused     bool
	blocksLock sync.RWMutex
}

func New() *Handler {
	h := &Handler{
		stdin:            os.Stdin,
		stdout:           os.Stdout,
		stopCh:           make(chan struct{}),
		paused:           false,
		registeredEvents: make(map[string][]func(*protocol.ClickEvent)),
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
		return nil
	}

	h.stdout.Write(append(b, '['))

	return h
}

func (h *Handler) RegisterClickEvent(name string, f func(*protocol.ClickEvent)) {
	h.LockBlocks()
	h.registeredEvents[name] = append(h.registeredEvents[name], f)
	h.UnlockBlocks()
}

func (h *Handler) RegisterBlock(b *protocol.Block) {
	h.LockBlocks()
	h.blocks = append(h.blocks, b)
	h.UnlockBlocks()
}

func (h *Handler) Tick() {
	if h.paused {
		return
	}

	h.blocksLock.RLock()

	b, err := json.Marshal(h.blocks)
	if err != nil {
		return
	}

	h.blocksLock.RUnlock()

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

func (h *Handler) LockBlocks() {
	h.blocksLock.Lock()
}

func (h *Handler) UnlockBlocks() {
	h.blocksLock.Unlock()
}

func (h *Handler) signalHandler() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGCONT, syscall.SIGSTOP, syscall.SIGKILL)

	for s := range sig {
		switch s {
		case syscall.SIGCONT:
			h.paused = false

		case syscall.SIGSTOP:
			h.paused = true

		case syscall.SIGKILL:
			os.Exit(1)

		default:
			continue
		}
	}
}
