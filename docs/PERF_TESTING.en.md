# JDK 25 Preset Testing

This document defines a reproducible method to evaluate JVM presets for `JDK 25` with target heap `8 GB`.

## Goal and ranking

- Goal: select a **balanced** strategy (latency + throughput + stability).
- Unified ranking formula:
  - `BalancedScore = 0.5 * Latency + 0.4 * Throughput + 0.1 * Stability`
- Safety gate:
  - if any preset in a PC class has `FullGC = 0`, presets with `FullGC > 0` are excluded from default selection.

## PC class matrix

- `low_end`: 4-6 cores, 8-16 GB RAM
- `mid_range`: 6-8 cores, 16-32 GB RAM
- `high_end`: 8-16+ cores, 32+ GB RAM

## KPI set

- `p95_pause_ms`, `p99_pause_ms`
- `avg_fps`, `low_1_fps`
- `full_gc_count`
- `startup_ms`
- `gc_cpu_pct`

## Telemetry quality and fallback

- `service.exe` writes a preset event for every normal wrapper run (`exit_code = 0`).
- If valid game metrics are present (`game_metrics.jsonl`), ranking prioritizes `fps/frame_time`.
- If game metrics are missing, soft fallback is applied using `avg_process_cpu_pct` and `wait_ms`.
- Confidence (`high/medium/low`) is downgraded for fallback-only data and shown in ranking output.

## Harness

Scripts:
- `scripts/perf/run-profiles.ps1`
- `scripts/perf/parse-results.py`
- `scripts/perf/report.py`

Run:

```powershell
pwsh ./scripts/perf/run-profiles.ps1 -Mode both
```

Artifacts:
- `artifacts/perf/<timestamp>/rows.csv`
- `artifacts/perf/<timestamp>/rows.json`
- `artifacts/perf/<timestamp>/winners.json`
- `artifacts/perf/<timestamp>/report.md`

## Real game runs

For `real` mode provide CSV input:

```csv
preset,pc_class,p95_pause_ms,p99_pause_ms,avg_fps,low_1_fps,full_gc_count,startup_ms,gc_cpu_pct
balanced,mid_range,72,96,130,99,0,2980,14
```

Default location:
- `artifacts/perf/real-runs.csv`

A template is auto-generated for each run:
- `real-runs.template.csv`

