# SPEC-AUTH-A — Specification Completeness Checklist

← [Index](SPEC-AUTH-00-index.md)

---

This checklist must be fully green before the spec is marked **Final** and
implementation begins. A single unchecked item blocks the spec from progressing.

| # | Criterion | Status |
|---|-----------|--------|
| 1 | Every acceptance criterion from Story 04 maps to ≥ 1 SPEC-AUTH-NNN ID | ✅ |
| 2 | Every SPEC-AUTH-NNN ID uses RFC-2119 keywords (`shall` / `must` / `should` / `may`) | ✅ |
| 3 | Every SPEC-AUTH-NNN ID maps to ≥ 1 named test in §9 | ✅ |
| 4 | Forward traceability: every test name cites ≥ 1 SPEC-ID | ✅ (see §9 table) |
| 5 | Backward traceability: every exported function signature cites its SPEC-IDs in the doc comment | ✅ (see §10 table) |
| 6 | All Go function signatures are present and complete in §6 | ✅ |
| 7 | SQL query in §6.6 matches `docs/architecture/auth.md` exactly | ✅ |
| 8 | No circular imports introduced (auth → users interface, not concrete) | ✅ (see §5 dependency graph) |
| 9 | Security: credentials never in source, `.gitignore` rule referenced | ✅ (SPEC-AUTH-024, §7.2) |
| 10 | Non-goals and out-of-scope items explicitly listed | ✅ (§3) |

**Spec status: ✅ FINAL**

---

## Acceptance criteria → SPEC-ID traceability matrix

| Story 04 AC | SPEC-IDs |
|-------------|----------|
| `firebase.go` initializes `*firebase.App` and `*auth.Client` once | SPEC-AUTH-001, SPEC-AUTH-002, SPEC-AUTH-003 |
| `auth.Verifier` interface with `Verify` method + two implementations | SPEC-AUTH-004, SPEC-AUTH-005, SPEC-AUTH-006 |
| `auth.User` type matches `docs/architecture/auth.md` | SPEC-AUTH-007 |
| `auth.Middleware` — missing header → 401 | SPEC-AUTH-008, SPEC-AUTH-009, SPEC-AUTH-012, SPEC-AUTH-023 |
| `auth.Middleware` — malformed header → 401 | SPEC-AUTH-010, SPEC-AUTH-012 |
| `auth.Middleware` — verifier error → 401 | SPEC-AUTH-011, SPEC-AUTH-012 |
| `auth.Middleware` — success → Touch + store in context | SPEC-AUTH-013, SPEC-AUTH-014, SPEC-AUTH-015 |
| `auth.UserFromContext` is the only retrieval path | SPEC-AUTH-016 |
| `internal/users/repo.go` implements Toucher against Postgres | SPEC-AUTH-017, SPEC-AUTH-018, SPEC-AUTH-019, SPEC-AUTH-020 |
| `GET /api/health` public → `{"status":"ok"}` | SPEC-AUTH-021 |
| `GET /api/health/me` auth-gated → `{"uid":…,"email":…}` | SPEC-AUTH-022 |
| Unit tests use `fakeVerifier`, no network calls | SPEC-AUTH-026, SPEC-AUTH-027 |
| DoD: `secrets/firebase-admin.json.example` committed | SPEC-AUTH-024 |
| DoD: `.env.example` has Firebase vars | SPEC-AUTH-025 |
