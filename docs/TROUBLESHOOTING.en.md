# Troubleshooting

## Quick diagnostics

1. Check `cli.exe --status`.
2. Confirm `logs/wrapper.log` is updating.
3. Confirm entries appear in `logs/presets/*.jsonl`.

## Expected messages

### `[START] Write param successful`

Normal signal that launch parameters were applied.

### `[START] Process exit with code 0`

Normal process termination (including manual game close).

## Common issues

### Game starts with no preset effect

- verify `Install` completed successfully;
- ensure antivirus did not remove/block `service.exe`;
- run `Uninstall` -> `Install` again.

### Benchmark says `not enough preset benchmark data`

- you need at least 2 successful runs per compared preset;
- use real game sessions, not only quick start/close cycles;
- enable `metrics-helper` if needed.

### No game metrics (`game_metrics_detected=false`)

- run `metrics-helper` and verify `game_metrics.jsonl` path;
- or set `STALART_GAME_METRICS_FILE`;
- make sure the file is actively appended during gameplay.

## When opening an issue

Attach:

- `logs/wrapper.log`;
- relevant `logs/presets/*.jsonl`;
- environment details (launcher/Steam/EGS, selected preset, reproduction steps).
