# SPEC-RECIPES — §1 AI Preamble

[← Index](SPEC-RECIPES-00-index.md)

## 1.1 Authorship and binding scope

This specification was authored under the `sdd-spec-author` skill. It is the
binding contract for Story 05 implementation. Any code that contradicts this
spec is a defect; any change of behavior shall amend this spec **first** in
the same PR.

## 1.2 Language

All requirements use RFC-2119 keywords (`SHALL`, `MUST`, `SHOULD`, `MAY`).
Bullet points without these keywords are explanatory context, not normative.

## 1.3 Traceability rules

- Every `SPEC-RECIPES-NNN` ID in §6 SHALL map to at least one named test in
  §9.
- Every public Go symbol introduced by this story SHALL have a doc comment
  citing the SPEC-RECIPES-NNN IDs it implements.
- Every commit body for Story 05 work SHALL cite at least one
  `SPEC-RECIPES-NNN` ID, per repository convention.

## 1.4 Constraints reaffirmed

- Modern Go 1.26 idioms only (see project guidelines).
- `internal/domain` SHALL remain pure: no `os`, `io`, `pgx`, `database/sql`,
  `net/http`, or any other I/O dependency may be imported.
- The loader SHALL aggregate errors and SHALL NOT fail-fast on the first
  invalid file.
- The store SHALL be transactional; partial loads are forbidden.
