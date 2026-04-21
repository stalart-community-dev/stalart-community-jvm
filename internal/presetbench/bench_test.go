package presetbench

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEvaluateSoftFallbackConfidence(t *testing.T) {
	t.Setenv("STALART_PRESETBENCH_LOG_DIR", t.TempDir())
	base := os.Getenv("STALART_PRESETBENCH_LOG_DIR")
	presetsDir := filepath.Join(base, "presets")
	if err := os.MkdirAll(presetsDir, 0o755); err != nil {
		t.Fatalf("mkdir presets dir: %v", err)
	}

	balancedRows := strings.Join([]string{
		`{"ts":"2026-04-22T10:00:00Z","preset_name":"balanced","exit_code":0,"wait_ms":120000,"avg_process_cpu_pct":11.2,"game_metrics_detected":true,"game_samples":120,"game_avg_fps":126.3,"game_avg_frame_time_ms":7.9}`,
		`{"ts":"2026-04-22T10:10:00Z","preset_name":"balanced","exit_code":0,"wait_ms":118000,"avg_process_cpu_pct":11.0,"game_metrics_detected":true,"game_samples":118,"game_avg_fps":124.1,"game_avg_frame_time_ms":8.1}`,
	}, "\n") + "\n"
	perfRows := strings.Join([]string{
		`{"ts":"2026-04-22T10:00:00Z","preset_name":"performance","exit_code":0,"wait_ms":125000,"avg_process_cpu_pct":12.6,"game_metrics_detected":false,"game_samples":0}`,
		`{"ts":"2026-04-22T10:10:00Z","preset_name":"performance","exit_code":0,"wait_ms":124000,"avg_process_cpu_pct":12.4,"game_metrics_detected":false,"game_samples":0}`,
	}, "\n") + "\n"

	if err := os.WriteFile(filepath.Join(presetsDir, "balanced.jsonl"), []byte(balancedRows), 0o644); err != nil {
		t.Fatalf("write balanced rows: %v", err)
	}
	if err := os.WriteFile(filepath.Join(presetsDir, "performance.jsonl"), []byte(perfRows), 0o644); err != nil {
		t.Fatalf("write performance rows: %v", err)
	}

	scores, err := Evaluate([]string{"balanced", "performance"}, 2)
	if err != nil {
		t.Fatalf("evaluate failed: %v", err)
	}
	if len(scores) != 2 {
		t.Fatalf("expected 2 scores, got %d", len(scores))
	}

	// balanced should be first due to valid game metrics and no fallback penalty.
	if scores[0].Preset != "balanced" {
		t.Fatalf("expected balanced to rank first, got %s", scores[0].Preset)
	}
	if scores[0].Confidence != "high" {
		t.Fatalf("expected balanced confidence=high, got %s", scores[0].Confidence)
	}
	if scores[1].Confidence != "low" {
		t.Fatalf("expected performance confidence=low, got %s", scores[1].Confidence)
	}
}
