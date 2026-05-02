# SPEC-BOOT-09 — Testing Specification

> Part of [SPEC-BOOT](SPEC-BOOT-00-index.md) — Story 01: Bootstrap Project Structure  
> SPEC-IDs covered: SPEC-BOOT-021, SPEC-BOOT-022, SPEC-BOOT-023

---

## 8. Testing Specification

### 8.1 Test Philosophy

Story 01 is a purely structural story. There is no business logic, no I/O, and no state to unit-test. However, the following verifications **shall** be performed as part of the Definition of Done:

| Verification | How |
|-------------|-----|
| **SPEC-BOOT-021**: All packages compile | `go build ./...` exits `0` |
| **SPEC-BOOT-022**: `main` prints and exits | `go run ./cmd/cooksense-server` outputs `cooksense-server starting` and exits `0` |
| **SPEC-BOOT-023**: Vet passes | `go vet ./...` exits `0` |

No `_test.go` files are required in Story 01. Test infrastructure (httptest, testcontainers) is introduced in Story 11.

### 8.2 Test Naming Convention (for future stories)

```
Test{What}_{Condition}_{ExpectedOutcome}
```

### 8.3 Coverage Thresholds

Story 01 introduces no testable logic. The coverage floor of **80%** (new code **90%**) applies from Story 03 onward.
