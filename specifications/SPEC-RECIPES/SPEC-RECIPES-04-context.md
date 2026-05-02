# SPEC-RECIPES — §4 System Context

[← Index](SPEC-RECIPES-00-index.md)

## 4.1 External dependencies (new in this story)

| Dependency | Version | Purpose |
|------------|---------|---------|
| `gopkg.in/yaml.v3` | latest stable | Parse `seed/recipes/*.yaml` files |

No other new third-party imports. The store reuses `github.com/jackc/pgx/v5`
already introduced in Story 03 (SPEC-DB).

## 4.2 Reused dependencies

| Dependency | Source story | Used here for |
|------------|--------------|----------------|
| `pgx/v5` + `pgxpool` | SPEC-DB | Connection pool injected into `seed.Store` |
| `internal/config` | SPEC-DB | Reads `DATABASE_URL` for the seed subcommand |
| `internal/db.Open` | SPEC-DB | Opens the pool from the seed subcommand |
| `log/slog` | stdlib | Structured logging from the loader/store |

## 4.3 Package boundary map

```
cmd/cooksense-server     ── dispatches "seed" subcommand
        │
        ▼
internal/seed            ── Load (YAML→[]domain.Recipe), Store (DB upsert)
        │                  Imports: domain, db, pgxpool, yaml.v3
        ▼
internal/domain          ── Recipe, RecipeIngredient, taxonomies, validators
                            Imports: stdlib only (regexp, errors, fmt)
```

Reverse imports are forbidden. `internal/domain` SHALL never import
`internal/seed`, `internal/db`, or any infrastructure package.

## 4.4 Filesystem layout introduced

```
internal/
  domain/
    recipe.go         (new — Recipe, RecipeIngredient, ValidateSlug, …)
    taxonomy.go       (new — canonical sets + Validate*)
  seed/
    recipes.go        (new — Load, LoadError)
    store.go          (new — Store)
    doc.go            (updated package comment)
seed/
  recipes/
    _sample.yaml      (new — minimal valid recipe for tests)
```

## 4.5 Out-of-band collaborators

- `migrations/0001_init.up.sql` — DDL prerequisite (already specified in
  SPEC-DB). Loader assumes the schema exists.
- `Makefile`'s `seed` target (SPEC-MAKE-009) — invokes the subcommand.
