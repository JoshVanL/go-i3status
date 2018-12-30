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

type Watcher struct {
	fd       int
	watching map[int32]chan struct{}
	sockets  map[string]net.Listener
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
		sockets:  make(map[string]net.Listener),
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

	l, err := net.Listen("unix", path)
	if err != nil {
		return err
	}

	w.mu.Lock()
	w.sockets[path] = l
	w.mu.Unlock()

	go func() {
		for {
			c, _ := l.Accept()
			ch <- struct{}{}

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
	for s, l := range w.sockets {
		l.Close()
		os.Remove(s)
	}
}
