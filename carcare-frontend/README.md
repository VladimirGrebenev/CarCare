# CARCARE Frontend

Frontend для CARCARE на Svelte 5 (SPA, PWA, TWA).

## Описание
Модерн-приложение для учёта расходов, заправок, ТО и штрафов. Поддержка мобильных устройств, PWA, генерация apk через TWA/Bubblewrap.

- Фреймворк: Svelte 5, SvelteKit
- Архитектура: SPA, модульная структура
- Тестирование: Vitest, Playwright
- CI/CD: GitHub Actions
- Контейнеризация: Docker

## Документация
- [PRD (Product Requirements Document)](../PRD_CARCARE.md)
- [Бэклог](../BACKLOG_CARCARE.md)

## Структура каталогов
- `public/` — статические файлы, manifest, service-worker
- `src/` — исходный код, маршруты, компоненты, виджеты, фичи
- `tests/` — unit и e2e тесты
- `.github/` — CI/CD workflows
- `docker/` — Dockerfile, nginx.conf

## Быстрый старт
(будет дополнено после инициализации package.json)
