# SPEC-DB-03 — Goals & Non-Goals

> Part of [SPEC-DB](SPEC-DB-00-index.md) — Story 03: Database pgx Pool, Migrations & Initial Schema

---

## 2. Goals & Non-Goals

### 2.1 Goals

| ID | Goal |
|----|------|
| G-1 | Provide a **managed connection pool** (`pgxpool`) that is initialized once at startup and shared via constructor injection. |
| G-2 | Provide **idempotent migrations** so that `make migrate` can be run any number of times safely. |
| G-3 | Deliver the **complete initial schema** (6 tables, 5 indexes, 1 enum) that all feature stories depend on. |
| G-4 | Fail **fast at startup** if `DATABASE_URL` or other required vars are missing — never silently proceed with zero values. |
| G-5 | Keep `internal/config` and `internal/db` **layer-pure** — no feature package imports, no domain logic. |

### 2.2 Non-Goals

| ID | Non-Goal |
|----|----------|
| NG-1 | **Per-package repositories** — feature repos (recipes, reactions, lessons) are introduced in stories 07–10. |
| NG-2 | **Read replicas / multi-region** — single primary Postgres for MVP. |
| NG-3 | **ORM** — we use `pgxpool` and raw SQL queries. No GORM, sqlboiler, or ent. |
| NG-4 | **`embed.FS` for migrations** — migrations are read from disk in MVP; embedded FS is a possible follow-up. |
| NG-5 | **Schema changes after 0001** — additional migrations belong to the feature stories that need them. |
