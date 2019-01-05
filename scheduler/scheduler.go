package scheduler

import (
	"sync"
	"time"
)

type Scheduler struct {
	modules map[time.Duration][]func()
	tick    chan struct{}
	mu      sync.Mutex
	wg      sync.WaitGroup

	update     chan struct{}
	updatingMu sync.RWMutex
	updating   bool
}

func New(tick chan struct{}) *Scheduler {
	return &Scheduler{
		tick:    tick,
		modules: make(map[time.Duration][]func()),
		update:  make(chan struct{}),
	}
}

func (s *Scheduler) Register(d time.Duration, update func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f := func() {
		update()
		s.wg.Done()
	}

	fs, ok := s.modules[d]
	if !ok {
		s.modules[d] = []func(){f}
	} else {
		s.modules[d] = append(fs, f)
	}
}

func (s *Scheduler) Run() {
	s.mu.Lock()

	for d, fs := range s.modules {
		go s.runTicker(time.NewTicker(d).C, fs)
	}

	for _, fs := range s.modules {
		s.wg.Add(len(fs))
		for _, f := range fs {
			go f()
		}
	}

	s.wg.Wait()

	go s.updater()
}

func (s *Scheduler) runTicker(ticker <-chan time.Time, fs []func()) {
	for {
		<-ticker
		s.wg.Add(len(fs))

		s.updatingMu.RLock()

		if !s.updating {
			s.updating = true
			s.update <- struct{}{}
		}

		s.updatingMu.RUnlock()

		for _, f := range fs {
			go f()
		}

	}
}

func (s *Scheduler) updater() {
	for {
		<-s.update
		s.wg.Wait()

		s.tick <- struct{}{}

		s.updatingMu.Lock()
		s.updating = false
		s.updatingMu.Unlock()
	}
}
