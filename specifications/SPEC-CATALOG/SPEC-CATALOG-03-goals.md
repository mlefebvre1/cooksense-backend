# SPEC-CATALOG §3 — Goals, Non-Goals, Constraints

[← back to index](SPEC-CATALOG-00-index.md)

## 3.1 Goals

- **G1 — Day-one usefulness.** From the very first time the app runs, a user
  shall encounter at least 15 distinct, varied, weeknight-friendly recipes.
- **G2 — Schema confidence.** Authoring 15+ real recipes shall stress-test the
  SPEC-RECIPES schema and taxonomies before any feature endpoint depends on it.
- **G3 — Diversity for swipe UX.** The catalog shall span enough cooking
  methods, proteins, and dietary profiles to make the swipe-discovery feed
  feel non-repetitive.
- **G4 — Reviewability.** Each recipe shall be small, atomic, and reviewable
  in isolation; PR diffs shall stay readable as YAML.
- **G5 — Reproducibility.** `make seed` against a freshly migrated database
  shall load the entire catalog deterministically and idempotently.

## 3.2 Non-goals

- **NG1** — User-submitted recipes (post-MVP, see `docs/product/scope-mvp.md`).
- **NG2** — Photos or media references in recipe files.
- **NG3** — Localized (FR/PT/ES) recipe text — MVP catalog is English-only.
- **NG4** — A CMS or admin UI for editing recipes — content lives in Git.
- **NG5** — Cross-recipe relationships (e.g. "see also") — not in MVP schema.
- **NG6** — Nutritional facts, calorie counts, allergen flags — post-MVP.

## 3.3 Constraints

- **C1** — All files shall be valid UTF-8 YAML 1.2 parseable by `gopkg.in/yaml.v3`.
- **C2** — All recipes shall satisfy `time_minutes ≤ 30` (active time only).
- **C3** — Slug uniqueness shall be guaranteed across the catalog
  (enforced by SPEC-RECIPES loader).
- **C4** — Ingredient names shall converge to a single canonical spelling
  per ingredient across the catalog (SPEC-CATALOG-015).
- **C5** — No third-party content (recipes copied verbatim from books, sites)
  may ship — recipes shall be original or traditional/public-domain
  formulations restated in the project's own words.
