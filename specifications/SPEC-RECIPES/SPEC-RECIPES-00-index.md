# SPEC-RECIPES — Recipe Domain + YAML Loader + `seed` Subcommand

**Story:** 05 — Recipe domain + YAML loader + `seed` subcommand
**Status:** Draft → Final
**Date:** 2026-05-02
**Authors:** sdd-spec-author

---

## Purpose

This specification governs the in-memory recipe domain model
(`internal/domain/recipe.go`), the canonical taxonomies and their validators
(`internal/domain/taxonomy.go`), the YAML loader with aggregated error
reporting (`internal/seed/recipes.go`), the transactional Postgres upsert
store (`internal/seed/store.go`), the `cooksense-server seed` subcommand, the
sample seed recipe required by tests, and all associated documentation and
tests.

Out of scope: HTTP handlers (stories 07–10), curated content of 15+ recipes
(story 06), lessons loader (story 10), photo handling (post-MVP).

---

## File map

| File | Contents |
|------|----------|
| [SPEC-RECIPES-01-preamble](SPEC-RECIPES-01-preamble.md) | AI constraints, authorship, traceability rules |
| [SPEC-RECIPES-02-introduction](SPEC-RECIPES-02-introduction.md) | Story relationship, decision references |
| [SPEC-RECIPES-03-goals](SPEC-RECIPES-03-goals.md) | Goals, non-goals, constraints |
| [SPEC-RECIPES-04-context](SPEC-RECIPES-04-context.md) | External dependencies, package boundary map |
| [SPEC-RECIPES-05-architecture](SPEC-RECIPES-05-architecture.md) | Pipeline diagram, dependency graph, transaction shape |
| [SPEC-RECIPES-06-packages](SPEC-RECIPES-06-packages.md) | All SPEC-RECIPES-NNN requirements: types, signatures, SQL |
| [SPEC-RECIPES-07-configuration](SPEC-RECIPES-07-configuration.md) | Seed directory path, sample recipe rules |
| [SPEC-RECIPES-08-build](SPEC-RECIPES-08-build.md) | Build, lint, vet, dependency rules |
| [SPEC-RECIPES-09-testing](SPEC-RECIPES-09-testing.md) | Test strategy, named tests, fixtures |
| [SPEC-RECIPES-10-documentation](SPEC-RECIPES-10-documentation.md) | Doc comments, README impact |
| [SPEC-RECIPES-A-checklist](SPEC-RECIPES-A-checklist.md) | Specification completeness checklist |
| [SPEC-RECIPES-B-tasks](SPEC-RECIPES-B-tasks.md) | Ordered implementation task list |

---

## SPEC-ID registry

| ID | Summary | Section |
|----|---------|---------|
| SPEC-RECIPES-001 | `internal/domain` imports no I/O packages (no `os`, `io`, `database/sql`, `pgx`) | §6.1 |
| SPEC-RECIPES-002 | `domain.Recipe` struct fields match the YAML schema | §6.1 |
| SPEC-RECIPES-003 | `domain.RecipeIngredient` struct fields (`Name`, `Category`, `Quantity`, `Unit`, `Optional`) | §6.1 |
| SPEC-RECIPES-004 | `domain.SlugPattern` exposes `^[a-z0-9]+(-[a-z0-9]+)*$` | §6.1 |
| SPEC-RECIPES-005 | `domain.ValidateSlug(s string) error` returns non-nil for any non-matching slug | §6.1 |
| SPEC-RECIPES-006 | `Recipe.Validate() error` enforces all invariants (≥ 2 ingredients, ≥ 1 step, `TimeMinutes > 0`, valid slug, at least one of `Tags` or `CookingMethods` non-empty) | §6.1 |
| SPEC-RECIPES-007 | `PassivePrepMinutes >= 0` enforced | §6.1 |
| SPEC-RECIPES-008 | `Steps` is `[]string` with length ≥ 1 | §6.1 |
| SPEC-RECIPES-009 | `CookingMethods`, `Tags`, `FlavorProfile` are `[]string` (taxonomy-validated separately) | §6.1 |
| SPEC-RECIPES-010 | `domain.CookingMethods` canonical set: `pan-sear`, `roast`, `braise`, `boil`, `steam`, `grill`, `slow-cook`, `pressure-cook`, `raw`, `bake` | §6.2 |
| SPEC-RECIPES-011 | `domain.FlavorProfiles` canonical set: `acid`, `fat`, `salt`, `sweet`, `bitter`, `umami`, `spicy`, `herbaceous` | §6.2 |
| SPEC-RECIPES-012 | `domain.IngredientCategories` canonical set: `protein`, `vegetable`, `fruit`, `starch`, `dairy`, `fat`, `spice`, `herb`, `condiment`, `liquid`, `other` | §6.2 |
| SPEC-RECIPES-013 | `domain.ValidateCookingMethod(s string) error` rejects values not in the canonical set | §6.2 |
| SPEC-RECIPES-014 | `domain.ValidateFlavorProfile(s string) error` rejects values not in the canonical set | §6.2 |
| SPEC-RECIPES-015 | `domain.ValidateCategory(s string) error` rejects values not in the canonical set | §6.2 |
| SPEC-RECIPES-016 | Taxonomy validators are case-sensitive (lowercase only) | §6.2 |
| SPEC-RECIPES-017 | Canonical sets are exposed as read-only iterables (not mutable slices) | §6.2 |
| SPEC-RECIPES-018 | `seed.Load(ctx, dir string) ([]domain.Recipe, error)` reads `*.yaml` files | §6.3 |
| SPEC-RECIPES-019 | YAML parsing uses `gopkg.in/yaml.v3` | §6.3 |
| SPEC-RECIPES-020 | Loader validates each recipe against domain invariants and taxonomies | §6.3 |
| SPEC-RECIPES-021 | Loader aggregates errors across all files; never fails on the first | §6.3 |
| SPEC-RECIPES-022 | Each error message includes filename and YAML line number when available | §6.3 |
| SPEC-RECIPES-023 | Slug uniqueness is enforced across the entire loaded set | §6.3 |
| SPEC-RECIPES-024 | `seed.LoadError` type wraps the aggregated errors and implements `errors.Unwrap` returning `[]error` | §6.3 |
| SPEC-RECIPES-025 | Empty directory returns `(nil, nil)` — no error | §6.3 |
| SPEC-RECIPES-026 | `seed.Store(ctx, pool *pgxpool.Pool, recipes []domain.Recipe) (loaded, ingredients int, err error)` | §6.4 |
| SPEC-RECIPES-027 | All writes happen in a single transaction (`pool.BeginTx`) | §6.4 |
| SPEC-RECIPES-028 | UPSERT recipes by `slug` (`ON CONFLICT (slug) DO UPDATE`) | §6.4 |
| SPEC-RECIPES-029 | UPSERT ingredients by `name` (`ON CONFLICT (name) DO UPDATE`) | §6.4 |
| SPEC-RECIPES-030 | `recipe_ingredients` rows are REPLACED for each affected recipe (DELETE-then-INSERT inside the tx) | §6.4 |
| SPEC-RECIPES-031 | Any error rolls back the transaction; no partial loads | §6.4 |
| SPEC-RECIPES-032 | Idempotent: re-running `Store` with the same input creates no duplicates and no orphan ingredients | §6.4 |
| SPEC-RECIPES-033 | `loaded` counts recipes upserted; `ingredients` counts distinct ingredients seen | §6.4 |
| SPEC-RECIPES-034 | Pool is passed via constructor argument; no global state | §6.4 |
| SPEC-RECIPES-035 | `cooksense-server seed` subcommand exists and is dispatched from `main.go` | §6.5 |
| SPEC-RECIPES-036 | Subcommand reads from `seed/recipes/` by default; `--dir` flag overrides | §6.5 |
| SPEC-RECIPES-037 | On success, prints `loaded N recipes, M ingredients` to stdout | §6.5 |
| SPEC-RECIPES-038 | On any error, prints to stderr and exits non-zero | §6.5 |
| SPEC-RECIPES-039 | Subcommand opens the pgx pool via `db.Open` (reuses SPEC-DB-NNN) | §6.5 |
| SPEC-RECIPES-040 | At least one valid sample recipe lives at `seed/recipes/_sample.yaml` for tests | §7.1 |
| SPEC-RECIPES-041 | Sample recipe is excluded from curated content (story 06) by underscore prefix | §7.1 |
| SPEC-RECIPES-042 | `internal/seed/doc.go` package comment describes Load + Store responsibilities | §10.1 |
| SPEC-RECIPES-043 | README updated with `make seed` usage | §10.1 |
| SPEC-RECIPES-044 | Test `Test_LoadRecipes_InvalidTaxonomy_ReturnsError` exists | §9.1 |
| SPEC-RECIPES-045 | Test `Test_LoadRecipes_DuplicateSlug_ReturnsError` exists | §9.1 |
| SPEC-RECIPES-046 | Test `Test_LoadRecipes_MissingRequiredField_ReturnsError` exists | §9.1 |
| SPEC-RECIPES-047 | Test `Test_LoadRecipes_ValidSample_ReturnsRecipes` exists | §9.1 |
| SPEC-RECIPES-048 | Test `Test_StoreRecipes_Idempotent_NoDuplicates` exists (integration, compose Postgres) | §9.1 |
