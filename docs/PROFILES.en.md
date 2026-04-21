# Presets

## Preset goals

Presets provide ready JVM tuning sets for different priorities:

- `compat` — maximum compatibility;
- `balanced` — stable baseline for most PCs;
- `performance` — more aggressive throughput;
- `ultra` — high-end focus.

## Practical selection flow

1. Start with `balanced`.
2. Apply `Apply Recommended Config` (or `cli.exe --autotune`).
3. Optionally compare nearby presets manually in the same in-game scenario.

## Where to inspect results

- `logs/wrapper.log`

## Notes

- Start with hardware-recommended preset first.
- Do manual fine tuning only after long real gameplay tests.
