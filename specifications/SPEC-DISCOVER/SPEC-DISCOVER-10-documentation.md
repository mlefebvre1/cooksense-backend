# SPEC-DISCOVER — §10 Documentation

[← Index](SPEC-DISCOVER-00-index.md)

## 10.1 Doc-comment requirements

Every exported symbol introduced by Story 07 SHALL carry a Go Doc comment
that:

1. Begins with the symbol name (`// Handler is …`, `// NewService …`).
2. References the SPEC-DISCOVER-NNN ID(s) it implements.
3. Documents pre-conditions, post-conditions, and error behavior where
   non-obvious (status codes, sentinel errors, side effects).

| File | Required doc-comment subjects | SPEC-IDs to cite |
|------|-------------------------------|------------------|
| `internal/recipes/doc.go` | Package overview: handler / service / repo trio for `/api/recipes/*` | 001, 035 |
| `internal/recipes/dto.go` | `RecipeBrief`, `RecipeFull`, `IngredientView` | 002, 003, 004, 005 |
| `internal/recipes/repo.go` | `Repo`, `ErrNotFound`, `PgRepo`, `NewPgRepo`, `(*PgRepo).Discover`, `(*PgRepo).GetBySlug` | 006–016 |
| `internal/recipes/service.go` | `Service`, `NewService`, `DefaultDiscoverLimit`, `MaxDiscoverLimit`, `(*Service).Discover`, `(*Service).GetBySlug` | 017–022 |
| `internal/recipes/handler.go` | `Handler`, `NewHandler`, `(*Handler).RegisterRoutes`, `(*Handler).Discover`, `(*Handler).GetBySlug` | 023–033 |
| `cmd/cooksense-server/main.go` | The wiring block (inline comment citing SPEC-DISCOVER-034) | 034 |

### SPEC-DISCOVER-035 — Package comment for `internal/recipes`

The `internal/recipes/doc.go` file SHALL contain a package comment of the
form:

```go
// Package recipes implements the recipe discovery and detail HTTP feature
// of CookSense. It is laid out in three layers — handler (HTTP),
// service (use cases + clamping), repo (Postgres) — and exposes the
// constructors required by cmd/cooksense-server to wire the routes.
//
// Implements: SPEC-DISCOVER-001 through SPEC-DISCOVER-034.
package recipes
```

The existing `doc.go` SHALL be updated, not duplicated.

## 10.2 README impact

### SPEC-DISCOVER-042 — README endpoint subsection

The top-level `README.md` SHALL gain (or update) a "Discover and detail
endpoints" subsection summarizing:

```markdown
## Recipe discovery and detail

Two authenticated endpoints power the swipe feed and the recipe detail
screen:

| Method | Path | Notes |
|--------|------|-------|
| `GET` | `/api/recipes/discover?limit=N` | `N` defaults to 10, max 25; excludes recipes the caller has already reacted to |
| `GET` | `/api/recipes/{slug}` | Returns the full recipe with sorted ingredients and stored-order steps |

Both endpoints require `Authorization: Bearer <Firebase ID token>`.
Response shapes are documented in `docs/architecture/api.md`.
```

The full curated quickstart lands in Story 12; SPEC-DISCOVER-042 only
requires the discovery paragraph to be present and accurate.

## 10.3 Architecture docs

`docs/architecture/api.md` already contains the canonical request/response
shapes (`<RecipeBrief>`, `<RecipeFull>`, error envelope). **One edit is
required**: the wording "Excludes `DISLIKE`d recipes by default (they can
never come back)" SHALL be amended to reflect the story's stronger rule
("excludes any recipe with any reaction"). This change SHALL ship in the
same PR that lands SPEC-DISCOVER-009/012, per the repo's "spec and code
ship together" rule.

## 10.4 Decision log

No new ADR is required. D-0002, D-0004, D-0007, and D-0008 already cover
the architectural choices. If implementation diverges from any of them
(e.g., introducing an HTTP framework, persisting reactions in a separate
table), a new ADR SHALL be added before merge.

## 10.5 Inline error catalog

The error codes used by Story-07 handlers SHALL match the catalog in
`docs/architecture/api.md`:

| Code | HTTP | Used by |
|------|------|---------|
| `NOT_FOUND` | 404 | `GET /api/recipes/{slug}` when slug is unknown |
| `INTERNAL` | 500 | Both handlers on unexpected errors |
| `UNAUTHENTICATED` | 401 | Emitted by the auth middleware (not by Story-07 code) |

If Story 07 needs a code not in the existing catalog, the catalog SHALL
be amended in the same PR.
