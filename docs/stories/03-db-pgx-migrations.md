# Story 03 — Database: pgx pool + migrations + 0001_init

Status: TODO
Estimate: M

## User story

As a developer, I want the server to talk to Postgres through a managed pool
and apply versioned SQL migrations so that schema changes are reproducible and
the runtime database access is fast and safe.

## Background

See `docs/architecture/data-model.md` for the full DDL of the initial schema
and `docs/architecture/overview.md` for the layering rules.

## Acceptance criteria

- [ ] `internal/config/config.go` loads `DATABASE_URL` and exposes it via a
      typed `Config` struct.
- [ ] `internal/db/pgx.go` exposes a `Open(ctx, cfg) (*pgxpool.Pool, error)`
      function that:
  - Sets sensible pool defaults (min 2, max 10, health check period 1m).
  - Pings the database on open and returns an error if unreachable.
- [ ] `internal/db/migrate.go` exposes `Up(ctx, dsn, migrationsDir) error`
      using `golang-migrate/migrate/v4` with the `pgx` driver and the `file://`
      source.
- [ ] A `migrate` subcommand on the server binary applies all up migrations:
      `cooksense-server migrate up` and `cooksense-server migrate down 1`.
- [ ] `migrations/0001_init.up.sql` matches the DDL in
      `docs/architecture/data-model.md` exactly.
- [ ] `migrations/0001_init.down.sql` drops everything created in
      `0001_init.up.sql` in reverse order, including the `reaction_kind` enum.
- [ ] `make up && make migrate` succeeds on a clean clone and produces all
      tables, indexes, and the enum type.

## Technical notes

- Use `github.com/jackc/pgx/v5/pgxpool` directly. Do not introduce an ORM.
- Migrations are read from the on-disk `migrations/` directory in MVP. Embedding
  via `embed.FS` is a possible follow-up but not required here.
- The pool must be passed via constructor injection to repositories — never
  reach for a global.
- Surface a single migration error type that includes the version number to
  ease debugging.

## Out of scope

- Per-package repositories — they are introduced in stories 07–10.
- Read replicas / multi-region.

## Dependencies

- depends on: 01, 02
- blocks: 04, 05, 07, 08, 09, 10, 11

## Definition of Done

- [ ] AC met.
- [ ] Smoke test: a Go test runs `db.Open` against the compose Postgres,
      pings it, and closes the pool.
- [ ] `make migrate` is idempotent (running it twice is a no-op).
