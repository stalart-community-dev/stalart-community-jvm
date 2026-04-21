# Stalart JVM Wrapper

[![eng](https://img.shields.io/badge/lang-English-blue)](README.en.md)
[![ru](https://img.shields.io/badge/lang-Russian-blue)](README.md)

Unofficial utility for optimizing JVM startup in Stalart on Windows.

## What this project does

- installs IFEO interception so `java.exe`/`javaw.exe` launches through `service.exe`;
- loads the active profile from `configs/*.json`;
- injects JVM flags for game launches;
- writes basic launch diagnostics to `logs/wrapper.log`;
- auto-selects recommended preset by hardware.

## How it works

1. User runs `cli.exe` and selects `Install`.
2. Windows registers `service.exe` as IFEO debugger for Java processes.
3. On game launch, Windows starts `service.exe` first.
4. `service.exe` validates launch context, applies profile, starts JVM, waits for exit.
5. On exit it writes process status and diagnostics.

Detailed flow: [docs/OVERVIEW.en.md](./docs/OVERVIEW.en.md).

## Quick start

1. Extract release files into `jvm_wrapper` next to Stalart launcher.
2. Run `cli.exe` as administrator.
3. Select `Install`.
4. Select `Apply Recommended Config`.
5. Launch the game.
6. Switch profiles later via `Select Config` if needed.

## Configs and presets

- Active profile key: `HKCU\\Software\\StalartWrapper`.
- Base profile: `configs/default.json`.
- Built-in presets: `compat`, `balanced`, `performance`, `ultra`.
- Default recommendation: `balanced`.
- Preset overview: [docs/PROFILES.en.md](./docs/PROFILES.en.md).
- JSON parameters: [docs/PARAMS.en.md](./docs/PARAMS.en.md).

## Hardware-based auto tune

- Global log: `logs/wrapper.log`.
- Auto-pick preset:
  - menu: `Apply Recommended Config`
  - CLI: `cli.exe --autotune`

Selection logic:
- `compat` for low-end systems,
- `balanced` for most systems,
- `performance` for stronger systems,
- `ultra` for high-end big-cache systems.

## Troubleshooting

See [docs/TROUBLESHOOTING.en.md](./docs/TROUBLESHOOTING.en.md).

## Build

```bash
go build -o build/cli.exe ./cmd/cli
go build -o build/service.exe ./cmd/service
```

## Disclaimer

- This utility is not affiliated with the game developer.
- Use at your own risk.
- Do not contact game support for this utility; use repository Issues instead.
