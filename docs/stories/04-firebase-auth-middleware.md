# Story 04 — Firebase ID token middleware + lazy user provisioning

Status: TODO
Estimate: M

## User story

As an authenticated mobile/web user, I want my Firebase ID token to be
verified by the backend so that my requests are tied to my identity and my
data is private.

## Background

See `docs/architecture/auth.md` for the full token flow, verification rules,
and threat model. Decision references: D-0004, D-0007.

## Acceptance criteria

- [ ] `internal/auth/firebase.go` initializes a `*firebase.App` and
      `*auth.Client` once at startup, given `FIREBASE_PROJECT_ID` and
      `GOOGLE_APPLICATION_CREDENTIALS`.
- [ ] An interface `auth.Verifier` with a single method
      `Verify(ctx, idToken string) (auth.User, error)` exists, with two
      implementations:
  - `firebaseVerifier` — wraps `auth.Client.VerifyIDToken`.
  - `fakeVerifier` — used by tests, accepts a map of `token → User`.
- [ ] `auth.User` is the type defined in `docs/architecture/auth.md`.
- [ ] `auth.Middleware(v Verifier, users users.Toucher) func(http.Handler) http.Handler`:
  - Returns `401 UNAUTHENTICATED` (standard error envelope) for missing or
    malformed `Authorization` headers.
  - Returns `401` for any verifier error (expired, invalid signature, wrong
    audience, revoked).
  - On success, calls `users.Toucher.Touch(ctx, user)` (lazy provisioning,
    UPSERT on `firebase_uid`) and stores the user in `r.Context()`.
- [ ] `auth.UserFromContext(ctx) (User, bool)` is the only way handlers
      retrieve the user.
- [ ] `internal/users/repo.go` implements `Toucher` against Postgres using the
      UPSERT shown in `docs/architecture/auth.md`.
- [ ] A public `GET /api/health` endpoint exists (no auth) and returns
      `{"status":"ok"}`.
- [ ] An auth-gated `GET /api/health/me` endpoint exists and returns
      `{"uid":"…","email":"…"}` for the authenticated caller.
- [ ] Unit tests cover: missing header, malformed header, expired token,
      wrong audience, happy path. They use `fakeVerifier` — no network calls.

## Technical notes

- Do not call `VerifyIDTokenAndCheckRevoked` for read-only endpoints; plain
  `VerifyIDToken` is enough and avoids a Firebase round-trip.
- The middleware must short-circuit before reading the request body.
- Set `WWW-Authenticate: Bearer realm="cooksense"` on `401` responses.
- Standard error envelope (see `docs/architecture/api.md`) — do not invent
  ad-hoc shapes.

## Out of scope

- Email verification gating, custom claims, role-based authorization (post-MVP).
- Refresh-token handling — Firebase clients handle that.

## Dependencies

- depends on: 01, 03
- blocks: 07, 08, 09, 10

## Definition of Done

- [ ] AC met.
- [ ] Tests pass.
- [ ] `secrets/firebase-admin.json.example` is committed; the real file is
      gitignored.
- [ ] `.env.example` updated with `FIREBASE_PROJECT_ID` and
      `GOOGLE_APPLICATION_CREDENTIALS`.
