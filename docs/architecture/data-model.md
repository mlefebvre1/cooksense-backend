# Data model

This is the canonical schema for the MVP. Any change requires a new migration
and a decision-log entry that supersedes the relevant DDL.

## Entity overview

```
users ─────────────┐
                   │ 1..N
user_reactions ────┘─────────► recipes ──────┐
                                              │ 1..N
                          recipe_ingredients ─┘─────► ingredients

lesson_articles  (independent, no FK)
```

## DDL — `migrations/0001_init.up.sql`

```sql
-- Enums

CREATE TYPE reaction_kind AS ENUM ('LIKE', 'DISLIKE', 'TRY_LATER');

-- Users (identity comes from Firebase; firebase_uid is the PK)

CREATE TABLE users (
    firebase_uid TEXT        PRIMARY KEY,
    display_name TEXT,
    email        TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Recipes

CREATE TABLE recipes (
    id                    BIGSERIAL   PRIMARY KEY,
    slug                  TEXT        UNIQUE NOT NULL,
    title                 TEXT        NOT NULL,
    concept               TEXT        NOT NULL,
    time_minutes          INT         NOT NULL CHECK (time_minutes > 0),
    passive_prep_minutes  INT         NOT NULL DEFAULT 0 CHECK (passive_prep_minutes >= 0),
    cooking_methods       TEXT[]      NOT NULL DEFAULT '{}',
    tags                  TEXT[]      NOT NULL DEFAULT '{}',
    flavor_profile        TEXT[]      NOT NULL DEFAULT '{}',
    steps                 JSONB       NOT NULL,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX recipes_tags_gin            ON recipes USING GIN (tags);
CREATE INDEX recipes_cooking_methods_gin ON recipes USING GIN (cooking_methods);

-- Ingredients

CREATE TABLE ingredients (
    id        BIGSERIAL PRIMARY KEY,
    name      TEXT      UNIQUE NOT NULL,
    category  TEXT      NOT NULL,           -- protein, vegetable, starch, dairy, spice, …
    aliases   TEXT[]    NOT NULL DEFAULT '{}'
);
CREATE INDEX ingredients_aliases_gin ON ingredients USING GIN (aliases);

-- Recipe ↔ Ingredient (M:N with quantities)

CREATE TABLE recipe_ingredients (
    recipe_id     BIGINT  NOT NULL REFERENCES recipes(id)     ON DELETE CASCADE,
    ingredient_id BIGINT  NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    quantity      NUMERIC(10,2),
    unit          TEXT,
    optional      BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (recipe_id, ingredient_id)
);
CREATE INDEX recipe_ingredients_ingredient_idx ON recipe_ingredients(ingredient_id);

-- User reactions (one row per user × recipe; kind is upserted)

CREATE TABLE user_reactions (
    firebase_uid TEXT          NOT NULL REFERENCES users(firebase_uid) ON DELETE CASCADE,
    recipe_id    BIGINT        NOT NULL REFERENCES recipes(id)         ON DELETE CASCADE,
    kind         reaction_kind NOT NULL,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT now(),
    PRIMARY KEY (firebase_uid, recipe_id)
);
CREATE INDEX user_reactions_uid_kind_idx ON user_reactions(firebase_uid, kind);

-- Lessons (Cooking School)

CREATE TABLE lesson_articles (
    slug       TEXT PRIMARY KEY,
    title      TEXT NOT NULL,
    body_md    TEXT NOT NULL,
    sort_order INT  NOT NULL DEFAULT 0
);
```

## Recipe YAML schema (input format for the seed loader)

```yaml
slug: pan-seared-chicken-thighs-lemon-cabbage   # required, kebab-case, unique
title: "Pan-seared chicken thighs, lemon-braised cabbage"
concept: |
  Sear bone-in thighs hard, then finish in a covered pan with cabbage that
  steams in chicken fat and lemon juice. One pan, 25 minutes.
time_minutes: 25
passive_prep_minutes: 0
cooking_methods: [pan-sear, braise]              # taxonomy: see below
tags: [weeknight, one-pan, gluten-free]
flavor_profile: [acid, fat, umami]               # taxonomy: see below
ingredients:
  - name: chicken thighs (bone-in, skin-on)
    category: protein
    quantity: 4
    unit: piece
  - name: green cabbage
    category: vegetable
    quantity: 0.5
    unit: head
  - name: lemon
    category: vegetable
    quantity: 1
    unit: piece
  - name: garlic
    category: vegetable
    quantity: 3
    unit: clove
  - name: olive oil
    category: fat
    quantity: 2
    unit: tbsp
    optional: true
steps:
  - "Salt the thighs heavily 10 min ahead."
  - "Sear skin-side down in a cold pan, render for 8 min until deeply golden."
  - "Flip, add sliced cabbage, garlic, lemon juice, cover and braise 12 min."
  - "Rest 3 min, slice, spoon pan juices over."
```

### Allowed values

- `cooking_methods`: `pan-sear`, `roast`, `braise`, `boil`, `steam`, `grill`,
  `slow-cook`, `pressure-cook`, `raw`, `bake`.
- `flavor_profile`: `acid`, `fat`, `salt`, `sweet`, `bitter`, `umami`, `spicy`,
  `herbaceous`.
- `category` (ingredients): `protein`, `vegetable`, `fruit`, `starch`, `dairy`,
  `fat`, `spice`, `herb`, `condiment`, `liquid`, `other`.

The loader **rejects** values outside these vocabularies (fail-fast). The
canonical lists are owned by `internal/domain/taxonomy.go`.

## Invariants

- `slug` is the public identifier in URLs; never expose the numeric `id`.
- A recipe must have **≥ 2 ingredients** and **≥ 1 step**.
- A recipe must have at least one of `tags` or `cooking_methods` non-empty.
- `time_minutes` ≤ 30 for evening recipes (the only kind we ship in MVP).
  Recipes with `passive_prep_minutes > 0` are allowed regardless.
