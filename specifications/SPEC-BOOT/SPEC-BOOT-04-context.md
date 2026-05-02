# SPEC-BOOT-04 — System Context & Dependencies

> Part of [SPEC-BOOT](SPEC-BOOT-00-index.md) — Story 01: Bootstrap Project Structure

---

## 3. System Context & Dependencies

### 3.1 Runtime Requirements

| Requirement | Specification |
|-------------|--------------|
| **Go** | `1.26.2` (exact, matching `go.mod` declaration) |
| **OS** | Linux, macOS |

### 3.2 Package Dependencies

No new third-party dependencies are introduced in Story 01.

The `go.mod` **shall** declare the module path `github.com/cooksense/cooksense-backend` (or the existing module path already in `go.mod` — do not change it) and the Go version `1.26.2`. No `require` block additions are needed.

### 3.3 External Systems & APIs

None. Story 01 is purely structural — no network calls, no I/O beyond stdout.
