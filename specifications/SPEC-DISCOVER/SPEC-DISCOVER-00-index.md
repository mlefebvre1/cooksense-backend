# SPEC-DISCOVER — Recipe Discovery & Detail HTTP Endpoints

**Story:** 07 — Recipe discovery and detail endpoints
**Status:** Draft → Final
**Date:** 2026-05-02
**Authors:** sdd-spec-author

---

## Purpose

This specification governs the two authenticated read-only HTTP endpoints
that power the swipe-feed and the recipe detail screen of the CookSense
mobile app:

- `GET /api/recipes/discover?limit=N` — randomized stream of recipes the
  caller has not yet reacted to.
- `GET /api/recipes/{slug}` — full recipe detail including ingredients
  (sorted by category then name) and steps (in stored order).

The spec covers the three layers introduced by the story
(`internal/recipes/{repo.go, service.go, handler.go}`), the DTOs they
produce, the router wiring in `cmd/cooksense-server/main.go`, and the
integration tests that exercise the endpoints against the compose Postgres
seeded by SPEC-RECIPES.

Out of scope: reactions endpoints (story 08, SPEC-REACTIONS), search by
ingredients (story 09), lessons (story 10), pagination (post-MVP),
"time-to-decision" telemetry (post-MVP), and any caching layer.

---

## File map

| File | Contents |
|------|----------|
| [SPEC-DISCOVER-01-preamble](SPEC-DISCOVER-01-preamble.md) | AI constraints, authorship, traceability rules |
| [SPEC-DISCOVER-02-introduction](SPEC-DISCOVER-02-introduction.md) | Story relationship, decision references |
| [SPEC-DISCOVER-03-goals](SPEC-DISCOVER-03-goals.md) | Goals, non-goals, constraints |
| [SPEC-DISCOVER-04-context](SPEC-DISCOVER-04-context.md) | Dependencies, package boundary map |
| [SPEC-DISCOVER-05-architecture](SPEC-DISCOVER-05-architecture.md) | Request/response flow, layering, SQL shape |
| [SPEC-DISCOVER-06-packages](SPEC-DISCOVER-06-packages.md) | All SPEC-DISCOVER-NNN normative requirements |
| [SPEC-DISCOVER-07-configuration](SPEC-DISCOVER-07-configuration.md) | Limits, defaults, environment variables |
| [SPEC-DISCOVER-08-build](SPEC-DISCOVER-08-build.md) | Build, lint, vet, dependency rules |
| [SPEC-DISCOVER-09-testing](SPEC-DISCOVER-09-testing.md) | Test strategy, named tests, fixtures |
| [SPEC-DISCOVER-10-documentation](SPEC-DISCOVER-10-documentation.md) | Doc comments, README impact |
| [SPEC-DISCOVER-A-checklist](SPEC-DISCOVER-A-checklist.md) | Spec completeness + AC traceability matrix |
| [SPEC-DISCOVER-B-tasks](SPEC-DISCOVER-B-tasks.md) | Ordered implementation task list |

---

## SPEC-ID registry

| ID | Summary | Section |
|----|---------|---------|
| SPEC-DISCOVER-001 | `internal/recipes` is split into `dto.go`, `repo.go`, `service.go`, `handler.go` (one concern per file) | §6.1 |
| SPEC-DISCOVER-002 | `RecipeBrief` DTO struct with JSON `snake_case` tags matching `docs/architecture/api.md` | §6.1 |
| SPEC-DISCOVER-003 | `RecipeFull` DTO embeds `RecipeBrief` and adds `Ingredients []IngredientView` and `Steps []string` | §6.1 |
| SPEC-DISCOVER-004 | `IngredientView` DTO struct with JSON tags `name`, `category`, `quantity`, `unit`, `optional` | §6.1 |
| SPEC-DISCOVER-005 | DTO marshaling SHALL emit `cooking_methods`, `tags`, `flavor_profile` as JSON arrays (never `null`) | §6.1 |
| SPEC-DISCOVER-006 | `Repo` interface exposes `Discover(ctx, uid string, limit int) ([]RecipeBrief, error)` and `GetBySlug(ctx, slug string) (RecipeFull, error)` | §6.2 |
| SPEC-DISCOVER-007 | `ErrNotFound` is the sentinel error returned by `Repo.GetBySlug` when no row matches | §6.2 |
| SPEC-DISCOVER-008 | `PgRepo` is the only `Repo` implementation; constructed via `NewPgRepo(pool *pgxpool.Pool)` | §6.2 |
| SPEC-DISCOVER-009 | Discover SQL selects from `recipes` excluding any row with a row in `user_reactions` for the caller (`NOT EXISTS`) | §6.2 |
| SPEC-DISCOVER-010 | Discover SQL randomizes results with `ORDER BY random()` and applies `LIMIT $2` | §6.2 |
| SPEC-DISCOVER-011 | Discover SQL parameter order: `$1 = firebase_uid`, `$2 = limit` | §6.2 |
| SPEC-DISCOVER-012 | Discover excludes recipes for **any** reaction kind (LIKE, DISLIKE, TRY_LATER) — story-07 simplification, supersedes the `DISLIKE-only` wording in `docs/architecture/api.md` for MVP | §6.2 |
| SPEC-DISCOVER-013 | `GetBySlug` SQL reads the recipe by slug, then loads `recipe_ingredients` joined to `ingredients` ordered by `ingredients.category ASC, ingredients.name ASC` | §6.2 |
| SPEC-DISCOVER-014 | `GetBySlug` returns `ErrNotFound` (wrapped via `fmt.Errorf("%w", ErrNotFound)`) when `pgx.ErrNoRows` is encountered for the recipe row | §6.2 |
| SPEC-DISCOVER-015 | `GetBySlug` returns the recipe's `steps` JSONB column unmarshalled into `[]string` preserving stored order | §6.2 |
| SPEC-DISCOVER-016 | The repository SHALL be the only layer that imports `pgx`/`pgxpool` for this feature | §6.2 |
| SPEC-DISCOVER-017 | `Service` struct holds a `Repo` and is constructed via `NewService(repo Repo) *Service` | §6.3 |
| SPEC-DISCOVER-018 | `Service.Discover(ctx, uid, limit) ([]RecipeBrief, error)` clamps `limit` to the `[1, MaxDiscoverLimit]` range; `limit <= 0` → `DefaultDiscoverLimit` | §6.3 |
| SPEC-DISCOVER-019 | `DefaultDiscoverLimit = 10` and `MaxDiscoverLimit = 25` are exported package-level constants | §6.3 |
| SPEC-DISCOVER-020 | `Service.GetBySlug(ctx, slug) (RecipeFull, error)` delegates to `Repo.GetBySlug` and propagates `ErrNotFound` unchanged (`errors.Is`-compatible) | §6.3 |
| SPEC-DISCOVER-021 | The service SHALL NOT import `net/http`, `pgx`, `pgxpool`, or any I/O package other than `context` | §6.3 |
| SPEC-DISCOVER-022 | The service SHALL log a single structured `slog.Info` line per request with attributes `route`, `uid` (no PII), `limit`, and `count` | §6.3 |
| SPEC-DISCOVER-023 | `Handler` struct holds a `*Service` and is constructed via `NewHandler(svc *Service) *Handler` | §6.4 |
| SPEC-DISCOVER-024 | `Handler.RegisterRoutes(mux *http.ServeMux, mw func(http.Handler) http.Handler)` registers `GET /api/recipes/discover` and `GET /api/recipes/{slug}` wrapped in the auth middleware | §6.4 |
| SPEC-DISCOVER-025 | The discover handler parses `?limit=` as an `int` via `strconv.Atoi`; missing/empty/invalid SHALL fall back to the default (no `400`) | §6.4 |
| SPEC-DISCOVER-026 | The discover handler retrieves the caller via `auth.UserFromContext`; absence yields `500 INTERNAL` (middleware contract violation) | §6.4 |
| SPEC-DISCOVER-027 | The discover handler responds `200` with body `{"recipes": [<RecipeBrief>, …]}`; the `recipes` field SHALL be `[]` (not `null`) when empty | §6.4 |
| SPEC-DISCOVER-028 | The detail handler extracts `slug` via `r.PathValue("slug")` (Go 1.22 mux pattern parameter) | §6.4 |
| SPEC-DISCOVER-029 | The detail handler responds `200` with the `RecipeFull` JSON object on success | §6.4 |
| SPEC-DISCOVER-030 | On `ErrNotFound`, the detail handler responds `404` with the standard error envelope `{"error":{"code":"NOT_FOUND","message":"recipe not found"}}` | §6.4 |
| SPEC-DISCOVER-031 | On any other error, handlers respond `500` with `{"error":{"code":"INTERNAL","message":"internal error"}}` and log the underlying error at `ERROR` level with the request id | §6.4 |
| SPEC-DISCOVER-032 | Handlers set `Content-Type: application/json; charset=utf-8` on every response (success and error) | §6.4 |
| SPEC-DISCOVER-033 | Handlers SHALL NOT touch the database directly; they SHALL only call the service | §6.4 |
| SPEC-DISCOVER-034 | `cmd/cooksense-server/main.go` wires `recipes.NewPgRepo` → `recipes.NewService` → `recipes.NewHandler` and calls `RegisterRoutes` once | §6.5 |
| SPEC-DISCOVER-035 | `internal/recipes/doc.go` package comment cites SPEC-DISCOVER-001 through SPEC-DISCOVER-034 | §10.1 |
| SPEC-DISCOVER-036 | Integration test `Test_Discover_RespectsLimit_ReturnsExpectedCount` exists | §9.1 |
| SPEC-DISCOVER-037 | Integration test `Test_Discover_ExcludesAlreadyReacted_NeverReturnsThem` exists | §9.1 |
| SPEC-DISCOVER-038 | Integration test `Test_GetRecipe_KnownSlug_ReturnsFullRecipe` exists | §9.1 |
| SPEC-DISCOVER-039 | Integration test `Test_GetRecipe_UnknownSlug_Returns404` exists | §9.1 |
| SPEC-DISCOVER-040 | Unit test `Test_ServiceDiscover_LimitClamping_AppliesBounds` exists | §9.1 |
| SPEC-DISCOVER-041 | Unit test `Test_HandlerDiscover_NegativeLimit_FallsBackToDefault` exists | §9.1 |
| SPEC-DISCOVER-042 | README updated with a short "Discover and detail endpoints" subsection | §10.2 |
