# SPEC-AUTH-B — Implementation Tasks

← [Index](SPEC-AUTH-00-index.md)

---

Ordered task list for Story 04. Each task is atomic — it can be reviewed and
merged independently if needed. Dependencies are noted explicitly.

| T# | Task | SPEC-IDs | Depends on | Description |
|----|------|----------|-----------|-------------|
| T-01 | Extend `config.Config` with Firebase fields | SPEC-AUTH-002, SPEC-AUTH-003 | Story 03 T-01 | Add `FirebaseProjectID` and `GoogleAppCredentials` to `internal/config/config.go`; add `os.Getenv` loading and required-field validation; update tests |
| T-02 | Add `firebase.google.com/go/v4` to `go.mod` | SPEC-AUTH-001 | T-01 | `go get firebase.google.com/go/v4@latest && go mod tidy` |
| T-03 | Create `internal/auth/user.go` | SPEC-AUTH-007 | T-02 | Define `User` struct and unexported `userKey{}` context key |
| T-04 | Create `internal/auth/verifier.go` | SPEC-AUTH-004, SPEC-AUTH-005, SPEC-AUTH-006 | T-03 | Define `Verifier` interface, `firebaseVerifier` (production), `FakeVerifier` (test helper) |
| T-05 | Create `internal/auth/firebase.go` | SPEC-AUTH-001, SPEC-AUTH-005 | T-04 | Implement `NewFirebaseApp` and `NewFirebaseVerifier`; use `option.WithCredentialsFile` and `firebase.Config{ProjectID}` |
| T-06 | Create `internal/users/repo.go` — `Toucher` interface | SPEC-AUTH-017 | T-03 | Define `Toucher` interface in `internal/users`; keep it in its own file or in `repo.go` top section |
| T-07 | Implement `users.Repo` and `Repo.Touch` | SPEC-AUTH-018, SPEC-AUTH-019 | T-06, Story 03 T-03 | Create `Repo` struct backed by `*pgxpool.Pool`; implement the UPSERT SQL exactly as specified |
| T-08 | Create `internal/auth/middleware.go` | SPEC-AUTH-008–SPEC-AUTH-016 | T-04, T-06 | Implement `Middleware` (header extraction, `Verify`, `Touch`, context injection, `WWW-Authenticate`) and `UserFromContext` |
| T-09 | Register `GET /api/health` (public) | SPEC-AUTH-021 | T-08 | Add handler in `main.go` or a dedicated `internal/httpx/health.go`; no auth, returns `{"status":"ok"}` |
| T-10 | Register `GET /api/health/me` (auth-gated) | SPEC-AUTH-022 | T-08, T-09 | Add handler wrapped in `auth.Middleware`; returns `{"uid":"…","email":"…"}` |
| T-11 | Wire auth in `main.go` | SPEC-AUTH-001, SPEC-AUTH-008 | T-05, T-07, T-09, T-10 | Call `NewFirebaseApp`, `NewFirebaseVerifier`, `users.NewRepo`; pass `Middleware` to protected routes |
| T-12 | Create `secrets/firebase-admin.json.example` | SPEC-AUTH-024 | — | Placeholder JSON with no real values; add `secrets/README.md` |
| T-13 | Verify `.env.example` has Firebase vars | SPEC-AUTH-025 | — | Confirm (already present); no code change needed |
| T-14 | Write middleware unit tests (5 required cases) | SPEC-AUTH-026, SPEC-AUTH-027 | T-08 | `internal/auth/middleware_test.go` — all 5 named cases from §9 using `FakeVerifier` and `noopToucher`/`errorToucher` |
| T-15 | Write additional coverage tests (health handlers) | SPEC-AUTH-021, SPEC-AUTH-022 | T-09, T-10, T-14 | `TestHealthHandler_Public_Returns200`, `TestHealthMeHandler_ValidToken_ReturnsUID`, `TestHealthMeHandler_MissingToken_Returns401` |
| T-16 | Run `go test ./...` — confirm ≥ 90% on `internal/auth` | SPEC-AUTH-026 | T-14, T-15 | `go test -cover ./internal/auth/...` must show ≥ 90% |
| T-17 | Run full validation checklist | all | T-16 | `go fmt ./...`, `golangci-lint run`, `go vet ./...`, `go build ./...` — all must pass |

---

## Critical path

```
T-01 → T-02 → T-03 → T-04 → T-05
                  └──→ T-06 → T-07
                  └──→ T-08 → T-09 → T-10 → T-11 → T-17
                             └──→ T-14 → T-15 → T-16
T-12 (independent)
T-13 (independent, confirm only)
```

Estimated total: **M** (matches Story 04 estimate). T-01 through T-08 form the
core; T-09 through T-11 wire everything into the running server; T-14–T-16
provide the test coverage gate.
