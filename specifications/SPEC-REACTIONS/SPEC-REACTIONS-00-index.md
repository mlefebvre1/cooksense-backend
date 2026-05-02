# SPEC-REACTIONS — Reactions API + "My recipes" listing

**Story:** 08 — Reactions and "My recipes"
**Status:** Draft → Final
**Date:** 2026-05-02
**Authors:** sdd-spec-author

---

## Purpose

This specification governs the reactions domain (`internal/domain/reaction.go`),
the reactions repository (`internal/reactions/repo.go`), the orchestration
service (`internal/reactions/service.go`), the HTTP handlers for
`POST /api/reactions`, `DELETE /api/reactions/{slug}`, and
`GET /api/me/recipes` (under `internal/api/`), the wiring in
`cmd/cooksense-server/main.go`, and all associated documentation and tests
required by Story 08.

Out of scope: reaction history (D-0008), bulk operations, exposure of the
`DISLIKE` kind in `GET /api/me/recipes` (story explicitly forbids it),
recipe discovery (story 07 — only consumed here), recipe detail (story 07),
search (story 09), lessons (story 10).

---

## File map

| File | Contents |
|------|----------|
| [SPEC-REACTIONS-01-preamble](SPEC-REACTIONS-01-preamble.md) | AI constraints, authorship, traceability rules |
| [SPEC-REACTIONS-02-introduction](SPEC-REACTIONS-02-introduction.md) | Story relationship, decision references |
| [SPEC-REACTIONS-03-goals](SPEC-REACTIONS-03-goals.md) | Goals, non-goals, constraints |
| [SPEC-REACTIONS-04-context](SPEC-REACTIONS-04-context.md) | External dependencies, package boundary map |
| [SPEC-REACTIONS-05-architecture](SPEC-REACTIONS-05-architecture.md) | Request flow, dependency graph, transaction shape |
| [SPEC-REACTIONS-06-packages](SPEC-REACTIONS-06-packages.md) | All SPEC-REACTIONS-NNN requirements: types, signatures, SQL, HTTP |
| [SPEC-REACTIONS-07-configuration](SPEC-REACTIONS-07-configuration.md) | Environment variables, route registration |
| [SPEC-REACTIONS-08-build](SPEC-REACTIONS-08-build.md) | Build, lint, vet, dependency rules |
| [SPEC-REACTIONS-09-testing](SPEC-REACTIONS-09-testing.md) | Test strategy, named tests, fixtures |
| [SPEC-REACTIONS-10-documentation](SPEC-REACTIONS-10-documentation.md) | Doc comments, README impact |
| [SPEC-REACTIONS-A-checklist](SPEC-REACTIONS-A-checklist.md) | Specification completeness checklist |
| [SPEC-REACTIONS-B-tasks](SPEC-REACTIONS-B-tasks.md) | Ordered implementation task list |

---

## SPEC-ID registry

| ID | Summary | Section |
|----|---------|---------|
| SPEC-REACTIONS-001 | `internal/domain` exposes `ReactionKind` string-typed enum with constants `ReactionLike`, `ReactionDislike`, `ReactionTryLater` | §6.1 |
| SPEC-REACTIONS-002 | `domain.ParseReactionKind(s string) (ReactionKind, error)` accepts only `"LIKE"`, `"DISLIKE"`, `"TRY_LATER"` | §6.1 |
| SPEC-REACTIONS-003 | `domain.ReactionKind.Valid() bool` returns true only for the three canonical values | §6.1 |
| SPEC-REACTIONS-004 | `domain.Reaction` struct fields: `FirebaseUID string`, `RecipeID int64`, `Kind ReactionKind`, `CreatedAt time.Time` | §6.1 |
| SPEC-REACTIONS-005 | `domain` package contains zero I/O imports (no `database/sql`, `pgx`, `net/http`, `os`) | §6.1 |
| SPEC-REACTIONS-006 | `domain.ReactionRepository` interface declares `Upsert`, `Delete`, `ListByKind` with the signatures of §6.2 | §6.1 |
| SPEC-REACTIONS-007 | `reactions.Repo` constructor `NewRepo(pool *pgxpool.Pool) *Repo` returns a struct that satisfies `domain.ReactionRepository` | §6.2 |
| SPEC-REACTIONS-008 | `Repo.Upsert(ctx, uid string, recipeID int64, kind domain.ReactionKind) (time.Time, error)` performs the INSERT … ON CONFLICT UPDATE statement of §6.2.2 and returns the resulting `created_at` | §6.2 |
| SPEC-REACTIONS-009 | The upsert SQL is exactly the statement specified in the story technical notes (`ON CONFLICT (firebase_uid, recipe_id) DO UPDATE SET kind = EXCLUDED.kind, created_at = now() RETURNING created_at`) | §6.2 |
| SPEC-REACTIONS-010 | `Repo.Delete(ctx, uid string, recipeID int64) error` issues `DELETE FROM user_reactions WHERE firebase_uid=$1 AND recipe_id=$2` and never returns `pgx.ErrNoRows` (deleting zero rows is a success) | §6.2 |
| SPEC-REACTIONS-011 | `Repo.ListByKind(ctx, uid string, kind domain.ReactionKind) ([]domain.RecipeBrief, error)` returns `RecipeBrief` rows joined with `recipes`, ordered by `user_reactions.created_at DESC` | §6.2 |
| SPEC-REACTIONS-012 | `ListByKind` SQL joins `user_reactions` with `recipes` on `recipes.id = user_reactions.recipe_id` and filters on `firebase_uid` and `kind` | §6.2 |
| SPEC-REACTIONS-013 | `ListByKind` returns an empty (non-nil) slice when no rows match | §6.2 |
| SPEC-REACTIONS-014 | The repository never accepts or stores anything other than the three canonical `ReactionKind` values; callers SHALL validate before calling | §6.2 |
| SPEC-REACTIONS-015 | `recipes` package exposes a `ResolveSlug(ctx, slug string) (int64, error)` lookup returning `domain.ErrRecipeNotFound` when missing | §6.3 |
| SPEC-REACTIONS-016 | `domain.ErrRecipeNotFound` sentinel error is exported and matched via `errors.Is` | §6.3 |
| SPEC-REACTIONS-017 | `reactions.Service` orchestrates slug → id resolution then delegates to `Repo.Upsert`/`Repo.Delete`/`Repo.ListByKind` | §6.4 |
| SPEC-REACTIONS-018 | `reactions.Service.SetReaction(ctx, uid, slug, kind)` returns `(time.Time, error)` and surfaces `domain.ErrRecipeNotFound` unchanged | §6.4 |
| SPEC-REACTIONS-019 | `reactions.Service.RemoveReaction(ctx, uid, slug)` returns `error` and surfaces `domain.ErrRecipeNotFound` unchanged | §6.4 |
| SPEC-REACTIONS-020 | `reactions.Service.MyRecipes(ctx, uid, kind)` returns `([]domain.RecipeBrief, error)` and rejects `kind=DISLIKE` with `domain.ErrInvalidMyRecipesKind` | §6.4 |
| SPEC-REACTIONS-021 | `domain.ErrInvalidMyRecipesKind` sentinel error is exported and matched via `errors.Is` | §6.4 |
| SPEC-REACTIONS-022 | `POST /api/reactions` is registered with the Go 1.22+ pattern router as `mux.HandleFunc("POST /api/reactions", h)` behind the auth middleware | §6.5 |
| SPEC-REACTIONS-023 | `POST /api/reactions` decodes a JSON body `{ "recipe_slug": string, "kind": string }`; missing or malformed body yields `400 INVALID_PAYLOAD` | §6.5 |
| SPEC-REACTIONS-024 | `POST /api/reactions` returns `404 NOT_FOUND` when `recipe_slug` is unknown | §6.5 |
| SPEC-REACTIONS-025 | `POST /api/reactions` returns `422 INVALID_REACTION_KIND` when `kind` is not in the enum | §6.5 |
| SPEC-REACTIONS-026 | `POST /api/reactions` returns `200` with body `{ "recipe_slug": "...", "kind": "...", "updated_at": "<RFC3339>" }` on success | §6.5 |
| SPEC-REACTIONS-027 | `POST /api/reactions` is idempotent: re-posting the same `(uid, slug, kind)` returns `200` and overwrites `created_at` (per §6.2.2 SQL) | §6.5 |
| SPEC-REACTIONS-028 | `POST /api/reactions` reads `firebase_uid` from the auth context (`auth.UIDFromContext`) and never trusts a body-supplied uid | §6.5 |
| SPEC-REACTIONS-029 | `DELETE /api/reactions/{slug}` is registered with `mux.HandleFunc("DELETE /api/reactions/{slug}", h)` behind the auth middleware | §6.6 |
| SPEC-REACTIONS-030 | `DELETE /api/reactions/{slug}` returns `204 No Content` when the slug exists, regardless of whether a reaction row was present | §6.6 |
| SPEC-REACTIONS-031 | `DELETE /api/reactions/{slug}` returns `404 NOT_FOUND` when the slug itself is unknown | §6.6 |
| SPEC-REACTIONS-032 | `DELETE /api/reactions/{slug}` writes no response body on success | §6.6 |
| SPEC-REACTIONS-033 | `GET /api/me/recipes` is registered with `mux.HandleFunc("GET /api/me/recipes", h)` behind the auth middleware | §6.7 |
| SPEC-REACTIONS-034 | `GET /api/me/recipes` defaults `kind` to `LIKE` when the query parameter is omitted or empty | §6.7 |
| SPEC-REACTIONS-035 | `GET /api/me/recipes` accepts only `kind=LIKE` or `kind=TRY_LATER`; `kind=DISLIKE` and any other value return `400 INVALID_PAYLOAD` | §6.7 |
| SPEC-REACTIONS-036 | `GET /api/me/recipes` response shape is `{ "recipes": [<RecipeBrief>, …] }` with `recipes` always a JSON array (never `null`) | §6.7 |
| SPEC-REACTIONS-037 | `GET /api/me/recipes` results are ordered by `user_reactions.created_at DESC` (delegated to `Repo.ListByKind`) | §6.7 |
| SPEC-REACTIONS-038 | All three handlers emit the standard error envelope `{"error":{"code":"...","message":"..."}}` for non-2xx responses | §6.8 |
| SPEC-REACTIONS-039 | All three handlers log lifecycle events with `slog` at `INFO` (success) or `WARN`/`ERROR` (failure) and never log Firebase tokens or full request bodies | §6.8 |
| SPEC-REACTIONS-040 | All three handlers thread `r.Context()` to the service layer; no `context.Background()` inside the request lifecycle | §6.8 |
| SPEC-REACTIONS-041 | Wiring in `cmd/cooksense-server/main.go` constructs `reactions.NewRepo(pool)`, `reactions.NewService(repo, recipesRepo)`, and registers the three routes | §6.9 |
| SPEC-REACTIONS-042 | No package-level mutable state holds the pool, service, or handler — all dependencies are injected via constructors | §6.9 |
| SPEC-REACTIONS-043 | `internal/reactions/doc.go` package comment describes Repo + Service responsibilities and cites SPEC-REACTIONS-007..021 | §10.1 |
| SPEC-REACTIONS-044 | API doc `docs/architecture/api.md` is updated only if behavior diverges from the current text; if amended, it ships in the same PR | §10.2 |
| SPEC-REACTIONS-045 | Test `Test_Reactions_Upsert_FromNoRow_InsertsRow` exists (integration) | §9.1 |
| SPEC-REACTIONS-046 | Test `Test_Reactions_Upsert_OverwritesPreviousKind` exists (integration) | §9.1 |
| SPEC-REACTIONS-047 | Test `Test_Reactions_Delete_RemovesExistingRow` exists (integration) | §9.1 |
| SPEC-REACTIONS-048 | Test `Test_Reactions_Delete_Idempotent_NoRow_Returns204` exists (integration) | §9.1 |
| SPEC-REACTIONS-049 | Test `Test_Reactions_Post_UnknownSlug_Returns404` exists (integration) | §9.1 |
| SPEC-REACTIONS-050 | Test `Test_Reactions_Post_InvalidKind_Returns422` exists (integration) | §9.1 |
| SPEC-REACTIONS-051 | Test `Test_Reactions_Post_MalformedJSON_Returns400` exists (integration) | §9.1 |
| SPEC-REACTIONS-052 | Test `Test_MyRecipes_DefaultsToLike_WhenKindOmitted` exists (integration) | §9.1 |
| SPEC-REACTIONS-053 | Test `Test_MyRecipes_OrderedByCreatedAtDesc` exists (integration) | §9.1 |
| SPEC-REACTIONS-054 | Test `Test_MyRecipes_DislikeKind_Returns400` exists (integration) | §9.1 |
| SPEC-REACTIONS-055 | Test `Test_MyRecipes_UnknownKind_Returns400` exists (integration) | §9.1 |
| SPEC-REACTIONS-056 | Test `Test_MyRecipes_EmptyResult_ReturnsEmptyArray` exists (integration) | §9.1 |
| SPEC-REACTIONS-057 | Test `Test_ParseReactionKind_RejectsLowercase` exists (unit) | §9.1 |
