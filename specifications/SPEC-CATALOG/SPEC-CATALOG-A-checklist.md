# SPEC-CATALOG §A — Specification Completeness Checklist

[← back to index](SPEC-CATALOG-00-index.md)

## A.1 Checklist (sdd-spec-author template)

| # | Item | Status |
|---|------|--------|
| 1 | Every story-06 acceptance criterion maps to ≥ 1 SPEC-CATALOG-NNN ID | ✅ (see §A.2) |
| 2 | Every SPEC-CATALOG-NNN ID has a verification artifact (test, reviewer step, or PR-description element) | ✅ (see §A.3) |
| 3 | RFC-2119 keywords (`shall`, `must`, `should`, `may`) used for every normative statement | ✅ |
| 4 | Scope boundary with SPEC-RECIPES is explicit and non-overlapping | ✅ §1.2 |
| 5 | All external dependencies enumerated; no implicit ones | ✅ §4.3 ("None") |
| 6 | Architecture diagram(s) present for non-trivial flows | ✅ §4.4, §5.3 |
| 7 | Configuration changes (env vars, files) listed exhaustively | ✅ §7 |
| 8 | Test strategy enumerates named tests with SPEC-ID coverage | ✅ §9.2 |
| 9 | Implementation tasks ordered with explicit dependencies | ✅ §B |
| 10 | Documentation impact (README, doc.go, PR templates) accounted for | ✅ §10 |

All 10 items green → spec **Final**.

---

## A.2 Acceptance criterion → SPEC-ID matrix

Source: `docs/stories/06-curated-recipes-content.md`.

| AC | Story-06 wording (paraphrased) | SPEC-IDs |
|----|----------------------------------|----------|
| AC-1 | `seed/recipes/` ≥ 15 YAML files; filenames match `slug` | 001, 002, 003, 004 |
| AC-2 | Each recipe valid: required fields, taxonomies, `time_minutes ≤ 30`, `passive_prep_minutes` allowed, ≥ 2 ingredients, ≥ 1 step | 005, 006, 007, 008, 009, 010 |
| AC-3 | Catalog diversity: ≥ 4 methods, ≥ 3 proteins, ≥ 2 vegetarian, ≥ 1 slow/pressure | 011, 012, 013, 014 |
| AC-4 | Every recipe read end-to-end by a reviewer | 024 |
| AC-5 | `make seed` loads all recipes without errors | 026 |

Technical notes (story-06):

| Note | SPEC-IDs |
|------|----------|
| Tone: terse, expert; no backstory; no "the secret is…" | 019, 020 |
| Steps describe what + why when non-obvious | 021 |
| Reuse ingredient names verbatim across recipes | 015 |
| Quantity/unit live on recipe-ingredient row, not on name | 016 |
| Aliases belong in story 09 | 018 |

DoD (story-06):

| DoD | SPEC-IDs |
|-----|----------|
| AC met | All |
| Catalog reviewed in single PR; reviewer comment per recipe | 024 |
| PR description summary table | 025 |

Out-of-scope reinforcements:

| OoS | SPEC-IDs |
|-----|----------|
| Photos | 023 |
| Translations (English only) | 022 |
| User-submitted content | (handled by NG1, no SPEC-ID; story 06 simply does not include UGC) |

---

## A.3 SPEC-ID → verification artifact

| SPEC-ID | Verification |
|---------|--------------|
| 001, 004 | Test #1 |
| 002, 005, 006 | Test #3 |
| 003 | Test #2 |
| 007 | Test #4 |
| 008 | Test #5 |
| 009 | SPEC-RECIPES loader test (inherited) |
| 010 | Test #6 |
| 011 | Test #7 |
| 012 | Test #8 |
| 013 | Test #9 |
| 014 | Test #10 |
| 015, 016, 017 | Test #11 |
| 018 | Manual review (no aliases column populated) + structural check inside Test #6 (taxonomy compliance) |
| 019 | Manual review (§9.4) |
| 020 | Manual review (§9.4) |
| 021 | Manual review (§9.4) |
| 022 | Manual review (§9.4) |
| 023 | Test #12 |
| 024 | PR review process (§9.4) |
| 025 | PR description (§9.4) |
| 026 | Manual quickstart, output pasted in PR (§9.5) |
| 027 | Tests #7–#10 collectively |

Every SPEC-ID has at least one verification artifact ✅.

---

## A.4 Sign-off

- **Drafted:** 2026-05-02 by sdd-spec-author
- **Status:** Final (checklist green)
- **Reviewer:** *to be filled at PR time*
- **Implementing story:** `docs/stories/06-curated-recipes-content.md`
- **Implementing PR:** *to be linked at PR time*
