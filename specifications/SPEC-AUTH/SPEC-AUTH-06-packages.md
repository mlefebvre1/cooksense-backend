# SPEC-AUTH-06 — Package Specifications

← [Index](SPEC-AUTH-00-index.md)

---

This file contains all normative SPEC-AUTH-NNN requirements with exact Go
signatures and SQL.

---

## 6.1 Firebase Admin SDK initialization — `internal/auth/firebase.go`

**SPEC-AUTH-001** — The Firebase Admin SDK (`*firebase.App` and
`*auth.Client`) shall be initialized exactly once, during server startup, by
calling `auth.NewFirebaseApp(ctx context.Context, cfg *config.Config)`.
The returned `*firebase.App` shall be stored and shared via dependency
injection; it shall never be created per-request.

**SPEC-AUTH-002** — `config.Config` shall expose `FirebaseProjectID string`.
If the value is empty at startup, `config.Load` shall return an error causing
the server to exit non-zero.

**SPEC-AUTH-003** — `config.Config` shall expose `GoogleAppCredentials string`
(the value of `GOOGLE_APPLICATION_CREDENTIALS`). If the value is empty at
startup, `config.Load` shall return an error causing the server to exit
non-zero.

```go
// NewFirebaseApp initializes the Firebase Admin SDK from cfg.
// SPEC-AUTH-001, SPEC-AUTH-002, SPEC-AUTH-003
func NewFirebaseApp(ctx context.Context, cfg *config.Config) (*firebase.App, error)

// NewFirebaseVerifier creates a Verifier backed by the given Firebase app.
// SPEC-AUTH-005
func NewFirebaseVerifier(app *firebase.App) (Verifier, error)
```

`NewFirebaseApp` shall use `option.WithCredentialsFile(cfg.GoogleAppCredentials)`
and set the project ID via `firebase.Config{ProjectID: cfg.FirebaseProjectID}`.

---

## 6.2 Verifier interface — `internal/auth/verifier.go`

**SPEC-AUTH-004** — The package shall define a `Verifier` interface with
exactly one method:

```go
// Verifier verifies a Firebase ID token and returns the caller's identity.
// SPEC-AUTH-004
type Verifier interface {
    Verify(ctx context.Context, idToken string) (User, error)
}
```

**SPEC-AUTH-005** — `firebaseVerifier` shall be the production implementation.
It wraps `auth.Client.VerifyIDToken` and maps the decoded `*auth.Token` fields
(`UID`, `Firebase.Identities["email"]`, claims) to a `User` value:

```go
// firebaseVerifier implements Verifier using the Firebase Admin SDK.
// SPEC-AUTH-005
type firebaseVerifier struct {
    client *auth.Client
}

func (v *firebaseVerifier) Verify(ctx context.Context, idToken string) (User, error)
```

`Verify` shall call `v.client.VerifyIDToken(ctx, idToken)`. Any error returned
by the SDK shall be propagated as-is; the middleware translates it to `401`.

**SPEC-AUTH-006** — `fakeVerifier` shall be an exported test helper (in
`internal/auth/fake_verifier_test.go` or an exported `_test` package) that
accepts a `map[string]User` and returns the matching `User` for a known token,
or an error for an unknown one:

```go
// FakeVerifier is an in-memory Verifier for use in tests only.
// SPEC-AUTH-006
type FakeVerifier struct {
    tokens map[string]User
}

func NewFakeVerifier(tokens map[string]User) *FakeVerifier
func (f *FakeVerifier) Verify(_ context.Context, idToken string) (User, error)
```

`FakeVerifier` shall return `fmt.Errorf("unknown token")` for tokens not in
the map.

---

## 6.3 User type — `internal/auth/user.go`

**SPEC-AUTH-007** — The package shall define `User` as a plain value struct
with no methods:

```go
// User represents an authenticated Firebase principal.
// SPEC-AUTH-007
type User struct {
    UID         string
    Email       string
    DisplayName string
}
```

The unexported context key shall be co-located in this file:

```go
type userKey struct{}
```

---

## 6.4 Middleware — `internal/auth/middleware.go`

**SPEC-AUTH-008** — The package shall expose `Middleware` with the following
signature:

```go
// Middleware returns an http.Handler wrapper that verifies Firebase ID tokens.
// SPEC-AUTH-008
func Middleware(v Verifier, toucher users.Toucher) func(http.Handler) http.Handler
```

`users.Toucher` is defined in `internal/users` (§6.6). `auth.Middleware`
accepts the interface, not the concrete type.

**SPEC-AUTH-009** — If the `Authorization` header is absent, the middleware
shall respond with HTTP `401` and the standard error envelope
(`code: "UNAUTHENTICATED"`) before calling `next.ServeHTTP`.

**SPEC-AUTH-010** — If the `Authorization` header is present but does not
match the pattern `Bearer <non-empty-token>`, the middleware shall respond
with HTTP `401` and the standard error envelope.

**SPEC-AUTH-011** — If `v.Verify(ctx, token)` returns any non-nil error, the
middleware shall respond with HTTP `401` and the standard error envelope. The
underlying error shall be logged at `DEBUG` level; it shall **not** be
included in the HTTP response body.

**SPEC-AUTH-012** — Every `401` response written by the middleware shall
include the response header:
```
WWW-Authenticate: Bearer realm="cooksense"
```

**SPEC-AUTH-013** — On a successful verification, the middleware shall call
`toucher.Touch(ctx, user)` before passing the request to `next`. This call
shall be synchronous (not in a goroutine).

**SPEC-AUTH-014** — On a successful verification, the middleware shall store
the `User` value in the request context using the unexported `userKey{}`:

```go
ctx := context.WithValue(r.Context(), userKey{}, user)
next.ServeHTTP(w, r.WithContext(ctx))
```

**SPEC-AUTH-015** — The middleware shall not read `r.Body` at any point. It
shall short-circuit (return the `401`) without consuming the request body.

---

## 6.5 Context accessor — `internal/auth/middleware.go`

**SPEC-AUTH-016** — The package shall expose `UserFromContext` as the sole
function for retrieving the authenticated user from a context:

```go
// UserFromContext retrieves the authenticated User from ctx.
// Returns (User{}, false) if no user is present.
// SPEC-AUTH-016
func UserFromContext(ctx context.Context) (User, bool)
```

No handler shall access `ctx.Value(userKey{})` directly.

---

## 6.6 Toucher interface and Postgres implementation

**SPEC-AUTH-017** — `internal/users` shall define the `Toucher` interface:

```go
// Toucher upserts the user row and updates last_seen_at.
// SPEC-AUTH-017
type Toucher interface {
    Touch(ctx context.Context, u auth.User) error
}
```

`Toucher` is defined here (consumer-side interface) to avoid importing
`internal/auth` from `internal/users`. The argument type `auth.User` is
acceptable because `internal/users` imports `internal/auth` (domain value).
No circular import is introduced.

**SPEC-AUTH-018** — `internal/users/repo.go` shall implement `Toucher` as
`Repo`, backed by a `*pgxpool.Pool` passed via the constructor:

```go
// Repo implements Toucher and other user-related persistence.
// SPEC-AUTH-018
type Repo struct {
    pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo
func (r *Repo) Touch(ctx context.Context, u auth.User) error
```

**SPEC-AUTH-019** — The SQL executed by `Touch` shall be exactly:

```sql
INSERT INTO users (firebase_uid, display_name, email)
VALUES ($1, $2, $3)
ON CONFLICT (firebase_uid) DO UPDATE
   SET last_seen_at  = now(),
       email         = COALESCE(EXCLUDED.email, users.email),
       display_name  = COALESCE(EXCLUDED.display_name, users.display_name);
```

No deviation from this query is permitted without a spec amendment.

**SPEC-AUTH-020** — If `Touch` returns a non-nil error, the middleware shall:
1. Log the error at `WARN` level with `slog.Warn`, including the `firebase_uid`
   as a structured attribute.
2. **Continue** serving the request — provisioning failure shall not block the
   authenticated user.

---

## 6.7 Health endpoints — registered in `main.go` or a dedicated handler file

**SPEC-AUTH-021** — `GET /api/health` shall be a public endpoint (no
middleware applied). It shall return HTTP `200` with:

```json
{ "status": "ok" }
```

No auth check. No database call.

**SPEC-AUTH-022** — `GET /api/health/me` shall be protected by
`auth.Middleware`. It shall return HTTP `200` with:

```json
{ "uid": "<user.UID>", "email": "<user.Email>" }
```

If the middleware blocks the request (missing/invalid token), the standard
`401` response applies.

---

## 6.8 Error envelope

**SPEC-AUTH-023** — Every error response produced by the middleware or health
handlers shall use the standard envelope defined in `docs/architecture/api.md`:

```json
{ "error": { "code": "<ERROR_CODE>", "message": "<human readable>" } }
```

Content-Type shall be `application/json`. The HTTP status shall be set before
writing the body. For authentication failures the error code shall be
`"UNAUTHENTICATED"`.
