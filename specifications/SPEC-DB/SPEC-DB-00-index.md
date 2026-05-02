# SPEC-DB — Database: pgx Pool, Migrations & Initial Schema — Index

> **Spec-Driven Development (SDD) — Story 03**
>
> **Story:** 03 — Database: pgx pool + migrations + 0001_init  
> **Status:** Final  
> **Version:** 1.0.0  
> **Authors:** CookSense Engineering  
> **License:** Proprietary — CookSense

> [!IMPORTANT]
> **The Contract** — This spec IS the source of truth. Code that contradicts
> the spec is a bug. A spec that contradicts reality must be updated first,
> then the code changed to match.

> [!TIP]
> **RFC-2119 Language**
> - **"shall" / "must"** = mandatory requirement (test must verify it)
> - **"should"** = strong recommendation (deviation needs justification)
> - **"may"** = optional

---

## Files in this specification

| File | Sections covered |
|------|-----------------|
| [SPEC-DB-00-index.md](SPEC-DB-00-index.md) | This index — metadata, SPEC-ID registry |
| [SPEC-DB-01-preamble.md](SPEC-DB-01-preamble.md) | §0 AI Steering Preamble |
| [SPEC-DB-02-introduction.md](SPEC-DB-02-introduction.md) | §1 Introduction |
| [SPEC-DB-03-goals.md](SPEC-DB-03-goals.md) | §2 Goals & Non-Goals |
| [SPEC-DB-04-context.md](SPEC-DB-04-context.md) | §3 System Context & Dependencies |
| [SPEC-DB-05-architecture.md](SPEC-DB-05-architecture.md) | §4 Architecture Overview |
| [SPEC-DB-06-packages.md](SPEC-DB-06-packages.md) | §5 Package Specifications (SPEC-DB-001 – 027) |
| [SPEC-DB-07-configuration.md](SPEC-DB-07-configuration.md) | §6 Configuration Specification |
| [SPEC-DB-08-build.md](SPEC-DB-08-build.md) | §7 Build, Tooling & Quality Specification |
| [SPEC-DB-09-testing.md](SPEC-DB-09-testing.md) | §8 Testing Specification |
| [SPEC-DB-10-documentation.md](SPEC-DB-10-documentation.md) | §9 Documentation Specification |
| [SPEC-DB-A-checklist.md](SPEC-DB-A-checklist.md) | Appendix A — Specification Checklist |
| [SPEC-DB-B-tasks.md](SPEC-DB-B-tasks.md) | Appendix B — Implementation Task Decomposition |

---

## SPEC-ID Registry

| SPEC-ID | Requirement Summary |
|---------|---------------------|
| SPEC-DB-001 | `config.Load()` reads all required env vars |
| SPEC-DB-002 | `config.Load()` fails fast on missing required vars |
| SPEC-DB-003 | `Config` struct is flat, typed, and immutable after construction |
| SPEC-DB-004 | `db.Open` accepts `ctx` and `cfg` |
| SPEC-DB-005 | Pool min 2, max 10, health check period 1 min |
| SPEC-DB-006 | `db.Open` pings the DB and returns error if unreachable |
| SPEC-DB-007 | `db.Open` accepts `pgxpool.Config` override for testing |
| SPEC-DB-008 | `db.Up` applies all pending migrations |
| SPEC-DB-009 | `db.Up` is idempotent (running twice is a no-op) |
| SPEC-DB-010 | `db.Down` rolls back N steps |
| SPEC-DB-011 | `MigrationError` type wraps version number |
| SPEC-DB-012 | `cooksense-server migrate up` subcommand |
| SPEC-DB-013 | `cooksense-server migrate down [N]` subcommand |
| SPEC-DB-014 | `reaction_kind` enum |
| SPEC-DB-015 | `users` table DDL |
| SPEC-DB-016 | `recipes` table DDL + GIN indexes |
| SPEC-DB-017 | `ingredients` table DDL + GIN index |
| SPEC-DB-018 | `recipe_ingredients` join table DDL + index |
| SPEC-DB-019 | `user_reactions` table DDL + index |
| SPEC-DB-020 | `lesson_articles` table DDL |
| SPEC-DB-021 | `0001_init.down.sql` drops all objects in reverse order |
| SPEC-DB-022 | Pool passed via constructor injection; no globals |
| SPEC-DB-023 | `internal/config` isolation from domain packages |
| SPEC-DB-024 | `internal/db` isolation from feature packages |
| SPEC-DB-025 | `go.mod` declares 4 new external dependencies |
| SPEC-DB-026 | `go build ./...` passes after adding dependencies |
| SPEC-DB-027 | `go vet ./...` produces zero violations |
