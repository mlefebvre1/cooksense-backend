# SPEC-DB-B — Implementation Task Decomposition

> Part of [SPEC-DB](SPEC-DB-00-index.md) — Story 03: Database pgx Pool, Migrations & Initial Schema

---

## Appendix B — Implementation Task Decomposition

Ordered list of atomic tasks for Story 03. Each task **shall** pass `go build ./...` + `go vet ./...` before moving to the next.

| Task | SPEC-IDs | Description | Dependencies |
|------|----------|-------------|--------------|
| T-01 | SPEC-DB-025 | Add 4 dependencies to `go.mod` via `go get` + `go mod tidy` | Story 01 merged |
| T-02 | SPEC-DB-001, SPEC-DB-002, SPEC-DB-003, SPEC-DB-023 | Implement `internal/config/config.go`: `Config` struct + `Load()` | T-01 |
| T-03 | SPEC-DB-001, SPEC-DB-002 | Write `internal/config/config_test.go`: 4 unit tests | T-02 |
| T-04 | SPEC-DB-004, SPEC-DB-005, SPEC-DB-006, SPEC-DB-007, SPEC-DB-024 | Implement `internal/db/pgx.go`: `Open()` with pool defaults + ping | T-02 |
| T-05 | SPEC-DB-008, SPEC-DB-009, SPEC-DB-010, SPEC-DB-011 | Implement `internal/db/migrate.go`: `Up()`, `Down()`, `MigrationError` | T-01 |
| T-06 | SPEC-DB-014, SPEC-DB-015, SPEC-DB-016, SPEC-DB-017, SPEC-DB-018, SPEC-DB-019, SPEC-DB-020 | Write `migrations/0001_init.up.sql` (full DDL) | None |
| T-07 | SPEC-DB-021 | Write `migrations/0001_init.down.sql` (reverse teardown) | T-06 |
| T-08 | SPEC-DB-012, SPEC-DB-013 | Update `cmd/cooksense-server/main.go` with `migrate up` / `migrate down N` subcommand dispatch | T-04, T-05 |
| T-09 | SPEC-DB-004, SPEC-DB-005, SPEC-DB-006 | Write `internal/db/pgx_test.go` (integration test: `TestOpen_WithRealPostgres_ReturnsPool`) | T-04 |
| T-10 | SPEC-DB-008, SPEC-DB-009, SPEC-DB-010, SPEC-DB-011 | Write `internal/db/migrate_test.go` (integration tests: idempotency, down, error wrapping) | T-05, T-07 |
| T-11 | SPEC-DB-026, SPEC-DB-027 | Verify `go build ./...` and `go vet ./...` pass | T-01..T-10 |
| T-12 | SPEC-DB-009, SPEC-DB-012 | `make up && make migrate && make migrate` — confirm idempotency | T-08, T-11 |
| T-13 | SPEC-DB-014..SPEC-DB-020 | Inspect tables with `psql -c '\dt'`; confirm 6 tables + `reaction_kind` enum | T-12 |
| T-14 | SPEC-DB-001, SPEC-DB-002 | Run `go test ./internal/config/...` — 4 unit tests pass | T-03 |
| T-15 | SPEC-DB-004..SPEC-DB-011 | Run `go test -tags=integration ./internal/db/...` — integration tests pass | T-09, T-10, T-12 |
| T-16 | — | Run `golangci-lint run`; fix any violations | T-11 |
| T-17 | SPEC-DB-026, SPEC-DB-027 | Final `go build ./...` + `go vet ./...` + coverage ≥ 80% for `internal/config` and `internal/db` | T-14, T-15, T-16 |

---

*End of SPEC-DB — Story 03 Database pgx Pool, Migrations & Initial Schema — v1.0.0*
