# SPEC-DISCOVER — §1 AI Preamble

[← Index](SPEC-DISCOVER-00-index.md)

## 1.1 Authorship and binding scope

This specification was authored under the `sdd-spec-author` skill. It is the
binding contract for Story 07 implementation. Any code that contradicts this
spec is a defect; any change of behavior shall amend this spec **first** in
the same PR.

## 1.2 Language

All requirements use RFC-2119 keywords (`SHALL`, `MUST`, `SHOULD`, `MAY`).
Bullet points without these keywords are explanatory context, not normative.

## 1.3 Traceability rules

- Every `SPEC-DISCOVER-NNN` ID in §6 SHALL map to at least one named test in
  §9.
- Every public Go symbol introduced by this story SHALL have a doc comment
  citing the SPEC-DISCOVER-NNN IDs it implements.
- Every commit body for Story 07 work SHALL cite at least one
  `SPEC-DISCOVER-NNN` ID, per repository convention.

## 1.4 Constraints reaffirmed

- Modern Go 1.26 idioms only (see project guidelines): `any`, `errors.Is`,
  `slices`/`maps`, `for i := range n`, `t.Context()`, `errors.AsType`,
  `wg.Go`, `omitzero`, `cmp.Or`, etc.
- Standard-library `net/http` only; **no** third-party HTTP framework.
- Clean Architecture: `handler → service → repo`. Reverse imports are
  forbidden.
- The repository SHALL be the **only** layer touching `pgx`/`pgxpool`
  (per AC #3 of the story).
- Handlers SHALL only handle JSON encode/decode and HTTP status mapping
  (per AC #5 of the story).
