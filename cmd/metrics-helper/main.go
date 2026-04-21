// Command metrics-helper writes runtime game metrics into JSONL so service.exe
// can ingest them and include aggregate telemetry in wrapper.log.
//
// Output format (one JSON object per line):
// {"fps":123.4,"frame_time_ms":8.1,"cpu_percent":42.0,"gpu_percent":67.5}
//
// FPS/frame_time/gpu_percent are read from PresentMon CSV when available.
// cpu_percent is always sampled from the game process.
package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"stalart-wrapper/internal/winapi"
)

var procGetProcessTimes = winapi.Kernel32.NewProc("GetProcessTimes")

const processQueryLimitedInformation = 0x1000

type filetime struct {
	LowDateTime  uint32
	HighDateTime uint32
}

type metric struct {
	FPS         float64 `json:"fps"`
	FrameTimeMS float64 `json:"frame_time_ms"`
	CPUPercent  float64 `json:"cpu_percent"`
	GPUPercent  float64 `json:"gpu_percent"`
}

func main() {
	outPath := flag.String("out", "logs/game_metrics.jsonl", "output JSONL path")
	presentMonCSV := flag.String("presentmon-csv", "", "optional PresentMon CSV path")
	interval := flag.Duration("interval", time.Second, "sampling interval")
	flag.Parse()

	if err := os.MkdirAll(filepath.Dir(*outPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "create output dir: %v\n", err)
		os.Exit(1)
	}
	f, err := os.OpenFile(*outPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open output: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	t := time.NewTicker(*interval)
	defer t.Stop()

	var prevProc100ns uint64
	var prevWall time.Time
	for range t.C {
		pid, err := detectGamePID()
		if err != nil || pid == 0 {
			continue
		}
		h, err := syscall.OpenProcess(processQueryLimitedInformation, false, pid)
		if err != nil {
			continue
		}
		now := time.Now()
		curProc100ns, err := processCPU100ns(h)
		if err != nil {
			_ = syscall.CloseHandle(h)
			continue
		}

		cpuPct := 0.0
		if !prevWall.IsZero() && prevProc100ns > 0 {
			dProcSec := float64(curProc100ns-prevProc100ns) / 1e7
			dWallSec := now.Sub(prevWall).Seconds()
			if dWallSec > 0 {
				cpuPct = (dProcSec / dWallSec) * 100.0
			}
		}
		prevProc100ns = curProc100ns
		prevWall = now
		_ = syscall.CloseHandle(h)

		m := metric{
			CPUPercent: round(cpuPct, 2),
		}
		if *presentMonCSV != "" {
			fps, ft, gpu, err := readPresentMonLast(*presentMonCSV)
			if err == nil {
				m.FPS = round(fps, 2)
				m.FrameTimeMS = round(ft, 2)
				m.GPUPercent = round(gpu, 2)
			}
		}

		b, err := json.Marshal(m)
		if err != nil {
			continue
		}
		if _, err := writer.WriteString(string(b) + "\n"); err != nil {
			continue
		}
		_ = writer.Flush()
	}
}

func detectGamePID() (uint32, error) {
	for _, n := range []string{"stalart.exe", "stalartw.exe"} {
		pid, err := pidByImageName(n)
		if err == nil && pid != 0 {
			return pid, nil
		}
	}
	return 0, errors.New("game process not found")
}

func pidByImageName(name string) (uint32, error) {
	cmd := exec.Command("tasklist", "/fo", "csv", "/nh", "/fi", "IMAGENAME eq "+name)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	r := csv.NewReader(strings.NewReader(string(out)))
	rec, err := r.Read()
	if err != nil || len(rec) < 2 {
		return 0, errors.New("no rows")
	}
	if strings.EqualFold(strings.TrimSpace(rec[0]), "INFO: No tasks are running which match the specified criteria.") {
		return 0, errors.New("not running")
	}
	pid, err := strconv.Atoi(strings.TrimSpace(rec[1]))
	if err != nil {
		return 0, err
	}
	return uint32(pid), nil
}

func processCPU100ns(h syscall.Handle) (uint64, error) {
	var create, exit, kernel, user filetime
	r, _, err := procGetProcessTimes.Call(
		uintptr(h),
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

func readPresentMonLast(path string) (fps, frameTime, gpu float64, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, 0, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	head, err := r.Read()
	if err != nil {
		return 0, 0, 0, err
	}
	var last []string
	for {
		rec, e := r.Read()
		if e != nil {
			break
		}
		last = rec
	}
	if len(last) == 0 {
		return 0, 0, 0, errors.New("no presentmon rows")
	}

	col := func(name string) int {
		for i, h := range head {
			if strings.EqualFold(strings.TrimSpace(h), name) {
				return i
			}
		}
		return -1
	}

	msIdx := col("msBetweenPresents")
	if msIdx < 0 {
		msIdx = col("msBetweenDisplayChange")
	}
	gpuIdx := col("gpuBusy")
	if gpuIdx < 0 {
		gpuIdx = col("gpu_time")
	}

	if msIdx >= 0 && msIdx < len(last) {
		frameTime, _ = strconv.ParseFloat(last[msIdx], 64)
		if frameTime > 0 {
			fps = 1000.0 / frameTime
		}
	}
	if gpuIdx >= 0 && gpuIdx < len(last) {
		gpu, _ = strconv.ParseFloat(last[gpuIdx], 64)
	}
	return fps, frameTime, gpu, nil
}

func round(v float64, places int) float64 {
	p := mathPow10(places)
	return float64(int64(v*p+0.5)) / p
}

func mathPow10(n int) float64 {
	p := 1.0
	for i := 0; i < n; i++ {
		p *= 10
	}
	return p
}
