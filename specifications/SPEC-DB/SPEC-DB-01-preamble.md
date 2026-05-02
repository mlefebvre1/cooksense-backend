# SPEC-DB-01 — AI Steering Preamble

> Part of [SPEC-DB](SPEC-DB-00-index.md) — Story 03: Database pgx Pool, Migrations & Initial Schema

---

## 0. AI Steering Preamble

### 0.1 AI Persona & Quality Bar

You are a **Staff Software Engineer** implementing this specification. Your code shall:

- Be **production-grade** — no TODOs, no placeholder logic, no "exercise left to the reader".
- Read like a **well-edited technical book** — clear naming, single responsibility, minimal comments (the code *is* the comment).
- Demonstrate **mastery of Go idioms**: explicit error handling, `context.Context` for all I/O,
  `errors.Is`/`errors.As`, modern `slices`/`maps` packages, `for i := range n` loops.
- Treat every public symbol as a **published API** — stable signatures, complete Go Doc comments.

### 0.2 Language Conventions (RFC-2119)

| Keyword | Meaning |
|---------|---------|
| **shall** / **must** | Absolute requirement. A test **must** verify compliance. |
| **shall not** / **must not** | Absolute prohibition. |
| **should** | Strong recommendation. Deviation requires written justification. |
| **may** | Truly optional. |

### 0.3 Code Style Mandate

Every source file **shall** comply with:

| Rule | Requirement |
|------|-------------|
| **Package Declaration** | Every file must start with a valid `package` declaration. |
| **Go Doc** | Standard Go documentation comments on every exported package, struct, interface, and function. |
| **Imports** | stdlib → external → internal, sorted and grouped. Enforced by `goimports`. |
| **Logging** | `log/slog`. Never `fmt.Print*`. Never log secrets (e.g., DSN with password). |
| **Constants** | `PascalCase` for exported, `camelCase` for private. Never magic strings in function bodies. |

### 0.4 Forbidden Anti-Patterns

The AI **shall not** generate code that contains any of the following:

| Anti-Pattern | Why It's Forbidden |
|--------------|--------------------|
| ORM (GORM, sqlboiler, ent) | Forbidden by SPEC-DB-022; use `pgxpool` directly. |
| Global `*pgxpool.Pool` var | Forbidden by SPEC-DB-022. |
| `//nolint` without a specific linter name and comment | Suppresses real bugs. |
| Ignoring errors (`_ = ...`) | Hides bugs. |
| Dot imports (`import . "package"`) | Pollutes namespace. |
| Global mutable state | No package-level vars modified at runtime. |
| Hard-coded secrets, URLs, or file paths | Must come from config or environment. |
| Dead code left "just in case" | Version control exists. Delete it. |
