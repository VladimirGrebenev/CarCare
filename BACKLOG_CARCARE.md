# CARCARE: Backlog по этапам Roadmap

## Этап 1: MVP

### Epic 1: Регистрация и авторизация
- **User Story:** Как пользователь, я хочу регистрироваться и входить через email, Яндекс, Google, чтобы быстро начать пользоваться сервисом.
  - **Acceptance Criteria:**
    - Регистрация/вход работают
    - Email подтверждается
    - Ошибки валидируются
    - Безопасность обеспечена
  - **Ответственные:** product-manager, javascript-expert, golang-expert
- **Tech Task:** Реализовать OAuth-интеграцию (Яндекс, Google)
  - **Acceptance Criteria:**
    - Успешный вход
    - Обработка ошибок
    - Логирование попыток
  - **Ответственные:** javascript-expert, golang-expert

### Epic 2: Профиль пользователя и авто
- **User Story:** Как пользователь, я хочу редактировать профиль и добавлять авто с фото, чтобы вести учёт расходов.
  - **Acceptance Criteria:**
    - Все поля валидируются
    - Фото ≤2 МБ, не более 3 фото
    - UX соответствует гайдлайнам
  - **Ответственные:** product-manager, javascript-expert, ui-ux-pro-max
- **Tech Task:** Ограничить бесплатный тариф одним авто, реализовать платную функцию добавления авто.
  - **Acceptance Criteria:**
    - Ограничение работает
    - Оплата активирует функцию
  - **Ответственные:** javascript-expert, golang-expert

### Epic 3: Основные разделы (Топливо, ТО, Штрафы, Отчёты)
- **User Story:** Как пользователь, я хочу вести расходы по топливу, ТО, штрафам и видеть отчёты.
  - **Acceptance Criteria:**
    - Добавление/редактирование записей
    - Фильтры, отчёты, визуализация, экспорт
  - **Ответственные:** product-manager, javascript-expert, golang-expert
- **Tech Task:** Реализовать кастомизацию виджетов на главной.
  - **Acceptance Criteria:**
    - Drag&drop, настройка порядка
    - Быстрый отклик
  - **Ответственные:** javascript-expert, ui-ux-pro-max

### Epic 4: CI/CD, тесты, документация
- **Tech Task:** Настроить CI/CD pipeline, покрыть тестами основные сценарии, подготовить документацию. Фронтенд реализовать на Svelte 5.x, SvelteKit 2.x, Vite 8.x, Node.js 20.19+, Vitest 4.x, Playwright, @sveltejs/vite-plugin-svelte.
  - **Acceptance Criteria:**
    - Все тесты проходят
    - Сборка apk
    - Документация актуальна
  - **Ответственные:** docker-expert, golang-expert, javascript-expert, python-expert

### Epic 5: Генерация apk для ручной установки
- **Tech Task:** Сборка apk через TWA/Bubblewrap, ручная установка. Использовать актуальные версии Svelte 5.x, SvelteKit 2.x, Vite 8.x, Node.js 20.19+, Vitest 4.x, Playwright, @sveltejs/vite-plugin-svelte для PWA.
  - **Acceptance Criteria:**
    - Apk доступен для тестирования
    - Инструкции по установке
  - **Ответственные:** javascript-expert

---

## Этап 2: Post-MVP

### Epic 1: Интеграции (Госуслуги, карты, оплата)
- **User Story:** Как пользователь, я хочу интеграцию с Госуслугами, картами и оплатой, чтобы автоматизировать учёт штрафов и оплат.
  - **Acceptance Criteria:**
    - Интеграции работают
    - Ошибки обрабатываются
    - Интерфейсы документированы
  - **Ответственные:** architecture-expert, golang-expert, javascript-expert
- **Tech Task:** Вынести интеграции в отдельные адаптеры (Hexagonal Architecture).
  - **Acceptance Criteria:**
    - Ядро не зависит от сервисов
    - Интеграции мокируются в тестах
  - **Ответственные:** architecture-expert, golang-expert

### Epic 2: Расширенные отчёты и фильтры
- **User Story:** Как пользователь, я хочу расширенные фильтры и новые виды отчётов.
  - **Acceptance Criteria:**
    - Фильтры по периоду, авто, типу расхода
    - Новые визуализации
  - **Ответственные:** product-manager, javascript-expert

### Epic 3: Публикация apk в RuStore
- **Tech Task:** Подготовить и опубликовать apk в RuStore.
  - **Acceptance Criteria:**
    - Apk опубликован
    - Проходит модерацию
    - Инструкции для пользователей
  - **Ответственные:** javascript-expert, product-manager

### Epic 4: Улучшения UI/UX, новые роли
- **User Story:** Как пользователь, я хочу улучшенный интерфейс и новые роли (например, семейный доступ).
  - **Acceptance Criteria:**
    - UI/UX соответствует гайдлайнам
    - Новые роли работают
  - **Ответственные:** javascript-expert, ui-ux-pro-max

### Epic 5: Performance/UX оптимизации
- **Tech Task:** Оптимизировать производительность и UX (lazy loading, кэширование, accessibility).
  - **Acceptance Criteria:**
    - Быстрый отклик
    - Отсутствие layout shift
    - Accessibility проверена
  - **Ответственные:** javascript-expert, ui-ux-pro-max

---

## Этап 3: Публикация (Google Play, масштабирование)

### Epic 1: Публикация apk в Google Play
- **Tech Task:** Подготовить apk к публикации, пройти модерацию Google Play.
  - **Acceptance Criteria:**
    - Apk опубликован
    - Соответствует требованиям Google
  - **Ответственные:** product-manager, javascript-expert

### Epic 2: Новые интеграции и расширение функционала
- **User Story:** Как пользователь, я хочу новые интеграции (например, страховые сервисы) и расширенный функционал.
  - **Acceptance Criteria:**
    - Интеграции реализованы через адаптеры
    - Тесты проходят
  - **Ответственные:** architecture-expert, golang-expert

### Epic 3: Масштабирование, мониторинг, SLA
- **Tech Task:** Обеспечить масштабирование, мониторинг, выполнение SLA.
  - **Acceptance Criteria:**
    - Система выдерживает нагрузку
    - Мониторинг работает
    - SLA выполняется
  - **Ответственные:** docker-expert, security-expert

---

**Назначение субагентов:**
- product-manager: детализация задач, приёмка, roadmap
- javascript-expert: frontend, PWA, apk, UI/UX
- golang-expert: backend, API, интеграции
- ui-ux-pro-max: дизайн, UX-ревью
- docker-expert: CI/CD, DevOps, масштабирование
- architecture-expert: архитектура, интеграции
- security-expert: безопасность, compliance
- python-expert: тесты, автоматизация
- code-review-excellence: code review, стандарты
