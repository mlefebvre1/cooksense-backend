# SPEC-AUTH — Firebase ID Token Middleware + Lazy User Provisioning

**Story:** 04 — Firebase ID token middleware + lazy user provisioning
**Status:** Draft → Final
**Date:** 2026-05-02
**Authors:** sdd-spec-author

---

## Purpose

This specification governs the authentication middleware, the Firebase Admin SDK
initialization, the `auth.Verifier` interface and its implementations, the lazy
user-provisioning path, the health-check endpoints, and all associated tests.

---

## File map

| File | Contents |
|------|----------|
| [SPEC-AUTH-01-preamble](SPEC-AUTH-01-preamble.md) | AI constraints, authorship, traceability rules |
| [SPEC-AUTH-02-introduction](SPEC-AUTH-02-introduction.md) | Story relationship, SPEC-ID registry |
| [SPEC-AUTH-03-goals](SPEC-AUTH-03-goals.md) | Goals, non-goals, constraints |
| [SPEC-AUTH-04-context](SPEC-AUTH-04-context.md) | External dependencies, decision references |
| [SPEC-AUTH-05-architecture](SPEC-AUTH-05-architecture.md) | Token flow, package dependency graph, middleware lifecycle |
| [SPEC-AUTH-06-packages](SPEC-AUTH-06-packages.md) | All SPEC-AUTH-NNN requirements: Go signatures, error taxonomy, SQL |
| [SPEC-AUTH-07-configuration](SPEC-AUTH-07-configuration.md) | Environment variables, credential file rules |
| [SPEC-AUTH-08-build](SPEC-AUTH-08-build.md) | Build, lint, vet constraints |
| [SPEC-AUTH-09-testing](SPEC-AUTH-09-testing.md) | Test strategy, named tests, coverage targets |
| [SPEC-AUTH-10-documentation](SPEC-AUTH-10-documentation.md) | Doc comments, example files, README impact |
| [SPEC-AUTH-A-checklist](SPEC-AUTH-A-checklist.md) | Specification completeness checklist |
| [SPEC-AUTH-B-tasks](SPEC-AUTH-B-tasks.md) | Ordered implementation task list |

---

## SPEC-ID registry

| ID | Summary | Section |
|----|---------|---------|
| SPEC-AUTH-001 | Firebase Admin SDK initialized once at startup | §6.1 |
| SPEC-AUTH-002 | `FIREBASE_PROJECT_ID` read from config; missing → startup failure | §6.1 |
| SPEC-AUTH-003 | `GOOGLE_APPLICATION_CREDENTIALS` read from config; missing → startup failure | §6.1 |
| SPEC-AUTH-004 | `auth.Verifier` interface with single `Verify` method | §6.2 |
| SPEC-AUTH-005 | `firebaseVerifier` wraps `auth.Client.VerifyIDToken` | §6.2 |
| SPEC-AUTH-006 | `fakeVerifier` for tests; accepts `map[token]User` | §6.2 |
| SPEC-AUTH-007 | `auth.User` struct with `UID`, `Email`, `DisplayName` fields | §6.3 |
| SPEC-AUTH-008 | `auth.Middleware` signature: accepts `Verifier` + `users.Toucher` | §6.4 |
| SPEC-AUTH-009 | Missing `Authorization` header → `401 UNAUTHENTICATED` | §6.4 |
| SPEC-AUTH-010 | Malformed header (not `Bearer <token>`) → `401 UNAUTHENTICATED` | §6.4 |
| SPEC-AUTH-011 | Any `Verify` error → `401 UNAUTHENTICATED` | §6.4 |
| SPEC-AUTH-012 | `WWW-Authenticate: Bearer realm="cooksense"` set on every `401` | §6.4 |
| SPEC-AUTH-013 | On success: `users.Toucher.Touch(ctx, user)` called before handler | §6.4 |
| SPEC-AUTH-014 | On success: `auth.User` stored in `r.Context()` via typed key | §6.4 |
| SPEC-AUTH-015 | Middleware short-circuits before reading the request body | §6.4 |
| SPEC-AUTH-016 | `auth.UserFromContext(ctx) (User, bool)` is the sole retrieval function | §6.5 |
| SPEC-AUTH-017 | `users.Toucher` interface with single `Touch(ctx, User) error` method | §6.6 |
| SPEC-AUTH-018 | `internal/users/repo.go` implements `Toucher` against Postgres | §6.6 |
| SPEC-AUTH-019 | UPSERT SQL matches the canonical query in `docs/architecture/auth.md` | §6.6 |
| SPEC-AUTH-020 | `Touch` error is logged and request continues (non-fatal) | §6.6 |
| SPEC-AUTH-021 | `GET /api/health` public endpoint returns `{"status":"ok"}` | §6.7 |
| SPEC-AUTH-022 | `GET /api/health/me` auth-gated endpoint returns `{"uid":"…","email":"…"}` | §6.7 |
| SPEC-AUTH-023 | Error responses use the standard envelope `{"error":{"code":"…","message":"…"}}` | §6.8 |
| SPEC-AUTH-024 | `secrets/firebase-admin.json.example` committed; real file gitignored | §7.1 |
| SPEC-AUTH-025 | `.env.example` contains `FIREBASE_PROJECT_ID` and `GOOGLE_APPLICATION_CREDENTIALS` | §7.1 |
| SPEC-AUTH-026 | Unit tests use `fakeVerifier` — no network calls | §9.1 |
| SPEC-AUTH-027 | Five named unit test cases required (see §9) | §9.1 |
