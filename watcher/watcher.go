package watcher

import (
	"sync"
	"syscall"
	"unsafe"
)

const (
	sys_IN_MOVED_TO    uint32 = syscall.IN_MOVED_TO
	sys_IN_CREATE      uint32 = syscall.IN_CREATE
	sys_IN_MOVED_FROM  uint32 = syscall.IN_MOVED_FROM
	sys_IN_ATTRIB      uint32 = syscall.IN_ATTRIB
	sys_IN_MODIFY      uint32 = syscall.IN_MODIFY
	sys_IN_DELETE_SELF uint32 = syscall.IN_DELETE_SELF
	sys_IN_DELETE      uint32 = syscall.IN_DELETE
	sys_IN_MOVE_SELF   uint32 = syscall.IN_MOVE_SELF

	sys_AGNOSTIC_EVENTS = sys_IN_MOVED_TO | sys_IN_MOVED_FROM | sys_IN_CREATE | sys_IN_ATTRIB | sys_IN_MODIFY | sys_IN_MOVE_SELF | sys_IN_DELETE | sys_IN_DELETE_SELF
)

type Watcher struct {
	fd       int
	watching map[int32]chan struct{}
	mu       sync.Mutex
}

func New() (*Watcher, error) {
	fd, err := syscall.InotifyInit()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		fd:       fd,
		watching: make(map[int32]chan struct{}),
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
