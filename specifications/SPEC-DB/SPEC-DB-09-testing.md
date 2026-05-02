# SPEC-DB-09 — Testing Specification

> Part of [SPEC-DB](SPEC-DB-00-index.md) — Story 03: Database pgx Pool, Migrations & Initial Schema

---

## 8. Testing Specification

### 8.1 Test Philosophy

Story 03 introduces I/O-heavy code (pool, migrations). Unit tests are insufficient — **integration tests against a real Postgres instance are mandatory**. The compose Postgres is the target for all integration tests in this story.

Integration tests **shall** use a build tag `//go:build integration` and be skipped when `DATABASE_URL` is not set (so they don't block `go test ./...` without the DB running).

### 8.2 Test Matrix

| Test Name | SPEC-IDs | File | Description |
|-----------|----------|------|-------------|
| `TestLoad_AllRequired_ReturnsConfig` | SPEC-DB-001, SPEC-DB-002, SPEC-DB-003 | `internal/config/config_test.go` | Sets all env vars; asserts `Config` fields match; ensures no error. |
| `TestLoad_MissingDatabaseURL_ReturnsError` | SPEC-DB-002 | `internal/config/config_test.go` | Unsets `DATABASE_URL`; asserts `Load()` returns non-nil error mentioning the var name. |
| `TestLoad_MissingFirebaseProjectID_ReturnsError` | SPEC-DB-002 | `internal/config/config_test.go` | Unsets `FIREBASE_PROJECT_ID`; asserts error. |
| `TestLoad_MissingGoogleCredentials_ReturnsError` | SPEC-DB-002 | `internal/config/config_test.go` | Unsets `GOOGLE_APPLICATION_CREDENTIALS`; asserts error. |
| `TestOpen_WithRealPostgres_ReturnsPool` | SPEC-DB-004, SPEC-DB-005, SPEC-DB-006 | `internal/db/pgx_test.go` (integration) | Calls `db.Open` with real DSN; asserts pool is non-nil; pings; closes. |
| `TestUp_IsIdempotent` | SPEC-DB-008, SPEC-DB-009 | `internal/db/migrate_test.go` (integration) | Calls `db.Up` twice; asserts second call returns nil (not an error). |
| `TestDown_OneStep_RollsBack` | SPEC-DB-010 | `internal/db/migrate_test.go` (integration) | Calls `db.Up` then `db.Down(1)`; asserts `schema_migrations` table shows correct version. |
| `TestMigrationError_WrapsVersion` | SPEC-DB-011 | `internal/db/migrate_test.go` | Constructs `MigrationError{Version:1, Err: errFake}`; asserts `Error()` contains version; asserts `errors.AsType[*MigrationError](err)` works. |

### 8.3 Test Helpers

A `testDB(t)` helper **shall** be provided in `internal/db/testhelper_test.go`:

```go
// testDB opens a real pool against DATABASE_URL, skips if DATABASE_URL is unset,
// and registers t.Cleanup to close the pool.
func testDB(t *testing.T) *pgxpool.Pool {
    t.Helper()
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        t.Skip("DATABASE_URL not set")
    }
    pool, err := db.Open(t.Context(), config.Config{DatabaseURL: dsn})
    if err != nil {
        t.Fatalf("db.Open: %v", err)
    }
    t.Cleanup(pool.Close)
    return pool
}
```

### 8.4 Coverage Thresholds

| Package | Required coverage |
|---------|-----------------|
| `internal/config` | ≥ 90% (unit-testable without DB) |
| `internal/db` | ≥ 80% (integration tests required) |
