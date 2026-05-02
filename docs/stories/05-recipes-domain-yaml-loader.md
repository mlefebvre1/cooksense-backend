# Story 05 — Recipe domain + YAML loader + `seed` subcommand

Status: TODO
Estimate: M

## User story

As a content owner, I want to author recipes as YAML files in the repo and
load them into the database with a single command so that recipes are
reviewable, version-controlled, and reproducible across environments.

## Background

See `docs/architecture/data-model.md` for the recipe YAML schema and the
allowed taxonomies. See D-0005 for the rationale.

## Acceptance criteria

- [ ] `internal/domain/recipe.go` defines the in-memory `Recipe`,
      `RecipeIngredient`, and supporting types matching the YAML schema.
- [ ] `internal/domain/taxonomy.go` defines the canonical sets for
      `cooking_methods`, `flavor_profile`, `category` (ingredient). Provides a
      `Validate*` function for each.
- [ ] `internal/seed/recipes.go` reads `seed/recipes/*.yaml`, parses with
      `gopkg.in/yaml.v3`, validates against the taxonomies and invariants
      (≥ 2 ingredients, ≥ 1 step, `time_minutes > 0`, `slug` kebab-case &
      unique across loaded set), and returns `[]domain.Recipe` plus an
      aggregated error listing all problems.
- [ ] `internal/seed/store.go` upserts recipes into Postgres in a single
      transaction:
  - UPSERT on `recipes.slug`.
  - UPSERT ingredients by `name`.
  - REPLACE the `recipe_ingredients` rows for each affected recipe.
- [ ] `cooksense-server seed` subcommand calls the loader on
      `seed/recipes/` and prints `loaded N recipes, M ingredients` on success.
- [ ] Re-running `seed` is idempotent (no duplicates, no orphan ingredients
      created).
- [ ] Unit tests cover: invalid taxonomy, duplicate slug, missing required
      field, valid sample.

## Technical notes

- Use `yaml.v3` (`gopkg.in/yaml.v3`).
- The loader **must aggregate errors** (don't fail on the first problem). A
  failed load reports every offending file with line numbers when possible.
- The transaction must roll back on any error — partial loads are forbidden.
- Use `pgx.CopyFrom` for `recipe_ingredients` if the catalog grows large; for
  MVP, batched `INSERT … ON CONFLICT` is fine.
- Slug regex: `^[a-z0-9]+(-[a-z0-9]+)*$`.

## Out of scope

- Lessons loader — story 10.
- Image/photo handling — post-MVP.
- Authoring UI — post-MVP (V1 user-submitted recipes).

## Dependencies

- depends on: 03
- blocks: 06, 07, 08, 09

## Definition of Done

- [ ] AC met.
- [ ] At least one valid sample recipe lives in `seed/recipes/` for tests
      (the curated content lands in story 06).
- [ ] `make seed` works end-to-end against the compose Postgres.
