# Автоматизация миграций в CarCare Backend

## Автоматическое применение миграций

Миграции хранятся в каталоге `internal/infrastructure/migration/`.

### Локально

Для применения миграций используйте [golang-migrate](https://github.com/golang-migrate/migrate):

1. Установите migrate:
   ```sh
   curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz
   sudo mv migrate /usr/local/bin/
   ```
2. Примените миграции:
   ```sh
   migrate -path internal/infrastructure/migration -database "postgres://user:password@localhost:5432/carcare?sslmode=disable" up
   ```

### В Docker

1. Добавьте в Dockerfile:
   ```dockerfile
   # ...
   COPY --from=builder /app/internal/infrastructure/migration ./migration
   ADD https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz /tmp/migrate.tar.gz
   RUN tar -xzf /tmp/migrate.tar.gz -C /usr/local/bin && rm /tmp/migrate.tar.gz
   # ...
   ```
2. В `docker-entrypoint.sh`:
   ```sh
   ./migrate -path ./migration -database "$DATABASE_URL" up
   ```

### В CI/CD (GitHub Actions)

В репозитории настроен workflow `.github/workflows/ci.yml`, который автоматически применяет миграции при каждом push/pull request:

```yaml
- name: Install golang-migrate
  run: |
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz
    sudo mv migrate /usr/local/bin/
- name: Run DB migrations
  run: |
    migrate -path carcare-backend/internal/infrastructure/migration -database "postgres://postgres:postgres@localhost:5432/carcare?sslmode=disable" up
```

## Best Practices
- Все миграции — в одном каталоге, версионируются.
- Применение миграций — до запуска приложения.
- Используйте инструменты автоматизации (golang-migrate/dbmate) для production/dev/CI.
- Не храните миграции вне VCS.
- Не запускайте приложение без актуальных миграций.

---

Подробнее — см. [README.md](../README.md) и [migration/README.md](internal/infrastructure/migration/README.md).
