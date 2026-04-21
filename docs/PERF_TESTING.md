# Тестирование пресетов JDK 25

Документ описывает воспроизводимую методику оценки пресетов JVM для `JDK 25` с целевой кучей `8 ГБ`.

## Цель и критерии

- Цель: выбрать пресет с **balanced** приоритетом (latency + throughput + стабильность).
- Единая формула ранжирования:
  - `BalancedScore = 0.5 * Latency + 0.4 * Throughput + 0.1 * Stability`
- Safety-gate:
  - если в классе ПК есть пресеты с `FullGC = 0`, пресеты с `FullGC > 0` не участвуют в выборе default.

## Матрица классов ПК

- `low_end`: 4-6 ядер, 8-16 ГБ RAM
- `mid_range`: 6-8 ядер, 16-32 ГБ RAM
- `high_end`: 8-16+ ядер, 32+ ГБ RAM

## KPI

- `p95_pause_ms`, `p99_pause_ms`
- `avg_fps`, `low_1_fps`
- `full_gc_count`
- `startup_ms`
- `gc_cpu_pct`

## Telemetry quality и fallback

- `service.exe` пишет событие пресета на каждый нормальный запуск (`exit_code = 0`).
- Если игровые метрики (`game_metrics.jsonl`) найдены и валидны, ранжирование использует `fps/frame_time` как основной сигнал.
- Если игровых метрик нет, включается soft fallback по `avg_process_cpu_pct` и `wait_ms`.
- Для таких результатов понижается confidence (`high/medium/low`) и это учитывается в итоговом выборе.

## Harness

Скрипты:
- `scripts/perf/run-profiles.ps1`
- `scripts/perf/parse-results.py`
- `scripts/perf/report.py`

Запуск:

```powershell
pwsh ./scripts/perf/run-profiles.ps1 -Mode both
```

Результат:
- `artifacts/perf/<timestamp>/rows.csv`
- `artifacts/perf/<timestamp>/rows.json`
- `artifacts/perf/<timestamp>/winners.json`
- `artifacts/perf/<timestamp>/report.md`

## Реальные игровые прогоны

Для режима `real` подготовьте CSV:

```csv
preset,pc_class,p95_pause_ms,p99_pause_ms,avg_fps,low_1_fps,full_gc_count,startup_ms,gc_cpu_pct
balanced,mid_range,72,96,130,99,0,2980,14
```

Путь по умолчанию:
- `artifacts/perf/real-runs.csv`

Шаблон автоматически создается в каждой папке прогона:
- `real-runs.template.csv`

## Принятый итог для JDK 25

Текущий winner для `low/mid/high` по balanced-стратегии:
- `balanced`

Причина:
- лучший суммарный балл среди безопасных (`FullGC = 0`) пресетов.
