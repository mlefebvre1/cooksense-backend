# SPEC-DB-06 — Package Specifications

> Part of [SPEC-DB](SPEC-DB-00-index.md) — Story 03: Database pgx Pool, Migrations & Initial Schema  
> SPEC-IDs covered: SPEC-DB-001 – SPEC-DB-027

---

## 5. Package Specifications

### 5.1 `internal/config` — Runtime Configuration

#### SPEC-DB-001: `config.Load()` reads all required env vars

`config.Load()` **shall** read the following environment variables:

| Env var | Field | Type | Required |
|---------|-------|------|----------|
| `APP_PORT` | `Config.AppPort` | `string` | no (default: `"8080"`) |
| `LOG_LEVEL` | `Config.LogLevel` | `string` | no (default: `"info"`) |
| `LOG_FORMAT` | `Config.LogFormat` | `string` | no (default: `"text"`) |
| `DATABASE_URL` | `Config.DatabaseURL` | `string` | **yes** |
| `FIREBASE_PROJECT_ID` | `Config.FirebaseProjectID` | `string` | **yes** |
| `GOOGLE_APPLICATION_CREDENTIALS` | `Config.GoogleAppCredentials` | `string` | **yes** |

#### SPEC-DB-002: `config.Load()` fails fast on missing required vars

If any required variable (`DATABASE_URL`, `FIREBASE_PROJECT_ID`, `GOOGLE_APPLICATION_CREDENTIALS`) is absent or empty, `config.Load()` **shall** return a non-nil error that describes which variable is missing. `main.go` **shall** log the error and exit `1`.

#### SPEC-DB-023: `internal/config` isolation

`internal/config` **shall not** import any `internal/domain` or feature packages. It **shall** only use stdlib packages (`os`, `fmt`, `errors`).

**Function signature:**

```go
// Load reads configuration from environment variables.
// It returns an error if any required variable is absent or empty.
func Load() (Config, error)
```

**Struct definition:**

```go
// Config holds the validated runtime configuration for cooksense-server.
type Config struct {
    AppPort              string
    LogLevel             string
    LogFormat            string
    DatabaseURL          string
    FirebaseProjectID    string
    GoogleAppCredentials string
}
```

#### SPEC-DB-003: `Config` struct immutability

The `Config` struct **shall** be a value type (not a pointer). After `Load()` returns, no field **shall** be mutated. The struct **shall** be passed by value to functions that need it.

---

### 5.2 `internal/db` — Database Pool & Migrations

#### SPEC-DB-004: `db.Open` signature

`db.Open` **shall** have the following signature:

```go
// Open creates and validates a PostgreSQL connection pool.
// It applies pool configuration defaults and pings the database before returning.
// Returns an error if the database is unreachable.
//
// SPEC-DB-004, SPEC-DB-005, SPEC-DB-006
func Open(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error)
```

#### SPEC-DB-005: Pool defaults

`db.Open` **shall** set the following pool defaults on the `pgxpool.Config` before calling `pgxpool.NewWithConfig`:

| Setting | Value |
|---------|-------|
| `MinConns` | `2` |
| `MaxConns` | `10` |
| `HealthCheckPeriod` | `1 * time.Minute` |

#### SPEC-DB-006: Ping on open

After creating the pool, `db.Open` **shall** call `pool.Ping(ctx)`. If the ping fails, `db.Open` **shall** close the pool and return the error.

#### SPEC-DB-007: Testing override

`db.Open` **shall** accept `pgxpool.Config` overrides. Recommended approach: an optional variadic `func(*pgxpool.Config)` parameter:

```go
func Open(ctx context.Context, cfg config.Config, overrides ...func(*pgxpool.Config)) (*pgxpool.Pool, error)
```

This allows tests to set `MaxConns = 2` or inject a test DSN without changing the production call site.

#### SPEC-DB-024: `internal/db` isolation

`internal/db` **shall not** import any feature packages (`internal/recipes`, `internal/reactions`, etc.). It **may** import `internal/config`.

---

### 5.3 Migration Functions

#### SPEC-DB-008: `db.Up` applies all pending migrations

```go
// Up applies all pending up migrations from migrationsDir using the pgx driver.
// It treats migrate.ErrNoChange as success (idempotent).
//
// SPEC-DB-008, SPEC-DB-009
func Up(ctx context.Context, dsn, migrationsDir string) error
```

`db.Up` **shall** apply all unapplied `.up.sql` files in `migrationsDir` in ascending version order using `golang-migrate` with the `pgx/v5` database driver and `file://` source.

#### SPEC-DB-009: Idempotency

`db.Up` **shall** treat `migrate.ErrNoChange` as success and return `nil`. Running `db.Up` when all migrations are already applied **shall** exit `0`.

#### SPEC-DB-010: `db.Down` rolls back N steps

```go
// Down rolls back the last n migration steps.
//
// SPEC-DB-010
func Down(ctx context.Context, dsn, migrationsDir string, n int) error
```

`db.Down` **shall** roll back exactly `n` migration steps. If `n <= 0`, it **shall** return an error (`"n must be positive"`).

#### SPEC-DB-011: `MigrationError` type

```go
// MigrationError is returned when a migration fails, carrying the version
// number of the failed migration for easier debugging.
//
// SPEC-DB-011
type MigrationError struct {
    Version uint
    Err     error
}

func (e *MigrationError) Error() string
func (e *MigrationError) Unwrap() error
```

When `golang-migrate` returns a `migrate.ErrDirty` or any migration execution error, `db.Up`/`db.Down` **shall** wrap it in `MigrationError`.

---

### 5.4 CLI Subcommand Dispatch

#### SPEC-DB-012: `cooksense-server migrate up`

Running `cooksense-server migrate up` **shall**:

1. Call `config.Load()`.
2. Call `db.Up(ctx, cfg.DatabaseURL, "migrations")`.
3. Log `"migrations applied"` at INFO level on success.
4. Exit `0` on success, `1` on error.

#### SPEC-DB-013: `cooksense-server migrate down [N]`

Running `cooksense-server migrate down 1` **shall**:

1. Call `config.Load()`.
2. Parse the `N` argument (default `1` if omitted).
3. Call `db.Down(ctx, cfg.DatabaseURL, "migrations", N)`.
4. Log `"migrations rolled back"` at INFO level on success.
5. Exit `0` on success, `1` on error.

---

### 5.5 `migrations/0001_init.up.sql` — Initial Schema DDL

The file `migrations/0001_init.up.sql` **shall** match the DDL in `docs/architecture/data-model.md` exactly.

#### SPEC-DB-014: `reaction_kind` enum

```sql
CREATE TYPE reaction_kind AS ENUM ('LIKE', 'DISLIKE', 'TRY_LATER');
```

#### SPEC-DB-015: `users` table

```sql
CREATE TABLE users (
    firebase_uid TEXT        PRIMARY KEY,
    display_name TEXT,
    email        TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### SPEC-DB-016: `recipes` table + GIN indexes

```sql
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
```

#### SPEC-DB-017: `ingredients` table + GIN index

```sql
CREATE TABLE ingredients (
    id        BIGSERIAL PRIMARY KEY,
    name      TEXT      UNIQUE NOT NULL,
    category  TEXT      NOT NULL,
    aliases   TEXT[]    NOT NULL DEFAULT '{}'
);
CREATE INDEX ingredients_aliases_gin ON ingredients USING GIN (aliases);
```

#### SPEC-DB-018: `recipe_ingredients` join table + index

```sql
CREATE TABLE recipe_ingredients (
    recipe_id     BIGINT  NOT NULL REFERENCES recipes(id)     ON DELETE CASCADE,
    ingredient_id BIGINT  NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    quantity      NUMERIC(10,2),
    unit          TEXT,
    optional      BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (recipe_id, ingredient_id)
);
CREATE INDEX recipe_ingredients_ingredient_idx ON recipe_ingredients(ingredient_id);
```

#### SPEC-DB-019: `user_reactions` table + index

```sql
CREATE TABLE user_reactions (
    firebase_uid TEXT          NOT NULL REFERENCES users(firebase_uid) ON DELETE CASCADE,
    recipe_id    BIGINT        NOT NULL REFERENCES recipes(id)         ON DELETE CASCADE,
    kind         reaction_kind NOT NULL,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT now(),
    PRIMARY KEY (firebase_uid, recipe_id)
);
CREATE INDEX user_reactions_uid_kind_idx ON user_reactions(firebase_uid, kind);
```

#### SPEC-DB-020: `lesson_articles` table

```sql
CREATE TABLE lesson_articles (
    slug       TEXT PRIMARY KEY,
    title      TEXT NOT NULL,
    body_md    TEXT NOT NULL,
    sort_order INT  NOT NULL DEFAULT 0
);
```

---

### 5.6 `migrations/0001_init.down.sql` — Reverse Teardown

#### SPEC-DB-021: Down migration drops all objects in reverse order

`migrations/0001_init.down.sql` **shall** drop all objects created in `0001_init.up.sql` in strict reverse order:

```sql
DROP TABLE IF EXISTS lesson_articles;
DROP TABLE IF EXISTS user_reactions;
DROP TABLE IF EXISTS recipe_ingredients;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS users;
DROP TYPE  IF EXISTS reaction_kind;
```

---

### 5.7 Constructor Injection Rule

#### SPEC-DB-022: No global pool

The `*pgxpool.Pool` **shall** be passed via constructor arguments to every repository. No package-level `var pool *pgxpool.Pool` **shall** exist anywhere in the codebase. This is enforced at code review.
