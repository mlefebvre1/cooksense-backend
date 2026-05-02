# SPEC-DISCOVER — §6 Package Specifications (Normative)

[← Index](SPEC-DISCOVER-00-index.md)

This section is normative. Every requirement uses RFC-2119 keywords.

---

## §6.1 `internal/recipes/dto.go`

### SPEC-DISCOVER-001 — File split
The `internal/recipes` package SHALL contain exactly one of each:
`dto.go`, `repo.go`, `service.go`, `handler.go`. No file SHALL exceed
500 lines (repo guideline). DTOs SHALL live only in `dto.go`.

### SPEC-DISCOVER-002 — `RecipeBrief` DTO
```go
type RecipeBrief struct {
    Slug                string   `json:"slug"`
    Title               string   `json:"title"`
    Concept             string   `json:"concept"`
    TimeMinutes         int      `json:"time_minutes"`
    PassivePrepMinutes  int      `json:"passive_prep_minutes"`
    CookingMethods      []string `json:"cooking_methods"`
    Tags                []string `json:"tags"`
    FlavorProfile       []string `json:"flavor_profile"`
}
```
Field shape and JSON tags SHALL be byte-for-byte identical to the
`<RecipeBrief>` example in `docs/architecture/api.md`.

### SPEC-DISCOVER-003 — `RecipeFull` DTO
```go
type RecipeFull struct {
    RecipeBrief                          // embedded; flattens JSON
    Ingredients []IngredientView `json:"ingredients"`
    Steps       []string         `json:"steps"`
}
```
Embedding `RecipeBrief` SHALL flatten its fields into the outer JSON
object (Go's default behavior with anonymous embedding). The two
additional fields SHALL appear after the brief fields in the wire format.

### SPEC-DISCOVER-004 — `IngredientView` DTO
```go
type IngredientView struct {
    Name     string  `json:"name"`
    Category string  `json:"category"`
    Quantity float64 `json:"quantity"`
    Unit     string  `json:"unit"`
    Optional bool    `json:"optional"`
}
```
Field types SHALL match the columns of `recipe_ingredients` joined to
`ingredients` (per `docs/architecture/data-model.md`). `Quantity` SHALL be
`float64` to accommodate the `NUMERIC(10,2)` column.

### SPEC-DISCOVER-005 — Empty arrays, never `null`
The repository SHALL initialize `CookingMethods`, `Tags`, `FlavorProfile`,
`Ingredients`, and `Steps` to non-nil empty slices (`[]string{}`,
`[]IngredientView{}`) when the database row contains no values, so that
JSON marshalling emits `[]` rather than `null`. This rule applies to all
slice fields on the response DTOs.

---

## §6.2 `internal/recipes/repo.go`

### SPEC-DISCOVER-006 — `Repo` interface
```go
type Repo interface {
    Discover(ctx context.Context, uid string, limit int) ([]RecipeBrief, error)
    GetBySlug(ctx context.Context, slug string) (RecipeFull, error)
}
```
The interface SHALL live in `repo.go` and SHALL be the contract the
service depends on. No other methods SHALL be added in Story 07.

### SPEC-DISCOVER-007 — `ErrNotFound` sentinel
```go
var ErrNotFound = errors.New("recipes: not found")
```
The package SHALL export this sentinel. `Repo.GetBySlug` SHALL return an
error matching `errors.Is(err, ErrNotFound) == true` when the recipe row
is missing.

### SPEC-DISCOVER-008 — `PgRepo` constructor
```go
type PgRepo struct {
    pool *pgxpool.Pool
}

func NewPgRepo(pool *pgxpool.Pool) *PgRepo
```
The pool SHALL be passed by argument; the package SHALL NOT hold any
package-level pool. `PgRepo` SHALL satisfy `Repo`.

### SPEC-DISCOVER-009 — Discover anti-join
The Discover query SHALL be exactly the SQL shown in
SPEC-DISCOVER-05 §5.4. The `NOT EXISTS` clause SHALL filter out every
recipe for which a `user_reactions` row exists for `$1`. No predicate on
`ur.kind` SHALL be present (per SPEC-DISCOVER-012).

### SPEC-DISCOVER-010 — Random ordering and limit
The query SHALL include `ORDER BY random()` and `LIMIT $2` literally. No
deterministic-seed ordering SHALL be introduced (the catalog is small
enough that test stability is achieved by counting and set membership,
not order — see §9).

### SPEC-DISCOVER-011 — Parameter order
The placeholders SHALL be `$1 = uid` (TEXT) and `$2 = limit` (INT). The
implementation SHALL pass them via `pool.Query(ctx, sql, uid, limit)`.

### SPEC-DISCOVER-012 — Any-reaction exclusion (story override)
For Story 07, "already reacted to" SHALL mean **any row in
`user_reactions`** for the calling `firebase_uid`, regardless of `kind`.
This deliberately overrides the `DISLIKE-only by default` wording in
`docs/architecture/api.md` for the MVP. The rationale is recorded in this
spec; future stories MAY refine the predicate without amending the
schema (LIKE/TRY_LATER could become re-discoverable later).

### SPEC-DISCOVER-013 — Detail ingredient ordering
`GetBySlug` SHALL return `Ingredients` ordered by `ingredients.category
ASC, ingredients.name ASC` (per AC #2 of the story). The ordering SHALL
be enforced by SQL (`ORDER BY i.category ASC, i.name ASC`), not by an
in-Go sort.

### SPEC-DISCOVER-014 — Not-found translation
When the first query (recipe row) returns `pgx.ErrNoRows`, `GetBySlug`
SHALL return `fmt.Errorf("recipes.repo: %w", ErrNotFound)`. No other path
SHALL produce `ErrNotFound`. The handler relies on `errors.Is` to
distinguish 404 from 500.

### SPEC-DISCOVER-015 — Steps preservation
The `recipes.steps` JSONB column SHALL be unmarshalled into `[]string`
preserving the stored order. The implementation SHALL use
`json.Unmarshal` on the raw bytes; no element SHALL be elided or
re-ordered.

### SPEC-DISCOVER-016 — Repo-only DB imports
`pgx`, `pgxpool`, and any database-driver package SHALL be imported only
by `internal/recipes/repo.go` and `internal/recipes/repo_test.go`.
`handler.go` and `service.go` SHALL NOT import them.

---

## §6.3 `internal/recipes/service.go`

### SPEC-DISCOVER-017 — `Service` constructor
```go
type Service struct {
    repo Repo
    log  *slog.Logger
}

func NewService(repo Repo) *Service
```
The constructor SHALL accept any value satisfying `Repo` (enabling
in-memory fakes for unit tests). The logger SHALL default to
`slog.Default()` if not injected; an exported variant
`NewServiceWithLogger(repo Repo, log *slog.Logger) *Service` MAY exist.

### SPEC-DISCOVER-018 — Limit clamping
```go
func (s *Service) Discover(ctx context.Context, uid string, limit int) ([]RecipeBrief, error)
```
SHALL clamp `limit` exactly as follows:

| Input `limit` | Effective `limit` |
|---------------|-------------------|
| `≤ 0` (incl. unset / negative) | `DefaultDiscoverLimit` (10) |
| `1..MaxDiscoverLimit` | unchanged |
| `> MaxDiscoverLimit` | `MaxDiscoverLimit` (25) |

The clamp SHALL happen **before** calling the repo. The repo SHALL trust
the value it receives and apply it verbatim to the `LIMIT $2` placeholder.

### SPEC-DISCOVER-019 — Limit constants
```go
const (
    DefaultDiscoverLimit = 10
    MaxDiscoverLimit     = 25
)
```
These SHALL be exported package-level `const` (PascalCase). The handler
and tests SHALL reference them by name; magic numbers SHALL NOT appear in
either.

### SPEC-DISCOVER-020 — `GetBySlug` propagation
```go
func (s *Service) GetBySlug(ctx context.Context, slug string) (RecipeFull, error)
```
SHALL delegate to `s.repo.GetBySlug(ctx, slug)` and return the result
unchanged. The `ErrNotFound` sentinel SHALL be propagated such that
`errors.Is(err, ErrNotFound)` remains true at the handler.

### SPEC-DISCOVER-021 — No I/O imports beyond `context`
`service.go` SHALL NOT import any of: `net/http`, `os`, `io`,
`database/sql`, `github.com/jackc/pgx/...`. Permitted standard-library
imports: `context`, `errors`, `fmt`, `log/slog`.

### SPEC-DISCOVER-022 — Structured request logging
Each successful service call SHALL emit exactly one `slog.Info` line:

| Call | Attributes |
|------|------------|
| `Discover` | `route="discover"`, `uid=<uid>`, `limit=<clamped>`, `count=<len>` |
| `GetBySlug` | `route="detail"`, `slug=<slug>` |

`uid` is **not** PII for the purposes of this log (it is a Firebase UID,
already opaque). Email and display name SHALL NOT be logged. Errors SHALL
be logged at `ERROR` level by the handler (single source of truth — see
SPEC-DISCOVER-031).

---

## §6.4 `internal/recipes/handler.go`

### SPEC-DISCOVER-023 — `Handler` constructor
```go
type Handler struct {
    svc *Service
    log *slog.Logger
}

func NewHandler(svc *Service) *Handler
```
The constructor SHALL accept the service. A logger may be injected via a
`NewHandlerWithLogger` variant; otherwise `slog.Default()` is used.

### SPEC-DISCOVER-024 — `RegisterRoutes`
```go
func (h *Handler) RegisterRoutes(
    mux *http.ServeMux,
    authMW func(http.Handler) http.Handler,
)
```
SHALL register exactly two patterns:

```go
mux.Handle("GET /api/recipes/discover", authMW(http.HandlerFunc(h.Discover)))
mux.Handle("GET /api/recipes/{slug}",   authMW(http.HandlerFunc(h.GetBySlug)))
```

The function SHALL panic on `nil` mux or `nil` middleware (programming
error, fail fast at startup).

### SPEC-DISCOVER-025 — Limit query parsing
The discover handler SHALL parse `r.URL.Query().Get("limit")` via
`strconv.Atoi`. On any parse error or when the parameter is absent, the
handler SHALL pass `0` to the service (which the service interprets as
"use default" per SPEC-DISCOVER-018). Malformed `?limit=` SHALL NOT yield
a `400` response.

### SPEC-DISCOVER-026 — User retrieval
The handler SHALL retrieve the caller via:

```go
user, ok := auth.UserFromContext(r.Context())
if !ok {
    h.writeError(w, r, http.StatusInternalServerError, "INTERNAL", "internal error", errors.New("auth: missing user in context"))
    return
}
```

A missing user is a middleware contract violation and SHALL produce
`500 INTERNAL`. The error SHALL be logged at `ERROR` level.

### SPEC-DISCOVER-027 — Discover response shape
On success, the handler SHALL respond:

```json
{ "recipes": [ <RecipeBrief>, … ] }
```

When the service returns `nil` or an empty slice, the handler SHALL
encode the field as `[]`, not `null`. The wrapper struct used for
encoding SHALL be:

```go
type discoverResponse struct {
    Recipes []RecipeBrief `json:"recipes"`
}
```

The handler SHALL pre-allocate `Recipes = make([]RecipeBrief, 0)` when
the result is `nil`.

### SPEC-DISCOVER-028 — Slug extraction
The detail handler SHALL extract the slug via:

```go
slug := r.PathValue("slug")
```

It SHALL NOT manually parse `r.URL.Path`. Empty slug SHALL be impossible
under the registered pattern, but the handler SHALL still treat it as
`404 NOT_FOUND` defensively.

### SPEC-DISCOVER-029 — Detail success response
On success, the handler SHALL encode the `RecipeFull` value directly as
the response body (no wrapper). The HTTP status SHALL be `200 OK`.

### SPEC-DISCOVER-030 — Detail not-found response
When the service returns an error matching
`errors.Is(err, ErrNotFound)`, the handler SHALL respond:

```
HTTP/1.1 404 Not Found
Content-Type: application/json; charset=utf-8

{"error":{"code":"NOT_FOUND","message":"recipe not found"}}
```

`code` SHALL be exactly `"NOT_FOUND"` (matching the API doc error catalog).
`message` SHALL be a stable, human-readable string suitable for client
display.

### SPEC-DISCOVER-031 — Internal-error response
For any error not matching `ErrNotFound`, the handler SHALL respond:

```
HTTP/1.1 500 Internal Server Error
Content-Type: application/json; charset=utf-8

{"error":{"code":"INTERNAL","message":"internal error"}}
```

The handler SHALL log the underlying error at `ERROR` level with
attributes `route`, `slug` (when present), `err` (the raw error string).
Underlying error details SHALL NOT leak into the JSON envelope.

### SPEC-DISCOVER-032 — Content-Type
Every response (success and error) SHALL set
`Content-Type: application/json; charset=utf-8` **before** writing the
status code or the body, ensuring `http.ResponseWriter.WriteHeader` does
not freeze the headers prematurely.

### SPEC-DISCOVER-033 — No DB in handler
`handler.go` SHALL NOT import `pgx`/`pgxpool`/`database/sql` and SHALL
NOT execute SQL directly. Its only dependency on data-access logic SHALL
be the `*Service` injected at construction.

---

## §6.5 `cmd/cooksense-server/main.go`

### SPEC-DISCOVER-034 — Wiring
The server bootstrap SHALL include the following block (or its
documented equivalent) **after** opening the pgx pool and constructing
the auth middleware:

```go
recipesRepo := recipes.NewPgRepo(pool)
recipesSvc  := recipes.NewService(recipesRepo)
recipesH    := recipes.NewHandler(recipesSvc)
recipesH.RegisterRoutes(mux, authMW)
```

The wiring SHALL appear exactly once. No package-level singleton SHALL be
introduced. The compose order SHALL be `Repo → Service → Handler →
RegisterRoutes`.

---

## §6.6 Concurrency & lifecycle

- All exported handler/service/repo methods accept `context.Context` as
  the first parameter (per repo guidelines).
- The repo SHALL release rows via `defer rows.Close()` after iteration.
- The handler SHALL respect client cancellation: `r.Context()` SHALL
  flow into the service and into pgx, so a closed connection aborts the
  query.
- No goroutines SHALL be spawned by Story 07 code paths.
