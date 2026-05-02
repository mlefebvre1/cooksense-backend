# SPEC-CATALOG §1 — Preamble

[← back to index](SPEC-CATALOG-00-index.md)

## 1.1 Authorship

This document was drafted by the `sdd-spec-author` skill on 2026-05-02 from
`docs/stories/06-curated-recipes-content.md`. It is a binding specification:
RFC-2119 `shall` / `must` / `should` / `may` keywords carry their formal meaning.

## 1.2 Scope boundary

SPEC-CATALOG is **content-centric**. It does NOT define:

- The recipe **schema** (owned by [SPEC-RECIPES §6.1–§6.2](../SPEC-RECIPES/SPEC-RECIPES-06-packages.md)).
- The **YAML loader** behavior, error aggregation, or DB upsert (owned by SPEC-RECIPES).
- The **DDL** for `recipes`, `ingredients`, `recipe_ingredients` (owned by [SPEC-DB](../SPEC-DB/SPEC-DB-00-index.md)).
- The **`seed` subcommand** dispatch (owned by SPEC-RECIPES §6.6).

SPEC-CATALOG **adds** rules about:

- Which recipes ship in the MVP catalog (count, diversity).
- How they are written (tone, naming conventions).
- How they are reviewed (peer review, PR table).
- How the catalog is verified end-to-end (`make seed`, diversity test).

If a SPEC-CATALOG rule contradicts a SPEC-RECIPES rule, the contradiction is a
spec bug — open an ADR before editing either file.

## 1.3 Traceability

Every SPEC-CATALOG-NNN ID has at least one named test in §9 or one verifiable
artifact in §6 (e.g. PR table, reviewer comment). The Appendix A matrix
demonstrates every story-06 acceptance criterion maps to ≥ 1 ID and every ID
maps to ≥ 1 verification artifact.

## 1.4 Authoring discipline

Recipe content is **immutable once merged** unless a follow-up PR explicitly
revises it with reviewer sign-off. Cosmetic edits (typos, capitalization)
follow the same review process as new recipes.
