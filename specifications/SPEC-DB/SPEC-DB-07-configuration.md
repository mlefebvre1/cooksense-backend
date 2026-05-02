# SPEC-DB-07 — Configuration Specification

> Part of [SPEC-DB](SPEC-DB-00-index.md) — Story 03: Database pgx Pool, Migrations & Initial Schema

---

## 6. Configuration Specification

### 6.1 Environment Variables

The full variable table is specified in [SPEC-DB-06 §5.1](SPEC-DB-06-packages.md#spec-db-001-configload-reads-all-required-env-vars) and mirrors `docs/architecture/infra.md`.

| Variable | Required | Default | Used by |
|----------|----------|---------|---------|
| `DATABASE_URL` | **yes** | — | `db.Open`, `db.Up`, `db.Down` |
| `FIREBASE_PROJECT_ID` | **yes** | — | Story 04 (auth middleware) |
| `GOOGLE_APPLICATION_CREDENTIALS` | **yes** | — | Story 04 (auth middleware) |
| `APP_PORT` | no | `"8080"` | Story 07 (HTTP server) |
| `LOG_LEVEL` | no | `"info"` | All packages |
| `LOG_FORMAT` | no | `"text"` | All packages |

### 6.2 `DATABASE_URL` Format

The DSN **shall** follow the `postgres://` URL scheme:

```
postgres://<user>:<password>@<host>:<port>/<dbname>?sslmode=disable
```

Example for local development (matching `docker-compose.yml` defaults):

```
postgres://cooksense:cooksense@localhost:5432/cooksense?sslmode=disable
```

### 6.3 Loading Order

`config.Load()` **shall** read values strictly from environment variables. It **shall not** read `.env` files directly — that is the responsibility of the shell (`-include .env` in the `Makefile`).

### 6.4 Sensitive Values

`DATABASE_URL` contains a password. It **shall never** be logged at any level. Log the database host/port only (masked DSN), never the full string.
