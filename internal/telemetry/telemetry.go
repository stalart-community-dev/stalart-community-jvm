package telemetry

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"stalart-wrapper/internal/winapi"
)

var (
	procGetProcessTimes      = winapi.Kernel32.NewProc("GetProcessTimes")
	procGetProcessMemoryInfo = syscall.NewLazyDLL("psapi.dll").NewProc("GetProcessMemoryInfo")
)

type filetime struct {
	LowDateTime  uint32
	HighDateTime uint32
}

type processMemoryCounters struct {
	CB                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uintptr
	WorkingSetSize             uintptr
	QuotaPeakPagedPoolUsage    uintptr
	QuotaPagedPoolUsage        uintptr
	QuotaPeakNonPagedPoolUsage uintptr
	QuotaNonPagedPoolUsage     uintptr
	PagefileUsage              uintptr
	PeakPagefileUsage          uintptr
}

type Snapshot struct {
	DurationSec          int64
	Samples              int
	AvgProcessCPUPercent float64
	PeakWorkingSetMB     float64

	GameSamples         int
	GameAvgFPS          float64
	GameAvgFrameTimeMS  float64
	GameAvgCPUPercent   float64
	GameAvgGPUPercent   float64
	GameMetricsSource   string
	GameMetricsDetected bool
}

type Sampler struct {
	handle syscall.Handle
	start  time.Time

	mu sync.Mutex
	sn Snapshot

	stopCh chan struct{}
	doneCh chan struct{}
}

func Start(handle syscall.Handle) *Sampler {
	s := &Sampler{
		handle: handle,
		start:  time.Now(),
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}
	go s.run()
	return s
}

func (s *Sampler) Stop() Snapshot {
	close(s.stopCh)
	<-s.doneCh
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sn.DurationSec = int64(time.Since(s.start).Seconds())
	return s.sn
}

func (s *Sampler) run() {
	defer close(s.doneCh)
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	lastWall := time.Now()
	lastProc, err := processCPUTime100ns(s.handle)
	if err != nil {
		return
	}
	numCPU := float64(runtime.NumCPU())

	for {
		select {
		case <-s.stopCh:
			return
		case <-t.C:
			now := time.Now()
			procNow, err := processCPUTime100ns(s.handle)
			if err != nil {
				continue
			}

			deltaProcSec := float64(procNow-lastProc) / 1e7
			deltaWallSec := now.Sub(lastWall).Seconds()
			if deltaWallSec <= 0 || numCPU <= 0 {
				lastWall = now
				lastProc = procNow
				continue
			}
			cpuPercent := (deltaProcSec / (deltaWallSec * numCPU)) * 100.0

			ws, peakWs, memErr := processWorkingSet(s.handle)
			if memErr != nil {
				ws, peakWs = 0, 0
			}

			s.mu.Lock()
			s.sn.Samples++
			n := float64(s.sn.Samples)
			s.sn.AvgProcessCPUPercent += (cpuPercent - s.sn.AvgProcessCPUPercent) / n
			peakMB := float64(peakWs) / (1024.0 * 1024.0)
			curMB := float64(ws) / (1024.0 * 1024.0)
			if peakMB > s.sn.PeakWorkingSetMB {
				s.sn.PeakWorkingSetMB = peakMB
			}
			if curMB > s.sn.PeakWorkingSetMB {
				s.sn.PeakWorkingSetMB = curMB
			}
			s.mu.Unlock()

			lastWall = now
			lastProc = procNow
		}
	}
}

type gameMetric struct {
	FPS         float64 `json:"fps"`
	FrameTimeMS float64 `json:"frame_time_ms"`
	CPUPercent  float64 `json:"cpu_percent"`
	GPUPercent  float64 `json:"gpu_percent"`
}

func (s *Sampler) MergeGameMetrics(exePath string) {
	path, ok := resolveGameMetricsPath(exePath)
	if !ok {
		return
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	// allow long lines
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1024*1024)

	var count int
	var fpsSum, ftSum, cpuSum, gpuSum float64
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var m gameMetric
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			continue
		}
		count++
		fpsSum += m.FPS
		ftSum += m.FrameTimeMS
		cpuSum += m.CPUPercent
		gpuSum += m.GPUPercent
	}
	if count == 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.sn.GameSamples = count
	s.sn.GameAvgFPS = fpsSum / float64(count)
	s.sn.GameAvgFrameTimeMS = ftSum / float64(count)
	s.sn.GameAvgCPUPercent = cpuSum / float64(count)
	s.sn.GameAvgGPUPercent = gpuSum / float64(count)
	s.sn.GameMetricsSource = path
	s.sn.GameMetricsDetected = true
}

func resolveGameMetricsPath(exePath string) (string, bool) {
	if p := strings.TrimSpace(os.Getenv("STALART_GAME_METRICS_FILE")); p != "" {
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	exeDir := filepath.Dir(exePath)
	candidates := []string{
		filepath.Join(exeDir, "logs", "game_metrics.jsonl"),
		filepath.Join(exeDir, "logs", "stalart_metrics.jsonl"),
		filepath.Join(exeDir, "game_metrics.jsonl"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	return "", false
}

func processCPUTime100ns(handle syscall.Handle) (uint64, error) {
	var create, exit, kernel, user filetime
	r, _, err := procGetProcessTimes.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&create)),
		uintptr(unsafe.Pointer(&exit)),
		uintptr(unsafe.Pointer(&kernel)),
		uintptr(unsafe.Pointer(&user)),
	)
	if r == 0 {
		return 0, err
	}
	k := (uint64(kernel.HighDateTime) << 32) | uint64(kernel.LowDateTime)
	u := (uint64(user.HighDateTime) << 32) | uint64(user.LowDateTime)
	return k + u, nil
}

func processWorkingSet(handle syscall.Handle) (workingSet uintptr, peakWorkingSet uintptr, err error) {
	var pmc processMemoryCounters
	pmc.CB = uint32(unsafe.Sizeof(pmc))
	r, _, callErr := procGetProcessMemoryInfo.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&pmc)),
		uintptr(pmc.CB),
	)
	if r == 0 {
		return 0, 0, callErr
	}
	return pmc.WorkingSetSize, pmc.PeakWorkingSetSize, nil
}

var ErrNoGameMetrics = errors.New("no game metrics found")
