# SPEC-BOOT-02 — Introduction

> Part of [SPEC-BOOT](SPEC-BOOT-00-index.md) — Story 01: Bootstrap Project Structure

---

## 1. Introduction

**cooksense-backend** is the Go-based HTTP API server for the CookSense recipe application. It **shall** operate as a standalone HTTP microservice that provides recipe discovery, ingredient-based search, user reactions, and cooking lessons to a Firebase-authenticated mobile/web frontend.

This specification covers Story 01: **the foundational project skeleton** that every subsequent story builds upon. It defines the package layout, `doc.go` convention, placeholder entry point, migration directory structure, `.gitignore` rules, and `go.mod` version declaration.

### 1.1 Scope

This document specifies:

- The canonical `internal/` package tree and the `doc.go` placeholder contract for each package.
- The placeholder `cmd/cooksense-server/main.go` behaviour (compile + print + exit 0).
- The `migrations/`, `seed/recipes/`, and `seed/lessons/` directory structure.
- The `.gitignore` rules protecting secrets, build artefacts, and local overrides.
- The `go.mod` module path and Go version constraint.

HTTP server construction uses Go's standard library (`net/http`) directly. No internal HTTP wrapper package is introduced.

### 1.2 Definitions

| Term | Definition |
|------|-----------|
| **doc.go** | A Go source file whose sole purpose is to carry the package-level documentation comment. Contains only the `package` declaration (and the comment above it). |
| **placeholder** | A file that satisfies the compiler and the package contract but contains no production logic. It will be replaced by a fully-featured implementation in a later story. |
| **SPEC-BOOT-NNN** | Requirement identifier for the Bootstrap story. NNN is a zero-padded three-digit number. |
| **internal package** | A Go package under `internal/` that cannot be imported by code outside this module. |
