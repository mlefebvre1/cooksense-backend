# SPEC-DISCOVER — §4 System Context

[← Index](SPEC-DISCOVER-00-index.md)

## 4.1 External dependencies (new in this story)

None. Story 07 introduces no new third-party imports.

## 4.2 Reused dependencies

| Dependency | Source story | Used here for |
|------------|--------------|----------------|
| `net/http` (stdlib) | SPEC-BOOT, D-0002 | `*http.ServeMux`, `r.PathValue`, handlers |
| `encoding/json` (stdlib) | — | DTO marshalling/unmarshalling |
| `pgx/v5` + `pgxpool` | SPEC-DB | Connection pool injected into `recipes.NewPgRepo` |
| `internal/auth` | SPEC-AUTH | `auth.UserFromContext` to obtain the caller |
| `internal/db.Open` | SPEC-DB | Wired in `cmd/cooksense-server/main.go` |
| `internal/config` | SPEC-DB | `DATABASE_URL` consumption (transitive) |
| `internal/domain` | SPEC-RECIPES | Taxonomy constants reused in DTO docs (no runtime import required) |
| `log/slog` (stdlib) | repo guideline | Structured logging in service and handler error paths |

## 4.3 Package boundary map

```
cmd/cooksense-server      ── wires NewPgRepo → NewService → NewHandler
        │                    and registers the /api/recipes/* routes.
        ▼
internal/recipes/handler  ── parses request, calls Service, encodes JSON.
        │                    Imports: net/http, encoding/json, log/slog,
        │                             internal/auth, internal/recipes/service.
        ▼
internal/recipes/service  ── clamps limit, delegates to Repo, logs request.
        │                    Imports: context, log/slog, internal/recipes/repo (interface).
        ▼
internal/recipes/repo     ── Postgres queries (Discover + GetBySlug).
                             Imports: context, errors, fmt,
                                      github.com/jackc/pgx/v5,
                                      github.com/jackc/pgx/v5/pgxpool,
                                      encoding/json (for steps JSONB).
```

Reverse imports are forbidden. The handler SHALL NOT import the repo.
The service SHALL NOT import `net/http`. The repo SHALL NOT import the
handler or the service.

## 4.4 Filesystem layout introduced

```
internal/
  recipes/
    doc.go        (updated package comment — SPEC-DISCOVER-035)
    dto.go        (new — RecipeBrief, RecipeFull, IngredientView)
    repo.go       (new — Repo interface, PgRepo, ErrNotFound, SQL)
    service.go    (new — Service, NewService, limit clamping)
    handler.go    (new — Handler, NewHandler, RegisterRoutes)
    repo_test.go        (new — integration tests against compose Postgres)
    service_test.go     (new — unit tests, in-memory fake Repo)
    handler_test.go     (new — integration tests via httptest + fakeVerifier)
cmd/
  cooksense-server/
    main.go       (updated — wires recipes feature into the mux)
```

## 4.5 Out-of-band collaborators

- `migrations/0001_init.up.sql` (SPEC-DB) — provides the `recipes`,
  `ingredients`, `recipe_ingredients`, and `user_reactions` tables.
- `seed/recipes/*.yaml` (SPEC-RECIPES + Story 06) — provides the data the
  endpoints serve. Integration tests SHALL run `make seed` (or the
  programmatic equivalent) before exercising endpoints.
- `internal/auth.Middleware` (SPEC-AUTH) — wraps the routes and injects
  `auth.User` into the request context.
