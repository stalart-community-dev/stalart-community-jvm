# Профили оптимизации

Ниже приведены встроенные профили с полными значениями в JSON.
Файлы лежат в `configs/` и сразу доступны в `Select Config`.

## Расширенная таблица

| Параметр | `compat` | `balanced` (`default`) | `performance` | `ultra` |
|---|---:|---:|---:|---:|
| `heap_size_gb` | 6 | 8 | 8 | 8 |
| `pre_touch` | false | true | true | true |
| `metaspace_mb` | 512 | 512 | 640 | 640 |
| `max_gc_pause_millis` | 60 | 50 | 35 | 30 |
| `g1_new_size_percent` | 23 | 23 | 30 | 35 |
| `g1_mixed_gc_count_target` | 3 | 3 | 4 | 4 |
| `initiating_heap_occupancy_percent` | 20 | 20 | 18 | 15 |
| `parallel_gc_threads` | 4 | 10 | 10 | 10 |
| `conc_gc_threads` | 2 | 5 | 5 | 5 |
| `soft_ref_lru_policy_ms_per_mb` | 25 | 25 | 40 | 50 |
| `reserved_code_cache_size_mb` | 400 | 400 | 512 | 512 |
| `max_inline_level` | 15 | 15 | 16 | 18 |
| `freq_inline_size` | 500 | 500 | 550 | 600 |
| `inline_small_code` | 4000 | 4000 | 4200 | 4500 |
| `max_node_limit` | 240000 | 240000 | 260000 | 280000 |
| `auto_box_cache_max` | 4096 | 4096 | 8192 | 8192 |
| `compile_threshold_scaling` | 0.5 | 0.5 | 0.6 | 0.75 |

## Полные значения профилей

- `configs/compat.json`
- `configs/balanced.json`
- `configs/default.json` (рабочий дефолт)
- `configs/performance.json`
- `configs/ultra.json`

## Как выбирать

- `compat` — если приоритет: запуск и стабильность.
- `balanced` — стандартный режим на каждый день.
- `performance` — упор на FPS при хорошем балансе риска.
- `ultra` — максимальный тюнинг для мощного ПК и тестирования.
