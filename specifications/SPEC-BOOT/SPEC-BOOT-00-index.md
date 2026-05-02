# SPEC-BOOT — Bootstrap Project Structure — Index

> **Spec-Driven Development (SDD) — Story 01**
>
> **Story:** 01 — Bootstrap project structure  
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
| [SPEC-BOOT-00-index.md](SPEC-BOOT-00-index.md) | This index — metadata, SPEC-ID registry |
| [SPEC-BOOT-01-preamble.md](SPEC-BOOT-01-preamble.md) | §0 AI Steering Preamble |
| [SPEC-BOOT-02-introduction.md](SPEC-BOOT-02-introduction.md) | §1 Introduction |
| [SPEC-BOOT-03-goals.md](SPEC-BOOT-03-goals.md) | §2 Goals & Non-Goals |
| [SPEC-BOOT-04-context.md](SPEC-BOOT-04-context.md) | §3 System Context & Dependencies |
| [SPEC-BOOT-05-architecture.md](SPEC-BOOT-05-architecture.md) | §4 Architecture Overview |
| [SPEC-BOOT-06-packages.md](SPEC-BOOT-06-packages.md) | §5 Package Specifications (SPEC-BOOT-001 – 018) |
| [SPEC-BOOT-07-configuration.md](SPEC-BOOT-07-configuration.md) | §6 Configuration Specification |
| [SPEC-BOOT-08-build.md](SPEC-BOOT-08-build.md) | §7 Build, Tooling & Quality Specification |
| [SPEC-BOOT-09-testing.md](SPEC-BOOT-09-testing.md) | §8 Testing Specification |
| [SPEC-BOOT-10-documentation.md](SPEC-BOOT-10-documentation.md) | §9 Documentation Specification |
| [SPEC-BOOT-A-checklist.md](SPEC-BOOT-A-checklist.md) | Appendix A — Specification Checklist |
| [SPEC-BOOT-B-tasks.md](SPEC-BOOT-B-tasks.md) | Appendix B — Implementation Task Decomposition |

---

## SPEC-ID Registry

| SPEC-ID | Requirement Summary |
|---------|---------------------|
| SPEC-BOOT-001 | `internal/config/doc.go` package declaration |
| SPEC-BOOT-002 | `internal/config` isolation (no internal imports) |
| SPEC-BOOT-003 | `internal/auth/doc.go` package declaration |
| SPEC-BOOT-004 | `internal/db/doc.go` package declaration |
| SPEC-BOOT-005 | `internal/domain/doc.go` package declaration |
| SPEC-BOOT-006 | `internal/domain` zero-import rule |
| SPEC-BOOT-007 | `internal/recipes/doc.go` package declaration |
| SPEC-BOOT-008 | `internal/reactions/doc.go` package declaration |
| SPEC-BOOT-009 | `internal/lessons/doc.go` package declaration |
| SPEC-BOOT-010 | `internal/users/doc.go` package declaration |
| SPEC-BOOT-011 | `internal/seed/doc.go` package declaration |
| SPEC-BOOT-012 | Placeholder `main` function (print + exit 0) |
| SPEC-BOOT-013 | `main.go` imports only `"fmt"` |
| SPEC-BOOT-014 | `migrations/` directory with `.gitkeep` |
| SPEC-BOOT-015 | `seed/recipes/` and `seed/lessons/` directories with `.gitkeep` |
| SPEC-BOOT-016 | `.gitignore` rules (`bin/`, `secrets/`, `.env`, `*.local.*`, etc.) |
| SPEC-BOOT-017 | `go.mod` declares `go 1.26.2` |
| SPEC-BOOT-018 | `go.mod` module path preserved |
| SPEC-BOOT-019 | `go build ./...` exits 0 |
| SPEC-BOOT-020 | `go vet ./...` exits 0 |
| SPEC-BOOT-021 | All packages compile |
| SPEC-BOOT-022 | `main` prints `cooksense-server starting` and exits 0 |
| SPEC-BOOT-023 | `go vet ./...` passes |
| SPEC-BOOT-024 | `go build ./...` + `go vet ./...` pass as Definition of Done |
