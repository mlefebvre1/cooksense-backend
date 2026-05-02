# SPEC-MAKE-08 — Build, Tooling & Quality Specification

> Part of [SPEC-MAKE](SPEC-MAKE-00-index.md) — Story 02: Makefile Targets  
> SPEC-IDs covered: SPEC-MAKE-018, SPEC-MAKE-019, SPEC-MAKE-020

---

## 7. Build, Tooling & Quality Specification

### 7.1 Verification Matrix

| SPEC-ID | Command | Expected outcome |
|---------|---------|-----------------|
| SPEC-MAKE-018 | `make help` | Exits `0`; all 10 target names and descriptions appear in stdout. |
| SPEC-MAKE-019 | `make build` | Exits `0`; `bin/cooksense-server` is created and is executable. |
| SPEC-MAKE-020 | `make lint` | Exits `0` on a clean tree (no Go vet violations). |

### 7.2 Toolchain Requirements

| Tool | Version | How verified |
|------|---------|-------------|
| GNU Make | ≥ 3.81 | `make --version` |
| Go | `1.26.2` | `go version` |
| docker compose (V2) | any | `docker compose version` |
| modd | any | `which modd` |
| golangci-lint | any (optional) | `which golangci-lint` |

### 7.3 Formatting

`Makefile` formatting **shall** use hard tabs for recipe indentation. The file **shall** not contain any trailing whitespace. A `diff --check` or `cat -A` check **should** confirm this before merging.
