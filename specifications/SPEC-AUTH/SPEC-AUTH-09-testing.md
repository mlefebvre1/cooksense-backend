# SPEC-AUTH-09 — Testing

← [Index](SPEC-AUTH-00-index.md)

---

## 1. Test strategy

**SPEC-AUTH-026** — All middleware unit tests shall use `FakeVerifier`
(SPEC-AUTH-006) and a `noopToucher` (implements `users.Toucher` with a
no-op `Touch`). No network calls to Firebase and no database shall be made
in unit tests.

**SPEC-AUTH-027** — The following five named test cases shall be implemented
in `internal/auth/middleware_test.go`:

| Test name | SPEC-IDs verified | Scenario | Expected outcome |
|-----------|------------------|----------|-----------------|
| `TestMiddleware_MissingAuthorizationHeader_Returns401` | SPEC-AUTH-009, SPEC-AUTH-012, SPEC-AUTH-023 | Request with no `Authorization` header | HTTP 401, `WWW-Authenticate` header present, body is standard error envelope with code `UNAUTHENTICATED` |
| `TestMiddleware_MalformedAuthorizationHeader_Returns401` | SPEC-AUTH-010, SPEC-AUTH-012, SPEC-AUTH-023 | Header present but not `Bearer <token>` (e.g. `"Basic abc"`, `"Bearer"`, `"bearer token"`) | HTTP 401, `WWW-Authenticate` header present |
| `TestMiddleware_ExpiredOrInvalidToken_Returns401` | SPEC-AUTH-011, SPEC-AUTH-012 | `FakeVerifier` returns an error for the given token | HTTP 401, downstream handler not called |
| `TestMiddleware_WrongAudienceToken_Returns401` | SPEC-AUTH-011 | `FakeVerifier` returns `fmt.Errorf("wrong audience")` | HTTP 401 |
| `TestMiddleware_ValidToken_CallsNextAndSetsContext` | SPEC-AUTH-013, SPEC-AUTH-014, SPEC-AUTH-016 | `FakeVerifier` returns a valid `User`; `Touch` is a no-op | HTTP 200 from downstream handler; `auth.UserFromContext` returns the expected user |

Additional tests that **should** be implemented (not strictly blocking DoD but
required to reach ≥ 90% coverage for new code):

| Test name | Covers |
|-----------|--------|
| `TestMiddleware_TouchError_ContinuesRequest` | SPEC-AUTH-020 — `Touch` returns error but request still reaches handler with 200 |
| `TestMiddleware_DoesNotReadBody` | SPEC-AUTH-015 — body untouched after a 401 |
| `TestHealthHandler_Public_Returns200` | SPEC-AUTH-021 — no token needed |
| `TestHealthMeHandler_ValidToken_ReturnsUID` | SPEC-AUTH-022 — returns uid + email |
| `TestHealthMeHandler_MissingToken_Returns401` | SPEC-AUTH-022, SPEC-AUTH-009 |

## 2. Test helpers pattern

```go
// noopToucher satisfies users.Toucher with no side effects.
type noopToucher struct{}
func (noopToucher) Touch(_ context.Context, _ auth.User) error { return nil }

// errorToucher returns a fixed error to verify SPEC-AUTH-020.
type errorToucher struct{ err error }
func (t errorToucher) Touch(_ context.Context, _ auth.User) error { return t.err }
```

Use `net/http/httptest.NewRecorder()` and `httptest.NewRequest()` for all
HTTP-layer assertions. Use `t.Context()` (Go 1.26) for any context needed
inside test functions.

## 3. Coverage targets

| Package | Minimum line coverage |
|---------|-----------------------|
| `internal/auth` | ≥ 90% |
| `internal/users` (new `repo.go`) | ≥ 80% (integration test via compose Postgres acceptable) |

## 4. Users repo integration test (optional for DoD, required for Story 11)

A smoke test for `users.Repo.Touch` may use the running compose Postgres
(started by `make up`). It shall live in `internal/users/repo_integration_test.go`
with build tag `//go:build integration`. Story 11 will formalize the
testcontainers pattern; for Story 04, running against the live compose
database is acceptable.
