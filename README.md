# Stalart JVM Wrapper

[![eng](https://img.shields.io/badge/lang-English-blue)](README.en.md)
[![ru](https://img.shields.io/badge/lang-Russian-blue)](README.md)

Неофициальная утилита для оптимизации запуска JVM в Stalart на Windows.

## Что делает проект

- устанавливает IFEO-перехват для запуска `java.exe`/`javaw.exe` через `service.exe`;
- загружает активный профиль из `configs/*.json`;
- подставляет JVM-флаги для игрового запуска;
- пишет базовый технический лог запуска в `logs/wrapper.log`;
- автоматически подбирает рекомендованный пресет по железу.

## Как это работает

**JVM (Java Virtual Machine)** — это среда выполнения, через которую работает [Stalart](https://stalart.net/).

1. Пользователь запускает `cli.exe` и выполняет `Install`.
2. Windows регистрирует `service.exe` как IFEO debugger для Java-процесса.
3. При запуске игры Windows сначала запускает `service.exe`.
4. `service.exe` проверяет сценарий запуска, применяет профиль, стартует JVM и ждет завершения процесса.
5. После выхода пишет код завершения и служебную диагностику.

Подробная схема: [docs/OVERVIEW.md](./docs/OVERVIEW.md).

## Быстрый старт

1. Распакуйте релиз в `jvm_wrapper` рядом с лаунчером Stalart.
2. Запустите `cli.exe` от имени администратора.
3. Выберите `Install`.
4. Выберите `Apply Recommended Config`.
5. Запустите игру.
6. При необходимости выберите другой профиль через `Select Config`.

## Конфиги и пресеты

- Активный профиль хранится в `HKCU\\Software\\StalartWrapper`.
- Базовый профиль: `configs/default.json`.
- Готовые пресеты: `compat`, `balanced`, `performance`, `ultra`.
- Рекомендованный по умолчанию: `balanced`.
- Таблица пресетов: [docs/PROFILES.md](./docs/PROFILES.md).
- Параметры JSON: [docs/PARAMS.md](./docs/PARAMS.md).

## Автоподбор по железу

- Общий лог: `logs/wrapper.log`.
- Автоподбор пресета:
  - меню: `Apply Recommended Config`
  - CLI: `cli.exe --autotune`

Логика выбора:
- `compat` для слабых систем (мало RAM/потоков),
- `balanced` для большинства конфигураций,
- `performance` для сильных систем,
- `ultra` для high-end с большим кешем/ресурсами.

## Troubleshooting

См. [docs/TROUBLESHOOTING.md](./docs/TROUBLESHOOTING.md).

## Сборка

```bash
go build -o build/cli.exe ./cmd/cli
go build -o build/service.exe ./cmd/service
```

## Дисклеймер

- Утилита не аффилирована с разработчиками игры.
- Используйте на свой риск.
- Не обращайтесь в техподдержку игры по вопросам этой утилиты — используйте Issues в репозитории.
