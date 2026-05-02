# Story 01 — Bootstrap project structure

Status: TODO
Estimate: S

## User story

As a developer, I want a clean and conventional project skeleton so that any
contributor can find code by feature without hunting.

## Background

`cmd/cooksense-server/main.go` exists but the rest of the layout is empty.
We need the directory tree described in `docs/architecture/overview.md` to
land before any feature work.

## Acceptance criteria

- [ ] The following directories exist and contain at least a placeholder
      `doc.go` file with the package comment:
  - `internal/config`
  - `internal/httpx`
  - `internal/auth`
  - `internal/db`
  - `internal/domain`
  - `internal/recipes`
  - `internal/reactions`
  - `internal/lessons`
  - `internal/users`
  - `internal/seed`
- [ ] `cmd/cooksense-server/main.go` compiles and prints "cooksense-server
      starting" then exits 0 (placeholder; real wiring lands in story #04/#07).
- [ ] `migrations/` and `seed/recipes/`, `seed/lessons/` directories exist
      (with `.gitkeep` if empty).
- [ ] `.gitignore` ignores `bin/`, `secrets/*.json` (except `*.example`),
      `.env`, and `*.local.*`.
- [ ] `go.mod` declares `go 1.26.2` (matching D-0001).
- [ ] `go build ./...` succeeds.

## Technical notes

- Package comments must be a single sentence (`// Package recipes provides …`).
- Do not import other internal packages from `internal/domain`.
- No third-party dependencies in this story; the next stories will add
  `pgx`, `firebase-admin`, `golang-migrate`, `yaml.v3`.

## Out of scope

- Real `main.go` wiring (config, db, server) — that lands in stories 03/04/07.
- Makefile — story 02.

## Dependencies

- depends on: —
- blocks: 02, 03, 04, 05

## Definition of Done

- [ ] AC met.
- [ ] `go vet ./...` clean.
- [ ] PR description includes a screenshot or `tree -L 2 internal/` output.
