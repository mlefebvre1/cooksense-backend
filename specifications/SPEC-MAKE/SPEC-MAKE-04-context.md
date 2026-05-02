# SPEC-MAKE-04 — System Context & Dependencies

> Part of [SPEC-MAKE](SPEC-MAKE-00-index.md) — Story 02: Makefile Targets

---

## 3. System Context & Dependencies

### 3.1 External Tooling Requirements

| Tool | Version constraint | How used |
|------|--------------------|----------|
| **GNU Make** | ≥ 3.81 (macOS ships 3.81; Linux ships ≥ 4.x) | Executes the `Makefile`. |
| **docker compose** | V2 plugin (`docker compose`, not `docker-compose`) | `up`, `down`, `clean` targets. |
| **Go** | `1.26.2` (matches `go.mod`) | `build`, `test`, `lint`, `migrate`, `seed`, `run`. |
| **modd** | any version on PATH | `run` target (hot reload). |
| **golangci-lint** | any version on PATH (optional) | `lint` target; graceful skip if absent. |

### 3.2 Files Consumed / Produced

| Artefact | Consumed / Produced | Notes |
|----------|---------------------|-------|
| `docker-compose.yml` | Consumed | Already in repo — not modified by this story. |
| `.env` | Consumed (optional) | Loaded if present; must be `.gitignore`d. |
| `.env.example` | Produced (updated) | Mirrors `docs/architecture/infra.md`. |
| `bin/cooksense-server` | Produced | Build output; under `bin/` (`.gitignore`d). |
| `go.sum` | Consumed | Must not be modified in this story. |
| `modd.conf` | Consumed | Already in repo — not modified. |

### 3.3 Dependencies on Other Stories

| Story | Kind | Notes |
|-------|------|-------|
| Story 01 | Hard pre-requisite | `cmd/cooksense-server/main.go` must compile. |
| Story 03 | Soft (runtime only) | `make migrate` calls `cooksense-server migrate up`; will print "not implemented" until story 03 merges. |
| Story 05 | Soft (runtime only) | `make seed` calls `cooksense-server seed`; will print "not implemented" until story 05 merges. |
