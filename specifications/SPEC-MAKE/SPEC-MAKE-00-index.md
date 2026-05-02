# SPEC-MAKE — Makefile Targets — Index

> **Spec-Driven Development (SDD) — Story 02**
>
> **Story:** 02 — Makefile targets wrapping `docker compose`, `modd`, and the Go toolchain  
> **Status:** Draft  
> **Version:** 1.0.0  
> **Authors:** CookSense Engineering  
> **License:** Proprietary — CookSense

> [!IMPORTANT]
> **The Contract** — This spec IS the source of truth. Code that contradicts
> the spec is a bug. A spec that contradicts reality must be updated first,
> then the code changed to match.

> [!TIP]
> **RFC-2119 Language**
> - **"shall" / "must"** = mandatory requirement (a verification check **must** confirm it)
> - **"should"** = strong recommendation (deviation needs justification)
> - **"may"** = optional

---

## Files in this specification

| File | Sections covered |
|------|-----------------|
| [SPEC-MAKE-00-index.md](SPEC-MAKE-00-index.md) | This index — metadata, SPEC-ID registry |
| [SPEC-MAKE-01-preamble.md](SPEC-MAKE-01-preamble.md) | §0 AI Steering Preamble |
| [SPEC-MAKE-02-introduction.md](SPEC-MAKE-02-introduction.md) | §1 Introduction |
| [SPEC-MAKE-03-goals.md](SPEC-MAKE-03-goals.md) | §2 Goals & Non-Goals |
| [SPEC-MAKE-04-context.md](SPEC-MAKE-04-context.md) | §3 System Context & Dependencies |
| [SPEC-MAKE-05-architecture.md](SPEC-MAKE-05-architecture.md) | §4 Architecture Overview |
| [SPEC-MAKE-06-artifacts.md](SPEC-MAKE-06-artifacts.md) | §5 Artifact Specifications (SPEC-MAKE-001 – 017) |
| [SPEC-MAKE-07-configuration.md](SPEC-MAKE-07-configuration.md) | §6 Configuration Specification |
| [SPEC-MAKE-08-build.md](SPEC-MAKE-08-build.md) | §7 Build, Tooling & Quality Specification |
| [SPEC-MAKE-09-testing.md](SPEC-MAKE-09-testing.md) | §8 Testing & Verification Specification |
| [SPEC-MAKE-10-documentation.md](SPEC-MAKE-10-documentation.md) | §9 Documentation Specification |
| [SPEC-MAKE-A-checklist.md](SPEC-MAKE-A-checklist.md) | Appendix A — Specification Checklist |
| [SPEC-MAKE-B-tasks.md](SPEC-MAKE-B-tasks.md) | Appendix B — Implementation Task Decomposition |

---

## SPEC-ID Registry

| SPEC-ID | Requirement Summary |
|---------|---------------------|
| SPEC-MAKE-001 | `Makefile` exists at repository root |
| SPEC-MAKE-002 | `help` is the default target; parses `##` comments |
| SPEC-MAKE-003 | Public target list (exactly 10 targets) |
| SPEC-MAKE-004 | `up` invokes `docker compose up -d` |
| SPEC-MAKE-005 | `up` waits until Postgres is healthy (≤ 30 s timeout) |
| SPEC-MAKE-006 | `down` does not drop the data volume |
| SPEC-MAKE-007 | `build` produces `bin/cooksense-server` |
| SPEC-MAKE-008 | `run` invokes `modd`; declares `up migrate` prerequisites |
| SPEC-MAKE-009 | `test` runs `go test ./...` |
| SPEC-MAKE-010 | `lint` runs `go vet`; optionally `golangci-lint` |
| SPEC-MAKE-011 | `clean` removes `bin/`; `CLEAN_VOLUMES=1` drops volume |
| SPEC-MAKE-012 | `migrate` and `seed` delegate to `cooksense-server` via `go run` |
| SPEC-MAKE-013 | Every non-file target declared `.PHONY` |
| SPEC-MAKE-014 | `.env` included if present, ignored if absent |
| SPEC-MAKE-015 | `Makefile` ≤ 80 lines; overflow goes to `scripts/` |
| SPEC-MAKE-016 | `.env.example` mirrors `infra.md` variable table |
| SPEC-MAKE-017 | `.env.example` contains no real secrets |
| SPEC-MAKE-018 | `make help` exits 0 with all targets listed |
| SPEC-MAKE-019 | `make build` exits 0, `bin/cooksense-server` is executable |
| SPEC-MAKE-020 | `make lint` exits 0 on a clean tree |
