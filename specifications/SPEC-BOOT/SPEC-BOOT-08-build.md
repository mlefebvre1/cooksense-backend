# SPEC-BOOT-08 — Build, Tooling & Quality Specification

> Part of [SPEC-BOOT](SPEC-BOOT-00-index.md) — Story 01: Bootstrap Project Structure  
> SPEC-IDs covered: SPEC-BOOT-019, SPEC-BOOT-020

---

## 7. Build, Tooling & Quality Specification

### 7.1 Build Verification

| Requirement | Command | Expected outcome |
|-------------|---------|-----------------|
| **SPEC-BOOT-019** | `go build ./...` | Exits `0`. No errors, no warnings. |
| **SPEC-BOOT-020** | `go vet ./...` | Exits `0`. Zero issues reported. |

### 7.2 Go Version

| Setting | Value |
|---------|-------|
| **Tool** | `go` |
| **Required version** | `1.26.2` |
| **Declared in** | `go.mod` (`go 1.26.2`) |

### 7.3 Linting

`golangci-lint run` **should** produce zero violations. Because Story 01 contains only `doc.go` placeholders and a trivial `main.go`, lint violations at this stage indicate a tooling misconfiguration that **shall** be fixed before marking the story done.

### 7.4 Formatting

All `.go` files **shall** be formatted by `go fmt ./...` with zero diffs. CI **shall** fail if any file would change.
