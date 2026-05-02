# SPEC-RECIPES — §2 Introduction

[← Index](SPEC-RECIPES-00-index.md)

## 2.1 Story relationship

| Source | Reference |
|--------|-----------|
| User story | `docs/stories/05-recipes-domain-yaml-loader.md` |
| Data model | `docs/architecture/data-model.md` (DDL + YAML schema + taxonomies) |
| Decision | D-0005 (curated YAML seed, no runtime AI) — `docs/product/decisions.md` |
| Depends on | Story 03 (SPEC-DB) — pgx pool, migrations, `db.Open` |
| Blocks | Stories 06 (curated content), 07 (discover API), 08 (reactions API), 09 (search API) |

## 2.2 Why this story exists

Without a domain model and a loader, the catalog cannot reach the database
and no API can return real data. Story 05 establishes:

1. The **canonical in-memory shape** of a recipe — used by every downstream
   handler/repo.
2. The **content pipeline**: YAML on disk → validated `domain.Recipe` →
   Postgres rows.
3. The **operator UX**: a single `make seed` command that ingests recipes
   idempotently across environments.

## 2.3 Decision references

- D-0001 — Go 1.26.2 (loader uses modern idioms only).
- D-0002 — Postgres 17 (UPSERT semantics, `ON CONFLICT`).
- D-0003 — `net/http` stdlib (no impact here; mentioned for context).
- D-0005 — Curated YAML seed instead of runtime AI generation. Drives the
  whole loader/store design.

## 2.4 Reading order

Readers SHOULD read in this order: §3 Goals → §5 Architecture → §6 Packages
(normative requirements) → §9 Testing → Appendix B Tasks.
