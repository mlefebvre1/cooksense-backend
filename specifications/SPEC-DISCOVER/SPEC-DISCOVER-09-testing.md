# SPEC-DISCOVER — §9 Testing

[← Index](SPEC-DISCOVER-00-index.md)

## 9.1 Required named tests

Test names SHALL follow `Test{What}_{Condition}_{ExpectedOutcome}` per
the project guideline. Each test SHALL document the SPEC-DISCOVER IDs it
verifies in its docstring (forward traceability).

| Test name | Layer | Verifies |
|-----------|-------|----------|
| `Test_RecipeBriefJSON_RoundTrip_MatchesAPIShape` | unit (dto) | SPEC-DISCOVER-002 |
| `Test_RecipeFullJSON_EmbedsBriefAndAddsIngredientsAndSteps` | unit (dto) | SPEC-DISCOVER-003 |
| `Test_IngredientViewJSON_FieldNames_MatchAPIShape` | unit (dto) | SPEC-DISCOVER-004 |
| `Test_DiscoverResponseJSON_NilSlice_EncodesEmptyArray` | unit (handler) | SPEC-DISCOVER-005, 027 |
| `Test_PgRepoDiscover_ExcludesAlreadyReacted_ReturnsRemainder` | integration (repo+pg) | SPEC-DISCOVER-009, 011, 012 |
| `Test_PgRepoDiscover_RandomOrderAndLimit_ReturnsAtMostLimit` | integration (repo+pg) | SPEC-DISCOVER-010 |
| `Test_PgRepoGetBySlug_KnownSlug_ReturnsFullRecipeWithSortedIngredients` | integration (repo+pg) | SPEC-DISCOVER-013, 015 |
| `Test_PgRepoGetBySlug_UnknownSlug_ReturnsErrNotFound` | integration (repo+pg) | SPEC-DISCOVER-007, 014 |
| `Test_ServiceDiscover_LimitClamping_AppliesBounds` | unit (service+fakeRepo) | SPEC-DISCOVER-018, 019 (SPEC-DISCOVER-040) |
| `Test_ServiceDiscover_DelegatesUidAndClampedLimitToRepo` | unit (service+fakeRepo) | SPEC-DISCOVER-006, 017 |
| `Test_ServiceGetBySlug_PropagatesErrNotFound_ErrorsIsTrue` | unit (service+fakeRepo) | SPEC-DISCOVER-020 |
| `Test_ServiceDiscover_LogsStructuredAttributes` | unit (service+fakeRepo+slog buffer) | SPEC-DISCOVER-022 |
| `Test_HandlerDiscover_NegativeLimit_FallsBackToDefault` | unit (handler+fakeService) | SPEC-DISCOVER-018, 025 (SPEC-DISCOVER-041) |
| `Test_HandlerDiscover_OverMaxLimit_ClampedToMax` | unit (handler+fakeService) | SPEC-DISCOVER-018, 025 |
| `Test_HandlerDiscover_MissingAuthUser_Returns500Internal` | unit (handler+httptest) | SPEC-DISCOVER-026, 031 |
| `Test_HandlerDiscover_HappyPath_Returns200WithEnvelope` | integration (handler+httptest+pg+fakeVerifier) | SPEC-DISCOVER-024, 027, 032 (SPEC-DISCOVER-036) |
| `Test_HandlerDiscover_ExcludesReactedRecipes_NeverReturnsThem` | integration (handler+httptest+pg+fakeVerifier) | SPEC-DISCOVER-009, 012 (SPEC-DISCOVER-037) |
| `Test_HandlerGetRecipe_KnownSlug_Returns200FullRecipe` | integration (handler+httptest+pg+fakeVerifier) | SPEC-DISCOVER-028, 029, 032 (SPEC-DISCOVER-038) |
| `Test_HandlerGetRecipe_UnknownSlug_Returns404NotFound` | integration (handler+httptest+pg+fakeVerifier) | SPEC-DISCOVER-014, 030 (SPEC-DISCOVER-039) |
| `Test_HandlerGetRecipe_RepoFails_Returns500Internal` | unit (handler+failingFakeService) | SPEC-DISCOVER-031 |
| `Test_RegisterRoutes_RegistersBothPatterns` | unit (handler) | SPEC-DISCOVER-024 |
| `Test_MainWiring_RecipesRoutesReachable_ReturnsNon404OnGet` | smoke (compose + binary) | SPEC-DISCOVER-034 |

The four AC #7 integration scenarios are mapped one-for-one in the table
above (see SPEC-DISCOVER-036..039).

## 9.2 Test fixtures

Fixtures live under `internal/recipes/testdata/`:

```
internal/recipes/testdata/
  seed/
    recipe-with-three-ingredients.yaml   # used by detail tests
    recipe-without-ingredients.yaml      # invalid by domain rules — NEVER seeded; reserved for negative tests
  fixtures/
    user-with-no-reactions.go            # helper builders, NOT YAML
    user-with-mixed-reactions.go         # helper builders
```

Integration tests SHALL:
- Load fixtures via `seed.Load` + `seed.Store` (SPEC-RECIPES) against the
  compose Postgres in `t.Cleanup`-wrapped transactions.
- Insert `users` and `user_reactions` rows via direct SQL helpers
  (`testhelpers.InsertReaction(t, pool, uid, slug, kind)`).
- Truncate `user_reactions, recipe_ingredients, ingredients, recipes,
  users RESTART IDENTITY CASCADE` in `t.Cleanup`.

Tests SHALL NOT mutate fixture YAML files at runtime.

## 9.3 Helpers

- `t.Helper()` SHALL be used in any helper that asserts on errors,
  status codes, or counts.
- `t.Context()` SHALL be used wherever the service/repo/handler under
  test needs a `context.Context` (per repo guideline).
- Integration tests SHALL connect to the compose Postgres via
  `DATABASE_URL` from the environment; if unset, they SHALL
  `t.Skip("DATABASE_URL not set")`.
- HTTP tests SHALL use `httptest.NewServer` plus the project's
  `auth.fakeVerifier` (SPEC-AUTH-006) seeded with one or more
  `(token → User)` pairs.
- Each request in HTTP tests SHALL set
  `Authorization: Bearer <fake-token>` exactly as the production client
  would.

## 9.4 Mocking policy

- The `Repo` interface (SPEC-DISCOVER-006) SHALL be the **only** seam
  that is faked for unit tests. A `fakeRepo` struct in `service_test.go`
  SHALL satisfy `Repo` with in-memory data.
- The pgx pool SHALL NOT be mocked; integration tests use the real
  compose Postgres (consistent with SPEC-RECIPES and SPEC-DB testing
  strategy).
- The Firebase verifier SHALL be replaced by `fakeVerifier` (SPEC-AUTH);
  no real network calls SHALL leave the test process.
- HTTP tests SHALL use `httptest`; `http.DefaultClient` SHALL NOT escape
  the test boundary.
- If a single test mocks more than three collaborators, the design is
  wrong — refactor before merging.

## 9.5 Coverage target

- `internal/recipes/dto.go`: ≥ 90 % line coverage (trivial; covered by
  the JSON round-trip tests).
- `internal/recipes/service.go`: ≥ 95 % line coverage (pure logic, no
  excuses).
- `internal/recipes/handler.go`: ≥ 90 % line coverage (httptest + table
  tests for the response branches).
- `internal/recipes/repo.go`: ≥ 85 % line coverage (integration tests
  exercise both happy paths and the not-found branch; rare
  internal-error branches MAY be exempt with a justified `//nolint`).
- Project floor remains 80 % per repo policy (SPECIFICATIONS §0.3).

## 9.6 AC traceability (story 07)

| Story AC | Verifying test(s) |
|----------|-------------------|
| AC #1 — `GET /api/recipes/discover?limit=N` returns up to N, default 10, max 25, excludes already-reacted, randomized, `<RecipeBrief>` shape | SPEC-DISCOVER-036, 037, plus 002/005/010 unit tests |
| AC #2 — `GET /api/recipes/{slug}` returns `<RecipeFull>` with sorted ingredients and stored-order steps; 404 on unknown | SPEC-DISCOVER-038, 039 |
| AC #3 — Repo is the only DB-touching layer; exposes `Discover` and `GetBySlug` | All SPEC-DISCOVER-006..016 tests |
| AC #4 — Service enforces clamping and orchestrates | SPEC-DISCOVER-040, plus the service unit tests |
| AC #5 — Handler only handles JSON/HTTP | SPEC-DISCOVER-041, plus all handler-layer tests |
| AC #6 — Integration test (compose Postgres, seeded data, `fakeVerifier`) covers the four scenarios | SPEC-DISCOVER-036, 037, 038, 039 |

## 9.7 Determinism guidance for `random()` ordering

The discover query is non-deterministic. Tests SHALL verify discover
behavior via:

1. **Cardinality** — `len(result) == expected`.
2. **Set membership** — `slices.ContainsFunc(result, …)` for inclusion;
   `!slices.ContainsFunc(reacted, r.Slug)` for exclusion.
3. **Limit ceiling** — `len(result) <= clampedLimit`.

Tests SHALL NOT assert on the specific order of returned recipes; doing
so makes the suite flaky and violates SPEC-DISCOVER-010.
