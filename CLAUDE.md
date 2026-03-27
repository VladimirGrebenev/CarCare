# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Running the Project

**One command to start everything:**
```bash
sudo docker-compose up --build
```

Services after startup:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- PostgreSQL: localhost:5432 (user: `postgres`, pass: `postgres`, db: `carcare`)

Other useful commands:
```bash
sudo docker-compose down -v        # Stop and remove all containers + volumes
sudo docker-compose logs -f backend  # Tail backend logs
sudo docker-compose ps             # Check service health
```

## Backend (Go)

**Module path:** `github.com/VladimirGrebenev/CarCare-backend`

**Build & run locally (outside Docker):**
```bash
cd carcare-backend
go mod download
go build -o carcare ./cmd/carcare
go run ./cmd/carcare
go test ./...
```

**Architecture: Clean Architecture (strict layer separation)**

```
cmd/carcare/main.go                    ← Entry point, DI wiring
internal/domain/{entity}/              ← Entities + Repository interfaces
internal/usecase/                      ← One struct per use case (e.g. AddCarUsecase)
internal/adapter/repository/           ← PostgreSQL implementations (database/sql, no ORM)
internal/adapter/rest/rest.go          ← HTTP handlers (net/http, no framework)
```

**DI pattern:** Wired bottom-up in `main.go`: `DB → Repositories → UsecaseContainer → Handlers → HTTP routes`.

**Handler pattern:** Each resource has a struct handler with injected use cases:
```go
type CarHandler struct {
    Add *usecase.AddCarUsecase
    // ...
}
func (h *CarHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { /* route by method */ }
```
New handlers must be registered with `http.Handle()` (not `HandleFunc`) since they implement `ServeHTTP`.

**Repository pattern:** One struct per domain, `db *sql.DB` field, raw SQL with `$1, $2...` placeholders.

**Use case pattern:** One struct per operation — `Execute()` method calls one repository method. No logic beyond delegation unless business rules require it.

**Database schema** is created automatically on first container start via `db/init.sql`. No migration tool — schema changes require recreating the volume (`docker-compose down -v`).

Tables: `users`, `cars`, `fuel_events`, `maintenance_events`, `fines`, `reports`.

## Frontend (Svelte 5 / SvelteKit)

**Build:**
```bash
cd carcare-frontend
npm ci
npm run dev          # Dev server with HMR
npm run build        # Production build (static output in /build)
npm run type-check   # svelte-check + TypeScript
npm run lint         # ESLint
```

**SvelteKit adapter:** `adapter-static` with `fallback: 'index.html'` — pure SPA, no SSR.

**Svelte 5 syntax only:** Use `$state()`, `$derived()`, `$effect()`, `$props()`, `{#snippet}`, `{@render}`. Never use Svelte 4 `export let`, `$:`, or `createEventDispatcher`.

**File-based routing:** Pages live in `src/routes/`. Protected pages must call `authGuard()` from `src/lib/authGuard.ts`.

**State management:** Svelte `writable`/`derived` stores in `src/stores/`. Follow the pattern in `src/stores/fuel.ts` — separate `list`, `loading`, `error`, `filters` writables + CRUD async functions.

**API layer:** All backend calls go through `src/lib/api.ts`. Uses `withAuthHeaders()` helper which reads token from `localStorage`. Returns empty arrays on error rather than throwing (graceful degradation).

**Design system:** Windows 11 / Fluent Design with glassmorphism. All tokens defined in `src/routes/+layout.css` as CSS variables (`--bg-base`, `--bg-layer`, `--accent`, `--border`, etc.). Dark theme is default; light theme applied via `.light` class on `<html>`. Toggle logic in `src/lib/theme.ts`.

**UI components** in `src/components/ui/`: Button (4 variants), Input, Card, Modal, Toast, Table, FAB, Loader, EmptyState. Use these — don't inline equivalent HTML.

**Auth:** Token stored in `localStorage` as `authToken`. `src/stores/auth.ts` exports `setAuth()`, `clearAuth()`, `bootstrapAuth()`. Auth is currently mock (backend does not validate tokens).

## Agent System

Specialized subagents in `.github/agents/`:

| Agent | Responsibility |
|---|---|
| `golang-expert` | Backend Go code |
| `javascript-expert` | Frontend Svelte/TS code |
| `database-expert` | SQL, schema, queries |
| `docker-expert` | Dockerfiles, Compose |
| `ui-ux-pro-max` | UI design, Svelte components |
| `architecture-expert` | Architectural decisions |
| `code-review-excellence` | Code review (Go, JS/TS, SQL) |
| `security-expert` | Security vulnerabilities |

**Always delegate tasks to the responsible subagent** rather than implementing directly. Each agent has defined skills in `.agents/skills/`.

## Known Issues (as of 2026-03-27)

- Auth is entirely mock: any email/password logs in, tokens are not validated on protected endpoints
- `POST /api/cars` requires `id` (UUID) from the client — backend doesn't generate IDs
- Empty list endpoints return `null` instead of `[]`
- `plate` field not persisted for cars (missing from SQL insert)
- No real JWT implementation
