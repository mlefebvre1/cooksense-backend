# Story 10 — Cooking School lessons API + 3 seed articles

Status: TODO
Estimate: S

## User story

As a user, I want to read short articles that teach me how to cook without
recipes (ratios, Maillard, flavor balance) so that I improve over time and
become less dependent on the catalog.

## Background

The "Cooking School" is the long-term differentiator (see
`docs/product/vision.md`). For the MVP we ship a minimal version: 3 static
Markdown articles served via 2 endpoints.

## Acceptance criteria

- [ ] `GET /api/lessons` (auth):
  - Returns `{"lessons": [{"slug","title","sort_order"}]}` ordered by
    `sort_order ASC`, `title ASC`.
- [ ] `GET /api/lessons/{slug}` (auth):
  - Returns `{"slug","title","body_md"}`.
  - `404 NOT_FOUND` for unknown slugs.
- [ ] Lesson loader (`internal/seed/lessons.go`) reads
      `seed/lessons/*.md` files. Each file starts with a YAML front-matter:
      ```
      ---
      slug: ratios-foundations
      title: "Foundational ratios"
      sort_order: 1
      ---
      # Foundational ratios
      …Markdown body…
      ```
      and is upserted into `lesson_articles`.
- [ ] The `seed` subcommand from story 05 also loads lessons (single
      transaction with recipes is acceptable).
- [ ] At least **3 lessons** are committed in `seed/lessons/`:
  - `ratios-foundations.md` — vinaigrette, batter, béchamel, stock ratios.
  - `maillard.md` — why and how to brown.
  - `flavor-balance.md` — acid/fat/salt/sweet/umami balancing at the bench.
- [ ] Unit tests cover front-matter parsing and slug validation.

## Technical notes

- Use `gopkg.in/yaml.v3` to parse the front-matter block. Body is the rest.
- Markdown is **not rendered** server-side; clients render. The backend is a
  pass-through for `body_md`.
- Slug validation: same regex as recipes (`^[a-z0-9]+(-[a-z0-9]+)*$`).

## Out of scope

- Quizzes, progression tracking (post-MVP V2).
- WYSIWYG editor / authoring UI.
- Image embedding.

## Dependencies

- depends on: 04, 05
- blocks: —

## Definition of Done

- [ ] AC met.
- [ ] Reviewer confirms each lesson is genuinely useful (not filler).
- [ ] Integration test hits both endpoints with a `fakeVerifier`.
