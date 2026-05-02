# Implementation stories — backlog

Each story is a self-contained spec ready to be picked up. Implement them in
the order below unless explicitly approved otherwise.

## Story template

```
# Story <NN> — <title>
Status: TODO | IN PROGRESS | DONE
Owner: <handle>
Estimate: S | M | L

## User story
As a … I want … so that …

## Background / Context
…

## Acceptance criteria
- [ ] …

## Technical notes
…

## Out of scope
…

## Dependencies
- depends on: #NN
- blocks: #NN

## Definition of Done
- [ ] AC met
- [ ] Tests pass
- [ ] go vet clean
- [ ] README/docs updated if a new env var or command was introduced
```

## Backlog (order = recommended implementation order)

| #  | Story                                                                                                | Size | Depends on | Critical path |
|----|------------------------------------------------------------------------------------------------------|------|------------|---------------|
| 01 | [Bootstrap project structure](./01-bootstrap-project-structure.md)                                   | S    | —          | ✅            |
| 02 | [Makefile targets wrapping docker-compose](./02-makefile-targets.md)                                 | S    | 01         | ✅            |
| 03 | [Database — pgx pool + migrations + 0001_init](./03-db-pgx-migrations.md)                            | M    | 01, 02     | ✅            |
| 04 | [Firebase ID token middleware + lazy user provisioning](./04-firebase-auth-middleware.md)            | M    | 01, 03     | ✅            |
| 05 | [Recipe domain + YAML loader + `seed` subcommand](./05-recipes-domain-yaml-loader.md)                | M    | 03         | ✅            |
| 06 | [Curate 15+ recipes as YAML content](./06-curated-recipes-content.md)                                | M    | 05         | ✅            |
| 07 | [`GET /api/recipes/discover` + `GET /api/recipes/{slug}`](./07-api-recipes-discover.md)              | M    | 04, 05     | ✅            |
| 08 | [`POST /api/reactions`, `DELETE /api/reactions/{slug}`, `GET /api/me/recipes`](./08-api-reactions-my-recipes.md) | M | 04, 05 | ✅ |
| 09 | [`GET /api/recipes/search?ingredients=…`](./09-api-recipes-search.md)                                | M    | 04, 05     | ✅            |
| 10 | [Lessons API + 3 seed articles](./10-api-lessons.md)                                                 | S    | 04         |               |
| 11 | [Integration tests skeleton](./11-integration-tests-skeleton.md)                                     | M    | 03, 04     |               |
| 12 | [README quickstart](./12-readme-quickstart.md)                                                       | S    | 02, 03, 06 |               |

## Dependency graph

```
01 ─► 02 ─► 03 ─► 04 ─┬─► 07 ─┐
                     ├─► 08 ─┼─► (MVP demo)
                     ├─► 09 ─┘
                     └─► 10
05 ─► 06 ─► (catalog ready) ─► 07/08/09 use it
03/04 ─► 11
02/03/06 ─► 12
```

## What "MVP demo ready" means

Stories 01 → 09 + 06 + 12 must be DONE. 10 and 11 are nice-to-haves that
should be done before merging to main but can slip the demo deadline by a
day if needed.
