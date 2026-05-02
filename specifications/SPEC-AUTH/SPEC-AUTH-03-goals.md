# SPEC-AUTH-03 — Goals, Non-Goals & Constraints

← [Index](SPEC-AUTH-00-index.md)

---

## 1. Goals

1. **Secure every non-public route.** All HTTP handlers that touch user data
   shall be unreachable without a valid, non-expired Firebase ID token.
2. **Single initialization point.** The Firebase Admin SDK client shall be
   created exactly once at server startup and shared via dependency injection.
   Re-creation per request is explicitly forbidden.
3. **Testable without Firebase.** The `auth.Verifier` interface shall allow
   unit tests to inject a `fakeVerifier` so that no network call to Firebase
   is required in the test suite.
4. **Lazy user provisioning.** A `users` row shall be upserted on every
   authenticated request; no separate registration step shall be required.
5. **Consistent error shape.** All `401` responses shall use the standard error
   envelope defined in `docs/architecture/api.md`.
6. **Zero secrets in the repo.** The real service-account JSON shall never be
   committed; only a `.example` placeholder shall be tracked by git.

## 2. Non-goals

- Email verification gating — post-MVP.
- Custom Firebase claims or role-based authorization — post-MVP.
- Refresh-token handling — the Firebase mobile SDK handles that transparently.
- Firebase Emulator integration test — optional post-MVP.
- Rate limiting — tracked as a post-MVP follow-up (threat model §5).
- Multi-tenancy or multiple Firebase projects.

## 3. Constraints

| Constraint | Source |
|-----------|--------|
| Go 1.26.2 | D-0001 |
| `net/http` stdlib only (no framework) | D-0003 |
| `firebase.google.com/go/v4` for SDK | `docs/architecture/auth.md` |
| JWT **never** parsed manually | `docs/architecture/auth.md` |
| No global mutable state — inject via constructors | CLAUDE.md |
| `context.Context` as first param on all I/O functions | CLAUDE.md |
| `log/slog` for logging — no `fmt.Print*` | CLAUDE.md |
| Functions ≤ 30 lines; files ≤ 500 lines | CLAUDE.md |
