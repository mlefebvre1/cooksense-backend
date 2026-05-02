# Authentication

## Source of truth

Identity is owned by **one Firebase project** shared between mobile and
backend. The backend never issues credentials, never stores passwords, and
never re-implements OAuth/OIDC.

## Token flow

```
Mobile/Web ── Firebase login ──► Firebase ID token (JWT, ~1h TTL)
            │
            ▼
       Authorization: Bearer <ID_TOKEN>
            │
            ▼
   CookSense backend  (auth middleware)
            │
            ▼
   firebase.google.com/go/v4/auth.VerifyIDToken(ctx, token)
            │
            ▼
   Decoded token  →  ctx user (UID, email)
            │
            ▼
   users table upsert (last_seen_at = now())
            │
            ▼
       Handler runs
```

## Verification rules

- Use `firebase.google.com/go/v4` and **never parse the JWT manually**.
- Initialize the Firebase Admin client **once** at startup, share via
  dependency injection. Re-init per request would crush latency.
- The middleware rejects with `401 UNAUTHENTICATED` when:
  - Header missing or not `Bearer …`.
  - Token signature invalid.
  - Token expired (`exp` past).
  - Token issued for a different audience/project.
  - `auth/revoked-token` (we call `VerifyIDTokenAndCheckRevoked` for sensitive
    actions; for read endpoints, plain `VerifyIDToken` is sufficient).
- The middleware **does not** call Firebase on every request beyond local
  signature verification (the SDK caches Google public keys).

## Lazy user provisioning

On first authenticated request from a UID we have not seen:

```sql
INSERT INTO users (firebase_uid, display_name, email)
VALUES ($1, $2, $3)
ON CONFLICT (firebase_uid) DO UPDATE
   SET last_seen_at = now(),
       email        = COALESCE(EXCLUDED.email, users.email),
       display_name = COALESCE(EXCLUDED.display_name, users.display_name);
```

This is performed by the auth middleware (or a thin `users.Service.Touch`
called from it). The latency cost is one UPSERT — acceptable for MVP.

## Context contract

- The middleware writes a typed user into `r.Context()`.
- Handlers retrieve it via `auth.UserFromContext(ctx) (auth.User, bool)`.
- `auth.User` is the only authenticated identity type used downstream:
  ```go
  type User struct {
      UID         string
      Email       string
      DisplayName string
  }
  ```

## Configuration

Required env vars:

- `FIREBASE_PROJECT_ID` — used to validate the token audience.
- `GOOGLE_APPLICATION_CREDENTIALS` — path to the Firebase Admin service
  account JSON (mounted from a secret in deployment, gitignored locally).

The repo never contains real credentials. A placeholder lives at
`secrets/firebase-admin.json.example` (created in story #4).

## Test strategy

- Unit-test the middleware with a **fake verifier** that returns a controlled
  `auth.Token`. Don't hit Firebase in unit tests.
- One integration test using a real-but-throwaway Firebase emulator is
  optional post-MVP. For MVP, the fake verifier is sufficient.

## Threat model (MVP)

- 🟢 Token forgery → blocked by signature verification.
- 🟢 Token replay across users → each request includes the UID; downstream
  authorization always uses `ctx.User.UID`.
- 🟡 Token theft on the client → out of backend scope (mobile responsibility).
- 🟡 Stolen service account JSON → keep it out of git, rotate on suspicion.
- 🔴 No rate limiting in MVP → tracked as a post-MVP follow-up.
