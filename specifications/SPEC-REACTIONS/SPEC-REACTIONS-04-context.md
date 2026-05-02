# SPEC-REACTIONS — §4 Context & Boundaries

[← Index](SPEC-REACTIONS-00-index.md)

## 4.1 Runtime requirements

| Requirement | Value |
|-------------|-------|
| Go | ≥ 1.26 |
| PostgreSQL | 17 (compose) |
| pgx | `github.com/jackc/pgx/v5` (already a project dependency) |
| HTTP router | stdlib `net/http` (Go 1.22+ pattern matching) |
| Logging | `log/slog` |

No new third-party dependency is introduced by this spec.

## 4.2 External dependencies (consumed)

| Dependency | Provider | Used for |
|------------|----------|----------|
| Auth context | SPEC-AUTH | `auth.UIDFromContext(ctx) (string, bool)` to extract `firebase_uid` |
| Pool | SPEC-DB | `*pgxpool.Pool` injected into `reactions.NewRepo` |
| `recipes` repo | SPEC-RECIPES + Story 07 | `recipes.Repo.ResolveSlug(ctx, slug) (int64, error)` |
| `domain.RecipeBrief` | SPEC-RECIPES | The shape returned by `ListByKind` |
| `domain.ErrRecipeNotFound` | SPEC-RECIPES (Story 07 introduces it) | Sentinel returned by `ResolveSlug` and surfaced as `404` |
| Error envelope | SPEC-AUTH / Story 07 | `api.WriteError(w, code, status, msg)` helper |

> **Note.** If Story 07 has not yet introduced `recipes.Repo.ResolveSlug` or
> `domain.ErrRecipeNotFound` at implementation time, Task T-04 of Appendix B
> SHALL add them in this same PR. Their addition is owned by SPEC-RECIPES;
> this spec only depends on their availability.

## 4.3 Package boundary map

| Package | Imports allowed | Imports forbidden |
|---------|-----------------|--------------------|
| `internal/domain` | stdlib only (`time`, `errors`, `fmt`) | `pgx`, `database/sql`, `net/http`, `os`, `io` |
| `internal/reactions` | `internal/domain`, `github.com/jackc/pgx/v5`, `github.com/jackc/pgx/v5/pgxpool`, stdlib | `internal/api`, `net/http` |
| `internal/api` (handlers added here) | `internal/domain`, `internal/reactions`, `internal/auth`, stdlib (`net/http`, `encoding/json`, `log/slog`) | `pgx` directly, `database/sql` |
| `cmd/cooksense-server` | all `internal/*` packages (wiring) | — |

Circular imports are forbidden by Go and would fail compilation. The
ordering above is checked manually during review.

## 4.4 Persistence context

The `user_reactions` table (from `migrations/0001_init.up.sql`):

```sql
CREATE TABLE user_reactions (
    firebase_uid TEXT          NOT NULL REFERENCES users(firebase_uid) ON DELETE CASCADE,
    recipe_id    BIGINT        NOT NULL REFERENCES recipes(id)         ON DELETE CASCADE,
    kind         reaction_kind NOT NULL,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT now(),
    PRIMARY KEY (firebase_uid, recipe_id)
);
CREATE INDEX user_reactions_uid_kind_idx ON user_reactions(firebase_uid, kind);
```

The `reaction_kind` enum:

```sql
CREATE TYPE reaction_kind AS ENUM ('LIKE', 'DISLIKE', 'TRY_LATER');
```

The `(firebase_uid, kind)` index is **already declared** in the migration —
SPEC-REACTIONS does NOT add an index. `ListByKind` queries are expected to
use it.

## 4.5 HTTP context

All three endpoints sit behind the auth middleware. The middleware:

- Verifies the Firebase ID token.
- Lazily provisions the `users` row (Story 04).
- Stores `firebase_uid` in the request context.
- Returns `401 UNAUTHENTICATED` on missing/invalid token (handled by the
  middleware; SPEC-REACTIONS handlers never see those requests).

The handlers added by this spec MUST extract `firebase_uid` via
`auth.UIDFromContext(ctx)`. If the function returns `ok=false`, that is a
programmer error (auth middleware not wired); the handler SHALL respond
`500 INTERNAL` and log at `ERROR`.
