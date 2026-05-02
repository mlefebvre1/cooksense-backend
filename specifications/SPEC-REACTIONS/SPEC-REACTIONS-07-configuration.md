# SPEC-REACTIONS — §7 Configuration

[← Index](SPEC-REACTIONS-00-index.md)

## 7.1 Environment variables

This spec introduces **no new environment variables**. It reuses what
SPEC-DB and SPEC-AUTH already define:

| Variable | Owner | Purpose for SPEC-REACTIONS |
|----------|-------|----------------------------|
| `DATABASE_URL` | SPEC-DB-007 | pgx pool connection used by `reactions.NewRepo` |
| `LOG_LEVEL` | SPEC-DB-009 | Handler/service `slog` levels |
| `LOG_FORMAT` | SPEC-DB-010 | Handler/service `slog` output format |
| `FIREBASE_*` | SPEC-AUTH | ID-token verification (consumed indirectly via auth middleware) |

`.env.example` SHALL NOT need updates beyond what SPEC-DB / SPEC-AUTH
already declare.

## 7.2 Route registration

The three routes added by this spec SHALL be registered against the same
`*http.ServeMux` instance that hosts the rest of the API. The auth
middleware SHALL wrap each handler using the helper `auth.Require(...)`
already provided by SPEC-AUTH.

| Method | Path | Auth | Handler |
|--------|------|------|---------|
| POST | `/api/reactions` | required | `api.PostReaction` |
| DELETE | `/api/reactions/{slug}` | required | `api.DeleteReaction` |
| GET | `/api/me/recipes` | required | `api.MyRecipes` |

## 7.3 Default behavior

| Setting | Default | Override |
|---------|---------|----------|
| `kind` query param on `GET /api/me/recipes` | `LIKE` | `?kind=TRY_LATER` (only) |
| Slug-resolution sentinel | `domain.ErrRecipeNotFound` | not configurable |
| JSON content type | `application/json; charset=utf-8` | not configurable |

## 7.4 Migration prerequisite

The `user_reactions` table and `reaction_kind` ENUM MUST already exist
(`migrations/0001_init.up.sql`). The server SHALL refuse to start (or
the routes SHALL return `500 INTERNAL` on first call) if the migration
has not been applied. SPEC-REACTIONS does NOT add a migration of its own.
