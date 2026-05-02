# SPEC-RECIPES — Appendix B · Implementation Tasks

[← Index](SPEC-RECIPES-00-index.md)

Each task is atomic, ordered, and traceable to SPEC-RECIPES-NNN IDs.
Implementation SHALL proceed in this order; reordering is allowed only
inside groups where dependencies do not cross.

| # | Task | SPEC-IDs | Depends on |
|---|------|----------|------------|
| T-01 | Add `gopkg.in/yaml.v3` to `go.mod` and run `go mod tidy`. | 019 | — |
| T-02 | Create `internal/domain/recipe.go` with `Recipe`, `RecipeIngredient`, `SlugPattern` const, compiled regex, `ValidateSlug`. | 001, 002, 003, 004, 005, 009 | T-01 |
| T-03 | Create `internal/domain/taxonomy.go` with `CookingMethods`, `FlavorProfiles`, `IngredientCategories`, and the three `Validate*` functions. | 010, 011, 012, 013, 014, 015, 016, 017 | T-02 |
| T-04 | Implement `Recipe.Validate()` in `recipe.go` aggregating with `errors.Join`. | 006, 007, 008 | T-02, T-03 |
| T-05 | Update `internal/domain/doc.go` package comment citing SPEC-RECIPES-001..017. | 042 | T-04 |
| T-06 | Write unit tests for the domain layer (slug, validate, taxonomies). | 044 | T-05 |
| T-07 | Create `internal/seed/recipes.go` with `LoadError`, `Load(ctx, dir)` reading `*.yaml`, parsing with `yaml.v3`, accumulating errors. | 018, 019, 020, 021, 022, 024, 025 | T-04 |
| T-08 | Add slug-uniqueness pass in `Load`. | 023 | T-07 |
| T-09 | Update `internal/seed/doc.go` package comment per SPEC-RECIPES-042. | 042 | T-07 |
| T-10 | Create `internal/seed/testdata/` with the 5 fixture scenarios from §9.2. | — | T-07 |
| T-11 | Write loader unit tests: `Test_LoadRecipes_*` (4 AC tests + aggregation + empty dir + line numbers). | 044, 045, 046, 047 | T-08, T-10 |
| T-12 | Create `internal/seed/store.go` with `Store(ctx, pool, recipes)`, transaction, recipe UPSERT, ingredient UPSERT, REPLACE recipe_ingredients. | 026, 027, 028, 029, 030, 031, 033, 034 | T-08 |
| T-13 | Verify idempotency: integration test `Test_StoreRecipes_Idempotent_NoDuplicates`. | 032, 048 | T-12 |
| T-14 | Add `runSeed(args)` helper in `cmd/cooksense-server/main.go`; wire `seed` subcommand dispatch. | 035, 036, 037, 038, 039 | T-12 |
| T-15 | Author `seed/recipes/_sample.yaml` (minimal valid recipe) with the documentation header. | 040, 041 | T-04 |
| T-16 | Update `README.md` with the "Seeding the catalog" subsection. | 043 | T-14 |
| T-17 | Run `make lint test`; verify `go build ./...`, coverage targets met (§9.5). Run `make up && make migrate && make seed` end-to-end and confirm output line. | All | T-15, T-16 |

## B.1 Critical-path graph

```
T-01 → T-02 → T-03 → T-04 ┬─► T-05 → T-06
                          │
                          ├─► T-07 → T-08 ─┬─► T-09
                          │                ├─► T-10 → T-11
                          │                └─► T-12 ─┬─► T-13
                          │                          └─► T-14 → T-16
                          └─► T-15 ──────────────────────────────► T-17
```

T-17 is the gate before PR submission and SHALL only run when every prior
task is complete.

## B.2 Estimated effort

| Group | Tasks | Effort |
|-------|-------|--------|
| Domain | T-01..T-06 | ~2 h |
| Loader | T-07..T-11 | ~3 h |
| Store + CLI | T-12..T-14 | ~3 h |
| Sample + docs | T-15, T-16 | ~30 min |
| Verification | T-17 | ~30 min |

Total estimate: ~9 h focused work, matching the "M" estimate on the story.

## B.3 PR strategy

The implementation MAY land in a single PR (`feat/recipes-loader`) or be
split into two:
1. `feat/recipes-domain` — T-01..T-06 (pure domain, no DB).
2. `feat/recipes-loader` — T-07..T-17 (loader, store, CLI, sample, README).

Either approach SHALL keep each commit body citing the SPEC-RECIPES IDs it
implements, per the repo's commit-message contract.
