# SPEC-CATALOG §2 — Introduction

[← back to index](SPEC-CATALOG-00-index.md)

## 2.1 Story relationship

| Artifact | Reference |
|----------|-----------|
| User story | `docs/stories/06-curated-recipes-content.md` |
| Persona | Marc — busy home cook (see `docs/product/personas.md`) |
| MVP scope | `docs/product/scope-mvp.md` — "≥ 15 curated recipes from day one" |
| Decisions | D-0005 (curated YAML catalog), D-0008 (English-only MVP catalog) |
| Upstream specs | SPEC-RECIPES (schema + loader), SPEC-DB (DDL) |
| Downstream stories | 07 (discover), 08 (reactions), 09 (search), 12 (README quickstart) |

## 2.2 Why this is a separate spec

Story 05 ships the **mechanism** (schema, loader, taxonomies, `seed` command).
Story 06 ships the **content** that mechanism consumes. They are decoupled so
that:

- The loader can be merged and tested with a single `_sample.yaml` fixture
  before the full catalog exists.
- Catalog diversity rules can evolve (more proteins, more methods) without
  touching loader code.
- Reviewers can focus on culinary correctness (tone, technique) in PRs that
  contain only YAML.

## 2.3 Out-of-band inputs

No external data sources. Recipes are authored by hand by humans, optionally
with LLM drafting offline followed by human edit and review. The result is
checked into Git in `seed/recipes/`.

## 2.4 Document conventions

- "Catalog" means the set of YAML files in `seed/recipes/` whose filename does
  NOT start with `_`. Underscore-prefixed files are reserved fixtures.
- "Recipe" refers to one YAML file inside the catalog.
- "Reviewer" means a human teammate who reads the recipe end-to-end and leaves
  at least one PR comment on it.
