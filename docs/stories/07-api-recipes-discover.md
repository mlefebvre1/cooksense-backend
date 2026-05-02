# Story 07 — Recipe discovery and detail endpoints

Status: TODO
Estimate: M

## User story

As a user opening the app, I want to receive a fresh stream of recipes I have
not seen yet, and to open the full detail of any of them, so that I can pick
what I'll cook tonight in seconds.

## Background

Implements the swipe-feed and the recipe detail pages. See
`docs/architecture/api.md` (sections "Recipes — discovery and detail") for the
contract. Decision references: D-0007, D-0008.

## Acceptance criteria

- [ ] `GET /api/recipes/discover?limit=N` (auth):
  - Returns up to `N` recipes (default 10, max 25).
  - Excludes recipes the caller has any reaction on (any of LIKE, DISLIKE,
    TRY_LATER) — this is the simplest "not yet reacted to" definition.
  - Result order is randomized per call (`ORDER BY random()` is acceptable
    in MVP given catalog size).
  - Response shape matches `<RecipeBrief>` in the API doc.
- [ ] `GET /api/recipes/{slug}` (auth):
  - Returns `<RecipeFull>` with ingredients (sorted by category then name)
    and steps (in stored order).
  - `404 NOT_FOUND` for unknown slugs.
- [ ] Repository (`internal/recipes/repo.go`) is the only DB-touching layer
      and exposes `Discover(ctx, uid, limit)` and `GetBySlug(ctx, slug)`.
- [ ] Service (`internal/recipes/service.go`) enforces limit clamping and
      orchestrates repo calls.
- [ ] Handler (`internal/recipes/handler.go`) only handles JSON
      encoding/decoding and HTTP status mapping.
- [ ] Integration test (against compose Postgres, with seeded data and a
      `fakeVerifier`):
  - Discover returns the expected number of recipes.
  - Discover never returns a recipe the user has already reacted to.
  - Detail returns the seeded recipe.
  - Detail returns 404 for an unknown slug.

## Technical notes

- SQL sketch for discover:
  ```sql
  SELECT r.* FROM recipes r
  WHERE NOT EXISTS (
      SELECT 1 FROM user_reactions ur
      WHERE ur.firebase_uid = $1 AND ur.recipe_id = r.id
  )
  ORDER BY random()
  LIMIT $2;
  ```
  Index `user_reactions(firebase_uid, …)` keeps the anti-join cheap.
- Ingredient hydration for detail uses one extra query (or `LEFT JOIN
  LATERAL` aggregating to JSONB) — either is fine; pick the simpler.
- JSON field names are `snake_case`, matching the API doc exactly.

## Out of scope

- Reactions endpoint — story 08.
- Search — story 09.
- "Time-to-decision" telemetry — post-MVP.

## Dependencies

- depends on: 04, 05, 06
- blocks: 12

## Definition of Done

- [ ] AC met.
- [ ] Integration test green against compose Postgres.
- [ ] `go vet ./...` clean.
