# SPEC-DB-08 — Build, Tooling & Quality Specification

> Part of [SPEC-DB](SPEC-DB-00-index.md) — Story 03: Database pgx Pool, Migrations & Initial Schema  
> SPEC-IDs covered: SPEC-DB-025, SPEC-DB-026, SPEC-DB-027

---

## 7. Build, Tooling & Quality Specification

### 7.1 Build Verification

| SPEC-ID | Command | Expected outcome |
|---------|---------|-----------------|
| SPEC-DB-025 | `go mod tidy` | `go.mod` and `go.sum` updated with the 4 new dependencies |
| SPEC-DB-026 | `go build ./...` | Exits `0`; no errors |
| SPEC-DB-027 | `go vet ./...` | Exits `0`; zero issues |

### 7.2 Dependencies Addition

After adding the new imports, run:

```bash
go get github.com/jackc/pgx/v5
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/pgx/v5
go get github.com/golang-migrate/migrate/v4/source/file
go mod tidy
```

### 7.3 Linting

`golangci-lint run` **should** produce zero violations after Story 03. Particularly watch for:

- `unused` — no exported symbols should be defined but never used.
- `errcheck` — all errors from `pgxpool` and `migrate` must be handled.
- `gosec` — DSN must not appear in log messages.

### 7.4 Integration Test (Smoke)

A `TestOpen_WithRealPostgres` integration test **shall** be present in `internal/db/pgx_test.go`. See [SPEC-DB-09](SPEC-DB-09-testing.md) for the full testing spec.
