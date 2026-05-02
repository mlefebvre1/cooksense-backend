# SPEC-RECIPES — §7 Configuration

[← Index](SPEC-RECIPES-00-index.md)

## 7.1 Seed directory and sample recipe

### SPEC-RECIPES-040 — Sample recipe presence
The repository SHALL contain at least one valid sample recipe at
`seed/recipes/_sample.yaml` so that:
- Story 05 unit tests have a real fixture to exercise.
- `make seed` produces non-zero output on a fresh checkout before Story 06
  curated content lands.

The sample SHALL satisfy every domain invariant (≥ 2 ingredients, ≥ 1 step,
`time_minutes > 0`, valid slug, valid taxonomies).

### SPEC-RECIPES-041 — Underscore-prefix exclusion
Story 06 (curated content) SHALL ignore filenames starting with `_`. The
sample is therefore named `_sample.yaml` so it does not pollute the curated
catalog count for the demo. The loader itself does NOT special-case the
underscore — it loads every `*.yaml` and the curated count is enforced at
the content layer.

## 7.2 Environment variables

The seed subcommand reuses environment variables already defined in the
project:

| Variable | Owner | Purpose for SPEC-RECIPES |
|----------|-------|--------------------------|
| `DATABASE_URL` | SPEC-DB-007 | Connect to Postgres for `seed.Store` |
| `LOG_LEVEL`    | SPEC-DB-009 | Loader/store structured logs |
| `LOG_FORMAT`   | SPEC-DB-010 | Loader/store structured logs |

No new environment variables are introduced by this story. `.env.example`
SHALL NOT need updates beyond what SPEC-DB already defined.

## 7.3 Filesystem expectations at runtime

- The working directory at the time `cooksense-server seed` runs SHALL
  contain `seed/recipes/` (relative path), unless `--dir` overrides.
- The directory SHALL be readable by the process user.
- The Postgres schema SHALL already be migrated (Story 03). The subcommand
  does NOT trigger migrations on its own.

## 7.4 Default values summary

| Setting | Default | Override |
|---------|---------|----------|
| Seed directory | `seed/recipes` | `--dir <path>` flag |
| Glob pattern | `*.yaml` | not configurable |
| YAML parser | `gopkg.in/yaml.v3` | not configurable |
