# SPEC-BOOT-03 — Goals & Non-Goals

> Part of [SPEC-BOOT](SPEC-BOOT-00-index.md) — Story 01: Bootstrap Project Structure

---

## 2. Goals & Non-Goals

### 2.1 Goals

| ID | Goal |
|----|------|
| G-1 | Provide a **complete, compilable project skeleton** so that any contributor can clone and run `go build ./...` on a clean machine without errors. |
| G-2 | Enforce the **Clean Architecture layer contract** through the package tree before any business logic is written. |
| G-3 | Protect secrets and build artefacts from accidental commits via a **complete `.gitignore`**. |
| G-4 | Anchor the **Go version** at `1.26.2` so all contributors and CI use identical language semantics. |
| G-5 | Provide a **navigable package tree** where any contributor can locate the right package by feature name, not by hunting. |

### 2.2 Non-Goals

| ID | Non-Goal |
|----|----------|
| NG-1 | This story **does not** implement any real wiring (config loading, database pool, HTTP server, Firebase auth). That is covered by stories 03, 04, 07. |
| NG-2 | This story **does not** add the `Makefile`. That is covered by story 02. |
| NG-3 | This story **does not** add third-party dependencies (`pgx`, `firebase-admin`, `golang-migrate`, `yaml.v3`). |
| NG-4 | This story **does not** write integration tests or seed data. |
