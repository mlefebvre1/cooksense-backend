# SPEC-DB-04 — System Context & Dependencies

> Part of [SPEC-DB](SPEC-DB-00-index.md) — Story 03: Database pgx Pool, Migrations & Initial Schema

---

## 3. System Context & Dependencies

### 3.1 Runtime Requirements

| Requirement | Specification |
|-------------|--------------|
| **Go** | `1.26.2` |
| **PostgreSQL** | `17-alpine` (via Docker — `docker-compose.yml`) |
| **OS** | Linux, macOS |

### 3.2 New External Dependencies (SPEC-DB-025)

The following packages **shall** be added to `go.mod` and `go.sum`:

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/jackc/pgx/v5` | latest stable | PostgreSQL driver + pool |
| `github.com/golang-migrate/migrate/v4` | latest stable | Migration runner |
| `github.com/golang-migrate/migrate/v4/database/pgx/v5` | same as above | pgx v5 driver adapter for migrate |
| `github.com/golang-migrate/migrate/v4/source/file` | same as above | `file://` source adapter |

No other new external dependencies are permitted in Story 03.

### 3.3 Internal Package Dependencies

| Package | May import |
|---------|-----------|
| `internal/config` | stdlib only (`os`, `fmt`, `errors`) |
| `internal/db` | stdlib + `github.com/jackc/pgx/v5/pgxpool` + `github.com/golang-migrate/migrate/v4` |
| `cmd/cooksense-server` | `internal/config`, `internal/db` (plus stdlib) |

`internal/config` and `internal/db` **shall not** import any `internal/domain` or feature packages.

### 3.4 External Systems & APIs

| System | Interface | Notes |
|--------|-----------|-------|
| PostgreSQL 17 | TCP (`DATABASE_URL`) | Accessed via `pgxpool`; must be reachable at `db.Open` time. |
