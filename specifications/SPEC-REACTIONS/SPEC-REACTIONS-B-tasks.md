# SPEC-REACTIONS — Appendix B · Implementation Tasks

[← Index](SPEC-REACTIONS-00-index.md)

Each task is atomic, ordered, and traceable to SPEC-REACTIONS-NNN IDs.
Implementation SHALL proceed in this order; reordering is allowed only
inside groups where dependencies do not cross.

| # | Task | SPEC-IDs | Depends on |
|---|------|----------|------------|
| T-01 | Create `internal/domain/reaction.go` with `ReactionKind`, the three constants, `Valid`, `Reaction` struct, and the two sentinel errors. | 001, 003, 004, 005, 006 | — |
| T-02 | Add `ParseReactionKind(s string) (ReactionKind, error)` to the same file. | 002 | T-01 |
| T-03 | Define `domain.ReactionRepository` interface (in `reaction.go` or a colocated `repository.go`) with the three methods. | 006 | T-01 |
| T-04 | (Conditional) If Story 07 has not landed `recipes.Repo.ResolveSlug` and `domain.ErrRecipeNotFound`, add them under `internal/recipes/repo.go` and `internal/domain/recipe.go` respectively, with their own doc comments. | 015, 016 | T-01 |
| T-05 | Write domain unit tests in `internal/domain/reaction_test.go` (table-driven `ParseReactionKind`, `Valid`, struct shape) and the import-purity test in `imports_test.go`. | 001–006, 057 | T-02, T-03 |
| T-06 | Create `internal/reactions/doc.go` with the package comment of SPEC-REACTIONS-043. | 043 | T-03 |
| T-07 | Create `internal/reactions/repo.go` with `Repo`, `NewRepo`, the compile-time interface assertion, `Upsert`, `Delete`, and `ListByKind`. SQL strings live as unexported `const` at the top of the file. | 007–014 | T-03, T-06 |
| T-08 | Author `internal/reactions/repo_integration_test.go` (`//go:build integration`) covering insert-from-no-row, overwrite, delete (existing + idempotent), list ordering, list empty result. Use `t.Context()` and a fresh schema (TRUNCATE) per test. | 045, 046, 047, 048, 053, 056 | T-07 |
| T-09 | Create `internal/reactions/service.go` with the unexported `recipeResolver` interface, `Service`, `NewService`, `SetReaction`, `RemoveReaction`, `MyRecipes`. | 015, 017, 018, 019, 020, 021 | T-04, T-07 |
| T-10 | Author `internal/reactions/service_test.go` with in-memory fakes (≤ 2 collaborators) covering unknown slug, dislike-rejection, like-delegation, `RemoveReaction` unknown slug. | 015, 016, 017, 018, 019, 020, 021 | T-09 |
| T-11 | Add the second seed fixture `seed/recipes/_sample2.yaml` (underscore-prefixed) for ordering tests. | §9.3 fixtures | T-08 |
| T-12 | Create `internal/api/reactions_post_handler.go` with `PostReaction(svc, log)` constructor returning `http.HandlerFunc`. Implement the algorithm of §6.5. | 022, 023, 024, 025, 026, 027, 028 | T-09 |
| T-13 | Create `internal/api/reactions_delete_handler.go` with `DeleteReaction(svc, log)`. Implement the algorithm of §6.6. | 029, 030, 031, 032 | T-09 |
| T-14 | Create `internal/api/me_recipes_handler.go` with `MyRecipes(svc, log)`. Implement the algorithm of §6.7, including default-to-LIKE and DISLIKE rejection. | 033, 034, 035, 036, 037 | T-09 |
| T-15 | Author handler tests `internal/api/reactions_post_handler_test.go`, `reactions_delete_handler_test.go`, `me_recipes_handler_test.go` using `httptest` and a stubbed auth middleware. | 045–056 | T-12, T-13, T-14, T-11 |
| T-16 | Author cross-handler tests for the error envelope, slog attrs, and request-context propagation. | 038, 039, 040 | T-15 |
| T-17 | Wire the three routes in `cmd/cooksense-server/main.go`'s bootstrap (constructing `reactionsRepo`, `reactionsSvc`, registering routes behind `auth.Require`). Add the wiring smoke test in `cmd/cooksense-server/main_test.go`. | 041, 042 | T-12, T-13, T-14 |
| T-18 | Run `make lint test`. Verify `go build ./...`, `go vet ./...`, `golangci-lint run`, coverage targets met (§9.5). Run a manual smoke against the compose Postgres for each endpoint (`curl` with a stub auth header in test mode). | All | T-11, T-15, T-16, T-17 |

## B.1 Critical-path graph

```
T-01 → T-02 → T-03 ┬─► T-04 (conditional)
                   │
                   ├─► T-05 (domain tests)
                   │
                   └─► T-06 → T-07 ┬─► T-08 (repo IT)
                                    │
                                    └─► T-09 ┬─► T-10 (service tests)
                                              │
                                              ├─► T-12 ┐
                                              ├─► T-13 ├─► T-15 → T-16 ┐
                                              ├─► T-14 ┘                ├─► T-18
                                              └─► T-17 ─────────────────┘
                                                              ↑
                                                         T-11 ┘
```

T-18 is the gate before PR submission and SHALL only run when every prior
task is complete.

## B.2 Estimated effort

| Group | Tasks | Effort |
|-------|-------|--------|
| Domain + interface | T-01..T-05 | ~1.5 h |
| Repo + integration | T-06..T-08, T-11 | ~2 h |
| Service | T-09, T-10 | ~1 h |
| Handlers | T-12..T-14 | ~2 h |
| Handler tests | T-15, T-16 | ~2 h |
| Wiring + smoke | T-17 | ~30 min |
| Verification | T-18 | ~30 min |

Total estimate: ~9.5 h focused work, matching the "M" estimate on the
story.

## B.3 PR strategy

The implementation MAY land in a single PR (`feat/reactions-my-recipes`) or
be split into two:

1. `feat/reactions-domain-repo` — T-01..T-08, T-11 (domain + repo +
   integration tests, no HTTP yet).
2. `feat/reactions-handlers` — T-09..T-10, T-12..T-18 (service, handlers,
   wiring, tests).

Either approach SHALL keep each commit body citing the SPEC-REACTIONS IDs
it implements, per the repo's commit-message contract.
