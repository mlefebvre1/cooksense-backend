# SPEC-REACTIONS — §1 AI Preamble

[← Index](SPEC-REACTIONS-00-index.md)

## 1.1 Authorship and binding scope

This specification was authored under the `sdd-spec-author` skill. It is the
binding contract for Story 08 implementation. Any code that contradicts this
spec is a defect; any change of behavior shall amend this spec **first** in
the same PR.

## 1.2 Language

All requirements use RFC-2119 keywords (`SHALL`, `MUST`, `SHOULD`, `MAY`).
Bulleted prose without one of those keywords is descriptive context, not a
requirement.

## 1.3 Traceability rules

- Every SPEC-REACTIONS-NNN ID listed in `SPEC-REACTIONS-00-index.md` SHALL be
  defined in §6 (or §7/§9/§10 where explicitly cross-referenced).
- Every SPEC-REACTIONS-NNN ID SHALL map forward to at least one named test in
  §9 (`SPEC-REACTIONS-09-testing.md`).
- Every non-trivial new function and exported symbol SHALL cite the
  SPEC-REACTIONS-NNN IDs it implements in its Go Doc comment (backward
  traceability).
- Every commit body SHALL cite the SPEC-REACTIONS-NNN IDs it implements,
  fixes, or deprecates, per the repo's commit-message contract.

## 1.4 Quality bar

The code produced from this spec is **production-grade**:

- No TODOs, no `fmt.Print*` for logging, no dead code.
- Functions ≤ 30 lines, structs ≤ 10 methods, files ≤ 500 lines.
- Doc comments on every exported symbol, in standard Go Doc form.
- Errors handled explicitly; never `_ = err`. Sentinel errors compared with
  `errors.Is`; typed errors extracted with `errors.AsType`.
- `context.Context` is the first argument of every function performing I/O.
- Logging uses `log/slog`. Firebase tokens, full request bodies, and PII
  SHALL NOT appear in any log message at any level.
