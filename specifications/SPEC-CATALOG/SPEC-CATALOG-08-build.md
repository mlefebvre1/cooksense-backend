# SPEC-CATALOG §8 — Build, Tooling, Quality

[← back to index](SPEC-CATALOG-00-index.md)

## 8.1 Go build impact

SPEC-CATALOG adds **no production Go code**. `go build ./...` and
`go vet ./...` outcomes are unchanged versus SPEC-RECIPES.

The single new test file (§9) compiles under the existing `go test ./...`
target and **shall not** require build tags.

## 8.2 YAML quality gates

### 8.2.1 Encoding

Catalog files **shall** be UTF-8, LF line endings (no CRLF, no BOM).

### 8.2.2 Indentation

Catalog files **shall** use 2-space indentation. Tabs are forbidden.

### 8.2.3 Trailing whitespace

Catalog files **shall not** contain trailing whitespace on any line.

### 8.2.4 Final newline

Catalog files **shall** end with exactly one trailing newline.

### 8.2.5 Optional yamllint configuration (recommended)

A `.yamllint.yaml` at repo root **may** be added with:

```yaml
extends: default
rules:
  document-start: disable
  line-length:
    max: 120
  truthy:
    allowed-values: ["true", "false"]
  comments:
    min-spaces-from-content: 1
```

If present, CI **should** invoke `yamllint seed/recipes/`. This is
non-blocking for SPEC-CATALOG (the diversity test in §9 is the
mandatory guard) but recommended.

## 8.3 Linting in `make lint`

`make lint` (SPEC-MAKE-008) **may** be extended to run yamllint on
`seed/recipes/` if yamllint is on `PATH`; it **shall** continue to pass
when yamllint is not installed.

## 8.4 No third-party dependencies added

The diversity test (§9) **shall** use only:

- Standard library (`os`, `path/filepath`, `strings`, `testing`).
- `gopkg.in/yaml.v3` (already pulled in by SPEC-RECIPES).

It **shall not** import any helper from `internal/seed` that performs DB
I/O — the test is filesystem-only.
