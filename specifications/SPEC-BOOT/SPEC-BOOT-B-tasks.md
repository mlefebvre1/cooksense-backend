# SPEC-BOOT-B — Implementation Task Decomposition

> Part of [SPEC-BOOT](SPEC-BOOT-00-index.md) — Story 01: Bootstrap Project Structure

---

## Appendix B — Implementation Task Decomposition

Ordered list of atomic tasks for Story 01. Each task **shall** be implemented, reviewed, and passing `go build ./...` + `go vet ./...` before the story is marked Done.

| Task | SPEC-IDs | Description | Dependencies |
|------|----------|-------------|--------------|
| T-01 | SPEC-BOOT-017, SPEC-BOOT-018 | Update `go.mod` to declare `go 1.26.2` | None |
| T-02 | SPEC-BOOT-016 | Create or update `.gitignore` with all required patterns | None |
| T-03 | SPEC-BOOT-001, SPEC-BOOT-002 | Create `internal/config/doc.go` | T-01 |
| T-04 | SPEC-BOOT-003 | Create `internal/auth/doc.go` | T-01 |
| T-05 | SPEC-BOOT-004 | Create `internal/db/doc.go` | T-01 |
| T-06 | SPEC-BOOT-005, SPEC-BOOT-006 | Create `internal/domain/doc.go` | T-01 |
| T-07 | SPEC-BOOT-007 | Create `internal/recipes/doc.go` | T-01 |
| T-08 | SPEC-BOOT-008 | Create `internal/reactions/doc.go` | T-01 |
| T-09 | SPEC-BOOT-009 | Create `internal/lessons/doc.go` | T-01 |
| T-10 | SPEC-BOOT-010 | Create `internal/users/doc.go` | T-01 |
| T-11 | SPEC-BOOT-011 | Create `internal/seed/doc.go` | T-01 |
| T-12 | SPEC-BOOT-012, SPEC-BOOT-013 | Update `cmd/cooksense-server/main.go` with placeholder `main` | T-01 |
| T-13 | SPEC-BOOT-014 | Create `migrations/.gitkeep` | None |
| T-14 | SPEC-BOOT-015 | Create `seed/recipes/.gitkeep` and `seed/lessons/.gitkeep` | None |
| T-15 | SPEC-BOOT-019, SPEC-BOOT-020 | Verify `go build ./...` and `go vet ./...` pass | T-02..T-14 |
| T-16 | SPEC-BOOT-021, SPEC-BOOT-022, SPEC-BOOT-023 | Run smoke verification: `go run ./cmd/cooksense-server` outputs `cooksense-server starting` | T-15 |

---

*End of SPEC-BOOT — Story 01 Bootstrap Project Structure — v1.0.0*
