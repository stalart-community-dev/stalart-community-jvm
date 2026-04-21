# Stalart JVM Wrapper

[![eng](https://img.shields.io/badge/lang-English-blue)](README.en.md)
[![ru](https://img.shields.io/badge/lang-Russian-blue)](README.md)

Неофициальная утилита для оптимизации запуска JVM в Stalart на Windows.

## Что делает проект

- устанавливает IFEO-перехват для запуска `java.exe`/`javaw.exe` через `service.exe`;
- загружает активный профиль из `configs/*.json`;
- подставляет JVM-флаги для игрового запуска;
- собирает телеметрию запуска и пишет ее в `logs/wrapper.log` и `logs/presets/*.jsonl`;
- позволяет выбрать лучший пресет через меню или `cli.exe --benchmark`.

## Как это работает

**JVM (Java Virtual Machine)** — это среда выполнения, через которую работает [Stalart](https://stalart.net/).

1. Пользователь запускает `cli.exe` и выполняет `Install`.
2. Windows регистрирует `service.exe` как IFEO debugger для Java-процесса.
3. При запуске игры Windows сначала запускает `service.exe`.
4. `service.exe` проверяет сценарий запуска, применяет профиль, стартует JVM и ждет завершения процесса.
5. После выхода пишет итоговые метрики и код завершения.

Подробная схема: [docs/OVERVIEW.md](./docs/OVERVIEW.md).

## Быстрый старт

1. Распакуйте релиз в `jvm_wrapper` рядом с лаунчером Stalart.
2. Запустите `cli.exe` от имени администратора.
3. Выберите `Install`.
4. Запустите игру.
5. При необходимости выберите другой профиль через `Select Config`.
6. При необходимости обычного старта игры выбирите `Uninstall`

## Конфиги и пресеты

- Активный профиль хранится в `HKCU\\Software\\StalartWrapper`.
- Базовый профиль: `configs/default.json`.
- Готовые пресеты: `compat`, `balanced`, `performance`, `ultra`.
- Таблица пресетов: [docs/PROFILES.md](./docs/PROFILES.md).
- Параметры JSON: [docs/PARAMS.md](./docs/PARAMS.md).

## Телеметрия и benchmark

- Общий лог: `logs/wrapper.log`.
- Логи по пресетам: `logs/presets/<preset>.jsonl`.
- Игра может отдавать дополнительные метрики через `logs/game_metrics.jsonl`.
- Автоподбор пресета:
  - меню: `Benchmark Presets (auto-select best)`
  - CLI: `cli.exe --benchmark`

Если игровых метрик пока нет, сравнение работает через soft fallback (CPU/wait), но с пониженным confidence.

Методика измерений: [docs/PERF_TESTING.md](./docs/PERF_TESTING.md).

## Troubleshooting

См. [docs/TROUBLESHOOTING.md](./docs/TROUBLESHOOTING.md).

## Сборка

```bash
go build -o build/cli.exe ./cmd/cli
go build -o build/service.exe ./cmd/service
go build -o build/metrics-helper.exe ./cmd/metrics-helper
```

## Дисклеймер

- Утилита не аффилирована с разработчиками игры.
- Используйте на свой риск.
- Не обращайтесь в техподдержку игры по вопросам этой утилиты — используйте Issues в репозитории или прямо обращайтесь к создателю (nyrokume).
