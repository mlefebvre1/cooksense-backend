# SPEC-REACTIONS — §10 Documentation

[← Index](SPEC-REACTIONS-00-index.md)

## 10.1 Doc-comment requirements

Every exported symbol introduced by Story 08 SHALL carry a Go Doc comment
that:

1. Begins with the symbol name (`// PostReaction is …`, `// Repo is …`).
2. References the SPEC-REACTIONS-NNN ID(s) it implements.
3. Documents pre-conditions, post-conditions, and error behavior.

| File | Required doc-comment subjects | SPEC-IDs to cite |
|------|-------------------------------|------------------|
| `internal/domain/reaction.go` | `ReactionKind`, constants, `ParseReactionKind`, `Valid`, `Reaction`, `ErrInvalidReactionKind`, `ErrInvalidMyRecipesKind`, `ReactionRepository` | 001–006 |
| `internal/reactions/doc.go` | Package overview: pgx-backed repo + slug-resolving service | 007, 017, 043 |
| `internal/reactions/repo.go` | `Repo`, `NewRepo`, `Upsert`, `Delete`, `ListByKind` | 007–014 |
| `internal/reactions/service.go` | `Service`, `NewService`, `SetReaction`, `RemoveReaction`, `MyRecipes` | 017–021 |
| `internal/api/reactions_post_handler.go` | `PostReaction` | 022–028 |
| `internal/api/reactions_delete_handler.go` | `DeleteReaction` | 029–032 |
| `internal/api/me_recipes_handler.go` | `MyRecipes` | 033–037 |
| `cmd/cooksense-server/main.go` | wiring section comment | 041, 042 |

### SPEC-REACTIONS-043 — Package comment for `internal/reactions`

The `internal/reactions/doc.go` file SHALL contain a package comment of the
form:

```go
// Package reactions implements the user-reactions persistence (Repo) and
// the slug-resolving orchestration layer (Service) for the
// POST /api/reactions, DELETE /api/reactions/{slug}, and
// GET /api/me/recipes endpoints.
//
// Implements: SPEC-REACTIONS-007 through SPEC-REACTIONS-021,
// SPEC-REACTIONS-041, SPEC-REACTIONS-042.
package reactions
```

## 10.2 Architecture documentation

### SPEC-REACTIONS-044 — `docs/architecture/api.md`

`docs/architecture/api.md` already specifies the three endpoints under the
"Reactions" and "My recipes" sections. **No edits to that document are
required by Story 08** — this spec implements the document verbatim.

If implementation diverges from the document for any reason (e.g., added
`updated_at` field, additional error code), the architecture doc SHALL be
amended **in the same PR** as the code, per the repo's
"spec and code ship together" rule.

## 10.3 Decision log

No new ADR is required. D-0008 (Reaction kinds) already covers the
high-level decision. If implementation choices in §6 diverge from D-0008
(e.g., introducing a history table, switching to a textual `kind` column),
a new ADR superseding D-0008 SHALL be added.

## 10.4 README impact

Story 12 owns the user-facing quickstart README. Story 08 SHALL NOT touch
`README.md` unless the implementation introduces a new environment variable
or required setup step (none expected — see §7).

## 10.5 OpenAPI / API reference

The repo does not currently publish an OpenAPI document. Should one be
introduced in a future story, the three SPEC-REACTIONS endpoints SHALL be
described in it; that future change is owned by that story, not by
SPEC-REACTIONS.
