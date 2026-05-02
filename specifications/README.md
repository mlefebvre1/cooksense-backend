# Specifications

This directory contains all **Spec-Driven Development (SDD)** contracts for
the `cooksense-backend` project.  
Each sub-directory holds one **SPEC-\*** module — a multi-file specification
that governs a vertical slice of the system from business intent through
implementation tasks.

> **The Contract rule:** code that contradicts a spec is a bug.  
> A spec that contradicts reality must be updated *first*, then the code
> changed to match.

---

## RFC-2119 language (applies to every spec)

| Keyword | Meaning |
|---------|---------|
| **shall / must** | Mandatory — a test must verify it |
| **should** | Strong recommendation — deviation requires justification |
| **may** | Optional |

---

## Specifications index

### SPEC-BOOT — Bootstrap Project Structure

| | |
|-|-|
| **Story** | 01 — Bootstrap project structure |
| **Status** | Final |
| **Version** | 1.0.0 |
| **SPEC-IDs** | SPEC-BOOT-001 – SPEC-BOOT-024 |

| File | Contents |
|------|----------|
| [SPEC-BOOT-00-index.md](SPEC-BOOT/SPEC-BOOT-00-index.md) | Index — metadata, SPEC-ID registry |
| [SPEC-BOOT-01-preamble.md](SPEC-BOOT/SPEC-BOOT-01-preamble.md) | §0 AI Steering Preamble |
| [SPEC-BOOT-02-introduction.md](SPEC-BOOT/SPEC-BOOT-02-introduction.md) | §1 Introduction |
| [SPEC-BOOT-03-goals.md](SPEC-BOOT/SPEC-BOOT-03-goals.md) | §2 Goals & Non-Goals |
| [SPEC-BOOT-04-context.md](SPEC-BOOT/SPEC-BOOT-04-context.md) | §3 System Context & Dependencies |
| [SPEC-BOOT-05-architecture.md](SPEC-BOOT/SPEC-BOOT-05-architecture.md) | §4 Architecture Overview |
| [SPEC-BOOT-06-packages.md](SPEC-BOOT/SPEC-BOOT-06-packages.md) | §5 Package Specifications |
| [SPEC-BOOT-07-configuration.md](SPEC-BOOT/SPEC-BOOT-07-configuration.md) | §6 Configuration Specification |
| [SPEC-BOOT-08-build.md](SPEC-BOOT/SPEC-BOOT-08-build.md) | §7 Build, Tooling & Quality |
| [SPEC-BOOT-09-testing.md](SPEC-BOOT/SPEC-BOOT-09-testing.md) | §8 Testing Specification |
| [SPEC-BOOT-10-documentation.md](SPEC-BOOT/SPEC-BOOT-10-documentation.md) | §9 Documentation Specification |
| [SPEC-BOOT-A-checklist.md](SPEC-BOOT/SPEC-BOOT-A-checklist.md) | Appendix A — Specification Checklist |
| [SPEC-BOOT-B-tasks.md](SPEC-BOOT/SPEC-BOOT-B-tasks.md) | Appendix B — Implementation Task Decomposition |

---

### SPEC-MAKE — Makefile Targets

| | |
|-|-|
| **Story** | 02 — Makefile targets wrapping `docker compose`, `modd`, and the Go toolchain |
| **Status** | Draft |
| **Version** | 1.0.0 |
| **SPEC-IDs** | SPEC-MAKE-001 – SPEC-MAKE-020 |

| File | Contents |
|------|----------|
| [SPEC-MAKE-00-index.md](SPEC-MAKE/SPEC-MAKE-00-index.md) | Index — metadata, SPEC-ID registry |
| [SPEC-MAKE-01-preamble.md](SPEC-MAKE/SPEC-MAKE-01-preamble.md) | §0 AI Steering Preamble |
| [SPEC-MAKE-02-introduction.md](SPEC-MAKE/SPEC-MAKE-02-introduction.md) | §1 Introduction |
| [SPEC-MAKE-03-goals.md](SPEC-MAKE/SPEC-MAKE-03-goals.md) | §2 Goals & Non-Goals |
| [SPEC-MAKE-04-context.md](SPEC-MAKE/SPEC-MAKE-04-context.md) | §3 System Context & Dependencies |
| [SPEC-MAKE-05-architecture.md](SPEC-MAKE/SPEC-MAKE-05-architecture.md) | §4 Architecture Overview |
| [SPEC-MAKE-06-artifacts.md](SPEC-MAKE/SPEC-MAKE-06-artifacts.md) | §5 Artifact Specifications |
| [SPEC-MAKE-07-configuration.md](SPEC-MAKE/SPEC-MAKE-07-configuration.md) | §6 Configuration Specification |
| [SPEC-MAKE-08-build.md](SPEC-MAKE/SPEC-MAKE-08-build.md) | §7 Build, Tooling & Quality |
| [SPEC-MAKE-09-testing.md](SPEC-MAKE/SPEC-MAKE-09-testing.md) | §8 Testing Specification |
| [SPEC-MAKE-10-documentation.md](SPEC-MAKE/SPEC-MAKE-10-documentation.md) | §9 Documentation Specification |
| [SPEC-MAKE-A-checklist.md](SPEC-MAKE/SPEC-MAKE-A-checklist.md) | Appendix A — Specification Checklist |
| [SPEC-MAKE-B-tasks.md](SPEC-MAKE/SPEC-MAKE-B-tasks.md) | Appendix B — Implementation Task Decomposition |

---

### SPEC-DB — Database: pgx Pool, Migrations & Initial Schema

| | |
|-|-|
| **Story** | 03 — Database: pgx pool + migrations + 0001_init |
| **Status** | Final |
| **Version** | 1.0.0 |
| **SPEC-IDs** | SPEC-DB-001 – SPEC-DB-027 |

| File | Contents |
|------|----------|
| [SPEC-DB-00-index.md](SPEC-DB/SPEC-DB-00-index.md) | Index — metadata, SPEC-ID registry |
| [SPEC-DB-01-preamble.md](SPEC-DB/SPEC-DB-01-preamble.md) | §0 AI Steering Preamble |
| [SPEC-DB-02-introduction.md](SPEC-DB/SPEC-DB-02-introduction.md) | §1 Introduction |
| [SPEC-DB-03-goals.md](SPEC-DB/SPEC-DB-03-goals.md) | §2 Goals & Non-Goals |
| [SPEC-DB-04-context.md](SPEC-DB/SPEC-DB-04-context.md) | §3 System Context & Dependencies |
| [SPEC-DB-05-architecture.md](SPEC-DB/SPEC-DB-05-architecture.md) | §4 Architecture Overview |
| [SPEC-DB-06-packages.md](SPEC-DB/SPEC-DB-06-packages.md) | §5 Package Specifications |
| [SPEC-DB-07-configuration.md](SPEC-DB/SPEC-DB-07-configuration.md) | §6 Configuration Specification |
| [SPEC-DB-08-build.md](SPEC-DB/SPEC-DB-08-build.md) | §7 Build, Tooling & Quality |
| [SPEC-DB-09-testing.md](SPEC-DB/SPEC-DB-09-testing.md) | §8 Testing Specification |
| [SPEC-DB-10-documentation.md](SPEC-DB/SPEC-DB-10-documentation.md) | §9 Documentation Specification |
| [SPEC-DB-A-checklist.md](SPEC-DB/SPEC-DB-A-checklist.md) | Appendix A — Specification Checklist |
| [SPEC-DB-B-tasks.md](SPEC-DB/SPEC-DB-B-tasks.md) | Appendix B — Implementation Task Decomposition |

---

### SPEC-AUTH — Firebase ID Token Middleware + Lazy User Provisioning

| | |
|-|-|
| **Story** | 04 — Firebase ID token middleware + lazy user provisioning |
| **Status** | Draft → Final |
| **Date** | 2026-05-02 |
| **SPEC-IDs** | SPEC-AUTH-001 – SPEC-AUTH-026 |

| File | Contents |
|------|----------|
| [SPEC-AUTH-00-index.md](SPEC-AUTH/SPEC-AUTH-00-index.md) | Index — metadata, SPEC-ID registry |
| [SPEC-AUTH-01-preamble.md](SPEC-AUTH/SPEC-AUTH-01-preamble.md) | AI constraints, authorship, traceability rules |
| [SPEC-AUTH-02-introduction.md](SPEC-AUTH/SPEC-AUTH-02-introduction.md) *(planned)* | Story relationship, SPEC-ID registry |
| [SPEC-AUTH-03-goals.md](SPEC-AUTH/SPEC-AUTH-03-goals.md) *(planned)* | Goals, non-goals, constraints |
| [SPEC-AUTH-04-context.md](SPEC-AUTH/SPEC-AUTH-04-context.md) *(planned)* | External dependencies, decision references |
| [SPEC-AUTH-05-architecture.md](SPEC-AUTH/SPEC-AUTH-05-architecture.md) *(planned)* | Token flow, package dependency graph, middleware lifecycle |
| [SPEC-AUTH-06-packages.md](SPEC-AUTH/SPEC-AUTH-06-packages.md) *(planned)* | All SPEC-AUTH-NNN requirements: Go signatures, error taxonomy, SQL |
| [SPEC-AUTH-07-configuration.md](SPEC-AUTH/SPEC-AUTH-07-configuration.md) *(planned)* | Environment variables, credential file rules |
| [SPEC-AUTH-08-build.md](SPEC-AUTH/SPEC-AUTH-08-build.md) *(planned)* | Build, lint, vet constraints |
| [SPEC-AUTH-09-testing.md](SPEC-AUTH/SPEC-AUTH-09-testing.md) *(planned)* | Test strategy, named tests, coverage targets |

---

## Directory structure

```
specifications/
├── README.md                   ← this file
├── SPEC-AUTH/                  ← Story 04 · Firebase Auth middleware
│   ├── SPEC-AUTH-00-index.md
│   └── SPEC-AUTH-01-preamble.md
├── SPEC-BOOT/                  ← Story 01 · Bootstrap project structure
│   ├── SPEC-BOOT-00-index.md
│   ├── SPEC-BOOT-01-preamble.md
│   ├── SPEC-BOOT-02-introduction.md
│   ├── SPEC-BOOT-03-goals.md
│   ├── SPEC-BOOT-04-context.md
│   ├── SPEC-BOOT-05-architecture.md
│   ├── SPEC-BOOT-06-packages.md
│   ├── SPEC-BOOT-07-configuration.md
│   ├── SPEC-BOOT-08-build.md
│   ├── SPEC-BOOT-09-testing.md
│   ├── SPEC-BOOT-10-documentation.md
│   ├── SPEC-BOOT-A-checklist.md
│   └── SPEC-BOOT-B-tasks.md
├── SPEC-DB/                    ← Story 03 · pgx pool, migrations, schema
│   ├── SPEC-DB-00-index.md
│   ├── SPEC-DB-01-preamble.md
│   ├── SPEC-DB-02-introduction.md
│   ├── SPEC-DB-03-goals.md
│   ├── SPEC-DB-04-context.md
│   ├── SPEC-DB-05-architecture.md
│   ├── SPEC-DB-06-packages.md
│   ├── SPEC-DB-07-configuration.md
│   ├── SPEC-DB-08-build.md
│   ├── SPEC-DB-09-testing.md
│   ├── SPEC-DB-10-documentation.md
│   ├── SPEC-DB-A-checklist.md
│   └── SPEC-DB-B-tasks.md
└── SPEC-MAKE/                  ← Story 02 · Makefile targets
    ├── SPEC-MAKE-00-index.md
    ├── SPEC-MAKE-01-preamble.md
    ├── SPEC-MAKE-02-introduction.md
    ├── SPEC-MAKE-03-goals.md
    ├── SPEC-MAKE-04-context.md
    ├── SPEC-MAKE-05-architecture.md
    ├── SPEC-MAKE-06-artifacts.md
    ├── SPEC-MAKE-07-configuration.md
    ├── SPEC-MAKE-08-build.md
    ├── SPEC-MAKE-09-testing.md
    ├── SPEC-MAKE-10-documentation.md
    ├── SPEC-MAKE-A-checklist.md
    └── SPEC-MAKE-B-tasks.md
```

---

## File naming conventions

```
SPEC-{MODULE}-{NN}-{slug}.md
```

| Segment | Rule |
|---------|------|
| `SPEC-{MODULE}` | Upper-case module token (e.g. `BOOT`, `DB`, `MAKE`, `AUTH`) |
| `{NN}` | Two-digit section number: `00` = index, `01`–`10` = ordered sections, `A`–`Z` = appendices |
| `{slug}` | Lower-case, hyphen-separated descriptor matching the section title |

### Standard section sequence

| Number | Typical content |
|--------|----------------|
| `00` | Index — metadata, file map, SPEC-ID registry |
| `01` | Preamble — AI steering, authorship, traceability |
| `02` | Introduction — scope, motivation, story relationship |
| `03` | Goals & Non-Goals |
| `04` | System Context & External Dependencies |
| `05` | Architecture Overview |
| `06` | Package / Artifact Specifications (SPEC-ID requirements) |
| `07` | Configuration Specification |
| `08` | Build, Tooling & Quality |
| `09` | Testing Specification |
| `10` | Documentation Specification |
| `A` | Appendix A — Specification Checklist |
| `B` | Appendix B — Implementation Task Decomposition |

---

## Adding a new specification

1. Create a directory `specifications/SPEC-{MODULE}/`.
2. Start from `SPEC-{MODULE}-00-index.md` — copy the preamble block from an
   existing index and fill in the metadata table.
3. Assign sequential `SPEC-{MODULE}-NNN` identifiers for every requirement
   (`NNN` starts at `001` and never reuses a retired ID).
4. Follow the standard section sequence above.
5. Register the new module in this README under **Specifications index**.
6. Every `SPEC-*` ID must have at least one test before the spec is
   considered **Final**.
