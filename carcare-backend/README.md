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
