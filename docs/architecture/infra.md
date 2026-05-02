# Infrastructure

## Local development

| Component        | Choice                                              |
|------------------|-----------------------------------------------------|
| OS               | macOS / Linux                                        |
| Go               | 1.26.2                                               |
| Postgres         | Docker — `postgres:17-alpine` via `docker-compose.yml` (already in repo) |
| Hot reload       | `modd` (config already in `modd.conf`)               |
| Migration tool   | `golang-migrate/migrate/v4`                          |
| Build tooling    | Plain `go build` + `Makefile`                        |

The repo's existing `docker-compose.yml` is the **only** way we run Postgres
locally. We do not duplicate it.

## Environment variables

Authoritative list (mirrored in `.env.example`):

| Var                              | Required | Default        | Description                                    |
|----------------------------------|----------|----------------|------------------------------------------------|
| `APP_PORT`                       | no       | `8080`         | HTTP listen port                                |
| `LOG_LEVEL`                      | no       | `info`         | `debug`/`info`/`warn`/`error`                   |
| `LOG_FORMAT`                     | no       | `text` (dev)   | `text` or `json`                                |
| `DATABASE_URL`                   | yes      | —              | `postgres://user:pass@host:port/db?sslmode=…`   |
| `POSTGRES_USER`                  | (compose)| `cooksense`    | used by docker-compose only                     |
| `POSTGRES_PASSWORD`              | (compose)| `cooksense`    | used by docker-compose only                     |
| `POSTGRES_DB`                    | (compose)| `cooksense`    | used by docker-compose only                     |
| `POSTGRES_PORT`                  | (compose)| `5432`         | host port mapping                               |
| `FIREBASE_PROJECT_ID`            | yes      | —              | Token audience                                  |
| `GOOGLE_APPLICATION_CREDENTIALS` | yes      | —              | Path to Firebase Admin SDK JSON                 |

Loading rules:
- The server reads env at startup.
- Missing required vars → log and exit non-zero.
- No fallback secrets in code.

## Make targets (canonical)

```
make up         # docker compose up -d  (starts Postgres)
make down       # docker compose down
make migrate    # apply all up migrations
make seed       # load YAML recipes + lessons into the DB
make run        # modd (hot reload) — depends on `up` + `migrate`
make build      # go build -o bin/cooksense-server ./cmd/cooksense-server
make test       # go test ./...
make lint       # go vet + (optional) golangci-lint run
make clean      # remove bin/, drop volumes optionally on confirm
```

## CI / CD (out of MVP, sketched)

- CI: GitHub Actions running `make lint test` against a service-container
  Postgres. Migrations applied as part of the test job.
- CD: out of scope until V1. MVP demo runs locally.

## Observability (MVP)

- Structured logs via `log/slog` JSON in production. Each request gets a
  `request_id` (UUID v4) injected by middleware and emitted on every log line
  for that request.
- No metrics or tracing in MVP. Stub: future story can add Prometheus.

## Backup / data lifecycle

- Local dev: Postgres data lives in the `postgres_data` Docker volume. It can
  be wiped at any time; `make seed` recreates the catalog.
- Production: out of scope for MVP.
