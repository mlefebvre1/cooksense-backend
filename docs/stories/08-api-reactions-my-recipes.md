# Story 08 — Reactions and "My recipes"

Status: TODO
Estimate: M

## User story

As a user, I want to mark recipes as liked, disliked, or "try later" and to
see a list of my liked / try-later recipes so that I can come back to them
when planning meals.

## Background

See `docs/architecture/api.md` (sections "Reactions" and "My recipes") and
D-0008. Reactions are stored as one row per `(uid, recipe_id)` with the
`reaction_kind` enum.

## Acceptance criteria

- [ ] `POST /api/reactions` (auth) with body
      `{"recipe_slug": "...", "kind": "LIKE|DISLIKE|TRY_LATER"}`:
  - UPSERT on `(firebase_uid, recipe_id)` — re-posting a different `kind`
    overwrites the previous one.
  - `200` on success with the new state.
  - `404 NOT_FOUND` if `recipe_slug` is unknown.
  - `422 INVALID_REACTION_KIND` if `kind` is outside the enum.
  - `400 INVALID_PAYLOAD` for malformed JSON or missing fields.
- [ ] `DELETE /api/reactions/{slug}` (auth):
  - Removes the reaction if any.
  - `204 No Content`, idempotent (also returns 204 if no reaction existed).
  - `404` if the recipe slug itself is unknown.
- [ ] `GET /api/me/recipes?kind=LIKE|TRY_LATER` (auth):
  - Defaults to `kind=LIKE` when omitted.
  - Returns `<RecipeBrief>` list, ordered by reaction `created_at DESC`.
  - `kind=DISLIKE` is **not** exposed — disliked recipes are private and not
    listable in MVP.
  - `400` for any other value of `kind`.
- [ ] Repository (`internal/reactions/repo.go`) exposes `Upsert(ctx, uid,
      recipeID, kind)`, `Delete(ctx, uid, recipeID)`, and
      `ListByKind(ctx, uid, kind)` returning `RecipeBrief` rows joined with
      `recipes`.
- [ ] Integration tests cover: upsert from no-row, upsert overwrite, delete,
      list ordering, 404 on unknown slug, 422 on invalid kind.

## Technical notes

- The handler resolves `recipe_slug` → `recipe_id` first; do not invent a
  variant API that takes IDs.
- SQL for upsert:
  ```sql
  INSERT INTO user_reactions (firebase_uid, recipe_id, kind)
  VALUES ($1, $2, $3)
  ON CONFLICT (firebase_uid, recipe_id) DO UPDATE
      SET kind = EXCLUDED.kind, created_at = now()
  RETURNING created_at;
  ```
- Watch the foreign key on `users(firebase_uid)`: the auth middleware must
  have provisioned the user before this query. That's the case for any
  authenticated request (story 04 handles it).

## Out of scope

- History of reactions (we keep only the latest).
- Bulk operations.

## Dependencies

- depends on: 04, 05, 06, 07
- blocks: 12

## Definition of Done

- [ ] AC met.
- [ ] Integration tests green.
- [ ] API doc unchanged or amended in the same PR if any deviation is needed.
