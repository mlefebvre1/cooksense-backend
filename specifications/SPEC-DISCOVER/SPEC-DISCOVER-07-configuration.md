# SPEC-DISCOVER — §7 Configuration

[← Index](SPEC-DISCOVER-00-index.md)

## 7.1 Limits and defaults

| Setting | Source | Value | Override |
|---------|--------|-------|----------|
| `DefaultDiscoverLimit` | Code constant (`internal/recipes/service.go`) | `10` | Compile-time only |
| `MaxDiscoverLimit` | Code constant (`internal/recipes/service.go`) | `25` | Compile-time only |
| `?limit=N` query parameter | HTTP request | clamped per SPEC-DISCOVER-018 | Client-controlled |

The story explicitly fixes both bounds; they SHALL NOT be made
runtime-configurable in MVP. A future ADR is required to introduce
environment-driven overrides.

## 7.2 Environment variables

The discover/detail feature reuses environment variables already defined
by upstream stories:

| Variable | Owner | Purpose for SPEC-DISCOVER |
|----------|-------|---------------------------|
| `DATABASE_URL` | SPEC-DB-007 | Connect to Postgres for `recipes.NewPgRepo` |
| `LOG_LEVEL` | SPEC-DB-009 | Service/handler structured logs |
| `LOG_FORMAT` | SPEC-DB-010 | Service/handler structured logs |
| `FIREBASE_PROJECT_ID` | SPEC-AUTH | Verifier wired into the auth middleware that protects the routes |

No new environment variables are introduced by Story 07. `.env.example`
SHALL NOT need updates beyond what SPEC-DB and SPEC-AUTH already define.

## 7.3 HTTP contract defaults

| Aspect | Value |
|--------|-------|
| Authentication | Required on both routes (`auth.Middleware` wraps them) |
| Methods | `GET` only — non-GET on either path yields `405 Method Not Allowed` (Go 1.22 mux behavior) |
| Response media type | `application/json; charset=utf-8` |
| Cache headers | None (post-MVP) |
| CORS | Out of scope for Story 07 |

## 7.4 Default values summary

| Setting | Default | Override |
|---------|---------|----------|
| Discover `limit` (no/empty/invalid query) | `10` | `?limit=N` (clamped 1..25) |
| Discover ordering | `random()` | not configurable |
| Detail ingredient ordering | `category ASC, name ASC` | not configurable |
| Detail steps ordering | as stored in JSONB | not configurable |
