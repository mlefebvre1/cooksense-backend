# Decisions log (ADR-lite)

Each decision is short, dated, and immutable. To change one, **add a new
entry** that supersedes the previous one — never edit history.

Format:
```
### D-NNNN — Title
- Date: YYYY-MM-DD
- Status: Accepted | Superseded by D-XXXX
- Context: …
- Decision: …
- Consequences: …
```

---

### D-0001 — Backend language and version
- Date: 2026-05-02
- Status: Accepted
- Context: Need a fast, simple, statically-typed backend.
- Decision: **Go 1.26.2** (the version pinned in this project's tooling).
- Consequences: We use modern Go idioms (see project guidelines): `any`,
  `errors.Is`, `slices`/`maps` packages, `for i := range n`, `t.Context()`,
  `b.Loop()`, `wg.Go`, `errors.AsType`, `omitzero`, etc.

### D-0002 — HTTP layer
- Date: 2026-05-02
- Status: Accepted
- Context: We want stability and minimal magic.
- Decision: **Standard library `net/http`** with the Go 1.22+ pattern-based
  router (`mux.HandleFunc("GET /api/...", h)`). No framework (no chi, echo, gin).
- Consequences: We write our own middleware chain. We trade a tiny bit of
  ergonomics for long-term simplicity and zero framework lock-in.

### D-0003 — Database
- Date: 2026-05-02
- Status: Accepted
- Context: Need arrays, JSONB, full-text/GIN indexes for ingredient matching.
- Decision: **PostgreSQL 17** (already pinned in `docker-compose.yml`).
- Consequences: We use `pgx/v5` as driver/pool. We use Postgres-specific
  features (arrays, GIN, ENUM) without apology.

### D-0004 — Authentication
- Date: 2026-05-02
- Status: Accepted
- Context: The mobile app already uses Firebase Auth.
- Decision: **One Firebase project shared between mobile and backend.** Backend
  verifies the Firebase ID token using the official `firebase.google.com/go/v4`
  Admin SDK. No anonymous mode.
- Consequences: Every authenticated request carries `Authorization: Bearer
  <ID_TOKEN>`. The backend lazily provisions a `users` row on first call.

### D-0005 — Recipe content pipeline
- Date: 2026-05-02
- Status: Accepted
- Context: We must avoid runtime LLM dependence and keep latency predictable.
- Decision: Recipes are authored as **YAML files in `seed/recipes/`** and
  loaded via a `seed` CLI command. Generation method (manual, LLM-assisted) is
  unconstrained, but **every published recipe is human-reviewed**.
- Consequences: Recipes ship with the binary's repo. Updates require a code
  release until V1 introduces user-submitted recipes.

### D-0006 — Migrations
- Date: 2026-05-02
- Status: Accepted
- Context: We need reproducible schema evolution.
- Decision: **`golang-migrate/migrate/v4`**, plain SQL files under
  `migrations/`. Migrations run via `make migrate` (or as a subcommand of the
  server binary).
- Consequences: No ORM, no schema drift. Every change is a new pair of
  `up.sql`/`down.sql` files.

### D-0007 — Identity model
- Date: 2026-05-02
- Status: Accepted
- Context: Identity comes from Firebase, but we still need joinable rows.
- Decision: The primary key for users is **`firebase_uid TEXT`** (no internal
  numeric id). All FK references use the same `TEXT` column.
- Consequences: When/if we ever migrate off Firebase, we add an internal id at
  that time. For MVP we avoid the indirection.

### D-0008 — Reaction kinds
- Date: 2026-05-02
- Status: Accepted
- Context: The product requires three swipe outcomes.
- Decision: A Postgres `ENUM('LIKE','DISLIKE','TRY_LATER')` named
  `reaction_kind`, with one row per `(user, recipe)` pair (UPSERT).
- Consequences: Simple model, no history. If we later need history, we add a
  `reaction_history` table without touching `user_reactions`.

### D-0009 — Out of MVP
- Date: 2026-05-02
- Status: Accepted
- Decision: Community photos, user-generated recipes, meal planning, and
  Cooking School quizzes are **explicitly out of MVP** (see `scope-mvp.md`).
- Consequences: No tables, no endpoints, no code paths for those features in
  Sprint 0–N until MVP ships.
