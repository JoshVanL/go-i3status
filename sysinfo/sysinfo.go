package sysinfo

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

type SysInfo struct {
	cpuIdleOld  uint64
	cpuIdleNew  uint64
	cpuTotalOld uint64
	cpuTotalNew uint64

	memUse    [2]uint64
	mu        sync.Mutex
	processes *Processes
}

func New() (*SysInfo, error) {
	var sysinfo_t unix.Sysinfo_t
	err := unix.Sysinfo(&sysinfo_t)
	if err != nil {
		return nil, err
	}

	s := &SysInfo{
		memUse:    [2]uint64{sysinfo_t.Freeram, sysinfo_t.Totalram},
		processes: newProcesses(),
	}

	go s.run()

	return s, nil
}

func (s *SysInfo) run() {
	ticker := time.NewTicker(time.Second * 3).C
	var sysinfo_t unix.Sysinfo_t

	for {
		s.updateCPU()

		err := unix.Sysinfo(&sysinfo_t)
		if err != nil {
			continue
		}

		s.mu.Lock()
		s.memUse = [2]uint64{sysinfo_t.Freeram, sysinfo_t.Totalram}
		s.mu.Unlock()

		<-ticker
	}
}

func (s *SysInfo) Memory() [3]float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, g := float64(s.memUse[0])/(1024*1024*1024), float64(s.memUse[1])/(1024*1024*1024)

	return [3]float64{u, g, (u / g) * 100}
}

func (s *SysInfo) CPULoad() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	idle := float64(s.cpuIdleNew - s.cpuIdleOld)
	total := float64(s.cpuTotalNew - s.cpuTotalOld)
	if total == 0 {
		return -1
	}

	return 100 * (total - idle) / total
}

func (s *SysInfo) updateCPU() {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}

	lines := bytes.Split(f, []byte{'\n'})
	if len(lines) < 1 {
		return
	}

	fields := strings.Fields(string(lines[0]))
	if len(fields) == 0 || fields[0] != "cpu" {
		return
	}

	var idle, total uint64
	for i := 1; i < len(fields); i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			return
		}

		total += val          // tally up all the numbers to get total ticks
		if i == 4 || i == 5 { // idle is the 5th field in the cpu line
			idle += val
		}
	}

	if s.cpuIdleOld == 0 {
		s.cpuIdleOld = idle
		s.cpuTotalOld = total
	} else {
		s.cpuIdleOld = s.cpuIdleNew
		s.cpuTotalOld = s.cpuTotalNew
	}

	s.cpuIdleNew = idle
	s.cpuTotalNew = total
}

func (s *SysInfo) Processes() *Processes {
	return s.processes
}
