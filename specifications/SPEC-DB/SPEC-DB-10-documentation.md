# SPEC-DB-10 — Documentation Specification

> Part of [SPEC-DB](SPEC-DB-00-index.md) — Story 03: Database pgx Pool, Migrations & Initial Schema

---

## 9. Documentation Specification

### 9.1 Go Doc Standards

Every exported symbol in `internal/config` and `internal/db` **shall** carry a Go Doc comment that:

1. Starts with the symbol name.
2. Cites the governing `SPEC-DB-NNN` IDs.
3. Describes the behaviour, not the implementation.

Example:

```go
// Open creates and validates a PostgreSQL connection pool using the DSN
// from cfg.DatabaseURL. It applies pool defaults (min 2, max 10, healthcheck 1m)
// and pings the database before returning.
// Returns an error if the database is unreachable.
//
// SPEC-DB-004, SPEC-DB-005, SPEC-DB-006
func Open(ctx context.Context, cfg config.Config, overrides ...func(*pgxpool.Config)) (*pgxpool.Pool, error)
```

### 9.2 SQL File Headers

Each migration file **shall** begin with a comment identifying the migration:

```sql
-- CookSense — migration 0001_init
-- Description: Creates the initial schema (users, recipes, ingredients,
--              recipe_ingredients, user_reactions, lesson_articles).
-- SPEC-DB-014 through SPEC-DB-020
```

### 9.3 README

Story 03 **shall not** modify `README.md` (deferred to Story 12). The PR description **shall** include the output of `make up && make migrate && psql -c '\dt'` confirming all 6 tables exist.
