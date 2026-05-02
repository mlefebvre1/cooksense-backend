-- CookSense — migration 0001_init (down)
-- Description: Reverses 0001_init.up.sql in strict reverse order.
-- SPEC-DB-021

DROP TABLE IF EXISTS lesson_articles;
DROP TABLE IF EXISTS user_reactions;
DROP TABLE IF EXISTS recipe_ingredients;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS users;
DROP TYPE  IF EXISTS reaction_kind;
