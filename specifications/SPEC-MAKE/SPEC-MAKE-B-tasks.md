# SPEC-MAKE-B — Implementation Task Decomposition

> Part of [SPEC-MAKE](SPEC-MAKE-00-index.md) — Story 02: Makefile Targets

---

## Appendix B — Implementation Task Decomposition

Ordered list of atomic tasks for Story 02. Each task **shall** be verified against the verification table in §8.2 before marking the story Done.

| Task | SPEC-IDs | Description | Dependencies |
|------|----------|-------------|--------------|
| T-01 | SPEC-MAKE-014 | Add `-include .env` guard at top of `Makefile` | None |
| T-02 | SPEC-MAKE-016, SPEC-MAKE-017 | Update `.env.example` with all 10 variables + inline `#` comments | None |
| T-03 | SPEC-MAKE-001 | Verify `Makefile` exists at repo root (or create it) | None |
| T-04 | SPEC-MAKE-013 | Add `.PHONY` line listing all 10 public targets | T-01, T-03 |
| T-05 | SPEC-MAKE-002 | Implement `help` as default target (parses `##` comments) | T-04 |
| T-06 | SPEC-MAKE-003 | Add variables `BINARY`, `CMD_DIR`, `BIN_DIR` | T-04 |
| T-07 | SPEC-MAKE-004, SPEC-MAKE-005 | Implement `up` target with healthcheck wait loop | T-04 |
| T-08 | SPEC-MAKE-006 | Implement `down` target | T-04 |
| T-09 | SPEC-MAKE-012 | Implement `migrate` target (`go run ... migrate up`) | T-07 |
| T-10 | SPEC-MAKE-012 | Implement `seed` target (`go run ... seed`) | T-07 |
| T-11 | SPEC-MAKE-008 | Implement `run` target (modd) with `up migrate` prerequisites | T-09, T-10 |
| T-12 | SPEC-MAKE-007 | Implement `build` target (`go build -o bin/...`) | T-04, T-06 |
| T-13 | SPEC-MAKE-009 | Implement `test` target (`go test ./...`) | T-04 |
| T-14 | SPEC-MAKE-010 | Implement `lint` target (go vet + optional golangci-lint) | T-04 |
| T-15 | SPEC-MAKE-011 | Implement `clean` target (rm bin/ + optional volume) | T-04 |
| T-16 | SPEC-MAKE-015 | Count Makefile lines; refactor to `scripts/` if > 80 | T-05..T-15 |
| T-17 | SPEC-MAKE-018, SPEC-MAKE-019, SPEC-MAKE-020 | Run verification table: `make help`, `make build`, `make lint` | T-16 |

---

*End of SPEC-MAKE — Story 02 Makefile Targets — v1.0.0*
