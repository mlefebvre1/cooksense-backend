# SPEC-DISCOVER — Appendix B · Implementation Tasks

[← Index](SPEC-DISCOVER-00-index.md)

Each task is atomic, ordered, and traceable to SPEC-DISCOVER-NNN IDs.
Implementation SHALL proceed in this order; reordering is allowed only
inside groups where dependencies do not cross.

| # | Task | SPEC-IDs | Depends on |
|---|------|----------|------------|
| T-01 | Create `internal/recipes/dto.go` with `RecipeBrief`, `RecipeFull` (embedding brief), and `IngredientView`, all with `snake_case` JSON tags. | 001, 002, 003, 004 | — |
| T-02 | Write DTO unit tests: `Test_RecipeBriefJSON_RoundTrip_MatchesAPIShape`, `Test_RecipeFullJSON_EmbedsBriefAndAddsIngredientsAndSteps`, `Test_IngredientViewJSON_FieldNames_MatchAPIShape`. | 002, 003, 004 | T-01 |
| T-03 | Create `internal/recipes/repo.go` skeleton: `Repo` interface, `ErrNotFound` sentinel, `PgRepo` struct, `NewPgRepo(pool *pgxpool.Pool) *PgRepo`. | 006, 007, 008, 016 | T-01 |
| T-04 | Implement `(*PgRepo).Discover(ctx, uid, limit)` with the SQL from §5.4; ensure empty slices on null array columns. | 005, 009, 010, 011, 012 | T-03 |
| T-05 | Implement `(*PgRepo).GetBySlug(ctx, slug)` with the two queries from §5.5; translate `pgx.ErrNoRows` to `fmt.Errorf("%w", ErrNotFound)`. | 005, 013, 014, 015 | T-03 |
| T-06 | Write repo integration tests against compose Postgres (`Test_PgRepoDiscover_*`, `Test_PgRepoGetBySlug_*`). Use `t.Skip("DATABASE_URL not set")` when env unset; truncate in `t.Cleanup`. | 007, 009–015 (see §9.1) | T-04, T-05 |
| T-07 | Create `internal/recipes/service.go` with `Service`, `NewService`, `DefaultDiscoverLimit`, `MaxDiscoverLimit`, and `(*Service).Discover` clamping per §6.3. | 017, 018, 019, 021 | T-03 |
| T-08 | Implement `(*Service).GetBySlug` (delegate + propagate). Add structured `slog.Info` lines per SPEC-DISCOVER-022. | 020, 022 | T-07 |
| T-09 | Write service unit tests with an in-memory `fakeRepo`: `Test_ServiceDiscover_LimitClamping_AppliesBounds`, `Test_ServiceDiscover_DelegatesUidAndClampedLimitToRepo`, `Test_ServiceGetBySlug_PropagatesErrNotFound_ErrorsIsTrue`, `Test_ServiceDiscover_LogsStructuredAttributes`. | 017–022, 040 | T-08 |
| T-10 | Create `internal/recipes/handler.go` with `Handler`, `NewHandler`, `(*Handler).RegisterRoutes`, `writeJSON`, `writeError` helpers. | 023, 024, 032 | T-07 |
| T-11 | Implement `(*Handler).Discover` with limit parsing, auth user retrieval, service delegation, `discoverResponse` wrapper, and empty-array guarantee. | 025, 026, 027, 031, 033 | T-10 |
| T-12 | Implement `(*Handler).GetBySlug` with `r.PathValue("slug")`, service delegation, and `errors.Is(err, ErrNotFound)` branch for 404. | 028, 029, 030, 031, 033 | T-10 |
| T-13 | Write handler unit tests with a fake service: `Test_HandlerDiscover_NegativeLimit_FallsBackToDefault`, `Test_HandlerDiscover_OverMaxLimit_ClampedToMax`, `Test_HandlerDiscover_MissingAuthUser_Returns500Internal`, `Test_HandlerGetRecipe_RepoFails_Returns500Internal`, `Test_RegisterRoutes_RegistersBothPatterns`, `Test_DiscoverResponseJSON_NilSlice_EncodesEmptyArray`. | 005, 023–033, 041 | T-11, T-12 |
| T-14 | Write handler integration tests against compose Postgres + `fakeVerifier`: `Test_HandlerDiscover_HappyPath_Returns200WithEnvelope`, `Test_HandlerDiscover_ExcludesReactedRecipes_NeverReturnsThem`, `Test_HandlerGetRecipe_KnownSlug_Returns200FullRecipe`, `Test_HandlerGetRecipe_UnknownSlug_Returns404NotFound`. | 036, 037, 038, 039 | T-13 |
| T-15 | Update `internal/recipes/doc.go` package comment per SPEC-DISCOVER-035 and SPEC-DISCOVER-001. Add Go Doc comments citing SPEC-DISCOVER IDs on every exported symbol introduced. | 035, 001 | T-12 |
| T-16 | Wire the feature in `cmd/cooksense-server/main.go` exactly as §6.5: `NewPgRepo` → `NewService` → `NewHandler` → `RegisterRoutes(mux, authMW)`. Add a smoke test `Test_MainWiring_RecipesRoutesReachable_ReturnsNon404OnGet`. | 034 | T-12, T-15 |
| T-17 | Update `docs/architecture/api.md`: replace "Excludes `DISLIKE`d recipes by default" with the any-reaction wording, citing SPEC-DISCOVER-012. | 012 (doc edit) | T-04 |
| T-18 | Update `README.md` with the "Recipe discovery and detail" subsection per SPEC-DISCOVER-042. | 042 | T-16 |
| T-19 | Run `make lint test`; verify `go vet ./...`, `go build ./...`, coverage targets met (§9.5). Run `make up && make migrate && make seed` then exercise both endpoints with `curl` + a fake bearer token (or via `Test_MainWiring_*`) end-to-end. | All | T-14, T-16, T-17, T-18 |

## B.1 Critical-path graph

```
T-01 ─► T-02
  │
  └─► T-03 ─► T-04 ─► T-06
        │      │
        │      └─► T-17
        │
        └─► T-05 ─► T-06
                  │
                  └─► T-07 ─► T-08 ─► T-09
                              │
                              └─► T-10 ─► T-11 ┐
                                          │     ├─► T-13 ─► T-14 ─► T-19
                                          T-12 ┘            │
                                            │               └─► T-16 ─► T-18 ─► T-19
                                            └─► T-15 ─────────► T-16
```

T-19 is the gate before PR submission and SHALL only run when every
prior task is complete.

## B.2 Estimated effort

| Group | Tasks | Effort |
|-------|-------|--------|
| DTOs | T-01, T-02 | ~1 h |
| Repo | T-03, T-04, T-05, T-06 | ~3 h |
| Service | T-07, T-08, T-09 | ~1.5 h |
| Handler | T-10, T-11, T-12, T-13 | ~2.5 h |
| Integration tests | T-14 | ~1.5 h |
| Wiring + docs | T-15, T-16, T-17, T-18 | ~1 h |
| Verification | T-19 | ~30 min |

Total estimate: ~11 h focused work, matching the "M" estimate on the
story.

## B.3 PR strategy

The implementation MAY land in a single PR (`feat/recipes-discover`) or be
split into two:

1. `feat/recipes-discover-domain` — T-01..T-09 (DTOs, repo, service +
   their tests). Lands a fully tested data-access layer with no HTTP.
2. `feat/recipes-discover-http` — T-10..T-19 (handler, wiring,
   integration tests, README/api.md updates).

Either approach SHALL keep each commit body citing the SPEC-DISCOVER IDs
it implements, per the repo's commit-message contract. PRs SHALL be
rebased onto `main` before merge (no merge commits).
