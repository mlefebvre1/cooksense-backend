# SPEC-CATALOG §7 — Configuration

[← back to index](SPEC-CATALOG-00-index.md)

## 7.1 Environment variables

SPEC-CATALOG introduces **no new environment variables**. The catalog is
loaded by the `seed` subcommand which already consumes:

- `DATABASE_URL` (SPEC-DB)
- `LOG_LEVEL`, `LOG_FORMAT` (SPEC-BOOT)

No additional config keys are added.

## 7.2 Configuration files added by this spec

| Path | Purpose | Format |
|------|---------|--------|
| `seed/recipes/*.yaml` | Catalog content (≥ 15 files) | YAML 1.2 |
| `seed/recipes/_animal_proteins.txt` | Animal-origin canonical names, one per line, used by SPEC-CATALOG-027 diversity test | UTF-8, `#` comments allowed |

### `_animal_proteins.txt` initial contents

The file **shall** ship with the following baseline (additions allowed via
PR; removals require an ADR):

```
# Animal-origin protein canonical names. Used by the catalog diversity test
# (SPEC-CATALOG-013). One canonical ingredient name per line. Lines starting
# with `#` are comments. Casing must match the canonical ingredient `name`.
chicken thighs
chicken thighs (bone-in, skin-on)
chicken breast
beef chuck
beef ground
pork shoulder
pork ground
salmon
cod
shrimp
eggs
bacon
anchovy
tuna
```

## 7.3 Excluded fixtures

`seed/recipes/_sample.yaml` is owned by SPEC-RECIPES (loader test fixture).
SPEC-CATALOG **shall not** modify or delete it; it is referenced here only
to clarify exclusion (SPEC-CATALOG-004).
