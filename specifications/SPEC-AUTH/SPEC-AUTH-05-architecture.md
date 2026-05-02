# SPEC-AUTH-05 — Architecture

← [Index](SPEC-AUTH-00-index.md)

---

## 1. Token flow (end-to-end)

```
Mobile / Web client
  │  Authorization: Bearer <Firebase ID token>
  ▼
net/http ServeMux
  │
  ▼
auth.Middleware(verifier, toucher)
  │
  ├─[missing or malformed header]──► 401 UNAUTHENTICATED (WWW-Authenticate header)
  │
  ├─[Verifier.Verify(ctx, token) error]
  │       └─[error]────────────────► 401 UNAUTHENTICATED (WWW-Authenticate header)
  │
  ├─[users.Toucher.Touch(ctx, user) error]
  │       └─[error]────────────────► slog.Warn (non-fatal, request continues)
  │
  ├─[ctx = context.WithValue(ctx, userKey{}, user)]
  │
  └─[next.ServeHTTP(w, r.WithContext(ctx))]
            │
            ▼
        Handler calls auth.UserFromContext(ctx)
```

## 2. Package dependency graph

```
internal/config ◄── internal/auth ──► (uses users.Toucher interface)
                                              │
                          internal/users ◄───┘ (implements Toucher)
                                  │
                          internal/db (pgxpool)
```

Clean-Architecture rule: `internal/auth` imports only `internal/config` and
the `users.Toucher` interface (passed in, not imported). It does **not** import
`internal/users` directly — this keeps the dependency direction clean and avoids
a circular import.

## 3. Middleware lifecycle

```
Server startup (main.go)
  1. config.Load() → cfg (reads FIREBASE_PROJECT_ID, GOOGLE_APPLICATION_CREDENTIALS)
  2. auth.NewFirebaseApp(ctx, cfg) → *firebase.App          [SPEC-AUTH-001]
  3. auth.NewFirebaseVerifier(app) → auth.Verifier          [SPEC-AUTH-005]
  4. users.NewRepo(pool) → users.Toucher                    [SPEC-AUTH-018]
  5. auth.Middleware(verifier, toucher) → http.Handler wrapper
  6. Wrap protected routes with the middleware

Per-request (protected route)
  1. Extract "Bearer <token>" from Authorization header
  2. verifier.Verify(ctx, token) → (User, error)
  3. toucher.Touch(ctx, user) — upsert + last_seen_at update
  4. r = r.WithContext(context.WithValue(r.Context(), userKey{}, user))
  5. next.ServeHTTP(w, r)
```

## 4. Context key design

The user value shall be stored using an unexported struct key to prevent
collisions with other packages:

```go
type userKey struct{}
```

`auth.UserFromContext` is the **only** exported function that reads this key.
Handlers shall never call `ctx.Value(userKey{})` directly.

## 5. Error response shape

All `401` responses shall conform to the standard error envelope:

```json
{
  "error": {
    "code": "UNAUTHENTICATED",
    "message": "missing or invalid Authorization header"
  }
}
```

HTTP status: `401`.
Header: `WWW-Authenticate: Bearer realm="cooksense"`.

## 6. Threat model (from `docs/architecture/auth.md`)

| Threat | Status | Mitigation |
|--------|--------|-----------|
| Token forgery | 🟢 Blocked | Firebase SDK signature verification |
| Token replay across users | 🟢 Blocked | UID from decoded token, never from caller |
| Expired token used | 🟢 Blocked | `exp` checked by `VerifyIDToken` |
| Wrong audience/project | 🟢 Blocked | `aud` checked against `FIREBASE_PROJECT_ID` |
| Stolen service account JSON | 🟡 Mitigated | Gitignored, `.example` only in repo |
| No rate limiting | 🔴 Open | Post-MVP follow-up |
