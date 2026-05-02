# cooksense-backend

## Database migrations

Migrations live in [migrations/](migrations/) and are managed with [tern](https://github.com/jackc/tern).

`tern` is pinned as a Go tool in [go.mod](go.mod), so you don't need to install it globally — invoke it via `go tool tern`.

### One-time setup

Start Postgres:

```
docker compose up -d
```

### Running migrations

Apply all pending migrations:

```
go tool tern migrate -m ./migrations -c ./migrations/tern.conf
```

`modd` also runs this automatically on every `.go` or `migrations/*.sql` change — see [modd.conf](modd.conf).

### Creating a new migration

```
go tool tern new -m ./migrations <name>
```

This creates `migrations/NNN_<name>.sql` with the up/down separator pre-filled:

```sql
-- write the "up" SQL above this line
---- create above / drop below ----
-- write the "down" SQL below this line
```

### Rolling back

Roll back the most recent migration:

```
go tool tern migrate -m ./migrations -c ./migrations/tern.conf --destination -1
```

Roll back to a specific version (e.g. version 0 = empty database):

```
go tool tern migrate -m ./migrations -c ./migrations/tern.conf --destination 0
```

### Status

Show current migration version:

```
go tool tern status -m ./migrations -c ./migrations/tern.conf
```
