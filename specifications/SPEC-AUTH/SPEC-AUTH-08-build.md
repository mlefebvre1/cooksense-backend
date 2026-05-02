# SPEC-AUTH-08 — Build, Lint & Quality

← [Index](SPEC-AUTH-00-index.md)

---

## 1. Build constraints

- `go build ./...` shall succeed after this story lands.
- No new build tags shall be introduced.
- The `firebase.google.com/go/v4` module shall be added to `go.mod` via
  `go get firebase.google.com/go/v4/...` followed by `go mod tidy`.

## 2. Lint rules

All existing `golangci-lint` rules apply without exception. Specific rules
relevant to this story:

| Rule | What it catches here |
|------|---------------------|
| `errcheck` | All `Touch` and `VerifyIDToken` errors must be handled |
| `contextcheck` | `context.Context` must be the first parameter of every I/O function |
| `gocritic` | No `interface{}` / `any` without justification |
| `revive` | Exported symbols must have Go Doc comments citing SPEC-IDs |
| `staticcheck` | No unused variables or dead code paths in middleware |

No `//nolint` directives shall be added without a bracketed SPEC-ID and a
justification comment.

## 3. `go vet` constraints

- `go vet ./...` shall produce zero warnings after this story.
- Particular attention: `copylocks` (the `userKey` struct is passed by value —
  correct); `httpmux` (route registration uses Go 1.22+ method+path syntax).

## 4. File size constraint

Each new file shall remain under 500 lines. The expected sizes:

| File | Estimated lines |
|------|----------------|
| `internal/auth/firebase.go` | ~40 |
| `internal/auth/verifier.go` | ~50 |
| `internal/auth/user.go` | ~20 |
| `internal/auth/middleware.go` | ~80 |
| `internal/users/repo.go` | ~50 |
| `internal/auth/middleware_test.go` | ~150 |

## 5. No `fmt.Print*` in production code

All log output shall use `log/slog`. This includes the `WARN` log in the
`Touch` error path (SPEC-AUTH-020) and the `DEBUG` log in the verifier error
path (SPEC-AUTH-011).
