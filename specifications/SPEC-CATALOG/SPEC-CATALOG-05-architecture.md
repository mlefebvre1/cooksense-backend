# SPEC-CATALOG §5 — Architecture

[← back to index](SPEC-CATALOG-00-index.md)

## 5.1 File layout

```
seed/
└── recipes/
    ├── _sample.yaml                              # SPEC-RECIPES fixture (NOT counted)
    ├── pan-seared-chicken-thighs-lemon-cabbage.yaml
    ├── miso-glazed-salmon-broccoli.yaml
    ├── ...                                       # ≥ 15 catalog files total
    └── pressure-cooker-beef-ragu.yaml            # ≥ 1 slow/pressure-cook recipe
```

- One recipe per file (SPEC-CATALOG-002).
- File basename equals `<slug>.yaml` (SPEC-CATALOG-003).
- Underscore-prefixed files are loader fixtures and **excluded** from catalog
  counts and diversity assertions (SPEC-CATALOG-004).

## 5.2 Authoring pipeline

```
1. Author drafts a recipe.yaml in a feature branch.
2. Local check:    yamllint + `make seed` against compose Postgres.
3. PR opened:      includes summary table (SPEC-CATALOG-025).
4. CI runs:        SPEC-RECIPES per-file validation + SPEC-CATALOG diversity test.
5. Human review:   ≥ 1 reviewer comment per recipe (SPEC-CATALOG-024).
6. Merge:          recipe becomes part of the catalog.
```

## 5.3 Verification flow

```
                 ┌────────────────────────────┐
                 │  seed/recipes/*.yaml       │
                 └────────────┬───────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌──────────────┐    ┌──────────────────┐    ┌────────────────────┐
│ Schema test  │    │ Diversity test   │    │ E2E `make seed`    │
│ (SPEC-RECIPES│    │ (SPEC-CATALOG-   │    │ against compose pg │
│   loader)    │    │  011..014, 027)  │    │ → exit 0           │
└──────┬───────┘    └────────┬─────────┘    └─────────┬──────────┘
       │                     │                         │
       └─────────────────────┴─────────────────────────┘
                             │
                       PR mergeable?
```

## 5.4 Ingredient name canonicalization

Because the loader UPSERTs ingredients by `name`, two recipes that say
`garlic` and `garlic clove` will produce **two distinct rows**, polluting
the table. SPEC-CATALOG-015 prevents this by mandating verbatim reuse:

| Canonical name | Allowed | Forbidden |
|----------------|---------|-----------|
| `garlic` | `garlic` | `garlic clove`, `garlic cloves`, `garlic (peeled)` |
| `olive oil` | `olive oil` | `extra-virgin olive oil`, `EVOO` |
| `lemon` | `lemon` | `lemons`, `whole lemon`, `lemon juice` (use a separate ingredient `lemon juice` if applicable) |
| `chicken thighs (bone-in, skin-on)` | qualified form *if* the cut materially changes the recipe | mixing `chicken thighs` and `chicken thighs (bone-in, skin-on)` interchangeably |

Quantities and units (`1`, `tbsp`, `clove`) live on the
`recipe_ingredients` row — never on the ingredient name (SPEC-CATALOG-016).

## 5.5 Diversity invariants — formal definitions

Let *C* = set of catalog recipes (filenames not starting with `_`).

- **D-METHODS** (SPEC-CATALOG-011): `|⋃_{r∈C} r.cooking_methods| ≥ 4`.
- **D-PROTEINS** (SPEC-CATALOG-012): the set of *primary proteins* across
  *C* has cardinality ≥ 3, where the primary protein of recipe *r* is the
  ingredient with `category: protein` of greatest `quantity` (ties broken
  by file order); for vegetarian recipes (no animal protein), the primary
  protein is the value `vegetarian`.
- **D-VEGETARIAN** (SPEC-CATALOG-013): the count of recipes with no
  ingredient whose `category: protein` matches a configured set of
  animal-origin canonical names (`chicken thighs`, `chicken breast`,
  `beef ground`, `beef chuck`, `pork`, `salmon`, `cod`, `shrimp`,
  `eggs`, `bacon`, …) shall be ≥ 2.
- **D-PASSIVE** (SPEC-CATALOG-014): there shall exist *r* ∈ *C* with
  `slow-cook ∈ r.cooking_methods ∨ pressure-cook ∈ r.cooking_methods`
  AND `r.passive_prep_minutes > 0`.

The animal-origin protein list is maintained in
`seed/recipes/_animal_proteins.txt` (one canonical name per line). The
diversity test reads this file at runtime so adding new proteins does not
require a code change.

## 5.6 No code path is added by SPEC-CATALOG

The only new code artifact is the **diversity test** (`internal/seed/
catalog_diversity_test.go` or equivalent under `tests/`), defined in §9.
Production code is untouched.
