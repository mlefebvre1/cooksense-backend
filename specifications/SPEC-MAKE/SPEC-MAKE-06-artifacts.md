# SPEC-MAKE-06 — Artifact Specifications

> Part of [SPEC-MAKE](SPEC-MAKE-00-index.md) — Story 02: Makefile Targets  
> SPEC-IDs covered: SPEC-MAKE-001 – SPEC-MAKE-017

---

## 5. Artifact Specifications

### 5.1 Makefile existence

#### SPEC-MAKE-001: `Makefile` at repository root

A file named `Makefile` **shall** exist at the repository root (same level as `go.mod`).

---

### 5.2 `help` — default target (self-documentation)

#### SPEC-MAKE-002: Self-documenting `help` target

`help` **shall** be the default target (first target defined, or declared via `.DEFAULT_GOAL`). It **shall** parse all `## <description>` trailing comments from the `Makefile` and print them as a formatted table to stdout.

Minimum accepted output format:

```
  help        Show this help message
  up          Start Postgres (docker compose up -d)
  down        Stop Postgres (docker compose down)
  migrate     Apply all pending SQL migrations
  seed        Load YAML seed data into the database
  run         Start server with hot-reload (modd)
  build       Compile the server binary to bin/cooksense-server
  test        Run the full test suite (go test ./...)
  lint        Run go vet (+ golangci-lint if installed)
  clean       Remove build artefacts (CLEAN_VOLUMES=1 to also drop DB)
```

#### SPEC-MAKE-003: Public target list

The `Makefile` **shall** define exactly the following 10 public targets: `help`, `up`, `down`, `migrate`, `seed`, `run`, `build`, `test`, `lint`, `clean`. No additional public targets are introduced in Story 02.

---

### 5.3 `up` — start Postgres

#### SPEC-MAKE-004: `up` invokes docker compose

`make up` **shall** invoke `docker compose up -d` to start the services defined in `docker-compose.yml` in detached mode.

#### SPEC-MAKE-005: `up` waits for Postgres health

`make up` **shall** wait until the Postgres container reports `healthy` (as defined by the `healthcheck` in `docker-compose.yml`) before returning. The wait **shall** time out after ≤ 30 seconds and print an error message if the container does not become healthy in time. The wait **shall** be implemented without external tools (a shell `until` loop with `docker inspect` or `docker compose ps` is acceptable).

---

### 5.4 `down` — stop Postgres

#### SPEC-MAKE-006: `down` preserves volume

`make down` **shall** invoke `docker compose down` (without `--volumes`). The `postgres_data` volume **shall not** be deleted by `make down`.

---

### 5.5 `build` — compile binary

#### SPEC-MAKE-007: `build` produces executable

`make build` **shall** invoke:

```sh
go build -o bin/cooksense-server ./cmd/cooksense-server
```

The resulting file `bin/cooksense-server` **shall** be executable (`chmod +x` if needed on Unix). `bin/` **shall** be created by the recipe if it does not exist.

---

### 5.6 `run` — hot-reload dev server

#### SPEC-MAKE-008: `run` uses modd

`make run` **shall** invoke `modd` to start the server with hot-reload using the existing `modd.conf`. The target **shall** declare `up migrate` as prerequisites.

---

### 5.7 `test` — full test suite

#### SPEC-MAKE-009: `test` runs all tests

`make test` **shall** invoke `go test ./...`. The exit code from `go test` **shall** be propagated by `make` (non-zero on failure).

---

### 5.8 `lint` — static analysis

#### SPEC-MAKE-010: `lint` runs go vet, optionally golangci-lint

`make lint` **shall** always run `go vet ./...`. If `golangci-lint` is found on `PATH` (via `which golangci-lint`), it **shall** additionally run `golangci-lint run`. If `golangci-lint` is not on PATH, the recipe **shall** print: `golangci-lint not found, skipping` and exit `0`.

---

### 5.9 `clean` — remove build artefacts

#### SPEC-MAKE-011: `clean` target contract

`make clean` **shall** remove the `bin/` directory (`rm -rf bin/`). It **shall not** drop the Docker volume unless `CLEAN_VOLUMES=1` is passed:

```sh
make clean CLEAN_VOLUMES=1   # removes bin/ AND drops the postgres_data volume
```

When `CLEAN_VOLUMES=1`, the recipe **shall** invoke `docker compose down --volumes` after `rm -rf bin/`.

---

### 5.10 `migrate` and `seed`

#### SPEC-MAKE-012: Subcommand delegation

`make migrate` **shall** invoke:

```sh
go run ./cmd/cooksense-server migrate up
```

`make seed` **shall** invoke:

```sh
go run ./cmd/cooksense-server seed
```

Both targets **shall** declare `up` as a prerequisite. Both targets **may** print "not implemented" and exit `0` until stories 03 and 05 land, as per the story scope.

---

### 5.11 Structural rules

#### SPEC-MAKE-013: `.PHONY` declaration

All 10 public targets **shall** be listed in a single `.PHONY` line near the top of the `Makefile`.

#### SPEC-MAKE-014: Conditional `.env` loading

The `Makefile` **shall** use `-include .env` (with the leading `-`) to silently load `.env` if present and skip it if absent.

#### SPEC-MAKE-015: Line count constraint

The `Makefile` **shall** be at most **80 lines** total. Any logic exceeding this limit **shall** be moved to a shell script under `scripts/` and called from the `Makefile`.

---

### 5.12 `.env.example`

#### SPEC-MAKE-016: Variable coverage

`.env.example` **shall** include every variable defined in `docs/architecture/infra.md`, grouped as:

- `APP_PORT`, `LOG_LEVEL`, `LOG_FORMAT`
- `DATABASE_URL`
- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT`
- `FIREBASE_PROJECT_ID`, `GOOGLE_APPLICATION_CREDENTIALS`

#### SPEC-MAKE-017: No real secrets

`.env.example` **shall** contain only placeholder values (e.g., `changeme`, `path/to/firebase-admin.json`). Real credentials **shall never** be committed to this file.
