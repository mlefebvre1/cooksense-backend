# SPEC-REACTIONS — §8 Build, Tooling, Quality

[← Index](SPEC-REACTIONS-00-index.md)

## 8.1 Build

- `go build ./...` SHALL succeed after Story 08 lands.
- `go vet ./...` SHALL be clean.
- `golangci-lint run` SHALL be clean. No `//nolint` SHALL appear without a
  bracketed code (e.g. `//nolint:errcheck`) and a justification comment, per
  the repo's CLAUDE.md rule.
- `go mod tidy` SHALL leave `go.mod`/`go.sum` in their tidy form. **No new
  third-party dependency** is introduced.

## 8.2 Module dependency rules

| Rule | Enforcement |
|------|-------------|
| `internal/domain` imports stdlib only | Manual review + lint forbid-list |
| `internal/reactions` may import `internal/domain`, `pgx`, `pgxpool`, stdlib; not `internal/api` or `net/http` | Architectural review |
| `internal/api` may import `internal/domain`, `internal/reactions`, `internal/auth`, stdlib; not `pgx` directly | Architectural review |
| `cmd/cooksense-server` may import all `internal/*` packages | Architectural review |
| No new third-party dependency | `go mod` audit |

## 8.3 Code-style guardrails specific to this story

- All new functions SHALL be ≤ 30 lines. Decomposition is mandatory if the
  body grows past that limit (see CLAUDE.md "Clean Code decomposition rules").
- All new structs SHALL have ≤ 10 methods.
- All new files SHALL be ≤ 500 lines.
- Doc comments on every exported symbol SHALL cite the SPEC-REACTIONS-NNN
  IDs the symbol implements (backward traceability).
- No `fmt.Print*` calls. Logging uses `log/slog`.
- No `time.Sleep` anywhere; cancellation flows through `r.Context()`.
- `interface{}` is forbidden; use `any` (Go 1.26 idiom).
- `errors.Is` / `errors.AsType[T]` for error inspection (never `err == sentinel`).

## 8.4 Forbidden patterns (auto-reject in review)

- `_ = err` or any silent error swallowing in repo/service/handler.
- Package-level mutable state holding pool, service, or repo.
- Body-supplied `firebase_uid` (the handler MUST use the auth context).
- A second API surface that takes a numeric `recipe_id` (NG-1).
- Returning `nil` slices from `ListByKind` (use `make([]…, 0)` so JSON is `[]`).
- Logging the Firebase token, the full request body, or the verbatim
  `firebase_uid` at any log level.
- `context.Background()` or `context.TODO()` inside the request lifecycle.

## 8.5 Performance budget

The MVP catalog is small (≤ 50 recipes per user). The endpoints SHALL meet,
on the local compose Postgres:

| Endpoint | p95 latency budget |
|----------|--------------------|
| `POST /api/reactions` | < 30 ms |
| `DELETE /api/reactions/{slug}` | < 30 ms |
| `GET /api/me/recipes` (≤ 50 rows) | < 50 ms |

These budgets are advisory for MVP; they MAY be enforced by future load
tests but are NOT required to gate the PR.
