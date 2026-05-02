# CookSense Backend — Documentation

This folder contains the product framing and the implementation specs for the
CookSense backend. It is the **single source of truth** for what we are building
and why. The original pitch lives in [`../project-proposal.md`](../project-proposal.md).

## Layout

```
docs/
├── README.md                  ← you are here
├── product/                   ← what & why
│   ├── vision.md
│   ├── personas.md
│   ├── scope-mvp.md
│   └── decisions.md           ← decision log (ADR-lite)
├── architecture/              ← how
│   ├── overview.md
│   ├── data-model.md
│   ├── api.md
│   ├── auth.md
│   └── infra.md
└── stories/                   ← implementation specs (one file per story)
    ├── README.md              ← backlog index + dependency graph
    ├── 01-bootstrap-project-structure.md
    ├── 02-makefile-targets.md
    ├── 03-db-pgx-migrations.md
    ├── 04-firebase-auth-middleware.md
    ├── 05-recipes-domain-yaml-loader.md
    ├── 06-curated-recipes-content.md
    ├── 07-api-recipes-discover.md
    ├── 08-api-reactions-my-recipes.md
    ├── 09-api-recipes-search.md
    ├── 10-api-lessons.md
    ├── 11-integration-tests-skeleton.md
    └── 12-readme-quickstart.md
```

## How to read this

1. **New to the project?** → Start with `product/vision.md`, then `product/scope-mvp.md`.
2. **Coming to implement?** → Read `architecture/overview.md`, then pick a story
   from `stories/` (follow the order in `stories/README.md`).
3. **Coming to review a PR?** → Open the matching story file; the *Acceptance
   Criteria* and *Definition of Done* sections are the review checklist.

## Conventions

- Stories are **numbered** and ordered by recommended implementation order.
- Each story is **self-contained**: a developer should be able to implement it
  without reading other stories (cross-references are explicit when needed).
- All product/architecture choices are recorded in `product/decisions.md`. If a
  story contradicts a decision, **the decision wins** until updated.
- Status of each story is tracked at the top of its file (`Status: TODO | IN
  PROGRESS | DONE`).
