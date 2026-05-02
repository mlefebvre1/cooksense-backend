# SPEC-CATALOG — Curated Recipe Catalog (≥15 YAML files)

**Story:** 06 — Curate 15+ recipes as YAML content
**Status:** Draft → Final
**Date:** 2026-05-02
**Authors:** sdd-spec-author

---

## Purpose

This specification governs the **content** that ships in `seed/recipes/` for the
MVP catalog: file naming, per-recipe schema compliance, catalog-wide diversity
constraints, editorial tone, ingredient-name normalization, peer-review process,
and the end-to-end verification that `make seed` loads the catalog without
errors.

This spec is intentionally narrow: the YAML *schema* and the *loader* are owned
by [SPEC-RECIPES](../SPEC-RECIPES/SPEC-RECIPES-00-index.md). SPEC-CATALOG only
adds normative rules about **which recipes ship and how they are written**.

---

## File map

| File | Contents |
|------|----------|
| [SPEC-CATALOG-01-preamble](SPEC-CATALOG-01-preamble.md) | AI constraints, scope boundary with SPEC-RECIPES |
| [SPEC-CATALOG-02-introduction](SPEC-CATALOG-02-introduction.md) | Story relationship, SPEC-ID registry |
| [SPEC-CATALOG-03-goals](SPEC-CATALOG-03-goals.md) | Goals, non-goals, constraints |
| [SPEC-CATALOG-04-context](SPEC-CATALOG-04-context.md) | Dependencies on SPEC-RECIPES & data-model |
| [SPEC-CATALOG-05-architecture](SPEC-CATALOG-05-architecture.md) | File layout, review pipeline, verification flow |
| [SPEC-CATALOG-06-packages](SPEC-CATALOG-06-packages.md) | All SPEC-CATALOG-NNN normative rules |
| [SPEC-CATALOG-07-configuration](SPEC-CATALOG-07-configuration.md) | No new env vars; pointers to existing config |
| [SPEC-CATALOG-08-build](SPEC-CATALOG-08-build.md) | Lint/format hooks for YAML, no Go build impact |
| [SPEC-CATALOG-09-testing](SPEC-CATALOG-09-testing.md) | Diversity test, schema test, end-to-end seed test |
| [SPEC-CATALOG-10-documentation](SPEC-CATALOG-10-documentation.md) | PR description template, reviewer checklist |
| [SPEC-CATALOG-A-checklist](SPEC-CATALOG-A-checklist.md) | Specification completeness checklist |
| [SPEC-CATALOG-B-tasks](SPEC-CATALOG-B-tasks.md) | Ordered authoring + verification task list |

---

## SPEC-ID registry

| ID | Summary | Section |
|----|---------|---------|
| SPEC-CATALOG-001 | `seed/recipes/` shall contain at least 15 distinct YAML files | §6.1 |
| SPEC-CATALOG-002 | Each catalog file shall be one recipe — no multi-document YAML | §6.1 |
| SPEC-CATALOG-003 | Filename shall equal `<slug>.yaml` (kebab-case) | §6.1 |
| SPEC-CATALOG-004 | Filenames starting with `_` (e.g. `_sample.yaml`) are reserved fixtures and not counted toward the catalog | §6.1 |
| SPEC-CATALOG-005 | Every catalog file shall validate against the SPEC-RECIPES schema | §6.2 |
| SPEC-CATALOG-006 | Every recipe shall declare `slug`, `title`, `concept`, `time_minutes`, `cooking_methods`, `tags`, `flavor_profile`, `ingredients`, `steps` | §6.2 |
| SPEC-CATALOG-007 | `time_minutes` shall be `≤ 30` for every recipe | §6.2 |
| SPEC-CATALOG-008 | Recipes whose total active+passive time exceeds 30 minutes shall set `passive_prep_minutes > 0` and keep `time_minutes ≤ 30` | §6.2 |
| SPEC-CATALOG-009 | Each recipe shall have ≥ 2 ingredients and ≥ 1 step (inherits SPEC-RECIPES invariant; restated for completeness) | §6.2 |
| SPEC-CATALOG-010 | All `cooking_methods`, `flavor_profile`, and ingredient `category` values shall be drawn from the SPEC-RECIPES taxonomies | §6.2 |
| SPEC-CATALOG-011 | The catalog shall represent ≥ 4 distinct `cooking_methods` across all recipes | §6.3 |
| SPEC-CATALOG-012 | The catalog shall include ≥ 3 distinct primary proteins (chicken, beef, fish, tofu, eggs, …) | §6.3 |
| SPEC-CATALOG-013 | The catalog shall include ≥ 2 vegetarian recipes (no ingredient with `category: protein` of animal origin) | §6.3 |
| SPEC-CATALOG-014 | The catalog shall include ≥ 1 recipe whose `cooking_methods` contains `slow-cook` or `pressure-cook` and whose `passive_prep_minutes > 0` | §6.3 |
| SPEC-CATALOG-015 | Ingredient `name` strings shall be reused verbatim across recipes (e.g. always `garlic`, never `garlic clove(s)`) | §6.4 |
| SPEC-CATALOG-016 | Quantity and unit shall live on the `recipe_ingredients` row, never on the ingredient `name` | §6.4 |
| SPEC-CATALOG-017 | Ingredient `name` shall be lowercase except where a proper noun is intrinsic (e.g. `Dijon mustard`) | §6.4 |
| SPEC-CATALOG-018 | Ingredient aliases (e.g. `boeuf` ↔ `beef`) shall NOT appear in story 06 — they are deferred to story 09 | §6.4 |
| SPEC-CATALOG-019 | `concept` shall be a single paragraph stating the dish idea and the single most important technique decision | §6.5 |
| SPEC-CATALOG-020 | `concept` and `steps` shall use a terse, expert tone — no backstory, no "the secret is…", no marketing copy | §6.5 |
| SPEC-CATALOG-021 | `steps` shall describe what to do and why when the reason is non-obvious; trivial mechanics shall be omitted | §6.5 |
| SPEC-CATALOG-022 | Catalog content shall be authored in English — translations are out of scope | §6.5 |
| SPEC-CATALOG-023 | Photos, image URLs, or external media references shall NOT appear in any recipe file | §6.5 |
| SPEC-CATALOG-024 | Every recipe file shall be reviewed by at least one human reviewer who leaves a PR comment on it | §6.6 |
| SPEC-CATALOG-025 | The PR description shall include a summary table listing each recipe's `slug`, `time_minutes`, `cooking_methods`, and primary protein | §6.6 |
| SPEC-CATALOG-026 | `make seed` shall complete with exit 0 against a freshly migrated database loaded with the catalog | §6.7 |
| SPEC-CATALOG-027 | A repository-level test shall enforce the diversity invariants (SPEC-CATALOG-011 through SPEC-CATALOG-014) | §9 |
