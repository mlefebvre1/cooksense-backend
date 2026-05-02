# CLAUDE.md

This file is the persistent operating manual for Claude Code when working in
**cooksense-backend**. Anything written here is loaded into every session as part of
Claude's working memory — keep it **short, specific, and high-signal**.
Prefer linking to `SPECIFICATIONS.md` / `STEERING.md` over duplicating rules.

> Philosophy: this codebase is written as if **Robert C. Martin** were
> reviewing every PR. We practice **Clean Code** (small, intent-revealing
> functions; no dead code; tests are the specification executed) and
> **Clean Architecture** (strict dependency direction, ports & adapters,
> business logic independent of databases, external frameworks, or transport layers).

---

## Always Do First — every session, no exceptions

1. **Read `SPECIFICATIONS.md` before writing or editing any production code.**
   The spec is the source of truth. Code that contradicts the spec is a bug;
   a spec that contradicts reality must be updated *first*, then the code
   changed to match.
2. **Read `STEERING.md` (or Appendix C of `SPECIFICATIONS.md`)** for the
   Always / Never / Architecture rules. They apply to every file you touch.
3. **Check `docs/features/in-progress/`** — if a spec file exists for the
   feature the user is asking about, resume from it; do not start over.
4. **Never generate code before you have located the relevant `SPEC-*` IDs**
   that govern the change. If none exist, you are in **Spec mode** (below).
5. **Every requirement shall use RFC-2119 language** — `shall` / `must` for
   mandatory, `should` for strong recommendation, `may` for optional.
   Prose without `shall` keywords is not a specification; reject it or
   rewrite it before treating it as binding.

---

## The SDD contract (read once, apply always)

Spec-Driven Development in this repo means **three non-negotiable loops**:

1. **Requirements → Design → Tasks.** A spec is not done when the business
   intent is captured; it is done when it lists `SPEC-*` IDs (§5–13), a
   design decision per ID, and an ordered task list (Appendix B of
   `SPECIFICATIONS.md`). Skipping the design or tasks step is a process bug.
2. **Traceability in both directions.** Every `SPEC-*` ID links forward to
   at least one test; every non-trivial function, test, and commit links
   back to one or more `SPEC-*` IDs. If a change has no `SPEC-*` reference,
   it is either a pure refactor (must say so) or a spec gap (stop and fix
   the spec first).
3. **Spec and code ship together.** If behavior changes, `SPECIFICATIONS.md`
   is edited in the **same PR** as the code. A code PR that contradicts the
   current spec without also amending it is rejected on sight.

---

## Skills available

| Skill                                                                 | When to invoke                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
|-----------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [`sdd-product-owner`](.claude/skills/sdd-product-owner/SKILL.md)      | User asks to act as a Product Owner, write/refine user stories, groom or prioritize backlog, slice epics, define acceptance criteria, or produce Definition of Ready/Done. The skill produces SDD-compliant artifacts (SPEC-*, RFC-2119, tasks table) — it does **not** write code and does **not** draft detailed `SPECIFICATIONS.md` contracts. When the story passes DoR, hand off to `sdd-spec-author`.                                                        |
| [`sdd-spec-author`](.claude/skills/sdd-spec-author/SKILL.md)          | User (typically a senior/staff engineer) asks to draft a `SPECIFICATIONS.md` section, expand SPEC-IDs into full contracts, write an ADR, or promote a DoR-passed story from `docs/features/in-progress/` into merged-truth spec. Produces RFC-2119 blocks with interfaces, error taxonomy, observability, security, test names, and a traceability matrix — it does **not** write code. When the spec is checklist-green, hand off to `sdd-design-document`.                  |
| [`sdd-design-document`](.claude/skills/sdd-design-document/SKILL.md)  | User asks to turn a checklist-green section of `SPECIFICATIONS.md` into a merge-ready `DESIGN-{AREA}-{NNN}.md` under `specifications/design/` — i.e. the concrete *how* (file layout, near-final artefact shape, integration call sites, observability lines, named tests, rollout, alternatives) for an Appendix B task. Produces a build artefact that maps every `SPEC-*` ID forward to ≥ 1 named test and backward to its user story; closed/archived once `T-{N}` merges. Does **not** write code, edit `SPECIFICATIONS.md`, or author ADRs. When checklist-green, hand off to the engineer flow below.        |

**Chain of ownership:** `sdd-product-owner` → `sdd-spec-author` → `sdd-design-document` → engineer flow (§"Feature requests — Spec mode first"). Each skill has a distinct input (raw intent / DoR-green story / checklist-green spec / checklist-green design) and a distinct output (story / spec section / design document / code). Do not collapse them.

---

## ⚠️ Feature requests — Spec mode first (MANDATORY)

When the user describes a feature idea or says *"start a new service"*,
*"add X"*, *"I want to build Y"*:

1. **Do NOT read code or scaffold files yet.**
2. **Do NOT create branches or files yet.**
3. Check `docs/features/in-progress/<feature-name>.md` — if present, read it
   and resume from where the previous session stopped.
4. If no spec exists, enter **Spec mode**:
   - Use `SPECIFICATIONS_TEMPLATE.md` (or §5 of `sdd-spec-author`) as the skeleton.
   - Ask clarifying questions **one group at a time**:
     a. *Business intent* — what data, what outcome, key stakeholders.
     b. *Domain model* — entities, value objects, invariants.
     c. *Logic & Rules* — business rules, validations, error cases.
     d. *Non-functional* — performance, security, observability, persistence.
     e. *Acceptance criteria* — how we prove the feature is correct in tests.
   - Produce a structured spec in **English**, with `SPEC-*` IDs for every
     requirement. Every requirement **shall** use RFC-2119 keywords (`shall` / `must` /
     `should` / `may`). Bulleted wishes without a keyword are not requirements.
   - Add a **Design** sub-section per module/component explaining *how* each
     `SPEC-*` will be satisfied (pattern chosen, collaborators, data shape).
   - Add a **Tasks** table (Appendix B of `SPECIFICATIONS.md` format):
     `Task | SPEC-IDs | Description | Dependencies`. No task without IDs.
   - Ask: *"Does this spec look correct? Should I proceed to implementation?"*
   - **Wait for explicit user approval before writing any code.**
5. After approval, save the spec to `docs/features/in-progress/<feature-name>.md`
   and tell the user they can close the session and resume later with:
   *"continue the feature from docs/features/in-progress/<feature-name>.md"*.
6. Only then: create a branch (`feat/<feature-name>`) and start implementation.
7. When the feature merges, move the file to `docs/features/done/`.

**Language rule:** when acting as the Spec Agent, always write the spec in
**English**, even if the user wrote in French or Portuguese.

---

## Project overview

**cooksense-backend** is a Go-based backend service for the CookSense application. It follows a clean architecture approach to provide a robust and maintainable codebase.

**Prerequisites:** Go ≥ 1.26, `make`, `golangci-lint`.

---

## Build, test, and run commands

Prefer `make` targets — they are the contract and match CI.

| Task                               | Command                                             |
|------------------------------------|-----------------------------------------------------|
| Install deps                       | `go mod download`                                   |
| Format                             | `go fmt ./...`                                      |
| Lint                               | `golangci-lint run`                                 |
| Full suite + coverage ≥ 80%        | `make test` (or `go test -cover ./...`)             |
| Build                              | `go build -o cooksense-server ./cmd/cooksense-server` |
| Run locally                        | `./cooksense-server`                                |

### Minimum validation checklist (run before every PR)

1. `go fmt ./...` — all files formatted.
2. `golangci-lint run` — zero violations.
3. `go test -v -cover ./...` — all tests green, coverage ≥ 80%.

**No `//nolint` without a specific linter name and justification, no `TODO`, no `fmt.Print*` for logging.**
If a suppression is unavoidable, it must carry a bracketed code and a
justification comment, per `SPECIFICATIONS.md` §0.3.

---

## Architecture — Clean Architecture in Go

```
cmd/
└── cooksense-server/
    └── main.go              # Entry point
internal/
├── domain/                  # Business logic (entities, interfaces)
├── service/                 # Use cases / orchestration
├── infra/                   # Adapters (DB, external APIs)
└── api/                     # HTTP/gRPC handlers
pkg/                         # Publicly shareable code
specifications/
└── design/                  # DESIGN-*.md files
docs/
└── features/                # SPEC-* in-progress and done
tests/                       # Integration and E2E tests
```

### Dependency flow (strict, no exceptions)

1. `api` -> `service` -> `domain`
2. `infra` -> `domain` (satisfying interfaces)
3. `cmd` -> all (for wiring)

Circular imports are **strictly forbidden** in Go and will result in compilation errors. See `SPECIFICATIONS.md` §4.3.

### Layers

| Layer | Who lives here | Rule |
|-------|----------------|------|
| **Presentation** | `api/`, `cmd/` | No business logic. Parse requests, delegate. |
| **Orchestration** | `service/` | Wires domain logic + infra. No low-level I/O details. |
| **Domain** | `domain/` | **Business rules live here.** No I/O, no external dependencies. |
| **Infrastructure** | `infra/` | External adapters. Satisfies domain interfaces. |

---

## Clean Code decomposition rules

- Functions **do one thing**. If you need the word *"and"* to name it, split it.
- Target: ≤ 30 lines per function, ≤ 10 methods per struct, ≤ 500 lines per module.
- **Intent-revealing names** — `NewUserRepository`, `ProcessPayment`, `ValidateOrder`. Never `HandleData`, `DoStuff`, `Helper`.
- Comments explain **why**, code explains **what**. Delete any comment that paraphrases its line.
- Business logic stays in the `domain` or `service` layers.
- Reject "god" packages. If a package grows too large, split it (e.g., `service/user`, `service/billing`).

---

## Go Idioms & Principles (MANDATORY)

This codebase follows **modern Go (1.26)** idioms.

### Go Idioms — what "Idiomatic Go" means here

- **Explicit error handling.** Always check `if err != nil`. Never ignore errors with `_`. Use `errors.Is` and `errors.As` for error inspection.
- **Interfaces for decoupling.** Keep interfaces small (often 1-3 methods). Interfaces should be defined by the consumer, not the producer, where possible.
- **`context.Context` for all I/O.** Every function that performs I/O or long-running work must accept `context.Context` as its first parameter.
- **Structs and composition.** Prefer composition over embedding. Avoid deep embedding hierarchies.
- **PascalCase for exported symbols.** PascalCase for public types, functions, and fields; camelCase for private ones.
- **Go Doc comments.** Use standard Go documentation comments on all exported symbols.
- **Logging with `log/slog`.** Use structured logging. Never log secrets.
- **Zero value is useful.** Design types such that their zero value is meaningful or safely unusable.
- **Avoid global state.** No package-level mutable variables. Inject dependencies via constructors.
- **Modern Go 1.26 features.** Use `slices`, `maps`, and `iter` packages. Use `for i := range n`. Use `new(val)` for pointers to literals.
- **Concurrency.** Use goroutines and channels judiciously. Always ensure goroutines have a clear lifecycle and can be cancelled via `context.Context`.

### Code review red flags (auto-reject)

- Ignoring errors (`_ = ...`) or missing `if err != nil`.
- Dot imports (`import . "package"`).
- Global mutable state.
- Hard-coded secrets, URLs, or file paths.
- `time.Sleep()` without a way to cancel (use `Context` or `Timer` with `select`).
- `Any` (interface{}) used without a strong justification.
- Large packages or files (> 500 lines).
- Lack of tests for new functionality.

---

## Service development workflow

Every new feature or service follows the Open-Closed Principle: **add files, do not edit existing ones** unless necessary.

1. **Spec** — add a section under `SPECIFICATIONS.md` §5 or create
   `docs/features/in-progress/<name>.md` with `SPEC-ID-<NNN>` IDs for:
   business intent, domain entities, rules, error handling, observability.
2. **Implement** logic in `internal/domain/` and `internal/service/`:
   - Define interfaces in `domain`.
   - Implement services that satisfy those interfaces.
   - Keep domain logic pure and free of I/O.
3. **Infrastructure** in `internal/infra/`:
   - Implement repository interfaces using databases (e.g., PostgreSQL, Redis).
   - Implement external API clients.
4. **API** in `internal/api/`:
   - Define HTTP handlers or gRPC services.
   - Wire dependencies in `cmd/cooksense-server/main.go`.
5. **Tests** `internal/.../*_test.go`:
   - Unit tests for domain logic.
   - Integration tests for infrastructure and services.
   - Naming: `Test{What}_{Condition}_{ExpectedOutcome}`.
6. **Docs** — update the README and API documentation.

---

## Configuration & environments

- Environment-specific values (DB connection, secrets, port) live in environment variables or a configuration file (not committed).
- The application should support loading configuration via a `Config` struct.
- Secrets come from secret management systems (e.g., Vault, AWS Secrets Manager, or environment variables). Never from committed files. See `SPECIFICATIONS.md` §12 for the full security spec.

---

## Testing policy (MANDATORY)

From `SPECIFICATIONS.md` §15 — the short version:

- **Every `SPEC-*` ID must have at least one test.** No test ⇒ the
  requirement does not exist.
- Test names are English sentences: `Test{What}_{Condition}_{ExpectedOutcome}`.
- Use `t.Context()` when a test function needs a context.
- Use `t.Helper()` in test helper functions.
- **Do mock**: External HTTP calls, databases (via interfaces), time-dependent operations.
- **Do not mock**: Pure business logic, struct construction.
- Coverage floor: **80% line coverage**, enforced by `make test` and CI.
  New code **should** be ≥ 90%.
- If you find yourself mocking more than 3 things in a test, the design is
  wrong — refactor the production code first.

---

## Observability rules

- Every package: use `slog` for logging.
- Use structured context with `slog.Attr` or key-value pairs.
- Levels: `DEBUG` = internal tracing, `INFO` = key lifecycle events, `WARN` = recoverable issues, `ERROR` = failures requiring attention.
- **Secrets never appear in a log message at any level.** Verify by inspection before merging.

---

## Branching & PR workflow

1. One focused PR per feature. Branch from latest `main`:
   `feat/<slug>`, `fix/<slug>`, `refactor/<slug>`, `docs/<slug>`.
2. Commits follow **Conventional Commits** via `commitizen`.
   The commit **body** (not the subject) **shall** cite every `SPEC-*` ID
   the change implements, fixes, or deprecates. A commit without a
   `SPEC-*` reference **must** include `Refs: refactor-only` or
   `Refs: chore-only` as the last line — anything else is rejected.
3. Before pushing: run the **Minimum validation checklist** above.
4. After opening the PR: request review. Fix all review comments, re-run
   validation, push again until the PR is clean.
5. **Always rebase onto `main` before merging.**
6. On merge: `cz bump` happens in the release workflow — do not bump
   versions manually.

---

## Non-obvious caveats

- Circular imports are a hard failure in Go. Design your package hierarchy carefully.
- Use `internal/` to hide implementation details from other modules.
- `context.Context` should never be stored in a struct; pass it as the first argument to functions.
- Avoid using `init()` for anything other than simple registrations.
- `go mod tidy` should be run after adding or removing dependencies.

---

## Review checklist (for Claude and humans)

Before you claim a task is done, verify:

- [ ] `golangci-lint run` passes.
- [ ] `go test ./...` passes — coverage ≥ 80%, new code ≥ 90%.
- [ ] Every new `SPEC-*` ID has at least one test.
- [ ] Every new or changed test **names or docstrings** the `SPEC-*` IDs it
      verifies (forward traceability).
- [ ] Every non-trivial new function **docstring** cites the `SPEC-*` IDs it
      implements (backward traceability).
- [ ] Every new public symbol has a Go Doc comment and follows PascalCase.
- [ ] No `fmt.Print*`, no `TODO`, no `FIXME`, no dead code, no `time.Sleep` without context.
- [ ] Go Idioms respected: explicit error handling, interfaces for decoupling, `context.Context` for I/O, `slog` for logging, no global mutable state.
- [ ] No secrets or hard-coded environment values in source.
- [ ] `SPECIFICATIONS.md` updated if the public behavior changed.
- [ ] README updated if the public API or config changed.
- [ ] Commit messages follow Conventional Commits and cite `SPEC-*` IDs.

If any box is unchecked, the work is not done — fix it before handing back
to the user.
