# SPEC-BOOT-06 — Package Specifications

> Part of [SPEC-BOOT](SPEC-BOOT-00-index.md) — Story 01: Bootstrap Project Structure  
> SPEC-IDs covered: SPEC-BOOT-001 – SPEC-BOOT-018

---

## 5. Package Specifications

> For Story 01 each package spec covers only the `doc.go` placeholder. Method and struct tables will be populated in the stories that implement each package's real functionality.

### 5.1 `internal/config` — Runtime configuration loading

#### SPEC-BOOT-001: `internal/config` package declaration

The system **shall** provide a file `internal/config/doc.go` with:

- A package-level comment: `// Package config provides runtime configuration loading and validation for cooksense-server.`
- Only the `package config` declaration below the comment.
- No imports.
- No exported symbols.

#### SPEC-BOOT-002: `internal/config` isolation

`internal/config` **shall not** import any other `internal/` package in Story 01.

---

### 5.2 `internal/auth` — Firebase authentication

#### SPEC-BOOT-003: `internal/auth` package declaration

The system **shall** provide a file `internal/auth/doc.go` with:

- A package-level comment: `// Package auth provides Firebase ID-token verification and authenticated-user context helpers.`
- Only the `package auth` declaration.
- No imports. No exported symbols.

---

### 5.3 `internal/db` — Database connection pool

#### SPEC-BOOT-004: `internal/db` package declaration

The system **shall** provide a file `internal/db/doc.go` with:

- A package-level comment: `// Package db provides the PostgreSQL connection pool and migration runner for cooksense-server.`
- Only the `package db` declaration.
- No imports. No exported symbols.

---

### 5.4 `internal/domain` — Business types and interfaces

#### SPEC-BOOT-005: `internal/domain` package declaration

The system **shall** provide a file `internal/domain/doc.go` with:

- A package-level comment: `// Package domain defines the core business entities, value objects, and repository interfaces for CookSense.`
- Only the `package domain` declaration.
- No imports. No exported symbols.

#### SPEC-BOOT-006: `internal/domain` zero-import rule

`internal/domain` **shall not** import any other `internal/` package — now or in future stories. Violations are rejected at code review.

---

### 5.5 `internal/recipes` — Recipe feature

#### SPEC-BOOT-007: `internal/recipes` package declaration

The system **shall** provide a file `internal/recipes/doc.go` with:

- A package-level comment: `// Package recipes implements the recipe discovery, detail, and ingredient-search feature of CookSense.`
- Only the `package recipes` declaration.
- No imports. No exported symbols.

---

### 5.6 `internal/reactions` — User reactions

#### SPEC-BOOT-008: `internal/reactions` package declaration

The system **shall** provide a file `internal/reactions/doc.go` with:

- A package-level comment: `// Package reactions implements the LIKE / DISLIKE / TRY_LATER reaction feature for CookSense recipes.`
- Only the `package reactions` declaration.
- No imports. No exported symbols.

---

### 5.7 `internal/lessons` — Cooking school

#### SPEC-BOOT-009: `internal/lessons` package declaration

The system **shall** provide a file `internal/lessons/doc.go` with:

- A package-level comment: `// Package lessons implements the Cooking School feature, serving curated culinary lesson articles.`
- Only the `package lessons` declaration.
- No imports. No exported symbols.

---

### 5.8 `internal/users` — User provisioning

#### SPEC-BOOT-010: `internal/users` package declaration

The system **shall** provide a file `internal/users/doc.go` with:

- A package-level comment: `// Package users provides lazy user provisioning and profile management backed by Firebase UID.`
- Only the `package users` declaration.
- No imports. No exported symbols.

---

### 5.9 `internal/seed` — Seed data loader

#### SPEC-BOOT-011: `internal/seed` package declaration

The system **shall** provide a file `internal/seed/doc.go` with:

- A package-level comment: `// Package seed loads curated recipe and lesson data from YAML files into the database at startup.`
- Only the `package seed` declaration.
- No imports. No exported symbols.

---

### 5.10 `cmd/cooksense-server/main.go` — Entry point

#### SPEC-BOOT-012: Placeholder `main` function

The file `cmd/cooksense-server/main.go` **shall** contain a `main` function that:

1. Prints the exact string `"cooksense-server starting"` to stdout using `fmt.Println`.
2. Exits with code `0` (normal return from `main`).

It **shall not** import any `internal/` package in Story 01 (real wiring arrives in stories 03/04/07).

#### SPEC-BOOT-013: `main.go` imports

`cmd/cooksense-server/main.go` **shall** import only `"fmt"` from the standard library.

---

### 5.11 `migrations/` — SQL migration files

#### SPEC-BOOT-014: `migrations/` directory

The directory `migrations/` **shall** exist at the project root. Because no SQL files are added in Story 01, the directory **shall** contain a single `.gitkeep` file to ensure it is tracked by Git.

---

### 5.12 `seed/recipes/` and `seed/lessons/` — Seed YAML directories

#### SPEC-BOOT-015: Seed sub-directories

The directories `seed/recipes/` and `seed/lessons/` **shall** exist at the project root. Because no YAML files are added in Story 01, each directory **shall** contain a single `.gitkeep` file.

---

### 5.13 `.gitignore`

#### SPEC-BOOT-016: `.gitignore` rules

The project-root `.gitignore` **shall** contain at minimum the following rules:

| Pattern | Rationale |
|---------|-----------|
| `bin/` | Compiled binaries produced by `go build`. |
| `secrets/` | Directory for local service-account JSON files (except `*.example`). |
| `secrets/*.json` | Firebase Admin SDK credentials and other JSON secrets. |
| `!secrets/*.example` | Explicitly un-ignores example/template credential files so they are committed. |
| `.env` | Local environment variable overrides. |
| `*.local.*` | Any file following the `name.local.ext` convention (e.g., `config.local.yaml`). |
| `coverage.out` | Go test coverage profiles. |
| `*.test` | Compiled test binaries. |

The `.gitignore` **should** also include standard Go and OS rules (e.g., `vendor/`, `.DS_Store`).

---

### 5.14 `go.mod` — Module declaration

#### SPEC-BOOT-017: Go version

The `go.mod` file **shall** declare `go 1.26.2`. Any lower version declaration is a violation of decision D-0001.

#### SPEC-BOOT-018: Module path

The `go.mod` **shall** preserve the existing module path. It **shall not** be changed in Story 01.
