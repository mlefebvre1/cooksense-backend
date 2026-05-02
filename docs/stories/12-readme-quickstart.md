# Story 12 — README quickstart

Status: TODO
Estimate: S

## User story

As a new contributor cloning the repo, I want a single page that tells me
exactly how to run the backend and exercise the API in under 5 minutes.

## Background

The current `README.md` is a stub. We finalize it once the moving pieces
exist (Make targets, migrations, seed, endpoints). It must reflect reality —
no aspirational instructions.

## Acceptance criteria

- [ ] The repo's top-level `README.md` contains the following sections, in
      order:
  1. **What is CookSense?** — one paragraph + link to
     `docs/product/vision.md`.
  2. **Tech stack** — Go 1.26.2, Postgres 17, Firebase Auth.
  3. **Prerequisites** — Go ≥ 1.26.2, Docker, `make`, `modd`.
  4. **Quickstart** —
     ```
     cp .env.example .env
     # edit .env to set FIREBASE_PROJECT_ID and GOOGLE_APPLICATION_CREDENTIALS
     make up           # start Postgres
     make migrate      # apply schema
     make seed         # load curated recipes + lessons
     make run          # start the server with hot reload
     ```
  5. **Sanity check** — `curl localhost:8080/api/health` should return
     `{"status":"ok"}`.
  6. **Common Make targets** — table generated from `make help`.
  7. **Project layout** — copy of the tree from
     `docs/architecture/overview.md` (kept in sync).
  8. **Where to look next** — pointers to `docs/`, with a short one-line
     description for each subdir.
- [ ] No section claims a feature that is not yet implemented. If a story is
      not done, either omit the feature or label it "🚧 in progress".
- [ ] All commands shown in the README are tested verbatim on a clean
      machine before merging.

## Technical notes

- Keep the README under ~150 lines; deeper content lives under `docs/`.
- Add badges only for things that exist (Go version, license). No fake build
  badges.

## Out of scope

- Deployment / production guide (post-MVP).
- API examples beyond the sanity check (link to `docs/architecture/api.md`).

## Dependencies

- depends on: 02, 03, 06
- blocks: —

## Definition of Done

- [ ] AC met.
- [ ] A teammate who has never seen the repo successfully runs the
      quickstart on their machine.
