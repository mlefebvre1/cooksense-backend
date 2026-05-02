# SPEC-DISCOVER — Appendix A · Specification Checklist

[← Index](SPEC-DISCOVER-00-index.md)

## A.1 Completeness checklist

| # | Criterion | Status |
|---|-----------|--------|
| 1 | Every AC and DoD item from `docs/stories/07-api-recipes-discover.md` maps to ≥ 1 SPEC-DISCOVER-NNN ID. | ✅ (matrix below) |
| 2 | Every SPEC-DISCOVER-NNN ID maps to ≥ 1 named test in §9. | ✅ |
| 3 | All requirements use RFC-2119 keywords (`SHALL` / `MUST` / `SHOULD` / `MAY`). | ✅ |
| 4 | Public Go signatures are stated for every new exported symbol. | ✅ |
| 5 | All SQL is shown verbatim (Discover anti-join + Detail two-query) and matches `docs/architecture/data-model.md`. | ✅ |
| 6 | Error taxonomy is defined (`ErrNotFound`, `INTERNAL`, `NOT_FOUND` codes). | ✅ |
| 7 | Configuration variables are documented (none new; reused vars listed). | ✅ |
| 8 | Observability/logging guidance present (`slog.Info` per request; `ERROR` on failures). | ✅ |
| 9 | Security/threat considerations addressed (auth required; uid only logged, no email/displayName). | ✅ |
| 10 | Appendix B lists ordered, atomic implementation tasks. | ✅ |

## A.2 Story AC → SPEC-ID traceability matrix

| Story acceptance criterion | SPEC-DISCOVER IDs |
|----------------------------|-------------------|
| AC #1 — `GET /api/recipes/discover?limit=N` (auth): up to N (default 10, max 25), excludes any-reaction, randomized, `<RecipeBrief>` shape | 002, 005, 009, 010, 011, 012, 018, 019, 024, 025, 027, 032 |
| AC #2 — `GET /api/recipes/{slug}` (auth): `<RecipeFull>` with sorted ingredients and ordered steps; 404 unknown slug | 003, 004, 013, 014, 015, 020, 028, 029, 030, 032 |
| AC #3 — Repo (`internal/recipes/repo.go`) is the only DB-touching layer with `Discover`, `GetBySlug` | 001, 006, 007, 008, 016 |
| AC #4 — Service (`internal/recipes/service.go`) clamps limit and orchestrates | 017, 018, 019, 020, 021, 022 |
| AC #5 — Handler (`internal/recipes/handler.go`) only does JSON / status mapping | 023, 024, 025, 026, 027, 028, 029, 030, 031, 032, 033 |
| AC #6 — Integration test (compose Postgres + seeded data + `fakeVerifier`) covers Discover count, Discover excludes reacted, Detail returns seeded recipe, Detail returns 404 | 036, 037, 038, 039 |

## A.3 Story DoD → SPEC-ID traceability matrix

| Story DoD item | SPEC-DISCOVER IDs |
|----------------|-------------------|
| AC met | All 042 IDs |
| Integration test green against compose Postgres | 036, 037, 038, 039 + §9.2/9.3 |
| `go vet ./...` clean | §8.1 |

## A.4 Technical-notes coverage

| Story technical note | SPEC-DISCOVER IDs |
|----------------------|-------------------|
| SQL sketch with `NOT EXISTS` against `user_reactions` | 009, 011 |
| Index `user_reactions(firebase_uid, …)` keeps anti-join cheap | §5.4 (references `user_reactions_uid_kind_idx` from SPEC-DB DDL) |
| Ingredient hydration uses one extra query (or `LEFT JOIN LATERAL`); pick simpler | §3.4, 013 |
| JSON field names are `snake_case`, matching API doc exactly | 002, 003, 004 |

## A.5 Out-of-scope items (story §61–65)

| Story note | Reflected in this spec |
|------------|------------------------|
| Reactions endpoint owned by story 08 | §3.2 |
| Search owned by story 09 | §3.2 |
| "Time-to-decision" telemetry post-MVP | §3.2 |

## A.6 Sign-off

- Spec reviewed against story file: ✅
- Spec reviewed against `docs/architecture/api.md`: ✅
- Spec reviewed against `docs/architecture/data-model.md`: ✅
- All SPEC-IDs referenced in §6 appear in the index registry of
  `SPEC-DISCOVER-00-index.md`: ✅
- Status transition: **Draft → Final**.
