# SPEC-DISCOVER — §2 Introduction

[← Index](SPEC-DISCOVER-00-index.md)

## 2.1 Story relationship

| Source | Reference |
|--------|-----------|
| User story | `docs/stories/07-api-recipes-discover.md` |
| API contract | `docs/architecture/api.md` (§ "Recipes — discovery and detail") |
| Data model | `docs/architecture/data-model.md` (`recipes`, `ingredients`, `recipe_ingredients`, `user_reactions`) |
| Decisions | D-0007 (identity model: `firebase_uid TEXT` PK), D-0008 (`reaction_kind` enum) |
| Depends on | Story 03 (SPEC-DB), Story 04 (SPEC-AUTH), Story 05 (SPEC-RECIPES), Story 06 (curated content) |
| Blocks | Story 12 (README quickstart) |

## 2.2 Why this story exists

Without Story 07, the catalog ingested by SPEC-RECIPES is invisible to the
mobile client — there is no swipe feed and no detail page. Story 07 is the
first feature endpoint of the API; it establishes the **layered request
pipeline** every subsequent story (reactions, search, lessons) reuses:

1. Handler — parses `*http.Request`, encodes JSON, maps statuses.
2. Service — orchestrates and enforces use-case invariants (limit clamping).
3. Repository — owns SQL and is the sole `pgx` consumer.

The story also crystallizes the **discovery semantics** for MVP: any
recipe the caller has reacted to (LIKE / DISLIKE / TRY_LATER) SHALL NOT
appear in the feed again. This is a deliberate simplification of the
"DISLIKE-only by default" wording in the API doc.

## 2.3 Decision references

- **D-0001** — Go 1.26.2 (modern idioms only; `r.PathValue`, `mux.HandleFunc("GET …")`).
- **D-0002** — `net/http` standard library router (no third-party framework).
- **D-0003** — PostgreSQL 17 + `pgx/v5` (`ORDER BY random()` is acceptable
  for MVP catalog size per the story's technical notes).
- **D-0004** — Firebase ID-token authentication; the auth middleware
  (SPEC-AUTH) supplies `auth.User` via `UserFromContext`.
- **D-0007** — Identity model: `firebase_uid TEXT` is the user PK and the
  FK target in `user_reactions`. The discover SQL filters on this column
  directly (no internal user id).
- **D-0008** — `reaction_kind` enum with three values. Story 07 treats all
  three uniformly — any present row hides the recipe from discover.

## 2.4 Reading order

Readers SHOULD read in this order: §3 Goals → §5 Architecture → §6 Packages
(normative requirements) → §9 Testing → Appendix B Tasks.
