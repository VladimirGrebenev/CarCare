# CARCARE Frontend

Frontend для CARCARE на Svelte 5 (SPA, PWA, TWA).

## Описание
Модерн-приложение для учёта расходов, заправок, ТО и штрафов. Поддержка мобильных устройств, PWA, генерация apk через TWA/Bubblewrap.

- **Стек:**
	- Node.js 20.19+
	- Svelte 5.x, SvelteKit 2.x
	- Vite 8.x
	- Vitest 4.x (unit/integration)
	- Playwright (e2e)
	- @sveltejs/vite-plugin-svelte
	- CI/CD: GitHub Actions
	- Docker

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
1. Установите зависимости:
	npm install
2. Запуск dev-сервера:
	npm run dev
3. Запуск unit/integration тестов:
	npm run test
	# или
	npx vitest run
4. Запуск e2e тестов:
	npx playwright test

## Автоматизация PWA-иконок и Lighthouse CI

- Скрипт: `generate-icons.js` (Node.js, canvas)
- Генерирует PNG-иконки (192x192, 512x512) для car, fuel, fines в стилях glassmorphism и minimalism
- Выход: `public/icons/{style}-maskable-{car|fuel|fines}-{size}x{size}.png`
- Использование:

```bash
npm install canvas
node generate-icons.js
```

## Lighthouse CI для PWA

- Workflow: `.github/workflows/lighthouse.yml`
- Проверяет installability, offline, manifest, maskable icons, splash, theme, accessibility, performance
- Конфиг: `lighthouserc.json` (жёсткие пороги для PWA)
- Запуск локально:

```bash
npm install -g @lhci/cli
npx lhci autorun --config=./lighthouserc.json
```

## Docker и CI/CD Best Practices

- Все автоматизации — в существующих каталогах, без новых .md-файлов
- Скрипты и workflow не засоряют корень репозитория
- Документация — только в этом README.md
- Для production-образов используйте multi-stage Dockerfile, .dockerignore, healthcheck, non-root user
- Для CI: кэшируйте node_modules, используйте артефакты, не храните секреты в коде

## Рекомендации по дальнейшей автоматизации

- Добавить генерацию splash screen (512x1024, 1024x1024)
- Интегрировать e2e-тесты Playwright в CI
- Проверять accessibility (axe, pa11y)
- Автоматизировать обновление service worker
- Публиковать Lighthouse-отчёты в GitHub Pages или S3
- Добавить секреты через GitHub Actions Secrets

## PWA (Progressive Web App)

### Возможности
- Установка на мобильные и десктопные устройства (installable)
- Оффлайн-режим (offline fallback)
- Кэширование статических ассетов
- Manifest с иконками (glassmorphism/minimalism, car/fuel/fines)
- Service Worker с автообновлением

### Как работает
1. Manifest подключён в public/manifest.webmanifest
2. Service Worker — public/service-worker.js (регистрируется автоматически)
3. Offline fallback — public/offline.html
4. Иконки — public/icons/

### Чек-лист PWA
- [x] Manifest с иконками, цветами, именем, start_url, display, lang
- [x] Service Worker: offline, кэширование, обновление
- [x] Offline fallback (offline.html)
- [x] Splash screen, theme color
- [x] Installability (add to home screen)
- [ ] Maskable icons (добавить SVG/PNG)
- [ ] Lighthouse 100% PWA

### Рекомендации по автоматизации
- Генерировать иконки из SVG (glassmorphism/minimalism, car/fuel/fines)
- Проверять PWA через Lighthouse (CI)
- Добавить тесты offline/online

## Structure
- `src/routes/` — Pages and endpoints
- `src/lib/` — Shared components, stores, utils
- `src/app.html` — HTML template
- `static/` — Static assets
- `tests/` — Unit and e2e tests (Vitest, Playwright)

## Development
- Install: `npm install`
- Dev: `npm run dev`
- Build: `npm run build`
- Preview: `npm run preview`
- Test: `npm run test`
- Lint: `npm run lint`

## TODO
- Implement all pages and API stubs
- Add PWA manifest and service worker
- Write unit and e2e tests

---

See [PRD_CARCARE.md](../PRD_CARCARE.md) for product requirements.
