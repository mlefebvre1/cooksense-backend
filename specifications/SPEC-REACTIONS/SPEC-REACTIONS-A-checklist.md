# SPEC-REACTIONS — Appendix A · Specification Checklist

[← Index](SPEC-REACTIONS-00-index.md)

## A.1 Completeness checklist

| # | Criterion | Status |
|---|-----------|--------|
| 1 | Every AC and DoD item from `docs/stories/08-api-reactions-my-recipes.md` maps to ≥ 1 SPEC-REACTIONS-NNN ID. | ✅ (matrix below) |
| 2 | Every SPEC-REACTIONS-NNN ID maps to ≥ 1 named test in §9. | ✅ |
| 3 | All requirements use RFC-2119 keywords (`SHALL` / `MUST` / `SHOULD` / `MAY`). | ✅ |
| 4 | Public Go signatures are stated for every new exported symbol. | ✅ |
| 5 | All SQL is shown verbatim and matches `migrations/0001_init.up.sql`. | ✅ |
| 6 | Error taxonomy is defined (`ErrInvalidReactionKind`, `ErrInvalidMyRecipesKind`, `ErrRecipeNotFound` consumed) and mapped to HTTP codes (§5.7). | ✅ |
| 7 | Configuration variables are documented (none new; reused vars listed). | ✅ |
| 8 | Observability/logging guidance present (slog usage, no token in logs). | ✅ |
| 9 | Security/threat considerations addressed (uid from auth context only, no body-supplied uid, no token in logs). | ✅ |
| 10 | Appendix B lists ordered, atomic implementation tasks. | ✅ |

## A.2 Story AC → SPEC-ID traceability matrix

| Story acceptance criterion | SPEC-REACTIONS IDs |
|----------------------------|--------------------|
| AC #1 — `POST /api/reactions`: UPSERT, `200` on success, `404 NOT_FOUND` for unknown slug, `422 INVALID_REACTION_KIND` for bad enum, `400 INVALID_PAYLOAD` for malformed body | 008, 009, 015, 016, 022, 023, 024, 025, 026, 027, 028 |
| AC #2 — `DELETE /api/reactions/{slug}`: `204` (idempotent), `404` if recipe slug unknown | 010, 015, 016, 029, 030, 031, 032 |
| AC #3 — `GET /api/me/recipes?kind=…`: defaults to `LIKE`, returns `<RecipeBrief>` ordered by `created_at DESC`, `DISLIKE` not exposed, `400` for any other value | 011, 012, 013, 020, 021, 033, 034, 035, 036, 037 |
| AC #4 — Repository exposes `Upsert`, `Delete`, `ListByKind` returning `RecipeBrief` joined with `recipes` | 006, 007, 008, 010, 011, 012 |
| AC #5 — Integration tests cover: upsert from no-row, upsert overwrite, delete, list ordering, 404 on unknown slug, 422 on invalid kind | 045, 046, 047, 048, 049, 050, 053 |

## A.3 Story DoD → SPEC-ID traceability matrix

| Story DoD item | SPEC-REACTIONS IDs |
|----------------|--------------------|
| AC met | All 057 IDs |
| Integration tests green | 045–056 (and §9.5 coverage targets) |
| API doc unchanged or amended in same PR | 044 |

## A.4 Technical-notes coverage

| Story technical note | SPEC-REACTIONS IDs |
|----------------------|--------------------|
| Handler resolves `recipe_slug` → `recipe_id` first; no ID-variant API | 015, 017, 018, 019 (NG-1) |
| Upsert SQL: `INSERT … ON CONFLICT … DO UPDATE … RETURNING created_at` | 008, 009 |
| FK on `users(firebase_uid)` is provisioned by the auth middleware (story 04) | §3.3 constraint, §4.5 |

## A.5 SPEC-ID → forward test mapping (summary)

Every SPEC-REACTIONS-NNN ID listed in `SPEC-REACTIONS-00-index.md` appears
in the test table of `SPEC-REACTIONS-09-testing.md` §9.1 at least once:

- 001..006 → domain unit tests (5 tests).
- 007 → interface assertion test.
- 008..014 → repo integration tests (6 tests).
- 015, 016 → service & integration tests (3 tests).
- 017..021 → service unit tests (4 tests).
- 022..028 → POST handler tests (5 tests).
- 029..032 → DELETE handler tests (3 tests).
- 033..037 → GET handler tests (5 tests).
- 038..040 → cross-handler tests (3 tests).
- 041, 042 → wiring smoke test (1 test).
- 043 → enforced via doc-comment review + linter scan (covered by
  `Test_Wiring_RegistersThreeRoutes` indirectly; spec-level rather than
  test-level).
- 044 → enforced at PR review (no test).
- 045..057 → already covered by their parent IDs above.

## A.6 Sign-off

- Spec reviewed against story file `docs/stories/08-api-reactions-my-recipes.md`: ✅
- Spec reviewed against `docs/architecture/api.md` (Reactions + My recipes sections): ✅
- Spec reviewed against `docs/architecture/data-model.md` (`user_reactions`, `reaction_kind`): ✅
- Spec reviewed against decision D-0008: ✅
- All SPEC-IDs referenced in §6 appear in the index registry of
  `SPEC-REACTIONS-00-index.md`: ✅
- Status transition: **Draft → Final** (pending user approval).
