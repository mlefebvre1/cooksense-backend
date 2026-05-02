# Story 11 — Integration tests skeleton

Status: TODO
Estimate: M

## User story

As a maintainer, I want a small but solid integration-test foundation so that
each new endpoint can be tested against a real Postgres without copy-pasting
boilerplate.

## Background

We rely on integration tests (real Postgres, real router, fake Firebase
verifier) more than on heavy mocking. This story provides the shared scaffold.

## Acceptance criteria

- [ ] `internal/testsupport/` package with:
  - `Postgres(t *testing.T) *pgxpool.Pool` — connects to the compose Postgres
    using `DATABASE_URL` from the env, runs migrations into a unique schema
    per test run (`cooksense_test_<timestamp>`), and registers a cleanup that
    drops the schema.
  - `Server(t *testing.T, pool, verifier) *httptest.Server` — wires the full
    router (using the real handlers from stories 07–10) with the provided
    pool and `fakeVerifier`. Returns the `*httptest.Server` and a base URL.
  - `Authed(t, baseURL, uid string) *http.Client` — returns a client whose
    `Authorization` header is preset to `Bearer <uid>` (the fake verifier
    maps token == uid in tests).
- [ ] One end-to-end smoke test (`internal/recipes/handler_e2e_test.go`)
      that:
  - Boots the server.
  - Seeds 2 recipes through the loader.
  - Calls `GET /api/recipes/discover` and asserts both come back.
  - Posts a `LIKE` reaction.
  - Calls `GET /api/me/recipes` and asserts the liked recipe appears.
- [ ] All tests use `t.Context()` for the request context.
- [ ] Tests skip with a clear message if `DATABASE_URL` is not set
      (so unit-test-only runs are not blocked).
- [ ] `make test` runs all tests; CI/local docs note that integration tests
      need `make up` first.

## Technical notes

- Per-test schema isolation is preferred over per-test database creation
  (faster and easier on Postgres).
- Avoid `testcontainers-go` for MVP — the existing compose Postgres is good
  enough and keeps deps lean. Reconsider when CI lands.
- Use `httptest.NewServer` rather than `httptest.NewRecorder` so the test
  exercises the real `http.Server` muxer behavior.
- Helpers must accept `t *testing.T` and call `t.Cleanup` for teardown — no
  global state.

## Out of scope

- Property-based tests.
- Performance/load tests.
- CI workflow files.

## Dependencies

- depends on: 03, 04
- blocks: —

## Definition of Done

- [ ] AC met.
- [ ] The smoke test runs green twice in a row (idempotency/cleanup proof).
- [ ] No leaked schemas after `make test` completes.
