# SPEC-CATALOG §B — Implementation Tasks

[← back to index](SPEC-CATALOG-00-index.md)

Atomic, ordered tasks for delivering story 06. Tasks are sized to be
individually reviewable; some are pure content, others are code/test.
Each task lists the SPEC-IDs it advances.

## B.1 Task list

| Task | Description | SPEC-IDs | Depends on |
|------|-------------|----------|------------|
| T-01 | Add `seed/recipes/_animal_proteins.txt` with the §7.2 baseline list | 013 | — |
| T-02 | Add `internal/seed/catalog_diversity_test.go` skeleton: file inventory + slug-filename test | 001, 002, 003, 004 | T-01, SPEC-RECIPES merged |
| T-03 | Extend the test with schema-validity case (calls `seed.Load` filesystem-only) | 005, 006 | T-02 |
| T-04 | Extend the test with `time_minutes` and `passive_prep_minutes` cases | 007, 008 | T-03 |
| T-05 | Extend the test with taxonomy-compliance case | 010 | T-03 |
| T-06 | Extend the test with diversity cases (methods, proteins, vegetarian, passive-cook) | 011, 012, 013, 014, 027 | T-01, T-03 |
| T-07 | Extend the test with ingredient-name discipline case | 015, 016, 017 | T-03 |
| T-08 | Extend the test with no-media-fields case | 023 | T-03 |
| T-09 | Author 5 chicken/poultry-based recipes (YAML) | 005–010, 015–021 | T-03 (loader green) |
| T-10 | Author 4 beef/pork/fish recipes (YAML) | 005–010, 015–021 | T-09 |
| T-11 | Author 4 vegetarian recipes (YAML) | 005–010, 013, 015–021 | T-10 |
| T-12 | Author ≥ 2 slow-cook / pressure-cook recipes (YAML) | 005–010, 014, 015–021 | T-11 |
| T-13 | Run `make up && make migrate && make seed`; capture output | 026 | T-09–T-12 |
| T-14 | Run `go test ./internal/seed/...`; ensure green | 001–017, 023, 027 | T-08, T-13 |
| T-15 | Compose PR description with summary table (§10.3 template) | 025 | T-14 |
| T-16 | Add README "Catalog" subsection (§10.1) | (doc) | T-14 |
| T-17 | Request reviewer round; ensure ≥ 1 inline comment per recipe | 024 | T-15 |

## B.2 Critical path

```
T-01 ──► T-02 ──► T-03 ──┬──► T-04 ──┐
                         ├──► T-05 ──┤
                         ├──► T-06 ──┤
                         ├──► T-07 ──┤
                         └──► T-08 ──┘
                                     │
                                     ▼
T-09 ──► T-10 ──► T-11 ──► T-12 ──► T-13 ──► T-14 ──► T-15 ──► T-17
                                                        │
                                                        └──► T-16
```

The diversity test (T-02..T-08) **shall** be merged before bulk recipe
authoring (T-09..T-12). This guarantees every authored recipe is
validated by the test as it lands, instead of accumulating violations.

## B.3 Effort estimate

| Bucket | Tasks | Hours |
|--------|-------|-------|
| Test code | T-01..T-08 | ~3 h |
| Recipe authoring | T-09..T-12 | ~6 h (≈ 25 min/recipe × 15) |
| Verification | T-13, T-14 | ~0.5 h |
| Documentation & PR | T-15, T-16 | ~0.5 h |
| Review cycle | T-17 | variable (out of author's control) |
| **Total (author side)** | | **~10 h** |

This matches the M estimate on the story.

## B.4 PR-split strategy

Recommended single PR for the whole story (the story explicitly says
"Catalog reviewed in a single PR"). However, the **diversity test**
(T-01..T-08) **may** ship in a precursor PR titled
`test(seed): catalog diversity guard (SPEC-CATALOG-027)` to:

- Land before recipe authors start.
- Make the diversity contract explicit and reviewed independently.
- Keep the recipe-PR diff smaller (YAML only).

Both strategies are acceptable; the precursor-PR approach is
recommended.

## B.5 Done criteria (recap)

- All tasks T-01..T-17 closed.
- `go test ./internal/seed/...` green.
- `make seed` succeeds; output pasted in PR.
- ≥ 15 catalog recipes merged.
- ≥ 1 reviewer comment per recipe in the PR thread.
- PR description contains the summary table (SPEC-CATALOG-025).
- This spec referenced in the story file (SPEC-CATALOG §10.5).
