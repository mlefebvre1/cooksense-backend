# SPEC-DISCOVER — §8 Build, Tooling, Quality

[← Index](SPEC-DISCOVER-00-index.md)

## 8.1 Build

- `go build ./...` SHALL succeed after Story 07 lands.
- `go vet ./...` SHALL be clean (story DoD requirement).
- `golangci-lint run` SHALL be clean (no `//nolint` without a bracketed
  code and justification, per repo guideline).
- `go mod tidy` SHALL leave `go.mod` and `go.sum` in their tidy form. No
  new direct dependencies are introduced by this story.

## 8.2 Module dependency rules

| Rule | Enforcement |
|------|-------------|
| `internal/recipes/handler.go` imports `net/http`, `encoding/json`, `errors`, `log/slog`, `strconv`, `internal/auth`, and the package's own service — but **not** `pgx`/`pgxpool` | Architectural review + import grep in CI (see §8.4) |
| `internal/recipes/service.go` imports only `context`, `errors`, `fmt`, `log/slog` from stdlib (plus the package's own `Repo` interface) — and **not** `pgx`/`pgxpool`/`net/http` | Architectural review |
| `internal/recipes/repo.go` is the only file allowed to import `pgx`/`pgxpool` for this feature | Architectural review |
| `internal/recipes/dto.go` imports nothing beyond stdlib | Architectural review |
| `cmd/cooksense-server/main.go` imports `internal/recipes` (for the wiring trio) but the wiring SHALL NOT pass any package-level value | Architectural review |

## 8.3 Code-style guardrails specific to this story

- All new functions SHALL be ≤ 30 lines (Clean Code rule from repo
  guidelines). Long handlers SHALL be split (e.g., `writeJSON`,
  `writeError` helpers in `handler.go`).
- All new structs SHALL have ≤ 10 methods.
- All new files SHALL be ≤ 500 lines.
- Doc comments on every exported symbol SHALL cite the
  SPEC-DISCOVER-NNN ID(s) the symbol implements (backward traceability).
- `fmt.Print*` SHALL NOT appear anywhere in `internal/recipes/*` —
  structured logging via `slog` only.

## 8.4 Forbidden patterns (auto-reject in review)

- `_ = err` or any silent error swallowing in handler/service/repo.
- Package-level mutable state holding pool, service, or handler.
- Using a third-party HTTP framework (chi, echo, gin, fiber) in any
  Story-07 file.
- Manual JWT inspection in handler.go (auth is the middleware's job).
- `interface{}` instead of `any` (modern-Go guideline).
- `for i := 0; i < n; i++` patterns (modern-Go guideline: `for i := range n`).
- `time.Sleep` (cancellation-unsafe) anywhere.
- Returning `null` for empty array fields (per SPEC-DISCOVER-005, 027).

## 8.5 Performance budget

- The discover query SHALL complete in **< 50 ms p95** wall-clock against
  the local compose Postgres for catalog sizes ≤ 50 recipes (Story 06's
  curated baseline). Implementations exceeding this budget SHALL be
  profiled before merge.
- The detail endpoint SHALL complete in **< 20 ms p95** under the same
  conditions (small read; one-or-two roundtrips).
- These are guidance numbers, not enforcement gates — they are spelt out
  here so future regressions are visible.

## 8.6 Lint hooks specific to this story

The following `golangci-lint` linters SHALL be active when reviewing
Story-07 code (already enabled by the repo config):

- `errcheck` — every `err` checked.
- `staticcheck` — modern-Go idiom enforcement.
- `gocyclo` (≤ 10) — handler complexity ceiling.
- `lll` (≤ 120) — line-length ceiling.
- `nilerr` — no `return nil` after a non-nil `err`.
- `bodyclose` — every `*http.Response.Body` closed (relevant to tests).
- `rowserrcheck` / `sqlclosecheck` — every `pgx.Rows` closed and `Err()`
  checked.
