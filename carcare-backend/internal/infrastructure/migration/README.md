# Миграции и Seed-скрипты CarCare

## Структура
- Все миграции и seed-скрипты размещаются в каталоге `internal/infrastructure/migration/`.
- Пример:
  - `20240322_create_reports_table.sql` — миграция для создания таблицы reports
  - `20240322_seed_reports.sql` — пример наполнения таблицы reports
  - `20260322_create_core_tables.sql` — миграция для production-ready схемы: users, cars, fuel_events, maintenance_events, fines

## Применение миграций (PostgreSQL)

1. Перейдите в каталог с миграциями:
   ```sh
   cd carcare-backend/internal/infrastructure/migration
   ```
2. Выполните миграции (замените параметры подключения на свои):
  ```sh
  psql "postgres://user:password@localhost:5432/carcare" -f 20240322_create_reports_table.sql
  psql "postgres://user:password@localhost:5432/carcare" -f 20260322_create_core_tables.sql
  # Повторите для других миграций по необходимости
  ```

## Запуск seed-скрипта

1. После применения миграций выполните:
   ```sh
   psql "postgres://user:password@localhost:5432/carcare" -f 20240322_seed_reports.sql
   ```

## Общая схема размещения файлов

```
carcare-backend/
  internal/
    infrastructure/
      migration/
        20240322_create_reports_table.sql
        20240322_seed_reports.sql
        20260322_create_core_tables.sql
        README.md
```

> Для автоматизации миграций рекомендуется использовать инструменты вроде [golang-migrate](https://github.com/golang-migrate/migrate) или [dbmate](https://github.com/amacneil/dbmate).
