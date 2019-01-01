package scheduler

import (
	"sync"
	"time"
)

type Scheduler struct {
	modules map[time.Duration][]func(sync.WaitGroup)
	tick    chan struct{}
	update  chan struct{}
	mu      sync.Mutex
	wg      sync.WaitGroup
}

func New(tick chan struct{}) *Scheduler {
	return &Scheduler{
		tick:    tick,
		modules: make(map[time.Duration][]func(sync.WaitGroup)),
		update:  make(chan struct{}),
	}
}

func (s *Scheduler) Register(d time.Duration, update func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f := func(wg sync.WaitGroup) {
		defer wg.Done()
		update()
	}

	fs, ok := s.modules[d]
	if !ok {
		s.modules[d] = []func(sync.WaitGroup){f}
	} else {
		s.modules[d] = append(fs, f)
	}
}

func (s *Scheduler) Run() {
	s.mu.Lock()

	go s.updater()

	for d, fs := range s.modules {
		ticker := time.NewTicker(d).C
		go s.runTicker(ticker, fs)
	}

	for _, fs := range s.modules {
		s.wg.Add(len(fs))
		for _, f := range fs {
			f(s.wg)
		}
	}
}

func (s *Scheduler) runTicker(ticker <-chan time.Time, fs []func(wg sync.WaitGroup)) {
	for {
		s.wg.Add(len(fs))
		s.update <- struct{}{}

		for _, f := range fs {
			go f(s.wg)
		}

		<-ticker
	}
}

func (s *Scheduler) updater() {
	for {
		s.update = make(chan struct{}, len(s.modules))
		<-s.update

		s.wg.Wait()
		s.tick <- struct{}{}
		close(s.tick)
	}
}
