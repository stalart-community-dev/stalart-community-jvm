# Параметры конфигурации (`configs/*.json`)

Ниже перечислены основные поля пресетов и их смысл.

## Базовые

- `heap_size_gb` — размер heap (`-Xms/-Xmx`).
- `metaspace_mb` — `Metaspace` лимит.
- `pre_touch` — включает `AlwaysPreTouch`.

## G1GC

- `max_gc_pause_millis` — целевая пауза GC.
- `initiating_heap_occupancy_percent` — порог старта concurrent mark.
- `g1_heap_region_size_mb` — размер G1 region.
- `g1_new_size_percent`, `g1_max_new_size_percent` — рамки young generation.
- `g1_reserve_percent` — резерв памяти под продвижение объектов.

## Потоки и компиляция

- `parallel_gc_threads`, `conc_gc_threads` — число GC-потоков.
- `reserved_code_cache_size_mb` — размер code cache.
- `compile_threshold_scaling` — агрессивность JIT-компиляции.

## Прочие

- `use_string_deduplication` — дедупликация строк.
- `use_large_pages` — large pages (требует корректной системной настройки).

## Рекомендации

- Для повседневного использования начинайте с `balanced`.
- Меняйте 1-2 параметра за раз и сравнивайте через benchmark.
- Используйте реальные игровые прогоны, а не только короткие запуски.
