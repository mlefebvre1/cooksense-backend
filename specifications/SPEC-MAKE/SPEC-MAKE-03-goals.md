# SPEC-MAKE-03 — Goals & Non-Goals

> Part of [SPEC-MAKE](SPEC-MAKE-00-index.md) — Story 02: Makefile Targets

---

## 2. Goals & Non-Goals

### 2.1 Goals

| ID | Goal |
|----|------|
| G-1 | Provide **one command per task** (`make up`, `make test`, etc.) so contributors never need to memorise docker/go flags. |
| G-2 | Make the `Makefile` **self-documenting** via a `make help` that parses `##` comments. |
| G-3 | **Forward `.env`** variables to every sub-process automatically when the file exists. |
| G-4 | Keep the `Makefile` **readable** (≤ 80 lines); any longer logic goes in `scripts/`. |
| G-5 | Protect the Postgres data volume from accidental deletion — `make clean` never drops it unless `CLEAN_VOLUMES=1`. |

### 2.2 Non-Goals

| ID | Non-Goal |
|----|----------|
| NG-1 | **CI/CD pipelines** — this is a developer DX tool only; GitHub Actions are out of MVP. |
| NG-2 | **Windows support** — targets rely on POSIX shell. |
| NG-3 | **`migrate` and `seed` implementation** — the targets exist and call `go run ./cmd/cooksense-server <subcmd>`, but the subcommand logic is wired in stories 03 and 05. |
| NG-4 | **Pinning external tool versions** — the `Makefile` assumes `docker compose` (V2 plugin), `modd`, `go`, and optionally `golangci-lint` are on PATH. |
