# SPEC-REACTIONS — §9 Testing

[← Index](SPEC-REACTIONS-00-index.md)

## 9.1 Named tests

Each test SHALL be named in the form
`Test{What}_{Condition}_{ExpectedOutcome}`.

| Test name | File | SPEC-IDs verified |
|-----------|------|-------------------|
| `Test_ParseReactionKind_AcceptsAllCanonicalValues` | `internal/domain/reaction_test.go` | 001, 002, 003 |
| `Test_ParseReactionKind_RejectsLowercase` | `internal/domain/reaction_test.go` | 002, 057 |
| `Test_ParseReactionKind_RejectsEmpty` | `internal/domain/reaction_test.go` | 002 |
| `Test_ReactionKind_Valid_TrueOnlyForCanonical` | `internal/domain/reaction_test.go` | 003 |
| `Test_DomainReaction_HasExpectedFields` | `internal/domain/reaction_test.go` | 004 |
| `Test_DomainImports_NoIO` | `internal/domain/imports_test.go` | 005 |
| `Test_Repo_SatisfiesInterface` | `internal/reactions/repo_test.go` | 006, 007 |
| `Test_Reactions_Upsert_FromNoRow_InsertsRow` | `internal/reactions/repo_integration_test.go` | 008, 009, 014, 045 |
| `Test_Reactions_Upsert_OverwritesPreviousKind` | `internal/reactions/repo_integration_test.go` | 008, 009, 046 |
| `Test_Reactions_Delete_RemovesExistingRow` | `internal/reactions/repo_integration_test.go` | 010, 047 |
| `Test_Reactions_Delete_Idempotent_NoRow_ReturnsNil` | `internal/reactions/repo_integration_test.go` | 010, 048 |
| `Test_Reactions_ListByKind_OrderedByCreatedAtDesc` | `internal/reactions/repo_integration_test.go` | 011, 012, 053 |
| `Test_Reactions_ListByKind_EmptyResult_ReturnsEmptySlice` | `internal/reactions/repo_integration_test.go` | 013, 056 |
| `Test_Service_SetReaction_UnknownSlug_ReturnsErrRecipeNotFound` | `internal/reactions/service_test.go` | 015, 016, 018 |
| `Test_Service_RemoveReaction_UnknownSlug_ReturnsErrRecipeNotFound` | `internal/reactions/service_test.go` | 019 |
| `Test_Service_MyRecipes_DislikeKind_ReturnsErrInvalidMyRecipesKind` | `internal/reactions/service_test.go` | 020, 021 |
| `Test_Service_MyRecipes_LikeKind_DelegatesToRepo` | `internal/reactions/service_test.go` | 017, 020 |
| `Test_Reactions_Post_HappyPath_Returns200WithUpdatedAt` | `internal/api/reactions_post_handler_test.go` | 022, 026, 027, 028 |
| `Test_Reactions_Post_MalformedJSON_Returns400` | `internal/api/reactions_post_handler_test.go` | 023, 051 |
| `Test_Reactions_Post_MissingSlug_Returns400` | `internal/api/reactions_post_handler_test.go` | 023 |
| `Test_Reactions_Post_UnknownSlug_Returns404` | `internal/api/reactions_post_handler_test.go` | 024, 049 |
| `Test_Reactions_Post_InvalidKind_Returns422` | `internal/api/reactions_post_handler_test.go` | 025, 050 |
| `Test_Reactions_Delete_Slug_Returns204` | `internal/api/reactions_delete_handler_test.go` | 029, 030, 032 |
| `Test_Reactions_Delete_Idempotent_NoRow_Returns204` | `internal/api/reactions_delete_handler_test.go` | 030, 048 |
| `Test_Reactions_Delete_UnknownSlug_Returns404` | `internal/api/reactions_delete_handler_test.go` | 031 |
| `Test_MyRecipes_DefaultsToLike_WhenKindOmitted` | `internal/api/me_recipes_handler_test.go` | 033, 034, 052 |
| `Test_MyRecipes_OrderedByCreatedAtDesc` | `internal/api/me_recipes_handler_test.go` | 037, 053 |
| `Test_MyRecipes_DislikeKind_Returns400` | `internal/api/me_recipes_handler_test.go` | 035, 054 |
| `Test_MyRecipes_UnknownKind_Returns400` | `internal/api/me_recipes_handler_test.go` | 035, 055 |
| `Test_MyRecipes_EmptyResult_ReturnsEmptyArray` | `internal/api/me_recipes_handler_test.go` | 036, 056 |
| `Test_Handlers_NonOK_WriteEnvelopedError` | `internal/api/error_envelope_test.go` | 038 |
| `Test_Handlers_LogAttrs_OmitTokenAndUID` | `internal/api/logging_test.go` | 039 |
| `Test_Handlers_PassRequestContext` | `internal/api/context_test.go` | 040 |
| `Test_Wiring_RegistersThreeRoutes` | `cmd/cooksense-server/main_test.go` | 041, 042 |

## 9.2 Test categories

### 9.2.1 Pure unit tests

- `internal/domain/reaction_test.go` — table-driven tests over
  `ParseReactionKind`, `Valid`, struct shape.
- `internal/reactions/service_test.go` — uses an in-memory fake
  `domain.ReactionRepository` and a fake `recipeResolver`. **Two collaborators
  is the maximum** (CLAUDE.md mocking rule).

### 9.2.2 Integration tests (compose Postgres)

- `internal/reactions/repo_integration_test.go` — real `*pgxpool.Pool`,
  build-tagged `//go:build integration` per the project convention. Uses
  `t.Context()` and `t.Cleanup` to truncate `user_reactions` between cases.
- `internal/api/*_handler_test.go` — spin up an `httptest.Server` mounted
  on the real mux with the auth middleware **stubbed** to set a known
  `firebase_uid` in the context. This avoids dragging Firebase into the
  test loop.

### 9.2.3 Wiring smoke

- `cmd/cooksense-server/main_test.go` — verifies that the three SPEC-REACTIONS
  routes are registered on the mux returned by the bootstrap helper.

## 9.3 Fixtures

| Fixture | Purpose |
|---------|---------|
| `seed/recipes/_sample.yaml` (existing) | Provides at least one valid recipe slug for handler/integration tests. |
| `seed/recipes/_sample2.yaml` (NEW under this story) | Second recipe so list-ordering tests have ≥ 2 distinct slugs to react to. SHALL be underscore-prefixed so it is excluded from the curated count (per SPEC-RECIPES-041). |
| Test user | Tests SHALL provision a user row directly via SQL (`INSERT INTO users(firebase_uid) VALUES ('test-uid-...')`) inside `t.Cleanup`-tracked transactions, since no Firebase exchange happens in tests. |

## 9.4 Mocking policy

- The pgx pool SHALL NOT be mocked. Integration tests use the real compose
  Postgres (consistent with SPEC-DB / SPEC-RECIPES testing strategy).
- The HTTP transport SHALL NOT be mocked beyond `httptest.NewRecorder` /
  `httptest.NewServer`.
- Pure business logic in `service.go` SHALL be tested without touching the
  database, using in-memory fakes that satisfy the small interfaces
  declared in §6.4.
- If a test needs more than three collaborators, the design is wrong —
  refactor the production code.

## 9.5 Coverage target

- `internal/domain/reaction.go`: ≥ 95 % line coverage (pure functions).
- `internal/reactions/service.go`: ≥ 95 % line coverage (pure orchestration).
- `internal/reactions/repo.go`: ≥ 90 % line coverage (integration).
- `internal/api/reactions_*_handler.go` and `me_recipes_handler.go`:
  ≥ 90 % line coverage.

The aggregated coverage SHALL be ≥ 80 % as enforced by `make test`.

## 9.6 Determinism

- Time-sensitive assertions (`updated_at` on POST) SHALL compare with
  tolerance: `time.Since(at) < 5*time.Second && at.After(start)`.
- List-ordering tests SHALL insert reactions with explicit `created_at`
  values via `UPDATE user_reactions SET created_at = $2 WHERE ...` after
  the upsert, so ordering is deterministic regardless of clock resolution.
