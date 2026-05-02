# SPEC-RECIPES — §3 Goals, Non-Goals, Constraints

[← Index](SPEC-RECIPES-00-index.md)

## 3.1 Goals

1. Provide a pure, I/O-free `domain.Recipe` model that every downstream
   feature (discover, search, reactions) consumes.
2. Provide canonical taxonomies and validators with a single source of truth,
   matching `docs/architecture/data-model.md` exactly.
3. Provide a YAML loader that produces actionable error reports — every bad
   file/line surfaced in one run, not one at a time.
4. Provide an idempotent transactional loader-to-Postgres pipeline so seeding
   is safe to repeat in any environment.
5. Provide an ergonomic `cooksense-server seed` subcommand wired through
   `make seed`.
6. Ship at least one valid sample recipe so subsequent stories and tests have
   real input data.

## 3.2 Non-goals

- Authoring 15+ curated recipes — owned by Story 06.
- HTTP handlers / API surface — Stories 07–10.
- Loading lessons (Cooking School) — Story 10.
- Image/photo handling — post-MVP.
- Recipe versioning, soft deletes, or audit trail — post-MVP.
- Hot-reload of recipes at runtime — out of scope.

## 3.3 Hard constraints

- Go 1.26.2 only; modern idioms required.
- `internal/domain` SHALL NOT import any I/O package.
- YAML library SHALL be `gopkg.in/yaml.v3` (only).
- All DB writes SHALL go through a single transaction in `seed.Store`.
- Slug regex SHALL be exactly `^[a-z0-9]+(-[a-z0-9]+)*$` (kebab-case).
- All SQL `ON CONFLICT` clauses SHALL match the unique constraints declared
  in `migrations/0001_init.up.sql` (per SPEC-DB).

## 3.4 Soft constraints

- Loader SHOULD report YAML line numbers when `yaml.v3` provides them.
- Store SHOULD prefer batched `INSERT … ON CONFLICT` over `pgx.CopyFrom` for
  the MVP catalog size; switching to `CopyFrom` is allowed when the catalog
  exceeds ~500 recipes.
- Subcommand stdout SHOULD be machine-parseable (one line on success).
