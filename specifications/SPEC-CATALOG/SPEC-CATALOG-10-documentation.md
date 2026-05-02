# SPEC-CATALOG §10 — Documentation

[← back to index](SPEC-CATALOG-00-index.md)

## 10.1 README impact

The repo `README.md` **shall** gain a "Catalog" subsection describing:

- Where catalog files live (`seed/recipes/`).
- How to add a new recipe (see §10.3 PR template).
- How to verify the catalog locally (`make seed`, `go test ./internal/seed/...`).
- A pointer to this spec for the full rules.

The subsection **shall** be ≤ 30 lines; full rules belong here, not in
the README.

## 10.2 Doc comments

No new exported Go symbols are introduced. Existing `internal/seed`
package-level doc comment (SPEC-RECIPES) **may** be extended with one
sentence pointing readers to SPEC-CATALOG for content rules:

```go
// Package seed loads YAML recipe files into Postgres. Schema and loader
// behavior are governed by SPEC-RECIPES; catalog content rules
// (diversity, naming, tone) are governed by SPEC-CATALOG.
```

## 10.3 PR template addition (recommended)

`.github/PULL_REQUEST_TEMPLATE/recipe.md` **may** be added to scaffold the
SPEC-CATALOG-025 summary table:

```markdown
## Summary

<!-- 1–2 sentences: what's added/changed in the catalog. -->

## Catalog table

| slug | time_minutes | passive_prep_minutes | cooking_methods | primary_protein |
|------|--------------|----------------------|-----------------|-----------------|
| ...  | ...          | ...                  | ...             | ...             |

## Verification

- [ ] `go test ./internal/seed/...` green (catalog diversity test passes).
- [ ] `make up && make migrate && make seed` succeeds; output pasted below.
- [ ] At least one reviewer comment on each new recipe file (SPEC-CATALOG-024).

```
seed output:
loaded 16 recipes, 87 ingredients
```

## SPEC traceability

Closes / advances: SPEC-CATALOG-NNN, ...
```

## 10.4 Reviewer checklist

The PR template **shall** include, or the spec page **shall** publish, a
reviewer checklist:

- [ ] Filename equals `<slug>.yaml`.
- [ ] All fields present; taxonomies respected.
- [ ] `time_minutes ≤ 30`. Passive time set if total > 30 min.
- [ ] Tone: terse, expert. No backstory, no marketing.
- [ ] Steps explain *why* when non-obvious.
- [ ] Ingredient names match canonical spelling already used in the catalog.
- [ ] No images, URLs, or non-schema fields.
- [ ] Reviewer left at least one substantive comment on the file.

## 10.5 Story documentation status

`docs/stories/06-curated-recipes-content.md` **shall** be updated to
reference this spec under "Background" once SPEC-CATALOG is merged:

> Implementation governed by [SPEC-CATALOG](../../specifications/SPEC-CATALOG/SPEC-CATALOG-00-index.md).
