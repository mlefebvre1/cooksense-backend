# SPEC-MAKE-01 — AI Steering Preamble

> Part of [SPEC-MAKE](SPEC-MAKE-00-index.md) — Story 02: Makefile Targets

---

## 0. AI Steering Preamble

### 0.1 AI Persona & Quality Bar

You are a **Staff Software Engineer** implementing this specification. The artefacts produced shall:

- Be **production-grade** — no TODOs, no placeholder logic except where this spec explicitly permits it.
- Read like a **well-edited technical book** — concise targets, intent-revealing names, no incantations.
- Demonstrate **mastery of POSIX make and shell idioms** (portable to GNU Make on Linux and macOS).
- Treat every public target as a **published API** — its name, exit code, and side-effects are stable.

### 0.2 Language Conventions (RFC-2119)

| Keyword | Meaning |
|---------|---------|
| **shall** / **must** | Absolute requirement. A verification step **must** confirm compliance. |
| **shall not** / **must not** | Absolute prohibition. |
| **should** | Strong recommendation. Deviation requires written justification. |
| **may** | Truly optional. |

### 0.3 Code Style Mandate (Makefile-specific)

| Rule | Requirement |
|------|-------------|
| **`.PHONY`** | Every non-file target **shall** be declared in a `.PHONY` line. |
| **Target docs** | Every public target **shall** carry a trailing `## <description>` comment that the `help` target parses. |
| **Tabs** | Recipe lines use a leading TAB, not spaces. |
| **Shell** | Recipes assume POSIX `sh`. macOS / Linux only. Windows is out of scope. |
| **Quoting** | Variable expansions used in shell **shall** be quoted: `"$(VAR)"`, never bare `$(VAR)`. |
| **Magic strings** | Hard-coded paths (`bin/`, `cmd/cooksense-server`) are acceptable in `Makefile` — they are part of the project layout fixed by SPEC-BOOT. |
| **Logging** | `make` recipes **shall not** log secrets. `.env` values are forwarded to subprocesses but never echoed. |

### 0.4 Forbidden Anti-Patterns

| Anti-Pattern | Why It's Forbidden |
|--------------|--------------------|
| Recipes longer than 5 lines of inline shell | Move to a script under `scripts/`. |
| Hidden side-effects (network calls, file deletion outside the target's contract) | Targets are advertised by their name; surprises break trust. |
| `make` calling `make` recursively without need | Slows builds and obscures dependency tracking. |
| Silent failure (recipe succeeds when subcommand fails) | Each step **must** propagate non-zero exit codes. |
| Real secret values committed to `.env.example` | Use placeholders only. |
