# SPEC-MAKE-A — Specification Checklist

> Part of [SPEC-MAKE](SPEC-MAKE-00-index.md) — Story 02: Makefile Targets

---

## Appendix A — Specification Checklist

Use this checklist before starting implementation. Every box must be checked.

- [x] **Section 0** — AI Preamble reviewed; Makefile-specific style and anti-patterns defined.
- [x] **Section 1** — Introduction written; scope and cross-story dependencies listed.
- [x] **Section 2** — Every goal is verifiable; non-goals (CI, Windows, pinned versions) explicit.
- [x] **Section 3** — External tooling listed with version constraints; files consumed/produced documented.
- [x] **Section 4** — Makefile structure defined; `.env` loading strategy shown; dependency graph explicit; 80-line budget stated.
- [x] **Section 5** — All 17 SPEC-MAKE-NNN IDs defined (001–017); every target's contract (command, exit code, side-effects) fully specified; `.env.example` content specified.
- [x] **Section 6** — Makefile variables named; `.env` forwarding strategy documented.
- [x] **Section 7** — Verification matrix covers the 3 primary build/lint checks.
- [x] **Section 8** — Verification table covers all 10 verification scenarios; manual vs automated distinction made.
- [x] **Section 9** — `##` comment format standard defined; `.env.example` documentation standard specified.
