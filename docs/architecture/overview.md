# Architecture overview

## High-level

```
                        Mobile / Web (PWA)
                                │  Authorization: Bearer <Firebase ID token>
                                ▼
                ┌──────────────────────────────────┐
                │  CookSense backend (Go 1.26.2)   │
                │                                  │
                │  net/http router (stdlib)        │
                │  ├── public:  /api/health        │
                │  └── auth:    /api/recipes/...   │
                │               /api/reactions     │
                │               /api/me/recipes    │
                │               /api/lessons/...   │
                │                                  │
                │  Middleware chain:               │
                │   recover → request id → log →   │
                │   firebase auth (when required)  │
                │                                  │
                │  Domain packages:                │
                │   recipes / reactions / lessons  │
                │   users (lazy provisioning)      │
                │                                  │
                │  Persistence: pgx/v5 → Postgres  │
                │  Migrations: golang-migrate      │
                │  Seed loader: YAML → DB          │
                └──────────────┬──────────┬────────┘
                               │          │
                ┌──────────────▼──┐   ┌───▼──────────────┐
                │  PostgreSQL 17   │   │  Firebase Auth    │
                │  (docker-compose)│   │  (token verify)   │
                └──────────────────┘   └───────────────────┘
```

## Project layout

```
cmd/cooksense-server/main.go        ← entrypoint (config, db, auth, router, server)

internal/
  config/         env loading and validation (no defaults that hide bugs)
  httpx/          HTTP server, middleware, error responses
  auth/           Firebase Admin SDK init, bearer-token middleware, ctx accessors
  db/             pgx pool init, migration runner
  domain/         pure types: Recipe, Ingredient, Reaction, Lesson, User
  recipes/        handler.go, service.go, repo.go (per package, hexagonal-lite)
  reactions/      same shape
  lessons/        same shape
  users/          lazy provisioning + last_seen update
  seed/           YAML loader and `seed` subcommand

migrations/       0001_init.up.sql, 0001_init.down.sql, …
seed/recipes/     one YAML per recipe
seed/lessons/     one Markdown per lesson

Makefile          up, down, migrate, run, seed, test, lint
docker-compose.yml  Postgres 17 (already present)
modd.conf         hot reload for dev (already present)
```

## Layering rules

- `domain/` has **no imports from other internal packages**. Pure types and
  business invariants only.
- `*/repo.go` is the **only** place that talks to the database.
- `*/handler.go` is the **only** place that touches `net/http`.
- `*/service.go` orchestrates: it depends on a repo interface (defined in the
  same package) and is the unit-test target.
- `cmd/cooksense-server/main.go` is the **only** place that wires concrete
  implementations together.

## Cross-cutting

- **Logging**: `log/slog` (stdlib), JSON handler in production, text in dev,
  driven by `LOG_LEVEL` env var.
- **Errors**: every public function returns `error`. Handlers map domain
  errors → HTTP status via `httpx.WriteError`. No panics in normal paths.
- **Context**: every repo/service method takes `ctx context.Context`. Tests
  use `t.Context()`.
- **Config**: all configuration is read once at startup from env. No hidden
  defaults; missing required vars → fail fast.
- **Graceful shutdown**: `*http.Server.Shutdown` on SIGINT/SIGTERM, with a
  10-second drain.

## Performance budget (MVP)

- p50 < 50 ms, p95 < 150 ms on read endpoints, local Postgres.
- All discover/search queries hit indexes — no sequential scans on `recipes`.
- Seed load on boot is **opt-in** (`seed` subcommand), not on every start.
