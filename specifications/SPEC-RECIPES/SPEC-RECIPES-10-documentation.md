# SPEC-RECIPES — §10 Documentation

[← Index](SPEC-RECIPES-00-index.md)

## 10.1 Doc-comment requirements

Every exported symbol introduced by Story 05 SHALL carry a Go Doc comment
that:

1. Begins with the symbol name (`// Recipe is …`, `// Load reads …`).
2. References the SPEC-RECIPES-NNN ID(s) it implements.
3. Documents pre-conditions, post-conditions, and error behavior.

| File | Required doc-comment subjects | SPEC-IDs to cite |
|------|-------------------------------|------------------|
| `internal/domain/doc.go` | Package overview: pure domain model, no I/O | SPEC-RECIPES-001 |
| `internal/domain/recipe.go` | `Recipe`, `RecipeIngredient`, `SlugPattern`, `ValidateSlug`, `Recipe.Validate` | 002–009 |
| `internal/domain/taxonomy.go` | Package taxonomy intent (read-only sets); each `Validate*` function | 010–017 |
| `internal/seed/doc.go` | Package overview: YAML→domain→Postgres pipeline | 018–034 |
| `internal/seed/recipes.go` | `Load`, `LoadError` | 018–025 |
| `internal/seed/store.go` | `Store` | 026–034 |
| `cmd/cooksense-server/main.go` | The `runSeed` helper (or inlined dispatch comment) | 035–039 |

### SPEC-RECIPES-042 — Package comment for `internal/seed`
The `internal/seed/doc.go` file SHALL contain a package comment of the form:

```go
// Package seed loads recipe YAML files from disk, validates them against
// the domain rules, and upserts the result into Postgres in a single
// transaction.
//
// Implements: SPEC-RECIPES-018 through SPEC-RECIPES-039.
package seed
```

## 10.2 README impact

### SPEC-RECIPES-043 — README quickstart update
The top-level `README.md` SHALL gain (or update) a "Seeding the catalog"
subsection explaining:

```markdown
## Seeding the catalog

After running `make up && make migrate`, populate the recipe catalog with:

    make seed

This loads every `*.yaml` file under `seed/recipes/` into Postgres in a
single transaction. Re-running is safe — `seed` is idempotent.

To load from a different directory:

    go run ./cmd/cooksense-server seed --dir path/to/yaml
```

The full curated quickstart lands in Story 12; SPEC-RECIPES-043 only
requires the seed paragraph to be present and accurate.

## 10.3 Architecture docs

`docs/architecture/data-model.md` already contains the canonical YAML
schema and taxonomy lists. **No edits to that document are required by
this story** — the spec borrows from it verbatim. If a future change
diverges, the architecture doc SHALL be updated first per the repo's
"spec and code ship together" rule.

## 10.4 Decision log

No new ADR is required. D-0005 (curated YAML seed) already covers the
high-level decision. If implementation choices in §6 diverge from D-0005
(e.g., switching to TOML), a new ADR superseding D-0005 SHALL be added.

## 10.5 Sample recipe documentation

The `seed/recipes/_sample.yaml` file SHALL include a top-of-file comment:

```yaml
# Sample recipe used by Story 05 fixtures and `make seed` smoke tests.
# Curated catalog content lives in non-underscore files (Story 06).
# See specifications/SPEC-RECIPES/ for the schema contract.
```
