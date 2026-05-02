# Story 06 — Curate 15+ recipes as YAML content

Status: TODO
Estimate: M (mostly content, low code)

## User story

As Marc (the busy home cook), I want a curated catalog of at least 15 fast
weeknight recipes so that the app is useful from day one, with no empty
states.

## Background

This story is **content work**, not engineering. It exercises the YAML schema
introduced in story 05 and validates the taxonomies before any feature
endpoint ships. See `docs/product/scope-mvp.md` and
`docs/architecture/data-model.md`.

## Acceptance criteria

- [ ] `seed/recipes/` contains **≥ 15** YAML files, one recipe per file.
      Filenames match the recipe `slug` (e.g.
      `pan-seared-chicken-thighs-lemon-cabbage.yaml`).
- [ ] Each recipe respects the schema and invariants from the data-model doc:
  - `slug`, `title`, `concept`, `time_minutes`, `cooking_methods`, `tags`,
    `flavor_profile`, `ingredients[]`, `steps[]` are present and valid.
  - All `cooking_methods`, `flavor_profile`, and ingredient `category` values
    are within the allowed taxonomies.
  - `time_minutes ≤ 30`, with `passive_prep_minutes` allowed if longer total.
  - At least 2 ingredients and 1 step.
- [ ] Catalog diversity:
  - At least 4 different `cooking_methods` represented across the catalog.
  - At least 3 distinct primary proteins (e.g. chicken, beef, fish, tofu, eggs).
  - At least 2 vegetarian recipes.
  - At least 1 slow-cooker / pressure-cooker recipe (uses
    `passive_prep_minutes`).
- [ ] Every recipe has been read end-to-end by a human reviewer (peer review
      tracked in PR comments).
- [ ] `make seed` loads all recipes without errors.

## Technical notes

- Tone: terse, expert. No backstory. No "the secret is…". One concept, key
  ingredients, essential steps.
- Steps should describe **what to do and why** when non-obvious; never
  describe trivial mechanics.
- Reuse ingredient names verbatim across recipes so the ingredients table
  consolidates correctly (e.g. always `garlic`, never `garlic clove(s)`).
  Quantity and unit live in the recipe-ingredient row, not the name.
- Aliases (e.g. `boeuf` ↔ `beef`) belong in story 09 (search), not here.

## Out of scope

- Photos.
- Translations (catalog is English in MVP; FR catalog is a future story).
- User-submitted content.

## Dependencies

- depends on: 05
- blocks: 07, 08, 09, 12

## Definition of Done

- [ ] AC met.
- [ ] Catalog reviewed in a single PR; reviewer leaves at least one comment
      on each recipe.
- [ ] Spreadsheet/table in the PR description summarizing each recipe's
      `time_minutes`, `cooking_methods`, and primary protein for at-a-glance
      diversity check.
