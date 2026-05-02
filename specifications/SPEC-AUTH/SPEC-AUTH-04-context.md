# SPEC-AUTH-04 — System Context

← [Index](SPEC-AUTH-00-index.md)

---

## 1. External dependencies introduced by this story

| Dependency | Version (go.mod) | Role |
|-----------|-----------------|------|
| `firebase.google.com/go/v4` | latest stable | Firebase Admin SDK — token verification, Admin client init |
| `google.golang.org/api` | (transitive) | OAuth2 credential loading for the Admin SDK |
| `github.com/jackc/pgx/v5` | already added (Story 03) | Postgres driver for the `users` UPSERT |

No additional direct dependencies beyond what Story 03 already introduced,
except the Firebase Admin SDK.

## 2. Runtime external services

| Service | Used for | MVP stance |
|---------|----------|-----------|
| Firebase Identity Platform | `VerifyIDToken` (offline, uses cached public keys) | Required — startup fails without valid credentials |
| Google key endpoint (`https://www.googleapis.com/...`) | SDK key cache refresh (done by the SDK automatically, not by our code) | Network access required; SDK handles retry |
| PostgreSQL 17 | `users` UPSERT | Required — from Story 03 |

## 3. Decision references

| ID | Decision | Impact |
|----|---------|--------|
| D-0004 | Firebase Auth is the sole identity provider | No homegrown JWT issuance; no other OAuth flows |
| D-0007 | Single Firebase project shared with mobile | `FIREBASE_PROJECT_ID` is the one audience to validate |
| D-0001 | Go 1.26.2 | Language version |
| D-0003 | `net/http` stdlib HTTP | No framework-level middleware wiring |

## 4. Package boundary map

```
cmd/cooksense-server/main.go
        │
        ├── internal/auth/        ← this story owns this package
        │     firebase.go         (SDK init, NewFirebaseVerifier)
        │     middleware.go       (Middleware, UserFromContext)
        │     verifier.go         (Verifier interface, fakeVerifier)
        │     user.go             (User type, contextKey)
        │
        ├── internal/users/       ← this story adds repo.go
        │     doc.go              (already exists)
        │     repo.go             (Toucher implementation)
        │
        └── internal/config/      ← this story reads FirebaseProjectID,
              config.go           GoogleAppCredentials (added to Config struct)
```

`internal/auth` shall **not** import `internal/users`. The dependency flows
inward through the `users.Toucher` interface, which is defined in
`internal/users` and accepted as a parameter by `auth.Middleware`.
