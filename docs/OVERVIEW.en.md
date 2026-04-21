# Architecture and Runtime Flow

## Components

- `cmd/cli` — interactive install/uninstall and profile selection UI.
- `cmd/service` — IFEO debugger, real entrypoint for game Java launch.
- `cmd/metrics-helper` — generator for `game_metrics.jsonl`.
- `internal/config` — profile loading/generation.
- `internal/jvm` — JVM flag building/filtering.
- `internal/telemetry` — process and game metrics collection.
- `internal/presetbench` — preset ranking and auto-selection.

## Launch flow

1. `cli.exe --install` registers IFEO interception.
2. Game starts Java runtime.
3. Windows launches `service.exe`.
4. `service.exe`:
   - detects game launch vs bootstrap;
   - loads active profile;
   - injects compatible flags;
   - starts child process;
   - collects telemetry until exit.
5. On exit:
   - writes `wrapper.log`;
   - appends JSONL event to `logs/presets/<preset>.jsonl`.

## Log safety

- raw JVM args are never written;
- user paths are redacted to `<user>`;
- `wrapper.log` size is bounded.

## Benchmark behavior

- requires at least 2 successful runs per preset;
- if game metrics exist (`fps/frame_time`), they are prioritized;
- otherwise soft fallback is used (`process cpu` + `wait_ms`);
- output includes confidence (`high/medium/low`).
