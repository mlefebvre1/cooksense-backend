# SPEC-AUTH-01 — AI Preamble

← [Index](SPEC-AUTH-00-index.md)

---

## 1. Authorship & constraints

This specification was produced by the **sdd-spec-author** skill operating
inside the cooksense-backend project. The following hard constraints apply to
every agent or engineer reading this document:

1. **Do not generate code before all SPEC-AUTH-NNN IDs governing the change
   have been located.** If a requirement is missing, open a spec amendment
   issue first.
2. **RFC-2119 keywords are binding.** `shall` / `must` = mandatory.
   `should` = strong recommendation. `may` = optional. Prose without a
   keyword is informational only.
3. **Every SPEC-ID must have at least one named test** in
   [SPEC-AUTH-09-testing](SPEC-AUTH-09-testing.md). A requirement without a
   test does not exist for the purposes of the Definition of Done.
4. **Traceability is bidirectional.** Every non-trivial function docstring
   shall cite the SPEC-ID(s) it implements; every test docstring or name shall
   cite the SPEC-ID(s) it verifies.
5. **This spec and its code ship in the same PR.** A code change that
   contradicts this document without a concurrent spec amendment is rejected.
6. **Security rules are non-negotiable.** Credentials shall never appear in
   source code, logs, or test fixtures. See §7 and the threat model in
   [SPEC-AUTH-05-architecture](SPEC-AUTH-05-architecture.md).

## 2. Revision history

| Date | Version | Author | Summary |
|------|---------|--------|---------|
| 2026-05-02 | 1.0 | sdd-spec-author | Initial Final draft from Story 04 |
