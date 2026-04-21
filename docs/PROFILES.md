# Профили оптимизации

Ниже приведены встроенные профили с полными значениями в JSON.
Файлы лежат в `configs/` и сразу доступны в `Select Config`.

## Расширенная таблица

| Параметр | `compat` | `balanced` (`default`) | `performance` | `ultra` | `jdk25_conservative_8g` |
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

## Полные значения профилей

- `configs/compat.json`
- `configs/balanced.json`
- `configs/default.json` (рабочий дефолт)
- `configs/performance.json`
- `configs/ultra.json`
- `configs/jdk25_conservative_8g.json`

## Результаты тестирования JDK 25 (heap 8 ГБ)

Источник данных: `artifacts/perf/latest/report.md`, расчет `balanced-score`  
Формула: `0.5 * latency + 0.4 * throughput + 0.1 * stability`.

| Класс ПК | Рекомендованный пресет | Balanced score | Комментарий |
|---|---|---:|---|
| `low_end` | `balanced` | 84.19 | Лучший безопасный баланс без Full GC. |
| `mid_range` | `balanced` | 84.80 | Наиболее устойчивый профиль для ежедневного использования. |
| `high_end` | `balanced` | 84.90 | Стабильная latency/throughput стратегия без агрессивных рисков. |

> [!NOTE]
> По synthetic-метрике `performance` иногда даёт выше throughput, но срабатывает safety-gate из-за `FullGC > 0`. Для default-стратегии JDK 25 выбран `balanced`.

## Стратегия флагов JDK 25

| Флаг/группа | Эффект | Риск | Решение |
|---|---|---|---|
| `UseG1GC`, `MaxGCPauseMillis`, `ParallelGCThreads`, `ConcGCThreads` | Базовое GC-управление | Низкий | Оставить |
| `InitiatingHeapOccupancyPercent`, `G1ReservePercent` | Стабилизация цикла маркировки | Средний | Оставить умеренные значения |
| `AlwaysPreTouch`, `UseStringDeduplication` | Предсказуемость памяти и экономия heap | Низкий/средний (startup) | Оставить |
| Жесткая фиксация `G1HeapRegionSize` | Может ухудшить авто-тюнинг JVM | Высокий | Не фиксировать без необходимости |
| Deep G1 refinement (`G1RSet...`, `G1ConcRefinement...`) | Нишевая оптимизация | Высокий, слабая переносимость | Не использовать по умолчанию |

## Как выбирать

- `compat` — если приоритет: запуск и стабильность.
- `balanced` — стандартный режим на каждый день.
- `performance` — упор на FPS при хорошем балансе риска.
- `ultra` — максимальный тюнинг для мощного ПК и тестирования.
