# SPEC-RECIPES — §6 Package Specifications (Normative)

[← Index](SPEC-RECIPES-00-index.md)

This section is normative. Every requirement uses RFC-2119 keywords.

---

## §6.1 `internal/domain/recipe.go`

### SPEC-RECIPES-001 — Purity
The `internal/domain` package SHALL NOT import any of: `os`, `io`, `io/fs`,
`net/http`, `database/sql`, `github.com/jackc/pgx/...`, `gopkg.in/yaml.v3`,
or any package outside the Go standard library other than `cmp` / `errors` /
`fmt` / `regexp` / `slices` / `strings`.

### SPEC-RECIPES-002 — `Recipe` struct
The package SHALL export the following struct, with YAML tags matching the
schema in `docs/architecture/data-model.md`:

```go
type Recipe struct {
    Slug                string              `yaml:"slug"`
    Title               string              `yaml:"title"`
    Concept             string              `yaml:"concept"`
    TimeMinutes         int                 `yaml:"time_minutes"`
    PassivePrepMinutes  int                 `yaml:"passive_prep_minutes"`
    CookingMethods      []string            `yaml:"cooking_methods"`
    Tags                []string            `yaml:"tags"`
    FlavorProfile       []string            `yaml:"flavor_profile"`
    Ingredients         []RecipeIngredient  `yaml:"ingredients"`
    Steps               []string            `yaml:"steps"`
}
```

### SPEC-RECIPES-003 — `RecipeIngredient` struct
```go
type RecipeIngredient struct {
    Name     string  `yaml:"name"`
    Category string  `yaml:"category"`
    Quantity float64 `yaml:"quantity"`
    Unit     string  `yaml:"unit"`
    Optional bool    `yaml:"optional"`
}
```

### SPEC-RECIPES-004 — Slug pattern
The package SHALL export `SlugPattern = "^[a-z0-9]+(-[a-z0-9]+)*$"` as a
package-level `const`. A `*regexp.Regexp` compiled from it SHALL be cached
package-private.

### SPEC-RECIPES-005 — `ValidateSlug`
```go
func ValidateSlug(s string) error
```
Returns a non-nil error of the form `"invalid slug %q: must match %s"` for
any input not matching `SlugPattern`. Returns `nil` for valid slugs.

### SPEC-RECIPES-006 — `Recipe.Validate`
```go
func (r Recipe) Validate() error
```
SHALL enforce, in this order, returning aggregated errors via `errors.Join`:
1. `ValidateSlug(r.Slug)` — slug shape.
2. `r.Title != ""`, `r.Concept != ""` — required strings.
3. `r.TimeMinutes > 0`.
4. `r.PassivePrepMinutes >= 0` (SPEC-RECIPES-007).
5. `len(r.Ingredients) >= 2`.
6. `len(r.Steps) >= 1` (SPEC-RECIPES-008).
7. At least one of `len(r.Tags) > 0` or `len(r.CookingMethods) > 0`.
8. For each `m` in `r.CookingMethods`: `ValidateCookingMethod(m)`.
9. For each `f` in `r.FlavorProfile`: `ValidateFlavorProfile(f)`.
10. For each ingredient `i`: `i.Name != ""` AND `ValidateCategory(i.Category)`.

Returns `nil` only if every check passes.

### SPEC-RECIPES-007 — Passive prep non-negative
Enforced by SPEC-RECIPES-006 step 4.

### SPEC-RECIPES-008 — Steps non-empty
Enforced by SPEC-RECIPES-006 step 6. The element type SHALL remain
`[]string` (not a struct) for MVP.

### SPEC-RECIPES-009 — String-slice taxonomy fields
`CookingMethods`, `Tags`, `FlavorProfile` SHALL be `[]string`. Tag values
are free-form (no taxonomy) for MVP; cooking methods and flavor profiles
SHALL be validated against §6.2.

---

## §6.2 `internal/domain/taxonomy.go`

### SPEC-RECIPES-010 — `CookingMethods` set
```go
var CookingMethods = []string{
    "pan-sear", "roast", "braise", "boil", "steam",
    "grill", "slow-cook", "pressure-cook", "raw", "bake",
}
```
The order SHALL be stable to enable test snapshotting.

### SPEC-RECIPES-011 — `FlavorProfiles` set
```go
var FlavorProfiles = []string{
    "acid", "fat", "salt", "sweet", "bitter", "umami", "spicy", "herbaceous",
}
```

### SPEC-RECIPES-012 — `IngredientCategories` set
```go
var IngredientCategories = []string{
    "protein", "vegetable", "fruit", "starch", "dairy",
    "fat", "spice", "herb", "condiment", "liquid", "other",
}
```

### SPEC-RECIPES-013 — `ValidateCookingMethod`
```go
func ValidateCookingMethod(s string) error
```
Returns non-nil iff `slices.Contains(CookingMethods, s) == false`. Error
message: `"invalid cooking_method %q: must be one of %v"`.

### SPEC-RECIPES-014 — `ValidateFlavorProfile`
Same shape as SPEC-RECIPES-013, against `FlavorProfiles`.

### SPEC-RECIPES-015 — `ValidateCategory`
Same shape as SPEC-RECIPES-013, against `IngredientCategories`.

### SPEC-RECIPES-016 — Case sensitivity
All three validators SHALL be case-sensitive. Inputs like `"Pan-sear"` or
`"ACID"` SHALL be rejected. Normalization (lowercasing) is the YAML
author's responsibility.

### SPEC-RECIPES-017 — Read-only iteration
The exported variables MAY be `[]string` for simplicity, but callers SHALL
NOT mutate them. The package SHALL document this in the `taxonomy.go` doc
comment. (Future tightening to an iterator-only API is allowed.)

---

## §6.3 `internal/seed/recipes.go`

### SPEC-RECIPES-018 — `Load` signature
```go
func Load(ctx context.Context, dir string) ([]domain.Recipe, error)
```

### SPEC-RECIPES-019 — YAML library
Parsing SHALL use `gopkg.in/yaml.v3`. No other YAML library is permitted.

### SPEC-RECIPES-020 — Per-recipe validation
For each successfully parsed recipe, `Load` SHALL call `recipe.Validate()`
and accumulate any returned error into the aggregate.

### SPEC-RECIPES-021 — Aggregate errors, never fail-fast
`Load` SHALL process every `*.yaml` file in `dir` before returning. A bad
file SHALL NOT abort iteration. The function SHALL return `(nil, LoadError)`
when any error occurred, and `(recipes, nil)` only when all files validated.

### SPEC-RECIPES-022 — File and line in error messages
Each error SHALL be prefixed with `<filename>:` and, when available from
`yaml.TypeError` or per-node `*yaml.Node.Line`, include `:<line>:` before
the message.

### SPEC-RECIPES-023 — Slug uniqueness
After per-file processing, `Load` SHALL check that every `Slug` in the
parsed set is unique. Each duplicate SHALL be reported as
`"duplicate slug %q in <file-a> and <file-b>"`.

### SPEC-RECIPES-024 — `LoadError` type
```go
type LoadError struct {
    Errs []error
}
func (e *LoadError) Error() string { /* concatenated, one per line */ }
func (e *LoadError) Unwrap() []error { return e.Errs }
```
`Unwrap` returning `[]error` SHALL be compatible with Go 1.20+
multi-error semantics (`errors.Is`/`errors.As` traversal).

### SPEC-RECIPES-025 — Empty directory
If `filepath.Glob(filepath.Join(dir, "*.yaml"))` returns no matches and no
filesystem error, `Load` SHALL return `(nil, nil)`.

---

## §6.4 `internal/seed/store.go`

### SPEC-RECIPES-026 — `Store` signature
```go
func Store(
    ctx context.Context,
    pool *pgxpool.Pool,
    recipes []domain.Recipe,
) (loaded int, ingredients int, err error)
```

### SPEC-RECIPES-027 — Single transaction
All writes SHALL execute inside one `pgx.Tx` obtained from `pool.BeginTx`.
The function SHALL `defer tx.Rollback(ctx)` and `Commit` only on the
happy path.

### SPEC-RECIPES-028 — Recipe UPSERT
The recipe upsert SHALL use the following shape (placeholders normalized):

```sql
INSERT INTO recipes
  (slug, title, concept, time_minutes, passive_prep_minutes,
   cooking_methods, tags, flavor_profile, steps)
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (slug) DO UPDATE SET
  title                = EXCLUDED.title,
  concept              = EXCLUDED.concept,
  time_minutes         = EXCLUDED.time_minutes,
  passive_prep_minutes = EXCLUDED.passive_prep_minutes,
  cooking_methods      = EXCLUDED.cooking_methods,
  tags                 = EXCLUDED.tags,
  flavor_profile       = EXCLUDED.flavor_profile,
  steps                = EXCLUDED.steps
RETURNING id;
```

`steps` SHALL be marshalled to JSONB (the column type per SPEC-DB DDL).

### SPEC-RECIPES-029 — Ingredient UPSERT
```sql
INSERT INTO ingredients (name, category, aliases)
VALUES ($1, $2, $3)
ON CONFLICT (name) DO UPDATE SET
  category = EXCLUDED.category,
  aliases  = EXCLUDED.aliases
RETURNING id;
```
`aliases` SHALL default to `'{}'` for MVP (no alias loader yet).

### SPEC-RECIPES-030 — Replace `recipe_ingredients`
For every recipe the function SHALL execute:
```sql
DELETE FROM recipe_ingredients WHERE recipe_id = $1;
```
followed by INSERTs for each `RecipeIngredient`:
```sql
INSERT INTO recipe_ingredients
  (recipe_id, ingredient_id, quantity, unit, optional)
VALUES ($1, $2, $3, $4, $5);
```
DELETE+INSERT SHALL happen inside the same transaction as the parent
recipe upsert.

### SPEC-RECIPES-031 — Atomic rollback
If any DB call returns a non-nil error, the function SHALL return
`(0, 0, fmt.Errorf("seed store: %w", err))`. The transaction SHALL be
rolled back via the deferred call. No row from the batch persists.

### SPEC-RECIPES-032 — Idempotency
Calling `Store` twice with the same `recipes` slice (against the same DB)
SHALL produce:
- Identical row counts in `recipes`, `ingredients`, and `recipe_ingredients`.
- No new `id` allocations for unchanged slugs/names (UPSERT semantics).
- No orphan ingredient rows beyond those created by the first call.

### SPEC-RECIPES-033 — Return counts
- `loaded` SHALL equal `len(recipes)` on success.
- `ingredients` SHALL equal the count of distinct `RecipeIngredient.Name`
  values across the input batch (case-sensitive).

### SPEC-RECIPES-034 — Constructor injection
The pool SHALL be passed by argument. The package SHALL NOT hold a
package-level `*pgxpool.Pool`. (Aligns with SPEC-DB-018.)

---

## §6.5 `cmd/cooksense-server` — `seed` subcommand

### SPEC-RECIPES-035 — Subcommand exists
Invoking `cooksense-server seed` SHALL dispatch to a `runSeed(args []string)`
function in `main.go` (or an internal helper).

### SPEC-RECIPES-036 — Default directory + `--dir` flag
The subcommand SHALL parse a `-dir` (or `--dir`) flag via `flag.NewFlagSet`,
defaulting to `"seed/recipes"`.

### SPEC-RECIPES-037 — Success output
On success the subcommand SHALL print exactly one line to stdout:
```
loaded N recipes, M ingredients
```
where `N` = `loaded` and `M` = `ingredients` from `Store`.

### SPEC-RECIPES-038 — Error exit
Any non-nil error from `Load` or `Store` SHALL be printed to stderr (full
multi-error contents in the case of `LoadError`) and the process SHALL
exit with status `1`.

### SPEC-RECIPES-039 — Pool lifecycle
The subcommand SHALL: (a) call `config.Load()`, (b) call `db.Open(ctx, cfg)`,
(c) `defer pool.Close()`, (d) call `seed.Load`, (e) call `seed.Store`,
(f) print the success line, (g) return exit code 0.

---

## §6.6 SQL prepared statements

Implementations MAY use `pgx.Batch` to send the per-recipe operations as a
batch. Implementations MAY use `pgx.CopyFrom` for `recipe_ingredients`
INSERTs once the catalog exceeds 500 recipes; the SPEC remains satisfied so
long as the resulting end-state matches §5.5.
