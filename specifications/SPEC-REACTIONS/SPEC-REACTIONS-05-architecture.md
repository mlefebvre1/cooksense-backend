# SPEC-REACTIONS — §5 Architecture

[← Index](SPEC-REACTIONS-00-index.md)

## 5.1 Layer diagram

```
┌─────────────────────────┐
│  internal/api           │   net/http handlers (POST/DELETE/GET)
│  reactions_handler.go   │   parses JSON, reads firebase_uid from ctx,
│  me_recipes_handler.go  │   calls reactions.Service, writes envelope
└──────────┬──────────────┘
           │
           ▼
┌─────────────────────────┐
│  internal/reactions     │   service.go     — slug → id resolution + delegation
│                         │   repo.go        — pgxpool SQL (UPSERT/DELETE/SELECT)
└──────────┬──────────────┘
           │
           ▼
┌─────────────────────────┐
│  internal/domain        │   reaction.go    — ReactionKind, Reaction,
│                         │                    sentinel errors, interface
└─────────────────────────┘
```

Dependency arrows point from outer layers to inner layers. The reverse is
forbidden.

## 5.2 Request flow — `POST /api/reactions`

```
Client → AuthMW → reactionsHandler.Post → reactions.Service.SetReaction
        ─────►  ───────────────────►   ┌─► recipes.Repo.ResolveSlug(slug) ──► recipes
                                        └─► reactions.Repo.Upsert(uid, id, kind) ──► user_reactions
        ◄────  ◄───────────────────  ◄───  returns created_at
        ↓
        200 { recipe_slug, kind, updated_at }
```

Error branches:
- `Decode` fails or required field empty → `400 INVALID_PAYLOAD`.
- `ParseReactionKind` fails → `422 INVALID_REACTION_KIND`.
- `ResolveSlug` returns `domain.ErrRecipeNotFound` → `404 NOT_FOUND`.
- Any other error → `500 INTERNAL` + `slog.Error`.

## 5.3 Request flow — `DELETE /api/reactions/{slug}`

```
Client → AuthMW → reactionsHandler.Delete → reactions.Service.RemoveReaction
                                            ┌─► recipes.Repo.ResolveSlug(slug)
                                            └─► reactions.Repo.Delete(uid, id)
        ↓
        204 (no body)
```

Error branches:
- `ResolveSlug` returns `domain.ErrRecipeNotFound` → `404 NOT_FOUND`.
- `Delete` affecting zero rows → still `204` (idempotent).
- Any other error → `500 INTERNAL`.

## 5.4 Request flow — `GET /api/me/recipes`

```
Client → AuthMW → meRecipesHandler.Get → reactions.Service.MyRecipes(uid, kind)
                                          └─► reactions.Repo.ListByKind(uid, kind)
        ↓
        200 { recipes: [<RecipeBrief>, …] }
```

Validation order in the handler:
1. Read `kind` query param. Empty → default `LIKE`.
2. `ParseReactionKind(value)` → if error → `400 INVALID_PAYLOAD`.
3. If `kind == DISLIKE` → `400 INVALID_PAYLOAD` (with code `INVALID_PAYLOAD`
   and message `"kind=DISLIKE is not exposed"`).
4. Service returns `[]RecipeBrief`. Empty slice → still `200` with empty
   array (never `null`).

## 5.5 Transaction shape

| Operation | Transaction? | Rationale |
|-----------|--------------|-----------|
| `Upsert` | Single statement | `INSERT … ON CONFLICT DO UPDATE` is atomic; no need for an explicit BEGIN/COMMIT. |
| `Delete` | Single statement | One DELETE statement; pgx auto-commits. |
| `ListByKind` | Single statement | Read-only SELECT JOIN. |

No multi-statement transactions are introduced by this spec.

## 5.6 Concurrency

- The pgx pool already serializes statements per connection.
- Two concurrent `POST /api/reactions` calls for the same `(uid, slug)`
  collapse into one row by the `(firebase_uid, recipe_id)` PK + ON CONFLICT
  clause; the last writer wins on `kind` and `created_at` (D-0008).
- No goroutines are spawned by SPEC-REACTIONS code outside the request
  handler; cancellation flows through `r.Context()`.

## 5.7 Error mapping (handler → HTTP)

| Source | Match | HTTP status | Envelope code |
|--------|-------|-------------|---------------|
| JSON decode error / missing field | `errors.Is(err, io.EOF)` or unmarshal error | 400 | `INVALID_PAYLOAD` |
| Empty `recipe_slug` | string check | 400 | `INVALID_PAYLOAD` |
| Bad reaction kind | `domain.ParseReactionKind` returns error | 422 | `INVALID_REACTION_KIND` |
| Unknown slug | `errors.Is(err, domain.ErrRecipeNotFound)` | 404 | `NOT_FOUND` |
| `MyRecipes` rejects `DISLIKE` | `errors.Is(err, domain.ErrInvalidMyRecipesKind)` | 400 | `INVALID_PAYLOAD` |
| Anything else | default | 500 | `INTERNAL` |
