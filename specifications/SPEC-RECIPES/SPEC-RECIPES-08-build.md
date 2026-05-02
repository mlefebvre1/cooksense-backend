# SPEC-RECIPES — §8 Build, Tooling, Quality

[← Index](SPEC-RECIPES-00-index.md)

## 8.1 Build

- `go build ./...` SHALL succeed after Story 05 lands.
- `go vet ./...` SHALL be clean.
- `golangci-lint run` SHALL be clean (no `//nolint` without a bracketed code
  and justification, per repo guideline).
- `go mod tidy` SHALL leave `go.mod` and `go.sum` in their tidy form. The
  only new direct dependency is `gopkg.in/yaml.v3`.

## 8.2 Module dependency rules

| Rule | Enforcement |
|------|-------------|
| `internal/domain` imports stdlib only | Manual review + lint forbid-list |
| `internal/seed` may import `internal/domain` and `pgxpool`, not `internal/api` | Architectural review |
| `cmd/cooksense-server` may import `internal/seed`, `internal/db`, `internal/config` | Architectural review |
| No third-party YAML library other than `gopkg.in/yaml.v3` | `go mod` audit |

## 8.3 Code-style guardrails specific to this story

- All new functions SHALL be ≤ 30 lines (Clean Code rule from repo guidelines).
- All new structs SHALL have ≤ 10 methods.
- All new files SHALL be ≤ 500 lines.
- Doc comments on every exported symbol (per Go Doc convention) SHALL cite
  the SPEC-RECIPES-NNN ID(s) the symbol implements (backward traceability).
- No `fmt.Print*` calls except in the `runSeed` subcommand's success/error
  output path (which is the user-facing CLI surface, not logging).

## 8.4 Forbidden patterns (auto-reject in review)

- `_ = err` or any silent error swallowing in loader/store.
- Package-level mutable state holding pool, config, or recipes.
- Building YAML errors with `errors.New(fmt.Sprintf(...))` instead of
  `fmt.Errorf`.
- Manual JWT-style error formatting (use `errors.Join` for aggregates).
- `time.Sleep` (cancellation-unsafe) anywhere.

## 8.5 Performance budget

The MVP catalog is small (≤ 50 recipes). The store SHALL complete in
< 500 ms wall-clock against the local compose Postgres. Implementations
that batch inserts via `pgx.SendBatch` are encouraged but not required.
