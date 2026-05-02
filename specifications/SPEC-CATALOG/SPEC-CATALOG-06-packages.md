# SPEC-CATALOG §6 — Normative Catalog Rules

[← back to index](SPEC-CATALOG-00-index.md)

This section enumerates every SPEC-CATALOG-NNN requirement. Each subsection
is grouped by concern. RFC-2119 keywords are normative.

---

## 6.1 File layout & inventory

### SPEC-CATALOG-001 — Minimum catalog size

`seed/recipes/` **shall** contain at least 15 catalog files (i.e. files whose
basename does not start with `_` and whose extension is `.yaml`).

### SPEC-CATALOG-002 — One recipe per file

Each catalog file **shall** contain exactly one YAML document defining one
recipe. Multi-document YAML (`---` separators) is forbidden.

### SPEC-CATALOG-003 — Filename equals slug

For every catalog file *f*, `basename(f) == r.slug + ".yaml"` where `r` is
the parsed recipe. Slug format is governed by SPEC-RECIPES
(`^[a-z0-9]+(-[a-z0-9]+)*$`).

### SPEC-CATALOG-004 — Reserved fixture prefix

Files in `seed/recipes/` whose basename starts with `_` are reserved as
loader fixtures (SPEC-RECIPES). They **shall** be excluded from
SPEC-CATALOG-001, the diversity invariants, and any catalog assertion.

---

## 6.2 Per-recipe schema compliance

### SPEC-CATALOG-005 — Schema validation

Every catalog file **shall** parse cleanly via the SPEC-RECIPES loader
(`internal/seed.Load`) without error.

### SPEC-CATALOG-006 — Required fields

Every recipe **shall** declare non-empty values for all of:
`slug`, `title`, `concept`, `time_minutes`, `cooking_methods`, `tags`,
`flavor_profile`, `ingredients`, `steps`.

### SPEC-CATALOG-007 — Active-time bound

`time_minutes` **shall** satisfy `0 < time_minutes ≤ 30`.

### SPEC-CATALOG-008 — Passive-time accommodation

If a recipe's total wall-clock duration exceeds 30 minutes, the recipe
**shall** model the excess as `passive_prep_minutes > 0` while keeping
`time_minutes ≤ 30`. The semantics: `time_minutes` covers active
attention; `passive_prep_minutes` covers unattended cooking (slow cooker,
marinade, oven idle).

### SPEC-CATALOG-009 — Minimum content

Every recipe **shall** declare ≥ 2 entries in `ingredients` and ≥ 1 entry in
`steps`. (Inherited from SPEC-RECIPES; restated for completeness.)

### SPEC-CATALOG-010 — Taxonomy compliance

All values in `cooking_methods`, `flavor_profile`, and per-ingredient
`category` **shall** be drawn from the canonical sets defined in
SPEC-RECIPES §6.3 / `internal/domain/taxonomy.go`. Values outside those
sets are loader-rejected.

---

## 6.3 Catalog-wide diversity

### SPEC-CATALOG-011 — Methods diversity

The union of `cooking_methods` across all catalog recipes **shall** have
cardinality ≥ 4.

### SPEC-CATALOG-012 — Protein diversity

The catalog **shall** include recipes whose primary proteins span ≥ 3
distinct canonical names. The primary-protein resolution rule is defined
in SPEC-CATALOG-05 §5.5. Vegetarian recipes contribute the literal
primary-protein value `vegetarian`.

### SPEC-CATALOG-013 — Vegetarian floor

The catalog **shall** include ≥ 2 recipes that contain no ingredient whose
`category` is `protein` AND whose `name` appears in the animal-origin list
at `seed/recipes/_animal_proteins.txt`.

### SPEC-CATALOG-014 — Passive-cook floor

The catalog **shall** include ≥ 1 recipe such that:

- `cooking_methods` contains `slow-cook` OR `pressure-cook`, AND
- `passive_prep_minutes > 0`.

---

## 6.4 Ingredient-name discipline

### SPEC-CATALOG-015 — Verbatim reuse

When two recipes refer to the same ingredient, the `name` string **shall**
be identical character-for-character (including spacing and parenthetical
qualifiers, if any).

### SPEC-CATALOG-016 — Quantity placement

Ingredient `name` **shall not** encode quantity, unit, or quantity-bearing
qualifiers (e.g. forbidden: `2 cloves garlic`, `garlic clove`,
`a pinch of salt`). All such information **shall** live in the
`quantity`/`unit` fields of the same `ingredients[]` entry.

### SPEC-CATALOG-017 — Casing

Ingredient `name` **shall** be lowercase, except where a proper noun is
intrinsic to the ingredient identity (e.g. `Dijon mustard`, `Parmigiano
Reggiano`). Capitalized brand names **should not** appear unless the
brand is functionally inseparable from the ingredient.

### SPEC-CATALOG-018 — No aliases here

Ingredient aliases (`boeuf` ↔ `beef`, `aubergine` ↔ `eggplant`)
**shall not** appear in story 06. Aliases are introduced by story 09
(search) and live in the `ingredients.aliases` column populated then.

---

## 6.5 Editorial tone

### SPEC-CATALOG-019 — Concept shape

`concept` **shall** be a single paragraph (one or two sentences) that
states the dish in plain terms and the single most important technique
decision (e.g. "sear cold-pan to render fat", "deglaze with cream after
the fond builds").

### SPEC-CATALOG-020 — Voice

`concept` and `steps` **shall** use a terse, expert tone. Forbidden
patterns include backstory ("growing up, my grandmother…"), suspense
("the secret is…"), marketing copy ("the best thing you'll eat all
week"), and superlatives without operational meaning.

### SPEC-CATALOG-021 — Step content

Each entry in `steps` **shall** describe what to do; it **should** add a
brief why-clause when the reason is non-obvious to a competent home
cook. Steps **shall not** describe trivial mechanics ("turn on the
stove", "open the bag").

### SPEC-CATALOG-022 — Language

All free-text fields (`title`, `concept`, `steps`, `tags`,
`flavor_profile`, ingredient `name`) **shall** be in English. Localized
catalogs are deferred.

### SPEC-CATALOG-023 — No media

Recipe files **shall not** contain image URLs, photo references, video
links, or any field outside the SPEC-RECIPES schema.

---

## 6.6 Review process

### SPEC-CATALOG-024 — Per-recipe reviewer comment

The PR introducing each recipe **shall** receive at least one inline
comment from a reviewer on each recipe file confirming a full read. A
single approving review without per-file comments does NOT satisfy
this requirement.

### SPEC-CATALOG-025 — PR description summary

The PR description **shall** include a Markdown table with one row per
recipe and the following columns:

| Column | Source |
|--------|--------|
| `slug` | YAML field |
| `time_minutes` | YAML field |
| `passive_prep_minutes` | YAML field |
| `cooking_methods` | comma-joined |
| `primary_protein` | per §5.5 (or `vegetarian`) |

The table **shall** be sortable by `cooking_methods` so the diversity
spread is visible at a glance.

---

## 6.7 End-to-end seed verification

### SPEC-CATALOG-026 — `make seed` succeeds

Against a freshly migrated database (i.e. `make up && make migrate` then
`make seed`), `make seed` **shall** exit 0 and the loader **shall** report
loaded counts whose recipe count equals the number of catalog files
(SPEC-CATALOG-001).

---

## 6.8 Automated diversity guard

### SPEC-CATALOG-027 — Diversity test

A repository test **shall** load every catalog file and assert all
diversity invariants of §6.3 (SPEC-CATALOG-011, 012, 013, 014). The test
**shall** fail with a clear message identifying which invariant was
violated and listing the offending catalog state. The test name is
specified in §9.
