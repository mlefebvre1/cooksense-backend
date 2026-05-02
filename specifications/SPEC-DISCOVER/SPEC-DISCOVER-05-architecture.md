# SPEC-DISCOVER — §5 Architecture

[← Index](SPEC-DISCOVER-00-index.md)

## 5.1 End-to-end request flow

### Discover

```
GET /api/recipes/discover?limit=10
        │  Authorization: Bearer <ID_TOKEN>
        ▼
auth.Middleware ──► verify token ──► UserFromContext(ctx) = auth.User{UID:"abc"}
        │
        ▼
recipes.Handler.Discover
        │   strconv.Atoi("?limit=") → raw int (or 0 on miss/invalid)
        ▼
recipes.Service.Discover(ctx, uid, raw)
        │   clamp(raw) → limit ∈ [1, 25]; raw≤0 → 10
        ▼
recipes.Repo.Discover(ctx, uid, limit)
        │   SELECT … WHERE NOT EXISTS (SELECT 1 FROM user_reactions …)
        │   ORDER BY random() LIMIT $2
        ▼
[]RecipeBrief                       ──► json.Encode → 200 {"recipes":[…]}
```

### Detail

```
GET /api/recipes/{slug}
        │
        ▼
auth.Middleware ──► UserFromContext(ctx)  (UID currently unused, but
                                           required for ownership in v2)
        │
        ▼
recipes.Handler.GetBySlug
        │   slug := r.PathValue("slug")
        ▼
recipes.Service.GetBySlug(ctx, slug)
        │
        ▼
recipes.Repo.GetBySlug(ctx, slug)
        │   1) SELECT … FROM recipes WHERE slug = $1
        │   2) SELECT i.name, i.category, ri.quantity, ri.unit, ri.optional
        │      FROM recipe_ingredients ri JOIN ingredients i ON …
        │      WHERE ri.recipe_id = $1
        │      ORDER BY i.category ASC, i.name ASC
        ▼
RecipeFull        ──► json.Encode → 200 <RecipeFull>
ErrNotFound       ──► 404 {"error":{"code":"NOT_FOUND",…}}
other error       ──► 500 {"error":{"code":"INTERNAL",…}}  (logged)
```

## 5.2 Dependency graph

```
cmd/cooksense-server
        │
        ▼
internal/recipes/handler  ──► internal/recipes/service ──► internal/recipes/repo
        │                                                       │
        ▼                                                       ▼
internal/auth                                          jackc/pgx/v5/pgxpool
```

The arrow direction is enforced by Clean Architecture and is verified
manually in code review. Reverse imports are a hard failure (Go compiler).

## 5.3 Layer responsibilities

| Layer | Owns | Forbidden |
|-------|------|-----------|
| **Handler** | HTTP parsing (`r.PathValue`, query params), JSON encode, HTTP status mapping, error envelope | SQL, business rules, `pgx` |
| **Service** | Use-case orchestration (limit clamping), structured logging, error propagation via `errors.Is` | `net/http`, `pgx`, `pgxpool` |
| **Repo** | SQL execution, row-to-DTO mapping, sentinel error translation (`pgx.ErrNoRows` → `ErrNotFound`) | HTTP, business rules |
| **DTO** | JSON shape; matches `docs/architecture/api.md` exactly | Mutability, validation |

## 5.4 SQL shape — Discover

```sql
SELECT r.slug, r.title, r.concept,
       r.time_minutes, r.passive_prep_minutes,
       r.cooking_methods, r.tags, r.flavor_profile
  FROM recipes r
 WHERE NOT EXISTS (
       SELECT 1
         FROM user_reactions ur
        WHERE ur.firebase_uid = $1
          AND ur.recipe_id    = r.id
       )
 ORDER BY random()
 LIMIT $2;
```

Exclusion is "any reaction kind" (per SPEC-DISCOVER-012). The index
`user_reactions_uid_kind_idx` on `(firebase_uid, kind)` (SPEC-DB DDL) keeps
the anti-join cheap; the planner uses the leading column for the
existence subquery.

## 5.5 SQL shape — Detail

The recipe row:

```sql
SELECT id, slug, title, concept,
       time_minutes, passive_prep_minutes,
       cooking_methods, tags, flavor_profile,
       steps
  FROM recipes
 WHERE slug = $1;
```

The ingredients (one extra roundtrip — see §3.4 for the alternative):

```sql
SELECT i.name, i.category,
       ri.quantity, ri.unit, ri.optional
  FROM recipe_ingredients ri
  JOIN ingredients i ON i.id = ri.ingredient_id
 WHERE ri.recipe_id = $1
 ORDER BY i.category ASC, i.name ASC;
```

`steps` is a JSONB column whose value SHALL be `json.Unmarshal`ed into a
`[]string` and emitted as-is in the response.

## 5.6 Error mapping

| Source | Service propagation | Handler response |
|--------|---------------------|------------------|
| `pgx.ErrNoRows` (recipe row) | wrapped as `fmt.Errorf("recipes.repo: %w", ErrNotFound)` | `404 NOT_FOUND` |
| Any other DB error | propagated unchanged | `500 INTERNAL` (logged) |
| Missing `auth.User` in context | `errors.New("missing auth user")` | `500 INTERNAL` (logged; middleware contract violation) |
| Empty result on Discover | NOT an error — returns `nil` slice | `200 {"recipes": []}` |

`errors.Is(err, recipes.ErrNotFound)` is the **only** authorized check the
handler performs to decide between `404` and `500`.

## 5.7 Route registration

```go
mux := http.NewServeMux()
authMW := authMiddleware(verifier, userTouch)

h := recipes.NewHandler(svc)
h.RegisterRoutes(mux, authMW)
// inside RegisterRoutes:
//   mux.Handle("GET /api/recipes/discover", authMW(http.HandlerFunc(h.Discover)))
//   mux.Handle("GET /api/recipes/{slug}",   authMW(http.HandlerFunc(h.GetBySlug)))
```

The router pattern `GET /api/recipes/discover` is more specific than
`GET /api/recipes/{slug}`, so the Go 1.22+ pattern matcher SHALL route
`/api/recipes/discover` to the discover handler — the literal `discover`
SHALL never be misinterpreted as a slug.
