# SPEC-DISCOVER — §3 Goals, Non-Goals, Constraints

[← Index](SPEC-DISCOVER-00-index.md)

## 3.1 Goals

1. Ship `GET /api/recipes/discover` returning a randomized stream of
   recipes the caller has not yet reacted to, capped by `limit`.
2. Ship `GET /api/recipes/{slug}` returning the full recipe payload
   (ingredients sorted by category then name, steps in stored order).
3. Establish a **three-layer pattern** (handler / service / repo) inside
   `internal/recipes/` that subsequent feature stories will mirror.
4. Encode the **standard error envelope** (`docs/architecture/api.md`
   error catalog) once, so future stories reuse the same shape.
5. Provide an integration-test harness (compose Postgres + `fakeVerifier`)
   that exercises both endpoints end-to-end.

## 3.2 Non-goals

- The reactions endpoints (`POST /api/reactions`, `DELETE /api/reactions/{slug}`)
  — owned by Story 08 (SPEC-REACTIONS).
- The search endpoint (`GET /api/recipes/search`) — owned by Story 09.
- The "my recipes" endpoint (`GET /api/me/recipes`) — owned by Story 08.
- Lessons (`GET /api/lessons*`) — owned by Story 10.
- Pagination (`?cursor=`) — post-MVP per `docs/architecture/api.md`.
- Caching, ETag, `If-None-Match` — post-MVP.
- Personalization beyond "exclude already-reacted" — post-MVP.
- Time-to-decision telemetry — post-MVP per the story.

## 3.3 Hard constraints

- Go 1.26.2 only; modern idioms required.
- Standard-library router only (`*http.ServeMux` with Go 1.22 method+path
  patterns and `r.PathValue`).
- `internal/recipes/handler.go` SHALL NOT import `pgx`/`pgxpool`.
- `internal/recipes/service.go` SHALL NOT import `pgx`/`pgxpool` or
  `net/http`.
- `internal/recipes/repo.go` is the only file that uses `pgxpool.Pool`.
- JSON field names are `snake_case`, matching `docs/architecture/api.md`
  exactly.
- `limit` clamping: default `10`, max `25`. Values outside the range
  (including negatives, zero, malformed) SHALL fall back to default and
  SHALL NOT yield `400`.
- The discover query SHALL use `ORDER BY random()` (acceptable for MVP
  catalog size per the story's technical notes).

## 3.4 Soft constraints

- The repo SHOULD execute the detail query in two roundtrips (recipe row
  + ingredients) for clarity. A `LEFT JOIN LATERAL … jsonb_agg` single-shot
  variant is permitted iff it produces an identical end-state and is no
  less readable.
- The handler SHOULD reject HTTP methods other than `GET` via the mux
  pattern itself (i.e., the route `GET /api/recipes/discover` causes the
  mux to return `405 Method Not Allowed` automatically for non-GET).
- Logs SHOULD include the request id (when available from request-scoped
  context) for cross-correlation with the auth middleware.
