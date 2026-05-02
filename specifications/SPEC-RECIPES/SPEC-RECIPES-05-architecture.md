# SPEC-RECIPES — §5 Architecture

[← Index](SPEC-RECIPES-00-index.md)

## 5.1 End-to-end pipeline

```
seed/recipes/*.yaml
        │
        ▼  (filepath.Glob "*.yaml")
   read each file ──► yaml.v3.Unmarshal ──► domain.Recipe (in-memory)
        │                                       │
        │                                       ▼
        │                                  Recipe.Validate()
        │                                       │
        ▼                                       ▼
  per-file errors  ◄─── aggregate ─────────  taxonomy validators
        │                                       │
        ▼                                       ▼
   seed.LoadError  (errors.Join)       []domain.Recipe (valid set)
                                                │
                                                ▼
                                         seed.Store(ctx, pool, recipes)
                                                │
                                                ▼
                                         BEGIN TRANSACTION
                                          ├── UPSERT recipes BY slug
                                          ├── UPSERT ingredients BY name
                                          ├── DELETE recipe_ingredients WHERE recipe_id IN (...)
                                          └── INSERT recipe_ingredients (...)
                                         COMMIT  (or ROLLBACK on any error)
```

## 5.2 Dependency graph

```
cmd/cooksense-server ──► internal/seed ──► internal/domain
                              │
                              └────────► internal/db (pool only)
                              │
                              └────────► gopkg.in/yaml.v3
```

The arrow direction is enforced by Clean Architecture: `domain` is a leaf
node and depends on nothing inside the project.

## 5.3 Error aggregation strategy

- The loader iterates every `*.yaml` file in the directory.
- For each file, it accumulates a `[]error` of issues:
  - Filesystem read error (1 entry, terminates that file).
  - YAML parse error (1 entry, terminates that file; line number when
    available via `yaml.TypeError`).
  - Per-recipe validation errors (multiple entries possible).
- After all files are processed, slug uniqueness is checked across the whole
  set; duplicates produce one error per duplicate group.
- The combined slice is wrapped in `seed.LoadError` and returned via
  `errors.Join` semantics. Callers can `errors.Unwrap() []error` to inspect.

## 5.4 Transaction shape

The `seed.Store` function follows this exact lifecycle:

1. `tx, err := pool.BeginTx(ctx, pgx.TxOptions{})` — fail closes nothing.
2. `defer tx.Rollback(ctx)` — best-effort; commit will short-circuit it.
3. UPSERT each recipe (returns the `id`).
4. UPSERT each distinct ingredient (returns the `id`).
5. For each recipe: `DELETE FROM recipe_ingredients WHERE recipe_id = $1`
   then `INSERT` the new rows.
6. `tx.Commit(ctx)` — on success.

If any step returns a non-nil error, the deferred rollback executes and
`Store` returns `(0, 0, err)`. No row from the input batch is persisted —
this guarantees the "no partial loads" rule from the story.

## 5.5 Idempotency proof sketch

- Recipes: UPSERT on `slug` → re-running with identical input updates the
  row in place; no new `id` is allocated for an existing slug.
- Ingredients: UPSERT on `name` → re-running creates no duplicates;
  `aliases` is overwritten with the same value (no-op).
- `recipe_ingredients`: DELETE+INSERT inside the tx → the final row set for
  each affected recipe is exactly what the YAML declares. No orphans accrue
  because `recipe_ingredients` is the only table referencing a join row.
- Ingredients no longer used by any recipe are NOT deleted (out of scope —
  catalog grows monotonically in MVP).

## 5.6 Subcommand dispatch

`cmd/cooksense-server/main.go` reuses the SPEC-DB-024 dispatch pattern:

```
switch os.Args[1] {
case "migrate": …            // SPEC-DB-024
case "seed":    runSeed(...) // SPEC-RECIPES-035
default:        runServer(...)
}
```

`runSeed` SHALL: (a) load config, (b) open the pgx pool, (c) call
`seed.Load`, (d) on success call `seed.Store`, (e) print the success line,
(f) close the pool, (g) exit 0. Any error path prints to stderr and exits
non-zero.
