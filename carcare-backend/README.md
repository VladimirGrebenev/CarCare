### Примеры curl для auth

#### Регистрация
```sh
curl -X POST http://localhost:8080/auth/register \
	-H "Content-Type: application/json" \
	-d '{"email":"user@example.com","password":"123456"}'
```

#### Подтверждение email
```sh
curl -X POST http://localhost:8080/auth/confirm \
	-H "Content-Type: application/json" \
	-d '{"token":"CONFIRM_TOKEN"}'
```

#### Повторная отправка подтверждения
```sh
curl -X POST http://localhost:8080/auth/resend \
	-H "Content-Type: application/json" \
	-d '{"email":"user@example.com"}'
```

#### Вход (login)
```sh
curl -X POST http://localhost:8080/auth/login \
	-H "Content-Type: application/json" \
	-d '{"email":"user@example.com","password":"123456"}'
```

#### Восстановление пароля (forgot password)
```sh
curl -X POST http://localhost:8080/auth/forgot \
	-H "Content-Type: application/json" \
	-d '{"email":"user@example.com"}'
```

#### Сброс пароля (reset password)
```sh
curl -X POST http://localhost:8080/auth/reset \
	-H "Content-Type: application/json" \
	-d '{"token":"RESET_TOKEN","new_password":"newpass123"}'
```

#### Обновление токена (refresh)
```sh
curl -X POST http://localhost:8080/auth/refresh \
	-H "Content-Type: application/json" \
	-d '{"refresh_token":"REFRESH_TOKEN"}'
```

#### Выход (logout)
```sh
curl -X POST http://localhost:8080/auth/logout \
	-H "Content-Type: application/json" \
	-d '{"token":"ACCESS_OR_REFRESH_TOKEN"}'
```
# CARCARE Backend

Современный backend для приложения CARCARE на Go (Clean Architecture, DDD, TDD).

## Описание
Backend реализует бизнес-логику, API, интеграции и асинхронные задачи для учёта расходов, заправок, техобслуживания и штрафов.

- Язык: Go
- Архитектура: Clean Architecture, DDD
- Тестирование: TDD, unit/integration/e2e
- Контейнеризация: Docker, Docker Compose
- Документация: OpenAPI/Swagger

## Документация
- [PRD (Product Requirements Document)](../PRD_CARCARE.md)
- [Бэклог](../BACKLOG_CARCARE.md)

## Автоматизация миграций

Миграции хранятся в `internal/infrastructure/migration/` и применяются автоматически:

- **Локально:**
	1. Установите [golang-migrate](https://github.com/golang-migrate/migrate):
		 ```sh
		 curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz
		 sudo mv migrate /usr/local/bin/
		 ```
	2. Примените миграции:
		 ```sh
		 migrate -path internal/infrastructure/migration -database "postgres://user:password@localhost:5432/carcare?sslmode=disable" up
		 ```
- **В Docker:**
	- Миграции и бинарь migrate добавляются в образ, применяются при старте контейнера через entrypoint.
- **В CI/CD:**
	- GitHub Actions workflow `.github/workflows/ci.yml` автоматически применяет миграции при push/pull request.

Подробнее: см. [README.migrations.md](README.migrations.md)

## Структура каталогов
- `cmd/` — точки входа (main.go)
- `internal/` — бизнес-логика, use cases, адаптеры, инфраструктура
- `pkg/` — переиспользуемые пакеты
- `api/` — спецификации API
- `scripts/` — скрипты для CI/CD
- `test/` — интеграционные/e2e тесты
- `deployments/` — Docker, Compose, Kubernetes
- `docs/` — документация

## Быстрый старт
(будет дополнено после инициализации go.mod)

## Development

- Go 1.21+
- Run tests: `go test ./...`
- Lint: `golangci-lint run`
- Build: `go build ./cmd/carcare`

## REST API

### Примеры curl для users

#### Получить список пользователей
```sh
curl -X GET http://localhost:8080/users
```

#### Получить пользователя по ID
```sh
curl -X GET http://localhost:8080/users/{id}
```

#### Создать пользователя
```sh
curl -X POST http://localhost:8080/users \
	-H "Content-Type: application/json" \
	-d '{
		"email": "user@example.com",
		"name": "Ivan Ivanov",
		"role": "user"
	}'
```

#### Обновить пользователя
```sh
curl -X PUT http://localhost:8080/users/{id} \
	-H "Content-Type: application/json" \
	-d '{
		"email": "user@example.com",
		"name": "Ivan Ivanov",
		"role": "admin"
	}'
```

#### Удалить пользователя
```sh
curl -X DELETE http://localhost:8080/users/{id}
```

REST endpoints реализованы для всех доменных сущностей:

- /cars
- /users
- /fuel
- /maintenance
- /fines
- /reports

OpenAPI спецификация: [api/openapi.yaml](api/openapi.yaml)

Примеры запросов:

```sh
# Проверка здоровья
curl http://localhost:8080/health

# Получить список машин
curl http://localhost:8080/cars

# Получить список пользователей
curl http://localhost:8080/users

# Получить список заправок
curl http://localhost:8080/fuel

# Получить список ТО
curl http://localhost:8080/maintenance

# Получить список штрафов
curl http://localhost:8080/fines

# Получить список отчетов
curl http://localhost:8080/reports
```

## Репозитории и тесты

Для каждой сущности реализован репозиторий (скелет) и unit-тест (TDD, проверка not implemented):
- internal/adapter/repository/repository.go
- test/unit/*_repository_test.go

## TODO
- Реализовать бизнес-логику и интеграцию с БД для всех репозиториев и REST endpoints
- Expand OpenAPI spec
- Add more unit and integration tests

---

See [PRD_CARCARE.md](../PRD_CARCARE.md) for product requirements.
