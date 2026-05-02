# SPEC-RECIPES — §9 Testing

[← Index](SPEC-RECIPES-00-index.md)

## 9.1 Required named tests

Test names SHALL follow `Test{What}_{Condition}_{ExpectedOutcome}` per the
project guideline. Each test SHALL document the SPEC-RECIPES IDs it
verifies in its docstring (forward traceability).

| Test name | Layer | Verifies |
|-----------|-------|----------|
| `Test_ValidateSlug_KebabCase_ReturnsNil` | unit (domain) | SPEC-RECIPES-004, 005 |
| `Test_ValidateSlug_UpperCase_ReturnsError` | unit (domain) | SPEC-RECIPES-005 |
| `Test_RecipeValidate_ValidInput_ReturnsNil` | unit (domain) | SPEC-RECIPES-006, 007, 008 |
| `Test_RecipeValidate_FewerThanTwoIngredients_ReturnsError` | unit (domain) | SPEC-RECIPES-006 |
| `Test_RecipeValidate_NoSteps_ReturnsError` | unit (domain) | SPEC-RECIPES-006, 008 |
| `Test_RecipeValidate_InvalidCookingMethod_ReturnsError` | unit (domain) | SPEC-RECIPES-006, 013 |
| `Test_RecipeValidate_InvalidFlavorProfile_ReturnsError` | unit (domain) | SPEC-RECIPES-014 |
| `Test_RecipeValidate_InvalidCategory_ReturnsError` | unit (domain) | SPEC-RECIPES-015 |
| `Test_ValidateCookingMethod_CaseSensitive_RejectsCapital` | unit (domain) | SPEC-RECIPES-016 |
| `Test_LoadRecipes_ValidSample_ReturnsRecipes` | unit (seed) | SPEC-RECIPES-018, 019, 020 (AC #7 happy path) |
| `Test_LoadRecipes_InvalidTaxonomy_ReturnsError` | unit (seed) | SPEC-RECIPES-020, 021, 024 (AC #7) |
| `Test_LoadRecipes_DuplicateSlug_ReturnsError` | unit (seed) | SPEC-RECIPES-023 (AC #7) |
| `Test_LoadRecipes_MissingRequiredField_ReturnsError` | unit (seed) | SPEC-RECIPES-020, 021 (AC #7) |
| `Test_LoadRecipes_AggregatesErrorsFromMultipleFiles` | unit (seed) | SPEC-RECIPES-021, 024 |
| `Test_LoadRecipes_EmptyDirectory_ReturnsNilNil` | unit (seed) | SPEC-RECIPES-025 |
| `Test_LoadRecipes_ErrorsIncludeFileAndLine` | unit (seed) | SPEC-RECIPES-022 |
| `Test_StoreRecipes_HappyPath_ReturnsCounts` | integration (seed+pg) | SPEC-RECIPES-026, 028, 029, 030, 033 |
| `Test_StoreRecipes_AnyErrorRollsBack_NoPartialState` | integration | SPEC-RECIPES-031 |
| `Test_StoreRecipes_Idempotent_NoDuplicates` | integration | SPEC-RECIPES-032 (AC #6) |
| `Test_SeedSubcommand_ValidSample_PrintsLoadedLine` | integration (CLI) | SPEC-RECIPES-035, 037, 039 |
| `Test_SeedSubcommand_InvalidYaml_ExitsNonZero` | integration (CLI) | SPEC-RECIPES-038 |

## 9.2 Test fixtures

Fixtures live under `internal/seed/testdata/`:

```
internal/seed/testdata/
  valid/
    breakfast.yaml          # minimal, valid
    lunch.yaml              # minimal, valid (different slug)
  invalid_taxonomy/
    bad-cooking-method.yaml # cooking_methods: [microwave]   ← not in canonical set
  duplicate_slug/
    a.yaml                  # slug: same-slug
    b.yaml                  # slug: same-slug
  missing_field/
    no-title.yaml           # missing `title:`
  malformed/
    parse-error.yaml        # broken YAML
```

Tests SHALL use `t.TempDir()` to copy fixtures or directly point `Load` at
`testdata/<scenario>/`. Tests SHALL NOT mutate fixture files at runtime.

## 9.3 Helpers

- `t.Helper()` SHALL be used in any helper that asserts on errors or counts.
- `t.Context()` SHALL be used wherever `Load`/`Store` need a `context.Context`.
- Integration tests SHALL connect to the compose Postgres via `DATABASE_URL`
  from the environment; if unset, tests SHALL `t.Skip("DATABASE_URL not set")`.
- Integration tests SHALL run inside their own database transaction or use
  `TRUNCATE recipes, ingredients, recipe_ingredients RESTART IDENTITY
  CASCADE` in `t.Cleanup()` to keep the DB clean.

## 9.4 Mocking policy

- Pure validation logic SHALL be tested without mocks.
- The pgx pool SHALL NOT be mocked; integration tests use the real compose
  Postgres (consistent with SPEC-DB testing strategy).
- The filesystem SHALL NOT be mocked; tests use real `t.TempDir()` /
  `testdata/` directories.

## 9.5 Coverage target

- `internal/domain/recipe.go` and `taxonomy.go`: ≥ 95 % line coverage
  (pure functions, no excuses).
- `internal/seed/recipes.go`: ≥ 90 % line coverage.
- `internal/seed/store.go`: ≥ 80 % line coverage (rollback paths exercised
  via fault injection — closing the pool mid-call).
- Project floor remains 80 % per repo policy.

## 9.6 AC traceability

The four AC #7 unit-test families are mapped one-for-one in the table
above:
- "invalid taxonomy" → `Test_LoadRecipes_InvalidTaxonomy_ReturnsError`
- "duplicate slug" → `Test_LoadRecipes_DuplicateSlug_ReturnsError`
- "missing required field" → `Test_LoadRecipes_MissingRequiredField_ReturnsError`
- "valid sample" → `Test_LoadRecipes_ValidSample_ReturnsRecipes`
