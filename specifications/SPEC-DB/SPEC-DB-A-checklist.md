# SPEC-DB-A — Specification Checklist

> Part of [SPEC-DB](SPEC-DB-00-index.md) — Story 03: Database pgx Pool, Migrations & Initial Schema

---

## Appendix A — Specification Checklist

Use this checklist before starting implementation. Every box must be checked.

- [x] **Section 0** — AI Preamble reviewed; DB-specific anti-patterns (ORM, global pool) listed.
- [x] **Section 1** — Introduction written; 3 deliverables named; story relationships mapped.
- [x] **Section 2** — Every goal is testable; non-goals (ORM, replicas, embed.FS) explicit.
- [x] **Section 3** — 4 new external dependencies named with module paths; internal dependency rules stated.
- [x] **Section 4** — Pool lifecycle sequence diagram; subcommand dispatch flowchart; idempotency guarantee explained.
- [x] **Section 5** — All 27 SPEC-DB-NNN IDs defined (001–027); `config.Load` contract + `Config` struct fully specified; `db.Open`, `db.Up`, `db.Down`, `MigrationError` signatures given; full DDL for 6 tables + 5 indexes + 1 enum; down migration stated.
- [x] **Section 6** — Variable table complete; DSN format example given; sensitive-value masking rule stated.
- [x] **Section 7** — Build commands listed; dependency addition steps shown; linter checks called out.
- [x] **Section 8** — 8 named tests with SPEC-ID traceability; `testDB` helper pattern specified; coverage thresholds set.
- [x] **Section 9** — Go Doc standard with SPEC-ID citation; SQL header format; README deferred to Story 12.
