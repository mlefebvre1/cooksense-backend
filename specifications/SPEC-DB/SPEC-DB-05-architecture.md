# SPEC-DB-05 — Architecture Overview

> Part of [SPEC-DB](SPEC-DB-00-index.md) — Story 03: Database pgx Pool, Migrations & Initial Schema

---

## 4. Architecture Overview

### 4.1 Design Patterns Applied

| Pattern | Usage in Story 03 |
|---------|-------------------|
| **Constructor Injection** | `db.Open` returns a `*pgxpool.Pool` that callers store and pass to repositories. No global pool. |
| **Fail-Fast Configuration** | `config.Load()` validates all required vars at startup; missing values → immediate `os.Exit(1)` (via log + return error to `main`). |
| **Command Pattern (CLI subcommand)** | `main.go` dispatches `os.Args[1]` to route `migrate up` / `migrate down N` to `db.Up` / `db.Down`. |

### 4.2 Dependency Graph

```mermaid
graph LR
    main["cmd/cooksense-server/main.go"] --> config["internal/config"]
    main --> db["internal/db"]
    db --> pgxpool["github.com/jackc/pgx/v5/pgxpool"]
    db --> migrate["github.com/golang-migrate/migrate/v4"]
    config --> stdlib["stdlib (os, fmt, errors)"]
```

### 4.3 Pool Lifecycle

```mermaid
sequenceDiagram
    participant main
    participant config
    participant db
    participant postgres

    main->>config: config.Load()
    config-->>main: Config{DatabaseURL, ...}
    main->>db: db.Open(ctx, cfg)
    db->>postgres: pgxpool.NewWithConfig(ctx, poolCfg)
    postgres-->>db: pool
    db->>postgres: pool.Ping(ctx)
    postgres-->>db: ok
    db-->>main: *pgxpool.Pool
    Note over main: pool injected into all repos
    main->>main: defer pool.Close()
```

### 4.4 Subcommand Dispatch

```mermaid
flowchart TD
    A["os.Args"] --> B{args[1]}
    B -- "migrate" --> C{args[2]}
    C -- "up" --> D["db.Up(ctx, dsn, migrationsDir)"]
    C -- "down" --> E["db.Down(ctx, dsn, migrationsDir, N)"]
    B -- other --> F["start HTTP server (future stories)"]
    D --> G["exit 0 / 1"]
    E --> G
```

### 4.5 Idempotency Guarantee

`golang-migrate` tracks applied migrations in a `schema_migrations` table (created automatically). Running `db.Up` when all migrations are already applied returns `migrate.ErrNoChange`, which **shall** be treated as success (exit `0`).
