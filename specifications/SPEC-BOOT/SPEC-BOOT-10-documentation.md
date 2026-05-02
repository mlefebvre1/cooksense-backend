# SPEC-BOOT-10 — Documentation Specification

> Part of [SPEC-BOOT](SPEC-BOOT-00-index.md) — Story 01: Bootstrap Project Structure

---

## 9. Documentation Specification

### 9.1 Package Docstring Standard

Every `doc.go` **shall** follow this template:

```go
// Package {name} {one-line description starting with a verb}.
package {name}
```

The description **shall** be a single sentence. It **shall** start with a verb in the present tense (e.g., "provides", "implements", "defines", "loads").

Examples of compliant comments:

```go
// Package config provides runtime configuration loading and validation for cooksense-server.
package config
```

```go
// Package domain defines the core business entities, value objects, and repository interfaces for CookSense.
package domain
```

### 9.2 `main.go` Documentation

`cmd/cooksense-server/main.go` **shall** carry a file-level comment stating the binary's purpose:

```go
// Command cooksense-server is the HTTP API server for the CookSense application.
package main
```

### 9.3 README Maintenance

Story 01 **shall not** modify `README.md` (that is deferred to Story 12). The PR description **shall** include the output of `tree -L 2 internal/` to prove the layout is correct.
