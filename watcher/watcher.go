package watcher

import (
	"net"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"unsafe"
)

const (
	socketDir = "/var/run/go-i3status"

	sys_IN_MODIFY uint32 = syscall.IN_MODIFY
)

type socketPair struct {
	chs []chan struct{}
	l   net.Listener
}

type Watcher struct {
	fd       int
	watching map[int32]chan struct{}
	sockets  map[string]*socketPair
	mu       sync.Mutex
}

func New() (*Watcher, error) {
	fd, err := syscall.InotifyInit()
	if err != nil {
		return nil, err
	}

	if err := clearSockets(); err != nil {
		return nil, err
	}

	w := &Watcher{
		fd:       fd,
		watching: make(map[int32]chan struct{}),
		sockets:  make(map[string]*socketPair),
	}

	go w.run()

	return w, nil
}

func (w *Watcher) run() {
	var buf [syscall.SizeofInotifyEvent * 1024]byte
	for {
		_, err := syscall.Read(w.fd, buf[:])
		if err != nil {
			continue
		}

		raw := (*syscall.InotifyEvent)(unsafe.Pointer(&buf))

		if (raw.Mask & sys_IN_MODIFY) != sys_IN_MODIFY {
			continue
		}

		w.mu.Lock()
		ch, ok := w.watching[raw.Wd]
		w.mu.Unlock()

		if !ok {
			continue
		}

		ch <- struct{}{}
	}
}

func (w *Watcher) AddSocket(module string, ch chan struct{}) error {
	path := filepath.Join(socketDir, module+".sock")

	w.mu.Lock()
	defer w.mu.Unlock()

	if p, ok := w.sockets[path]; ok {
		p.chs = append(p.chs, ch)
		return nil
	}

	l, err := net.Listen("unix", path)
	if err != nil {
		return err
	}

	w.sockets[path] = &socketPair{[]chan struct{}{ch}, l}

	go func() {
		for {
			c, _ := l.Accept()

			w.mu.Lock()

			p, ok := w.sockets[path]
			if !ok {
				w.mu.Unlock()
				continue
			}

			for _, ch := range p.chs {
				ch <- struct{}{}
			}

			w.mu.Unlock()

			if c != nil {
				c.Close()
			}
		}
	}()

	return nil
}

func (w *Watcher) AddFile(path string, ch chan struct{}) error {
	wd, err := syscall.InotifyAddWatch(w.fd, path, sys_IN_MODIFY)
	if err != nil {
		return err
	}

	w.mu.Lock()
	w.watching[int32(wd)] = ch
	w.mu.Unlock()

	return nil
}

func clearSockets() error {
	files, err := filepath.Glob(filepath.Join(socketDir, "*"))
	if err != nil {
		return err
	}

	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Watcher) Kill() {
	// we don't unlock because we don't want more watchers
	w.mu.Lock()
	for s, p := range w.sockets {
		p.l.Close()
		os.Remove(s)
	}
}
