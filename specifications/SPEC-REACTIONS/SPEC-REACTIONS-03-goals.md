# SPEC-REACTIONS — §3 Goals, Non-Goals, Constraints

[← Index](SPEC-REACTIONS-00-index.md)

## 3.1 Goals

| ID | Goal |
|----|------|
| G-1 | Provide an idempotent `POST /api/reactions` endpoint with UPSERT semantics on `(firebase_uid, recipe_id)`. |
| G-2 | Provide an idempotent `DELETE /api/reactions/{slug}` endpoint that returns `204` whether or not a row was previously present. |
| G-3 | Provide a `GET /api/me/recipes` endpoint that lists the caller's `LIKE` or `TRY_LATER` recipes ordered by reaction recency. |
| G-4 | Keep the domain layer free of I/O — `internal/domain` defines `ReactionKind`, `Reaction`, sentinel errors, and the `ReactionRepository` interface. |
| G-5 | Keep the repository layer free of HTTP — `internal/reactions` operates on `(uid, recipeID)` only and never sees a slug or HTTP request. |
| G-6 | Keep handler logic free of SQL — handlers call `reactions.Service`, which orchestrates slug resolution and repo calls. |
| G-7 | Surface a stable error catalog (`NOT_FOUND`, `INVALID_PAYLOAD`, `INVALID_REACTION_KIND`, `UNAUTHENTICATED`) consistent with `docs/architecture/api.md`. |
| G-8 | Maintain `≥ 90 %` line coverage on new code, including HTTP-level integration tests against the compose Postgres. |

## 3.2 Non-Goals

| ID | Non-Goal |
|----|----------|
| NG-1 | This spec **does not** introduce a numeric-id API surface for reactions. The recipe MUST be addressed by `slug`. |
| NG-2 | This spec **does not** persist a history of reactions. Re-posting overwrites `kind` and `created_at` (D-0008). |
| NG-3 | This spec **does not** expose `DISLIKE` reactions through `GET /api/me/recipes`. Disliked recipes are private. |
| NG-4 | This spec **does not** introduce bulk endpoints. One reaction per request. |
| NG-5 | This spec **does not** introduce caching of reactions. The MVP catalog is small enough to query directly. |
| NG-6 | This spec **does not** modify the auth middleware (story 04 / SPEC-AUTH) or the recipes repository beyond the addition of `ResolveSlug` already required by stories 07/09. |

## 3.3 Constraints

- **Schema is fixed.** The `user_reactions` table and the `reaction_kind`
  enum are defined in `migrations/0001_init.up.sql` (SPEC-DB / SPEC-RECIPES
  precursor). This spec MUST NOT alter the schema.
- **Auth is fixed.** The `firebase_uid` is read from the request context
  populated by SPEC-AUTH middleware. This spec MUST NOT introduce alternate
  identity sources.
- **Foreign-key integrity.** Inserting a reaction requires that the user row
  already exists. SPEC-AUTH guarantees lazy provisioning on every
  authenticated request (Story 04), so the FK never fires under normal
  flow. The handler MUST NOT attempt to create the user row.
- **Sort key.** `GET /api/me/recipes` MUST sort by
  `user_reactions.created_at DESC` (most recent reaction first). No
  pagination in MVP.
- **Slug-only.** Every public payload references the recipe by its `slug`.
  The numeric `recipe_id` MUST NOT leak into responses or error messages.
