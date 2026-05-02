# SPEC-RECIPES — Appendix A · Specification Checklist

[← Index](SPEC-RECIPES-00-index.md)

## A.1 Completeness checklist

| # | Criterion | Status |
|---|-----------|--------|
| 1 | Every AC and DoD item from `docs/stories/05-recipes-domain-yaml-loader.md` maps to ≥ 1 SPEC-RECIPES-NNN ID. | ✅ (matrix below) |
| 2 | Every SPEC-RECIPES-NNN ID maps to ≥ 1 named test in §9. | ✅ |
| 3 | All requirements use RFC-2119 keywords (`SHALL` / `MUST` / `SHOULD` / `MAY`). | ✅ |
| 4 | Public Go signatures are stated for every new exported symbol. | ✅ |
| 5 | All SQL is shown verbatim (UPSERTs + DELETE/INSERT) and matches `migrations/0001_init.up.sql`. | ✅ |
| 6 | Error taxonomy is defined (`LoadError`, wrapped store errors). | ✅ |
| 7 | Configuration variables are documented (none new; reused vars listed). | ✅ |
| 8 | Observability/logging guidance present (slog usage, no fmt.Print except CLI). | ✅ |
| 9 | Security/threat considerations addressed (no secrets in YAML; pool injected). | ✅ |
| 10 | Appendix B lists ordered, atomic implementation tasks. | ✅ |

## A.2 Story AC → SPEC-ID traceability matrix

| Story acceptance criterion | SPEC-RECIPES IDs |
|----------------------------|------------------|
| AC #1 — `internal/domain/recipe.go` defines `Recipe`, `RecipeIngredient`, supporting types | 001, 002, 003, 004, 005, 006, 007, 008, 009 |
| AC #2 — `internal/domain/taxonomy.go` defines canonical sets + `Validate*` for each | 010, 011, 012, 013, 014, 015, 016, 017 |
| AC #3 — `internal/seed/recipes.go` reads YAML, parses with `yaml.v3`, validates, aggregates errors | 018, 019, 020, 021, 022, 023, 024, 025 |
| AC #4 — `internal/seed/store.go` upserts in a single tx (recipes by slug, ingredients by name, REPLACE recipe_ingredients) | 026, 027, 028, 029, 030, 031, 033, 034 |
| AC #5 — `cooksense-server seed` subcommand prints `loaded N recipes, M ingredients` | 035, 036, 037, 039 |
| AC #6 — Re-running `seed` is idempotent (no dups, no orphans) | 032 |
| AC #7 — Unit tests cover invalid taxonomy, duplicate slug, missing required field, valid sample | 044, 045, 046, 047 |

## A.3 Story DoD → SPEC-ID traceability matrix

| Story DoD item | SPEC-RECIPES IDs |
|----------------|------------------|
| AC met | All 048 IDs |
| At least one valid sample recipe lives in `seed/recipes/` | 040, 041 |
| `make seed` works end-to-end against the compose Postgres | 035–039, 048 (integration test) |

## A.4 Technical-notes coverage

| Story technical note | SPEC-RECIPES IDs |
|----------------------|------------------|
| Use `yaml.v3` (`gopkg.in/yaml.v3`) | 019 |
| Loader must aggregate errors; never fail on first | 021, 024 |
| Errors must report filenames + line numbers | 022 |
| Transaction must roll back on any error | 027, 031 |
| `pgx.CopyFrom` allowed for large catalogs (else batched ON CONFLICT) | §6.6 |
| Slug regex: `^[a-z0-9]+(-[a-z0-9]+)*$` | 004, 005 |

## A.5 Sign-off

- Spec reviewed against story file: ✅
- Spec reviewed against `docs/architecture/data-model.md`: ✅
- All SPEC-IDs referenced in §6 appear in the index registry of
  `SPEC-RECIPES-00-index.md`: ✅
- Status transition: **Draft → Final**.
