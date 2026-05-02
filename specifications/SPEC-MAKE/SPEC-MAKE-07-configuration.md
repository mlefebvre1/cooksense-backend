# SPEC-MAKE-07 — Configuration Specification

> Part of [SPEC-MAKE](SPEC-MAKE-00-index.md) — Story 02: Makefile Targets

---

## 6. Configuration Specification

### 6.1 Makefile Variables

The `Makefile` **shall** declare the following variables near the top, after the `-include .env` line:

| Variable | Default value | Description |
|----------|--------------|-------------|
| `BINARY` | `cooksense-server` | Name of the compiled binary. |
| `CMD_DIR` | `./cmd/cooksense-server` | Path to the main package. |
| `BIN_DIR` | `bin` | Output directory for compiled binaries. |
| `CLEAN_VOLUMES` | (unset) | Set to `1` to also drop the Postgres volume on `make clean`. |

Variables **shall** be declared using simple (`=`) assignment, not recursive expansion (`:=` is acceptable for variables that expand once). No variable **shall** be hard-coded inside a recipe; all paths go through the variable.

### 6.2 Environment Variables Forwarded to Sub-processes

When `.env` is included, all variables in it are automatically exported to sub-process environments by `make`. The `Makefile` **shall not** `export` them explicitly (redundant). The `.env` file **shall not** be committed; only `.env.example` is tracked.
