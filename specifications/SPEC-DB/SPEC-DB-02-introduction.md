# SPEC-DB-02 — Introduction

> Part of [SPEC-DB](SPEC-DB-00-index.md) — Story 03: Database pgx Pool, Migrations & Initial Schema

---

## 1. Introduction

This specification covers Story 03: **wiring the PostgreSQL connection pool, migration runner, and initial schema** so that every subsequent story can talk to a real database. It defines three concrete deliverables:

1. **`internal/config`** — a `Config` struct and `Load()` function that reads env vars at startup.
2. **`internal/db`** — a `db.Open()` function wrapping `pgxpool`, and `db.Up()`/`db.Down()` migration helpers backed by `golang-migrate`.
3. **`migrations/0001_init.up.sql` / `0001_init.down.sql`** — the canonical DDL for the MVP schema.

### 1.1 Scope

This document specifies:

- The `config.Config` struct fields and `Load()` contract.
- The `db.Open`, `db.Up`, `db.Down`, and `MigrationError` API.
- The `migrate` CLI subcommand wired in `cmd/cooksense-server/main.go`.
- Full DDL for the 6 MVP tables, 5 explicit indexes, and the `reaction_kind` enum.
- The `0001_init.down.sql` reverse teardown.

### 1.2 Definitions

| Term | Definition |
|------|-----------|
| **DSN** | Data Source Name — the `DATABASE_URL` connection string in `postgres://...` format. |
| **pgxpool** | The managed connection pool from `github.com/jackc/pgx/v5/pgxpool`. |
| **golang-migrate** | The `github.com/golang-migrate/migrate/v4` migration runner. |
| **SPEC-DB-NNN** | Requirement identifier for the Database story. NNN is a zero-padded three-digit number. |
| **up migration** | SQL that creates / alters schema objects (applied in forward order). |
| **down migration** | SQL that reverses an up migration (applied in reverse order). |

### 1.3 Relationship to Other Stories

| Story | Relationship |
|-------|-------------|
| SPEC-BOOT | Story 01 must be merged first; `internal/db/doc.go` and `internal/config/doc.go` already exist. |
| SPEC-MAKE | Story 02 provides `make up` (starts Postgres) and `make migrate` (calls this code). |
| Stories 04–10 | All depend on this story for the pool and schema. |
