# SPEC-CATALOG §9 — Testing Strategy

[← back to index](SPEC-CATALOG-00-index.md)

## 9.1 Test files introduced

| File | Purpose |
|------|---------|
| `internal/seed/catalog_diversity_test.go` | Filesystem-only test enforcing SPEC-CATALOG-001, -003, -004, -011, -012, -013, -014, -015, -016, -017 |

The test resides under `internal/seed` so the `seed` package can re-use
its YAML structs. It **shall** open `seed/recipes/` from the repo root via
`runtime.Caller` or the test's working directory + `../../seed/recipes`.

## 9.2 Required named test cases

All names follow the project convention `Test{What}_{Condition}_{ExpectedOutcome}`.

| # | Test name | SPEC-IDs |
|---|-----------|----------|
| 1 | `TestCatalog_FilesystemInventory_AtLeastFifteenRecipes` | SPEC-CATALOG-001, 004 |
| 2 | `TestCatalog_FilenameMatchesSlug_AllRecipes` | SPEC-CATALOG-003 |
| 3 | `TestCatalog_SchemaValid_AllRecipesParse` | SPEC-CATALOG-002, 005, 006 |
| 4 | `TestCatalog_TimeBound_ActiveLessThanOrEqualThirty` | SPEC-CATALOG-007 |
| 5 | `TestCatalog_PassivePrep_PositiveWhenTotalExceedsThirty` | SPEC-CATALOG-008 |
| 6 | `TestCatalog_TaxonomyCompliance_AllValuesInCanonicalSets` | SPEC-CATALOG-010 |
| 7 | `TestCatalog_MethodsDiversity_AtLeastFourDistinct` | SPEC-CATALOG-011 |
| 8 | `TestCatalog_ProteinDiversity_AtLeastThreeDistinct` | SPEC-CATALOG-012 |
| 9 | `TestCatalog_VegetarianFloor_AtLeastTwoRecipes` | SPEC-CATALOG-013 |
| 10 | `TestCatalog_PassiveCookFloor_AtLeastOneSlowOrPressureRecipe` | SPEC-CATALOG-014 |
| 11 | `TestCatalog_IngredientNamesVerbatimReuse_NoCloseDuplicates` | SPEC-CATALOG-015, 016, 017 |
| 12 | `TestCatalog_NoMediaFields_NoForbiddenKeys` | SPEC-CATALOG-023 |

The diversity test (SPEC-CATALOG-027) is satisfied by the union of cases
7–10. Cases 1–6, 11, 12 enforce per-file rules that complement the
SPEC-RECIPES loader unit tests.

## 9.3 Test-helper expectations

- `TestCatalog_*` **shall** call `t.Helper()` in any internal helper.
- Each test **shall** be independent (no shared mutable state).
- Tests **shall** print, on failure, the offending file path(s) and the
  exact assertion that failed (e.g. `cooking_methods union = [pan-sear,
  bake] (size 2 < 4)`).

## 9.4 Manual-only verifications

The following SPEC-IDs have **no automated test** and **shall** be verified
manually during PR review:

| SPEC-ID | Reason | Reviewer artifact |
|---------|--------|--------------------|
| SPEC-CATALOG-019 | Editorial concept shape | Reviewer comment per file |
| SPEC-CATALOG-020 | Voice / tone | Reviewer comment per file |
| SPEC-CATALOG-021 | Step content (why-clause judgement) | Reviewer comment per file |
| SPEC-CATALOG-022 | English language | Reviewer comment / scan |
| SPEC-CATALOG-024 | Reviewer comment per file | The PR's review thread itself |
| SPEC-CATALOG-025 | PR description summary table | The PR description |

## 9.5 End-to-end seed run

SPEC-CATALOG-026 is verified by the project's existing manual quickstart
(`make up && make migrate && make seed`). It **shall** be re-run by the
PR author and the result pasted into the PR description ("Seed output:
loaded N recipes, M ingredients").

A future story **may** automate this in CI via testcontainers; that is
out of scope for SPEC-CATALOG.

## 9.6 Coverage targets

SPEC-CATALOG does not change Go coverage thresholds. The diversity test
itself is the contract; coverage of test code is not measured.

## 9.7 Negative tests (loader path)

Negative-path tests (invalid taxonomy, duplicate slug, missing required
field) are owned by SPEC-RECIPES §9. SPEC-CATALOG **shall not** duplicate
them. If a catalog file violates the schema, the SPEC-RECIPES loader
test will fail first; the catalog tests assume schema validity for the
diversity assertions.
