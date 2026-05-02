# SPEC-MAKE-10 — Documentation Specification

> Part of [SPEC-MAKE](SPEC-MAKE-00-index.md) — Story 02: Makefile Targets

---

## 9. Documentation Specification

### 9.1 Makefile Inline Documentation

Every public target **shall** carry a trailing `##` comment on its declaration line. These comments serve two purposes:

1. They are parsed by the `help` target and printed to developers.
2. They make the Makefile self-explanatory without requiring external docs.

Required format:
```makefile
targetname: ## Short description of what this target does.
	recipe
```

### 9.2 README

Story 02 **shall not** update `README.md`. The README is updated in Story 12 which references and documents all `make` targets.

### 9.3 `.env.example` as Documentation

`.env.example` **shall** include brief inline comments for each variable (using `#`) to explain what it controls and whether it is required:

```
DATABASE_URL=postgres://cooksense:cooksense@localhost:5432/cooksense?sslmode=disable  # required
FIREBASE_PROJECT_ID=changeme  # required
GOOGLE_APPLICATION_CREDENTIALS=path/to/firebase-admin.json  # required
```
