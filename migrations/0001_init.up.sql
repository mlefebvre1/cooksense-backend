-- CookSense — migration 0001_init
-- Description: Creates the initial schema (users, recipes, ingredients,
--              recipe_ingredients, user_reactions, lesson_articles).
-- SPEC-DB-014 through SPEC-DB-020

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
    category  TEXT      NOT NULL,
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
