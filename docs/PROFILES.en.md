# Presets

## Preset goals

Presets provide ready JVM tuning sets for different priorities:

- `compat` — maximum compatibility;
- `balanced` — stable baseline for most PCs;
- `performance` — more aggressive throughput;
- `ultra` — high-end focus.

## Practical selection flow

1. Start with `balanced`.
2. Run at least 2 real sessions per candidate preset.
3. Run `cli.exe --benchmark`.
4. Check ranking confidence.

## Where to inspect results

- `logs/presets/<preset>.jsonl`
- `logs/wrapper.log`

## Notes

- Short bootstrap launches do not represent real FPS behavior.
- Without game metrics, ranking still works but is less precise.
