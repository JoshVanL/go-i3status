package sysinfo

import (
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

type SysInfo struct {
	cpuLoads [3]uint64
	mu       sync.Mutex
}

func New() (*SysInfo, error) {
	var sysinfo_t unix.Sysinfo_t
	err := unix.Sysinfo(&sysinfo_t)
	if err != nil {
		return nil, err
	}

	s := &SysInfo{
		cpuLoads: sysinfo_t.Loads,
	}
	go s.run()

	return s, nil
}

func (s *SysInfo) run() {
	ticker := time.NewTicker(time.Second * 3).C
	var sysinfo_t unix.Sysinfo_t

	for {
		<-ticker

		err := unix.Sysinfo(&sysinfo_t)
		if err != nil {
			continue
		}

		s.mu.Lock()
		s.cpuLoads = sysinfo_t.Loads
		s.mu.Unlock()
	}
}

//const loadScale = 65536.0 // LINUX_SYSINFO_LOADS_SCALE
const loadScale = 1.0 // LINUX_SYSINFO_LOADS_SCALE

func (s *SysInfo) CPULoads() [3]float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return [3]float64{
		float64(s.cpuLoads[0]) / loadScale,
		float64(s.cpuLoads[1]) / loadScale,
		float64(s.cpuLoads[2]) / loadScale,
	}
}
