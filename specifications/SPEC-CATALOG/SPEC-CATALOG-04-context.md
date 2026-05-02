# SPEC-CATALOG §4 — System Context

[← back to index](SPEC-CATALOG-00-index.md)

## 4.1 Upstream dependencies (must be merged before this spec ships)

| Spec | Provides | Used for |
|------|----------|----------|
| SPEC-BOOT | `seed/recipes/` directory exists | File destination |
| SPEC-DB | DDL for `recipes`, `ingredients`, `recipe_ingredients` | Idempotent upsert target |
| SPEC-RECIPES | YAML schema, taxonomies, loader, `seed` subcommand | Validation + DB load |
| SPEC-MAKE | `make seed` target | E2E verification harness |

## 4.2 Decision references

| Decision | Title | Why it matters here |
|----------|-------|---------------------|
| D-0005 | Curated YAML catalog (no runtime LLM) | Justifies the very existence of `seed/recipes/` |
| D-0008 | English-only MVP catalog | SPEC-CATALOG-022 |
| D-0006 | Slug as public identifier | SPEC-CATALOG-003 (filename = slug) |

## 4.3 External dependencies introduced by this spec

**None.** SPEC-CATALOG introduces no new Go modules, no new env vars, no new
infrastructure. It only adds content files and one repository-level diversity
test that uses the existing `gopkg.in/yaml.v3` (already pulled in by
SPEC-RECIPES) and the standard library.

## 4.4 Boundary diagram

```
                ┌─────────────────────────────────┐
                │  Author writes YAML + opens PR  │
                └──────────────┬──────────────────┘
                               │
        ┌──────────────────────┼──────────────────────────┐
        │                      │                          │
        ▼                      ▼                          ▼
┌───────────────┐     ┌──────────────────┐     ┌────────────────────┐
│ SPEC-RECIPES  │     │  SPEC-CATALOG    │     │  Human reviewer    │
│ schema test   │     │  diversity test  │     │  (PR comment per   │
│ (per file)    │     │  (catalog-wide)  │     │   recipe)          │
└──────┬────────┘     └─────────┬────────┘     └─────────┬──────────┘
       │                        │                         │
       └────────────┬───────────┴─────────────────────────┘
                    ▼
         ┌─────────────────────────┐
         │  make seed (E2E load)   │
         │  against compose Postgres│
         └─────────────────────────┘
```

## 4.5 Out-of-scope dependencies

- Lessons content (`seed/lessons/`) — story 10, separate spec.
- Recipe photos — post-MVP.
- Aliases / synonyms — story 09 (search), separate spec.
