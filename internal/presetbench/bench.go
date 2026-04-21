package presetbench

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type runEvent struct {
	TS                 string  `json:"ts"`
	PresetName         string  `json:"preset_name"`
	WaitMS             int64   `json:"wait_ms"`
	ExitCode           int     `json:"exit_code"`
	GameMetricsDetected bool   `json:"game_metrics_detected"`
	GameSamples        int     `json:"game_samples"`
	GameAvgFPS         float64 `json:"game_avg_fps"`
	GameAvgFrameTimeMS float64 `json:"game_avg_frame_time_ms"`
	GameAvgCPUPercent  float64 `json:"game_avg_cpu_pct"`
	GameAvgGPUPercent  float64 `json:"game_avg_gpu_pct"`
	AvgProcessCPU      float64 `json:"avg_process_cpu_pct"`
	PeakWorkingSetMB   float64 `json:"peak_working_set_mb"`
}

type PresetScore struct {
	Preset         string
	Runs           int
	AvgFPS         float64
	AvgFrameTimeMS float64
	AvgCPUPercent  float64
	AvgGPUPercent  float64
	AvgWaitMS      float64
	BalancedScore  float64
}

var ErrNotEnoughData = errors.New("not enough preset benchmark data")

func logDir() string {
	self, err := os.Executable()
	if err != nil {
		return filepath.Join(".", "logs")
	}
	return filepath.Join(filepath.Dir(self), "logs")
}

func presetsDir() string {
	return filepath.Join(logDir(), "presets")
}

func readPresetEvents(preset string) ([]runEvent, error) {
	path := filepath.Join(presetsDir(), preset+".jsonl")
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []runEvent
	sc := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var e runEvent
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			continue
		}
		// consider only successful runs
		if e.ExitCode != 0 {
			continue
		}
		// keep only recent enough runs (14 days)
		if e.TS != "" {
			if ts, err := time.Parse(time.RFC3339, e.TS); err == nil {
				if time.Since(ts) > 14*24*time.Hour {
					continue
				}
			}
		}
		out = append(out, e)
	}
	return out, sc.Err()
}

func avg(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

func normalizeHigher(v, lo, hi float64) float64 {
	if hi <= lo {
		return 1
	}
	return (v - lo) / (hi - lo)
}

func normalizeLower(v, lo, hi float64) float64 {
	if hi <= lo {
		return 1
	}
	return (hi - v) / (hi - lo)
}

// Evaluate computes balanced scores for presets and returns ranking desc.
func Evaluate(presets []string, minRuns int) ([]PresetScore, error) {
	var scores []PresetScore
	for _, p := range presets {
		events, err := readPresetEvents(p)
		if err != nil || len(events) < minRuns {
			continue
		}
		fpsVals := make([]float64, 0, len(events))
		ftVals := make([]float64, 0, len(events))
		cpuVals := make([]float64, 0, len(events))
		gpuVals := make([]float64, 0, len(events))
		waitVals := make([]float64, 0, len(events))
		for _, e := range events {
			if e.GameMetricsDetected && e.GameSamples > 0 {
				if e.GameAvgFPS > 0 {
					fpsVals = append(fpsVals, e.GameAvgFPS)
				}
				if e.GameAvgFrameTimeMS > 0 {
					ftVals = append(ftVals, e.GameAvgFrameTimeMS)
				}
				if e.GameAvgCPUPercent > 0 {
					cpuVals = append(cpuVals, e.GameAvgCPUPercent)
				} else if e.AvgProcessCPU > 0 {
					cpuVals = append(cpuVals, e.AvgProcessCPU)
				}
				if e.GameAvgGPUPercent > 0 {
					gpuVals = append(gpuVals, e.GameAvgGPUPercent)
				}
			} else {
				if e.AvgProcessCPU > 0 {
					cpuVals = append(cpuVals, e.AvgProcessCPU)
				}
			}
			if e.WaitMS > 0 {
				waitVals = append(waitVals, float64(e.WaitMS))
			}
		}
		scores = append(scores, PresetScore{
			Preset:         p,
			Runs:           len(events),
			AvgFPS:         avg(fpsVals),
			AvgFrameTimeMS: avg(ftVals),
			AvgCPUPercent:  avg(cpuVals),
			AvgGPUPercent:  avg(gpuVals),
			AvgWaitMS:      avg(waitVals),
		})
	}
	if len(scores) < 2 {
		return nil, ErrNotEnoughData
	}

	loFPS, hiFPS := scores[0].AvgFPS, scores[0].AvgFPS
	loFT, hiFT := scores[0].AvgFrameTimeMS, scores[0].AvgFrameTimeMS
	loCPU, hiCPU := scores[0].AvgCPUPercent, scores[0].AvgCPUPercent
	loWait, hiWait := scores[0].AvgWaitMS, scores[0].AvgWaitMS
	for _, s := range scores[1:] {
		if s.AvgFPS < loFPS {
			loFPS = s.AvgFPS
		}
		if s.AvgFPS > hiFPS {
			hiFPS = s.AvgFPS
		}
		if s.AvgFrameTimeMS < loFT {
			loFT = s.AvgFrameTimeMS
		}
		if s.AvgFrameTimeMS > hiFT {
			hiFT = s.AvgFrameTimeMS
		}
		if s.AvgCPUPercent < loCPU {
			loCPU = s.AvgCPUPercent
		}
		if s.AvgCPUPercent > hiCPU {
			hiCPU = s.AvgCPUPercent
		}
		if s.AvgWaitMS < loWait {
			loWait = s.AvgWaitMS
		}
		if s.AvgWaitMS > hiWait {
			hiWait = s.AvgWaitMS
		}
	}

	for i := range scores {
		fpsScore := normalizeHigher(scores[i].AvgFPS, loFPS, hiFPS)
		ftScore := normalizeLower(scores[i].AvgFrameTimeMS, loFT, hiFT)
		cpuScore := normalizeLower(scores[i].AvgCPUPercent, loCPU, hiCPU)
		waitScore := normalizeLower(scores[i].AvgWaitMS, loWait, hiWait)
		latency := ftScore
		throughput := fpsScore
		stability := 0.6*cpuScore + 0.4*waitScore
		scores[i].BalancedScore = 100 * (0.5*latency + 0.4*throughput + 0.1*stability)
	}

	// sort descending
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].BalancedScore > scores[i].BalancedScore {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}
	return scores, nil
}

func (s PresetScore) String() string {
	return fmt.Sprintf(
		"%s: score=%.2f runs=%d fps=%.2f frame_ms=%.2f cpu=%.2f gpu=%.2f wait_ms=%.0f",
		s.Preset, s.BalancedScore, s.Runs, s.AvgFPS, s.AvgFrameTimeMS, s.AvgCPUPercent, s.AvgGPUPercent, s.AvgWaitMS,
	)
}
