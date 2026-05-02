# HTTP API

All payloads are JSON (UTF-8). Times are RFC 3339 (`time.RFC3339`).

## Conventions

- **Authentication**: every route except `/api/health` requires
  `Authorization: Bearer <Firebase ID token>`. Missing/invalid → `401`.
- **Errors**: standardized envelope:
  ```json
  { "error": { "code": "INVALID_INGREDIENT", "message": "human readable" } }
  ```
  Codes are stable strings. HTTP status is set appropriately
  (`400` validation, `401` auth, `403` forbidden, `404` not found, `409`
  conflict, `422` semantic, `500` internal).
- **IDs**: recipes are addressed by `slug`. Numeric `id` never leaves the
  backend.
- **Pagination**: not required for MVP (catalog is small). When introduced:
  `?limit=` + `?cursor=` (opaque base64).
- **Idempotency**: `POST /api/reactions` is idempotent (UPSERT on
  `(uid, recipe_id)`).

## Endpoints

### Health

```
GET /api/health                                          (public)
200 → { "status": "ok", "version": "<git sha>" }
```

```
GET /api/health/me                                       (auth)
200 → { "uid": "abc123", "email": "marc@example.com" }
401 → standard error envelope
```

### Recipes — discovery and detail

```
GET /api/recipes/discover?limit=10                       (auth)
```
- Returns recipes the caller has **not yet reacted to**, randomized.
- Excludes `DISLIKE`d recipes by default (they can never come back).
- `limit` defaults to 10, max 25.
- `200 → { "recipes": [<RecipeBrief>, …] }`.

```
GET /api/recipes/{slug}                                  (auth)
200 → <RecipeFull>
404 → if the slug is unknown
```

`<RecipeBrief>`:
```json
{
  "slug": "pan-seared-chicken-thighs-lemon-cabbage",
  "title": "Pan-seared chicken thighs, lemon-braised cabbage",
  "concept": "Sear bone-in thighs hard, …",
  "time_minutes": 25,
  "passive_prep_minutes": 0,
  "cooking_methods": ["pan-sear", "braise"],
  "tags": ["weeknight", "one-pan"],
  "flavor_profile": ["acid", "fat", "umami"]
}
```

`<RecipeFull>` = `<RecipeBrief>` plus:
```json
{
  "ingredients": [
    { "name": "chicken thighs", "category": "protein",
      "quantity": 4, "unit": "piece", "optional": false }
  ],
  "steps": ["…", "…"]
}
```

### Recipes — search by ingredients

```
GET /api/recipes/search?ingredients=tomato,beef&scope=all|liked   (auth)
```
- `ingredients`: comma-separated list, **2 to 5 entries**, case- and
  accent-insensitive, matched against `ingredients.name` and `aliases`.
- `scope=all` (default): all recipes minus the caller's `DISLIKE`s.
- `scope=liked`: only recipes the caller has reacted with `LIKE`.
- Results are ranked by **number of matched ingredients DESC**, then
  `time_minutes ASC`.
- `200 → { "recipes": [<RecipeBrief>, …], "matched": { "<slug>": 3 } }`
- `400` if fewer than 2 or more than 5 ingredients.

### Reactions

```
POST /api/reactions                                       (auth)
Body: { "recipe_slug": "…", "kind": "LIKE" | "DISLIKE" | "TRY_LATER" }
```
- UPSERT: subsequent calls with a different `kind` overwrite the previous one.
- `200 → { "recipe_slug": "…", "kind": "LIKE", "updated_at": "…" }`
- `404` if `recipe_slug` is unknown.
- `422` if `kind` is not in the enum.

```
DELETE /api/reactions/{recipe_slug}                       (auth)
204 (also if no reaction existed — DELETE is idempotent)
```

### My recipes

```
GET /api/me/recipes?kind=LIKE|TRY_LATER                   (auth)
```
- Defaults to `kind=LIKE`. Returns recipes the caller reacted to with that
  kind, sorted by reaction `created_at DESC`.
- `200 → { "recipes": [<RecipeBrief>, …] }`

### Lessons (Cooking School)

```
GET /api/lessons                                          (auth)
200 → { "lessons": [{ "slug":"…","title":"…","sort_order":1 }] }
```

```
GET /api/lessons/{slug}                                   (auth)
200 → { "slug":"…", "title":"…", "body_md":"…" }
404 → unknown slug
```

## Error code catalog (initial)

| Code                       | HTTP | When                                          |
|----------------------------|------|-----------------------------------------------|
| `UNAUTHENTICATED`          | 401  | Missing/invalid/expired Firebase token         |
| `FORBIDDEN`                | 403  | Token valid but caller not allowed             |
| `NOT_FOUND`                | 404  | Resource does not exist                        |
| `INVALID_PAYLOAD`          | 400  | Malformed JSON or missing required field       |
| `INVALID_INGREDIENT_COUNT` | 400  | `ingredients` query param outside [2,5]        |
| `INVALID_REACTION_KIND`    | 422  | `kind` outside enum                            |
| `INTERNAL`                 | 500  | Unhandled error (logged with request id)       |
