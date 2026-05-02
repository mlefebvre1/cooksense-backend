# Story 02 — Makefile targets wrapping docker-compose

Status: TODO
Estimate: S

## User story

As a developer, I want a single `make` entrypoint for every common dev task so
that the workflow does not depend on memorizing `docker compose` flags or
`go run` paths.

## Background

`docker-compose.yml` already provides Postgres 17. We want `make` to be the
canonical UX. See `docs/architecture/infra.md` for the full target list.

## Acceptance criteria

- [ ] A top-level `Makefile` exists.
- [ ] `make help` (default target) prints a list of available targets with a
      one-line description per target.
- [ ] Targets implemented: `up`, `down`, `migrate`, `seed`, `run`, `build`,
      `test`, `lint`, `clean`, `help`.
- [ ] `make up` runs `docker compose up -d` and waits until the Postgres
      healthcheck reports healthy before returning (use the existing
      healthcheck — wait via `docker compose wait` or a small loop).
- [ ] `make down` runs `docker compose down` (data volume is preserved).
- [ ] `make build` produces `bin/cooksense-server`.
- [ ] `make run` invokes `modd` (using the existing `modd.conf`).
- [ ] `make test` runs `go test ./...`.
- [ ] `make lint` runs `go vet ./...`. If `golangci-lint` is on PATH, run it
      too; otherwise just print a notice.
- [ ] `make clean` removes `bin/`. It does **not** drop the Postgres volume
      unless `CLEAN_VOLUMES=1` is set.
- [ ] An `.env.example` is present and lists every variable from
      `docs/architecture/infra.md`.

## Technical notes

- Use `.PHONY` for every non-file target.
- Forward env vars from `.env` if it exists (`include .env`-style is fine but
  guard against missing file).
- Keep the Makefile under 80 lines; complex logic belongs in scripts under
  `scripts/`, not in Make recipes.

## Out of scope

- CI workflows.
- Real implementations of `migrate`/`seed` subcommands — they will be wired
  by stories 03 and 05. For this story it is acceptable that `make migrate`
  and `make seed` invoke `go run ./cmd/cooksense-server <subcmd>` even if the
  subcommand currently prints "not implemented" and exits 0.

## Dependencies

- depends on: 01
- blocks: 03, 12

## Definition of Done

- [ ] AC met.
- [ ] `make help` output reviewed for clarity.
- [ ] `.env.example` reviewed against the env-vars table in
      `docs/architecture/infra.md`.
