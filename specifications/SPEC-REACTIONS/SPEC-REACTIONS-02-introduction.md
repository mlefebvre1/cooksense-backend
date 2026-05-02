# SPEC-REACTIONS — §2 Introduction

[← Index](SPEC-REACTIONS-00-index.md)

## 2.1 Story relationship

This specification implements `docs/stories/08-api-reactions-my-recipes.md`
("Reactions and 'My recipes'"), an MVP story estimated **M** that depends on
stories 04 (auth middleware), 05 (recipe domain + loader), 06 (curated
content), and 07 (recipe discovery / detail endpoints). It blocks story 12
(quickstart README) because the curated demo flow requires functioning
reactions.

## 2.2 Decision references

| Decision | Title | Relevance to SPEC-REACTIONS |
|----------|-------|-----------------------------|
| D-0001 | Backend language: Go 1.26 | All code MUST use modern Go 1.26 idioms |
| D-0002 | HTTP layer: stdlib `net/http` with Go 1.22 pattern router | Routes registered with `mux.HandleFunc("METHOD /path", h)` |
| D-0003 | Database: PostgreSQL 17 + `pgx/v5` | Repo uses `pgxpool.Pool` and pgx ENUM types |
| D-0004 | Authentication: Firebase ID tokens | Auth middleware provides `firebase_uid` via context |
| D-0007 | Identity: `firebase_uid TEXT` is the FK across all user-owned tables | `user_reactions.firebase_uid` is `TEXT` |
| D-0008 | Reaction kinds: `ENUM('LIKE','DISLIKE','TRY_LATER')`; one row per `(user, recipe)`, UPSERT semantics, no history | The entire data shape and upsert behavior |

## 2.3 Architecture documents consumed

- `docs/architecture/api.md` — sections "Reactions" and "My recipes".
- `docs/architecture/data-model.md` — `user_reactions` table definition,
  `reaction_kind` enum, FK to `recipes(id)` and `users(firebase_uid)`.
- `docs/architecture/auth.md` — `firebase_uid` context propagation.

## 2.4 Specifications consumed

- **SPEC-AUTH** — provides the auth middleware that populates the
  request-scoped `firebase_uid`.
- **SPEC-DB** — provides `pgxpool.Pool` construction and `DATABASE_URL`
  configuration.
- **SPEC-RECIPES** — provides the `recipes` table, `slug` uniqueness, and
  the `domain.RecipeBrief` shape that `GET /api/me/recipes` returns.

## 2.5 Specifications produced

This spec produces no new public Go module. It introduces the
`internal/reactions` package and the `domain.ReactionKind` /
`domain.Reaction` / `domain.ReactionRepository` types, all internal to the
binary.

## 2.6 Out of scope

- Reaction history (D-0008 explicitly defers it; if needed later, a
  `reaction_history` table is added without touching `user_reactions`).
- Bulk reactions (`POST /api/reactions/bulk` is **not** introduced).
- Listing `DISLIKE` reactions in `/api/me/recipes` (Story 08 forbids it).
- Variants of the API that take a numeric `recipe_id` — the handler MUST
  resolve `recipe_slug` first.
