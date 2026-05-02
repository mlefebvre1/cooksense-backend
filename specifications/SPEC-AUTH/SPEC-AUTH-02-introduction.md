# SPEC-AUTH-02 — Introduction

← [Index](SPEC-AUTH-00-index.md)

---

## 1. Story relationship

| Field | Value |
|-------|-------|
| Story | 04 — Firebase ID token middleware + lazy user provisioning |
| Story file | `docs/stories/04-firebase-auth-middleware.md` |
| Estimate | M |
| Depends on | Story 01 (bootstrap), Story 03 (db + migrations) |
| Blocks | Stories 07, 08, 09, 10 |
| Decision refs | D-0004 (Firebase Auth as sole identity provider), D-0007 (single shared Firebase project) |
| Architecture refs | `docs/architecture/auth.md`, `docs/architecture/api.md` |

## 2. Business intent

Every request that carries user-owned data (reactions, personalized recipe
feeds) must be tied to a verified Firebase identity. This story introduces
the middleware that performs that verification, provisions the local `users`
row on first contact, and exposes two health endpoints (public and
auth-gated) that unblock all downstream feature stories.

## 3. Scope summary

**In scope:**
- Firebase Admin SDK initialization (`internal/auth/firebase.go`).
- `auth.Verifier` interface + `firebaseVerifier` + `fakeVerifier`.
- `auth.User` value type and context contract.
- `auth.Middleware` HTTP middleware.
- `users.Toucher` interface + Postgres implementation.
- `GET /api/health` (public) and `GET /api/health/me` (auth-gated).
- Credential placeholder file and `.gitignore` rule.

**Out of scope (see §3 goals file):**
- Email verification gating, custom claims, role-based authorization.
- Refresh-token handling.
- Firebase Emulator integration tests.

## 4. SPEC-ID summary table

See the full registry in [SPEC-AUTH-00-index](SPEC-AUTH-00-index.md). The 27
IDs break down by area:

| Area | IDs | Count |
|------|-----|-------|
| Firebase SDK init | 001–003 | 3 |
| Verifier interface | 004–006 | 3 |
| User type | 007 | 1 |
| Middleware | 008–015 | 8 |
| Context accessor | 016 | 1 |
| Toucher / repo | 017–020 | 4 |
| Health endpoints | 021–022 | 2 |
| Error envelope | 023 | 1 |
| Credential files | 024–025 | 2 |
| Test rules | 026–027 | 2 |
| **Total** | | **27** |
