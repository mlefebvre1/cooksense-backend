# SPEC-AUTH-10 — Documentation

← [Index](SPEC-AUTH-00-index.md)

---

## 1. Go Doc comments (mandatory on all exported symbols)

Every exported symbol introduced in this story shall carry a Go Doc comment.
The comment shall cite the governing SPEC-AUTH-NNN ID(s).

| Symbol | File | Required comment opening |
|--------|------|--------------------------|
| `Verifier` | `verifier.go` | `// Verifier verifies a Firebase ID token … // SPEC-AUTH-004` |
| `FakeVerifier` | `fake_verifier.go` or test file | `// FakeVerifier is an in-memory Verifier for use in tests … // SPEC-AUTH-006` |
| `User` | `user.go` | `// User represents an authenticated Firebase principal. // SPEC-AUTH-007` |
| `Middleware` | `middleware.go` | `// Middleware returns an http.Handler wrapper … // SPEC-AUTH-008` |
| `UserFromContext` | `middleware.go` | `// UserFromContext retrieves the authenticated User … // SPEC-AUTH-016` |
| `NewFirebaseApp` | `firebase.go` | `// NewFirebaseApp initializes the Firebase Admin SDK … // SPEC-AUTH-001` |
| `NewFirebaseVerifier` | `firebase.go` | `// NewFirebaseVerifier creates a Verifier … // SPEC-AUTH-005` |
| `Toucher` | `internal/users/repo.go` | `// Toucher upserts the user row … // SPEC-AUTH-017` |
| `Repo` | `internal/users/repo.go` | `// Repo implements Toucher … // SPEC-AUTH-018` |
| `NewRepo` | `internal/users/repo.go` | `// NewRepo returns a new Repo backed by pool.` |

## 2. `secrets/firebase-admin.json.example`

The example credential file (SPEC-AUTH-024) shall include a header comment
explaining its purpose:

```
# This file is a non-functional placeholder.
# Copy it to secrets/firebase-admin.json and replace with your real
# Firebase Admin SDK service account credentials.
# NEVER commit the real file — it is gitignored.
```

(Note: JSON does not support comments; the block above should be placed as a
`README` in the `secrets/` directory or as a comment in the `.gitignore` entry.)
The actual JSON file shall be valid JSON with placeholder values.

## 3. `secrets/README.md`

A `secrets/README.md` file shall be created explaining the `secrets/` directory
convention:

```markdown
# secrets/

This directory holds credential files that are **never committed** to git.

| File | How to obtain |
|------|--------------|
| `firebase-admin.json` | Download from Firebase Console → Project Settings → Service Accounts → Generate new private key |

A non-functional example is provided at `firebase-admin.json.example`.
```

## 4. README impact

The top-level `README.md` shall be updated (or the quickstart story 12 will
handle it) to reference the `secrets/` setup step. For this story, a note in
`secrets/README.md` is sufficient — the full quickstart lives in Story 12.

## 5. Architecture doc update

`docs/architecture/auth.md` shall not be modified by this story — it is
already authoritative. If a discrepancy is found between the spec and
`auth.md`, `auth.md` is amended first in the same PR.
