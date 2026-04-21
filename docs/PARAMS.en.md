# Configuration Parameters (`configs/*.json`)

Key preset fields and what they control.

## Core

- `heap_size_gb` — heap size (`-Xms/-Xmx`).
- `metaspace_mb` — metaspace limit.
- `pre_touch` — enables `AlwaysPreTouch`.

## G1GC

- `max_gc_pause_millis` — target GC pause.
- `initiating_heap_occupancy_percent` — concurrent mark trigger threshold.
- `g1_heap_region_size_mb` — G1 region size.
- `g1_new_size_percent`, `g1_max_new_size_percent` — young gen range.
- `g1_reserve_percent` — memory reserve for promotions.

## Threads and JIT

- `parallel_gc_threads`, `conc_gc_threads` — GC thread counts.
- `reserved_code_cache_size_mb` — code cache size.
- `compile_threshold_scaling` — JIT aggressiveness.

## Other

- `use_string_deduplication` — string dedup.
- `use_large_pages` — large pages mode (requires correct OS setup).

## Recommendations

- Start with `balanced` for daily use.
- Change 1-2 parameters at a time and compare by benchmark.
- Prefer real gameplay runs over short startup-only runs.
