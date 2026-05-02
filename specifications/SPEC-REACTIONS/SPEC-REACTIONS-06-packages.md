# SPEC-REACTIONS — §6 Package Specifications

[← Index](SPEC-REACTIONS-00-index.md)

## 6.1 `internal/domain/reaction.go`

### SPEC-REACTIONS-001 — `ReactionKind` enum type

```go
// ReactionKind is the canonical enum of user reactions on a recipe.
//
// Implements: SPEC-REACTIONS-001, SPEC-REACTIONS-003.
type ReactionKind string

const (
    ReactionLike     ReactionKind = "LIKE"
    ReactionDislike  ReactionKind = "DISLIKE"
    ReactionTryLater ReactionKind = "TRY_LATER"
)
```

The values match the Postgres `reaction_kind` enum verbatim (D-0008).

### SPEC-REACTIONS-002 — `ParseReactionKind`

```go
// ParseReactionKind returns the canonical ReactionKind for s, or
// ErrInvalidReactionKind if s is not one of LIKE/DISLIKE/TRY_LATER.
//
// Comparison is case-sensitive: "like" SHALL be rejected.
//
// Implements: SPEC-REACTIONS-002.
func ParseReactionKind(s string) (ReactionKind, error)
```

The function SHALL return `ErrInvalidReactionKind` (a sentinel exported by
the same file) and the zero value of `ReactionKind`.

### SPEC-REACTIONS-003 — `Valid`

```go
// Valid reports whether r is one of the canonical ReactionKind values.
//
// Implements: SPEC-REACTIONS-003.
func (r ReactionKind) Valid() bool
```

### SPEC-REACTIONS-004 — `Reaction` struct

```go
// Reaction is the in-memory projection of a user_reactions row.
//
// Implements: SPEC-REACTIONS-004.
type Reaction struct {
    FirebaseUID string
    RecipeID    int64
    Kind        ReactionKind
    CreatedAt   time.Time
}
```

### SPEC-REACTIONS-005 — Domain purity

The `internal/domain` package SHALL import only stdlib packages
(`time`, `errors`, `fmt`). It SHALL NOT import `pgx`, `database/sql`,
`net/http`, `os`, or `io`. Reviewers SHALL verify by inspecting the file's
import block.

### SPEC-REACTIONS-006 — Repository interface

```go
// ReactionRepository is the persistence contract for user reactions.
//
// Implementations live in internal/reactions.
//
// Implements: SPEC-REACTIONS-006.
type ReactionRepository interface {
    Upsert(ctx context.Context, uid string, recipeID int64, kind ReactionKind) (time.Time, error)
    Delete(ctx context.Context, uid string, recipeID int64) error
    ListByKind(ctx context.Context, uid string, kind ReactionKind) ([]RecipeBrief, error)
}
```

Sentinel errors exported alongside `ReactionKind`:

```go
var (
    // ErrInvalidReactionKind is returned by ParseReactionKind when s is not
    // one of the canonical enum values.
    ErrInvalidReactionKind = errors.New("invalid reaction kind")

    // ErrInvalidMyRecipesKind is returned by reactions.Service.MyRecipes
    // when the caller asks for kind=DISLIKE.
    ErrInvalidMyRecipesKind = errors.New("kind not exposed by /api/me/recipes")
)
```

`ErrRecipeNotFound` is owned by SPEC-RECIPES; this spec only consumes it.

## 6.2 `internal/reactions/repo.go`

### 6.2.1 Constructor — SPEC-REACTIONS-007

```go
// Repo is the pgx-backed implementation of domain.ReactionRepository.
//
// Implements: SPEC-REACTIONS-007.
type Repo struct {
    pool *pgxpool.Pool
}

// NewRepo returns a Repo bound to the given pool. The pool MUST be ready
// (pinged) by the caller; Repo does not validate connectivity.
//
// Implements: SPEC-REACTIONS-007.
func NewRepo(pool *pgxpool.Pool) *Repo
```

A compile-time assertion SHALL appear in the same file:

```go
var _ domain.ReactionRepository = (*Repo)(nil)
```

### 6.2.2 `Repo.Upsert` — SPEC-REACTIONS-008, 009, 014

```go
// Upsert inserts a new reaction or overwrites the existing one for the
// (uid, recipeID) pair. It returns the resulting created_at, which is
// always now() on the server because re-upserts reset the timestamp
// (D-0008).
//
// kind MUST be one of domain.ReactionLike, ReactionDislike, ReactionTryLater.
// Callers SHALL validate before invocation; Repo passes kind through to
// Postgres unchanged.
//
// Implements: SPEC-REACTIONS-008, SPEC-REACTIONS-009, SPEC-REACTIONS-014.
func (r *Repo) Upsert(ctx context.Context, uid string, recipeID int64, kind domain.ReactionKind) (time.Time, error)
```

SQL (verbatim from the story technical notes):

```sql
INSERT INTO user_reactions (firebase_uid, recipe_id, kind)
VALUES ($1, $2, $3)
ON CONFLICT (firebase_uid, recipe_id) DO UPDATE
    SET kind = EXCLUDED.kind, created_at = now()
RETURNING created_at;
```

`kind` is bound as `string(kind)`; pgx serializes it to the Postgres
`reaction_kind` ENUM via the standard text encoder.

### 6.2.3 `Repo.Delete` — SPEC-REACTIONS-010

```go
// Delete removes the reaction row for (uid, recipeID), if any.
// Deleting zero rows is NOT an error; the operation is idempotent.
//
// Implements: SPEC-REACTIONS-010.
func (r *Repo) Delete(ctx context.Context, uid string, recipeID int64) error
```

SQL:

```sql
DELETE FROM user_reactions
WHERE firebase_uid = $1 AND recipe_id = $2;
```

The implementation SHALL NOT inspect the `RowsAffected()` count for the
purpose of returning an error. It MAY log the count at `DEBUG`.

### 6.2.4 `Repo.ListByKind` — SPEC-REACTIONS-011, 012, 013

```go
// ListByKind returns the caller's recipes that have a reaction of the
// given kind, ordered by reaction created_at DESC.
//
// On no rows, the returned slice is empty (length 0) and not nil.
//
// Implements: SPEC-REACTIONS-011, SPEC-REACTIONS-012, SPEC-REACTIONS-013.
func (r *Repo) ListByKind(ctx context.Context, uid string, kind domain.ReactionKind) ([]domain.RecipeBrief, error)
```

SQL (column list MUST match `domain.RecipeBrief`):

```sql
SELECT
    r.slug,
    r.title,
    r.concept,
    r.time_minutes,
    r.passive_prep_minutes,
    r.cooking_methods,
    r.tags,
    r.flavor_profile
FROM user_reactions ur
JOIN recipes r ON r.id = ur.recipe_id
WHERE ur.firebase_uid = $1
  AND ur.kind         = $2
ORDER BY ur.created_at DESC;
```

The arrays (`cooking_methods`, `tags`, `flavor_profile`) are scanned with
`pgx`'s native `[]string` mapping for `text[]`. The slice returned to the
caller MUST be initialized empty (`make([]domain.RecipeBrief, 0)`) so JSON
serialization yields `[]`, not `null`.

## 6.3 `internal/recipes/repo.go` — `ResolveSlug`

### SPEC-REACTIONS-015, 016

This spec relies on a thin lookup the recipes package provides:

```go
// ResolveSlug returns the numeric primary key for the given slug, or
// domain.ErrRecipeNotFound if no recipe matches.
//
// Implements (across SPEC-RECIPES + SPEC-REACTIONS): SPEC-REACTIONS-015.
func (r *Repo) ResolveSlug(ctx context.Context, slug string) (int64, error)
```

```go
// ErrRecipeNotFound is the sentinel for "no recipe with this slug".
//
// Implements: SPEC-REACTIONS-016.
var ErrRecipeNotFound = errors.New("recipe not found")
```

If these symbols already exist (Story 07), SPEC-REACTIONS does not
duplicate them. If not, Task T-04 of Appendix B SHALL add them.

## 6.4 `internal/reactions/service.go`

### SPEC-REACTIONS-017 — `Service` struct

```go
// recipeResolver is the narrow interface Service uses to resolve a slug.
// Defined here (the consumer) per the repo's "interfaces at the consumer"
// rule (Go idioms, CLAUDE.md).
type recipeResolver interface {
    ResolveSlug(ctx context.Context, slug string) (int64, error)
}

// Service orchestrates slug resolution and reaction persistence.
// Service is stateless and goroutine-safe: a single instance SHALL be
// shared across all HTTP requests.
//
// Implements: SPEC-REACTIONS-017.
type Service struct {
    reactions domain.ReactionRepository
    recipes   recipeResolver
}

func NewService(reactions domain.ReactionRepository, recipes recipeResolver) *Service
```

### SPEC-REACTIONS-018 — `SetReaction`

```go
// SetReaction upserts the given reaction and returns the resulting
// created_at. It returns domain.ErrRecipeNotFound when slug is unknown.
//
// Implements: SPEC-REACTIONS-018.
func (s *Service) SetReaction(ctx context.Context, uid, slug string, kind domain.ReactionKind) (time.Time, error)
```

Behavior:

1. `id, err := s.recipes.ResolveSlug(ctx, slug)` — surface the error.
2. `return s.reactions.Upsert(ctx, uid, id, kind)`.

### SPEC-REACTIONS-019 — `RemoveReaction`

```go
// RemoveReaction deletes the reaction (if any) for the given slug.
// Returns domain.ErrRecipeNotFound when slug is unknown.
// Returns nil when the slug exists, regardless of whether a row was
// actually removed.
//
// Implements: SPEC-REACTIONS-019.
func (s *Service) RemoveReaction(ctx context.Context, uid, slug string) error
```

### SPEC-REACTIONS-020, 021 — `MyRecipes`

```go
// MyRecipes returns the caller's recipes for the given kind, sorted by
// reaction created_at DESC.
//
// kind MUST be ReactionLike or ReactionTryLater. ReactionDislike yields
// domain.ErrInvalidMyRecipesKind. Any non-canonical value is the caller's
// programming error and yields ErrInvalidReactionKind.
//
// Implements: SPEC-REACTIONS-020, SPEC-REACTIONS-021.
func (s *Service) MyRecipes(ctx context.Context, uid string, kind domain.ReactionKind) ([]domain.RecipeBrief, error)
```

Validation:

```go
switch kind {
case domain.ReactionLike, domain.ReactionTryLater:
    return s.reactions.ListByKind(ctx, uid, kind)
case domain.ReactionDislike:
    return nil, domain.ErrInvalidMyRecipesKind
default:
    return nil, domain.ErrInvalidReactionKind
}
```

## 6.5 `internal/api/reactions_post_handler.go`

### SPEC-REACTIONS-022..028

```go
// PostReaction is the http.HandlerFunc for POST /api/reactions.
//
// Implements: SPEC-REACTIONS-022..028.
func PostReaction(svc *reactions.Service, log *slog.Logger) http.HandlerFunc
```

Request body:

```go
type postReactionRequest struct {
    RecipeSlug string `json:"recipe_slug"`
    Kind       string `json:"kind"`
}
```

Response body on success (`200 OK`):

```go
type postReactionResponse struct {
    RecipeSlug string    `json:"recipe_slug"`
    Kind       string    `json:"kind"`
    UpdatedAt  time.Time `json:"updated_at"` // RFC3339 via json.Marshal
}
```

Algorithm:

1. `uid, ok := auth.UIDFromContext(ctx)` — `ok=false` → `500 INTERNAL`.
2. Decode the JSON body with `json.NewDecoder(r.Body).Decode(&req)`. On
   error → `400 INVALID_PAYLOAD`.
3. `req.RecipeSlug == ""` → `400 INVALID_PAYLOAD`.
4. `kind, err := domain.ParseReactionKind(req.Kind)` — error →
   `422 INVALID_REACTION_KIND`.
5. `at, err := svc.SetReaction(ctx, uid, req.RecipeSlug, kind)`. Map error
   per §5.7.
6. Write `200` + JSON `postReactionResponse{req.RecipeSlug, string(kind), at}`.

The handler SHALL register itself in `cmd/cooksense-server/main.go` as:

```go
mux.Handle("POST /api/reactions", auth.Require(api.PostReaction(svc, log)))
```

## 6.6 `internal/api/reactions_delete_handler.go`

### SPEC-REACTIONS-029..032

```go
// DeleteReaction is the http.HandlerFunc for DELETE /api/reactions/{slug}.
//
// Implements: SPEC-REACTIONS-029..032.
func DeleteReaction(svc *reactions.Service, log *slog.Logger) http.HandlerFunc
```

Algorithm:

1. `uid, ok := auth.UIDFromContext(ctx)` — `ok=false` → `500`.
2. `slug := r.PathValue("slug")` — empty → `400 INVALID_PAYLOAD`
   (defense-in-depth; the router pattern guarantees presence).
3. `err := svc.RemoveReaction(ctx, uid, slug)`. Map error per §5.7.
4. On success → `w.WriteHeader(http.StatusNoContent)` and return without
   writing a body.

Route registration:

```go
mux.Handle("DELETE /api/reactions/{slug}", auth.Require(api.DeleteReaction(svc, log)))
```

## 6.7 `internal/api/me_recipes_handler.go`

### SPEC-REACTIONS-033..037

```go
// MyRecipes is the http.HandlerFunc for GET /api/me/recipes.
//
// Implements: SPEC-REACTIONS-033..037.
func MyRecipes(svc *reactions.Service, log *slog.Logger) http.HandlerFunc
```

Response body:

```go
type myRecipesResponse struct {
    Recipes []domain.RecipeBrief `json:"recipes"`
}
```

Algorithm:

1. `uid, ok := auth.UIDFromContext(ctx)` — `ok=false` → `500`.
2. `q := r.URL.Query().Get("kind")`.
3. If `q == ""` → `kind := domain.ReactionLike`. Else
   `kind, err := domain.ParseReactionKind(q)` — error →
   `400 INVALID_PAYLOAD`.
4. `recipes, err := svc.MyRecipes(ctx, uid, kind)`. Map error per §5.7.
5. If `recipes == nil` → `recipes = make([]domain.RecipeBrief, 0)`.
6. Write `200` + JSON `myRecipesResponse{recipes}`.

Route registration:

```go
mux.Handle("GET /api/me/recipes", auth.Require(api.MyRecipes(svc, log)))
```

## 6.8 Cross-handler concerns

### SPEC-REACTIONS-038 — Error envelope

Every non-2xx response SHALL be written via the shared helper:

```go
api.WriteError(w, http.StatusXxx, "ERROR_CODE", "human readable")
```

producing the body:

```json
{ "error": { "code": "ERROR_CODE", "message": "human readable" } }
```

### SPEC-REACTIONS-039 — Logging

Each handler SHALL emit:

- `INFO` on success with attrs `op`, `uid_hash`, `slug` (no full `uid`,
  no token, no body).
- `WARN` on 4xx with attrs `op`, `code`, `slug`.
- `ERROR` on 5xx with attrs `op`, `err` (the `error` is wrapped, not the
  full request).

Where `uid_hash` is `fmt.Sprintf("%x", sha256.Sum256([]byte(uid)))[:12]` —
short, non-reversible, useful for log correlation. (This helper MAY live
in `internal/auth` if not already present; that addition is owned by
SPEC-AUTH.)

### SPEC-REACTIONS-040 — Context plumbing

Every call to the service layer SHALL pass `r.Context()`. Handlers MUST
NOT use `context.Background()` or `context.TODO()`.

## 6.9 Wiring — `cmd/cooksense-server/main.go`

### SPEC-REACTIONS-041, 042

The server's `run` function (or equivalent bootstrap) SHALL contain:

```go
recipesRepo := recipes.NewRepo(pool)
reactionsRepo := reactions.NewRepo(pool)
reactionsSvc := reactions.NewService(reactionsRepo, recipesRepo)

mux.Handle("POST /api/reactions",
    auth.Require(api.PostReaction(reactionsSvc, log)))
mux.Handle("DELETE /api/reactions/{slug}",
    auth.Require(api.DeleteReaction(reactionsSvc, log)))
mux.Handle("GET /api/me/recipes",
    auth.Require(api.MyRecipes(reactionsSvc, log)))
```

No package-level mutable state SHALL hold the pool, repository, service,
or handler. All instances are owned by the local stack of `run()`.
