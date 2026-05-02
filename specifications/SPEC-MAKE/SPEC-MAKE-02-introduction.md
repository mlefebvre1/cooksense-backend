# SPEC-MAKE-02 — Introduction

> Part of [SPEC-MAKE](SPEC-MAKE-00-index.md) — Story 02: Makefile Targets

---

## 1. Introduction

The `Makefile` is the **canonical developer entrypoint** for the cooksense-backend project. It abstracts away the specific invocations of `docker compose`, `modd`, `go build`, `go test`, `go vet`, and `golangci-lint` so that contributors do not need to memorise command-line flags or path structure.

This specification covers Story 02: **authoring and validating all 10 `make` targets** that wrap the existing `docker-compose.yml` and Go toolchain.

### 1.1 Scope

This document specifies:

- The content and behaviour of each of the 10 public `make` targets.
- The Makefile structural rules (`.PHONY`, `##` doc comments, line limit).
- The `.env.example` variable list.
- The `.env` loading strategy (include-if-present).

### 1.2 Relationship to Other Stories

| Story | Relationship |
|-------|-------------|
| Story 01 | Provides the directory layout and `cmd/cooksense-server` binary path that the `Makefile` hard-codes. **Must merge first.** |
| Story 03 | `make migrate` will invoke `cooksense-server migrate up` — the subcommand is wired in story 03, but the make target ships in story 02. |
| Story 05 | `make seed` will invoke `cooksense-server seed` — wired in story 05. |
| Story 12 | README documents the `make` targets — depends on this story. |

### 1.3 Definitions

| Term | Definition |
|------|-----------|
| **target** | A named recipe in `Makefile` that `make <target>` executes. |
| **`.PHONY` target** | A target whose name does not correspond to a real file; always runs when invoked. |
| **`##` comment** | Inline comment after a target declaration used by `help` for self-documentation. |
| **SPEC-MAKE-NNN** | Requirement identifier for the Makefile story. |
