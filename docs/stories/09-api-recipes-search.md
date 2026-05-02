# Story 09 — Search recipes by ingredients

Status: TODO
Estimate: M

## User story

As a user with specific ingredients on hand, I want to enter 2 to 5
ingredients and get the recipes that use the most of them, so that I cook
with what I have instead of buying more.

## Background

See `docs/architecture/api.md` ("Recipes — search by ingredients") for the
contract. The matching must be case- and accent-insensitive and respect
ingredient aliases.

## Acceptance criteria

- [ ] `GET /api/recipes/search?ingredients=tomato,beef&scope=all|liked` (auth):
  - Parses `ingredients` as a comma-separated list. **Trims** whitespace,
    lower-cases, removes diacritics.
  - Returns `400 INVALID_INGREDIENT_COUNT` for fewer than 2 or more than 5
    entries.
  - Matches each input against `ingredients.name` and `ingredients.aliases`
    (an input matches if it equals any of those after normalization).
  - `scope=all` (default): excludes the caller's `DISLIKE`d recipes.
  - `scope=liked`: only recipes the caller has reacted to with `LIKE`.
  - Ranks results by **count of matched ingredients DESC**, then
    `time_minutes ASC`.
  - Response shape: `{"recipes": [<RecipeBrief>, …], "matched": {"<slug>": N}}`
    where `N` is the number of matched ingredients for that recipe.
- [ ] Normalization helper (`internal/recipes/normalize.go`) with unit tests
      for: case folding, diacritic stripping (e.g. `bœuf` → `boeuf`),
      whitespace trimming.
- [ ] Repository method `Search(ctx, uid, normalizedNames, scope)` — single
      SQL query, no N+1.
- [ ] Integration tests:
  - `boeuf` matches a recipe whose ingredient name is `beef` and aliases
    contain `boeuf`.
  - Ranking is stable: a recipe matching 3 inputs comes before one matching 2.
  - `scope=liked` excludes recipes the user has not liked.
  - Boundary checks: 1 ingredient → 400; 6 ingredients → 400.

## Technical notes

- Normalization must happen on **both sides**:
  - Inputs from the query string (always normalized).
  - At seed time, the loader populates a normalized projection (e.g. a
    `lower(unaccent(name))` expression index) so SQL matches are fast.
- Suggested DDL addition (in this story, as `0002_search.up.sql`):
  ```sql
  CREATE EXTENSION IF NOT EXISTS unaccent;
  CREATE INDEX ingredients_norm_idx
      ON ingredients (lower(unaccent(name)));
  CREATE INDEX ingredients_aliases_norm_gin
      ON ingredients USING GIN ((array(SELECT lower(unaccent(a)) FROM unnest(aliases) a)));
  ```
- One reasonable query approach: CTE that resolves input strings to
  `ingredient_id`s, then aggregates counts per recipe.

## Out of scope

- Full-text search on recipe title / concept (post-MVP).
- Suggesting *substitute* ingredients.
- Anti-waste optimization (post-MVP V2).

## Dependencies

- depends on: 04, 05, 06, 07
- blocks: 12

## Definition of Done

- [ ] AC met.
- [ ] New migration `0002_search.up.sql` and matching `down.sql` committed.
- [ ] Integration tests green.
