# Stalart JVM Wrapper

[![eng](https://img.shields.io/badge/lang-English-blue)](README.en.md)
[![ru](https://img.shields.io/badge/lang-Russian-blue)](README.md)

Unofficial utility for optimizing JVM startup in Stalart on Windows.

## What this project does

- installs IFEO interception so `java.exe`/`javaw.exe` launches through `service.exe`;
- loads the active profile from `configs/*.json`;
- injects JVM flags for game launches;
- records launch telemetry into `logs/wrapper.log` and `logs/presets/*.jsonl`;
- can auto-select the best preset via menu or `cli.exe --benchmark`.

## How it works

1. User runs `cli.exe` and selects `Install`.
2. Windows registers `service.exe` as IFEO debugger for Java processes.
3. On game launch, Windows starts `service.exe` first.
4. `service.exe` validates launch context, applies profile, starts JVM, waits for exit.
5. On exit it writes metrics and process exit status.

Detailed flow: [docs/OVERVIEW.en.md](./docs/OVERVIEW.en.md).

## Quick start

1. Extract release files into `jvm_wrapper` next to Stalart launcher.
2. Run `cli.exe` as administrator.
3. Select `Install`.
4. Launch the game.
5. Switch profiles later via `Select Config` if needed.

## Configs and presets

- Active profile key: `HKCU\\Software\\StalartWrapper`.
- Base profile: `configs/default.json`.
- Built-in presets: `compat`, `balanced`, `performance`, `ultra`.
- Preset overview: [docs/PROFILES.en.md](./docs/PROFILES.en.md).
- JSON parameters: [docs/PARAMS.en.md](./docs/PARAMS.en.md).

## Telemetry and benchmark

- Global log: `logs/wrapper.log`.
- Per-preset logs: `logs/presets/<preset>.jsonl`.
- Optional game metrics input: `logs/game_metrics.jsonl`.
- Auto-pick preset:
  - menu: `Benchmark Presets (auto-select best)`
  - CLI: `cli.exe --benchmark`

If game metrics are missing, ranking still works using soft fallback (CPU/wait) with lower confidence.

Benchmark methodology: [docs/PERF_TESTING.en.md](./docs/PERF_TESTING.en.md).

## Troubleshooting

See [docs/TROUBLESHOOTING.en.md](./docs/TROUBLESHOOTING.en.md).

## Build

```bash
go build -o build/cli.exe ./cmd/cli
go build -o build/service.exe ./cmd/service
go build -o build/metrics-helper.exe ./cmd/metrics-helper
```

## Disclaimer

- This utility is not affiliated with the game developer.
- Use at your own risk.
- Do not contact game support for this utility; use repository Issues instead.
