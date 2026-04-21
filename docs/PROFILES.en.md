# Optimization Profiles

Below are built-in optimization profiles with full JSON values.
Files are stored in `configs/` and available via `Select Config`.

## Extended table

| Parameter | `compat` | `balanced` (`default`) | `performance` | `ultra` | `jdk25_conservative_8g` |
|---|---:|---:|---:|---:|---:|
| `heap_size_gb` | 6 | 8 | 8 | 8 | 8 |
| `pre_touch` | false | true | true | true | true |
| `metaspace_mb` | 512 | 512 | 640 | 640 | 512 |
| `max_gc_pause_millis` | 60 | 50 | 35 | 30 | 100 |
| `g1_new_size_percent` | 23 | 23 | 30 | 35 | 5 |
| `g1_mixed_gc_count_target` | 3 | 3 | 4 | 4 | 3 |
| `initiating_heap_occupancy_percent` | 20 | 20 | 18 | 15 | 30 |
| `parallel_gc_threads` | 4 | 10 | 10 | 10 | 8 |
| `conc_gc_threads` | 2 | 5 | 5 | 5 | 4 |
| `soft_ref_lru_policy_ms_per_mb` | 25 | 25 | 40 | 50 | 25 |
| `reserved_code_cache_size_mb` | 400 | 400 | 512 | 512 | 512 |
| `max_inline_level` | 15 | 15 | 16 | 18 | 15 |
| `freq_inline_size` | 500 | 500 | 550 | 600 | 500 |
| `inline_small_code` | 4000 | 4000 | 4200 | 4500 | 4000 |
| `max_node_limit` | 240000 | 240000 | 260000 | 280000 | 240000 |
| `auto_box_cache_max` | 4096 | 4096 | 8192 | 8192 | 4096 |
| `compile_threshold_scaling` | 0.5 | 0.5 | 0.6 | 0.75 | 0.5 |

## Full profile values

- `configs/compat.json`
- `configs/balanced.json`
- `configs/default.json` (active default profile)
- `configs/performance.json`
- `configs/ultra.json`
- `configs/jdk25_conservative_8g.json`

## JDK 25 Benchmark Results (8 GB heap)

Source: `artifacts/perf/latest/report.md`, `balanced-score` ranking  
Formula: `0.5 * latency + 0.4 * throughput + 0.1 * stability`.

| PC Class | Recommended preset | Balanced score | Notes |
|---|---|---:|---|
| `low_end` | `balanced` | 84.19 | Best safe balance without Full GC events. |
| `mid_range` | `balanced` | 84.80 | Most consistent profile for day-to-day use. |
| `high_end` | `balanced` | 84.90 | Stable latency/throughput strategy without aggressive risk. |

> [!NOTE]
> `performance` can produce slightly better synthetic throughput, but it is filtered out by the safety gate (`FullGC > 0`). The default JDK 25 strategy remains `balanced`.

## JDK 25 Flag Strategy

| Flag/group | Impact | Risk | Decision |
|---|---|---|---|
| `UseG1GC`, `MaxGCPauseMillis`, `ParallelGCThreads`, `ConcGCThreads` | Core GC control | Low | Keep |
| `InitiatingHeapOccupancyPercent`, `G1ReservePercent` | Marking-cycle stability | Medium | Keep with moderate values |
| `AlwaysPreTouch`, `UseStringDeduplication` | Predictable memory + lower heap pressure | Low/Medium (startup) | Keep |
| Fixed `G1HeapRegionSize` | Can fight JVM auto-tuning | High | Avoid by default |
| Deep G1 refinement (`G1RSet...`, `G1ConcRefinement...`) | Niche tuning only | High portability risk | Do not use by default |

## Quick selection guide

- `compat` — highest launch compatibility/stability.
- `balanced` — default day-to-day profile.
- `performance` — FPS-oriented with moderate risk.
- `ultra` — maximum tuning for high-end hardware and testing.
