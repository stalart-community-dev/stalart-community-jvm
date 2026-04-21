# Performance Testing

## Goal

Compare JVM presets and choose the best one for a specific PC using a balanced view of:

- latency (frametime/pause),
- throughput (FPS),
- stability (CPU/wait behavior).

## Data sources

- `logs/presets/<preset>.jsonl` — primary ranking input;
- `logs/game_metrics.jsonl` — optional game metrics;
- `logs/wrapper.log` — launch diagnostics.

## Minimal test protocol

1. Pick 2-4 presets.
2. Run at least 2 full gameplay sessions per preset.
3. Execute `cli.exe --benchmark`.
4. Record top-1 result and confidence.

## Confidence interpretation

- `high` — enough game-metric coverage;
- `medium` — mixed data (partial fallback);
- `low` — mostly process-only fallback.

## Measurement hygiene

- same map/route scenario for each run;
- same graphics settings;
- close heavy background processes;
- avoid comparing based on very short sessions.
