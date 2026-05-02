# SPEC-MAKE-09 — Testing & Verification Specification

> Part of [SPEC-MAKE](SPEC-MAKE-00-index.md) — Story 02: Makefile Targets

---

## 8. Testing & Verification Specification

### 8.1 Verification Approach

Story 02 delivers a `Makefile`, not Go code — therefore there are no Go unit tests for the Makefile itself. Verification is done by **running each `make` target** and observing its exit code and side-effects.

### 8.2 Verification Table

| Test name | Target | Precondition | Steps | Expected result |
|-----------|--------|-------------|-------|-----------------|
| `Verify_MakeHelp_PrintsAllTargets` | `help` | clean repo | `make help` | Exit 0; stdout contains 10 target names |
| `Verify_MakeBuild_ProducesBinary` | `build` | story 01 merged | `make build` | Exit 0; `bin/cooksense-server` exists and is executable |
| `Verify_MakeTest_PassesOnCleanTree` | `test` | story 01 merged | `make test` | Exit 0 |
| `Verify_MakeLint_PassesOnCleanTree` | `lint` | story 01 merged | `make lint` | Exit 0 |
| `Verify_MakeUp_StartsPostgres` | `up` | Docker running | `make up` | Exit 0; Postgres container healthy |
| `Verify_MakeDown_StopsPostgres` | `down` | `make up` ran | `make down` | Exit 0; container stopped; volume present |
| `Verify_MakeClean_RemovesBin` | `clean` | `make build` ran | `make clean` | Exit 0; `bin/` removed; volume present |
| `Verify_MakeCleanVolumes_DropsVolume` | `clean CLEAN_VOLUMES=1` | `make up` ran | `make clean CLEAN_VOLUMES=1` | Exit 0; `bin/` removed; `postgres_data` volume absent |
| `Verify_MakeLint_GolangciMissing_Skips` | `lint` | `golangci-lint` absent | `PATH=... make lint` | Exit 0; "skipping" message printed |
| `Verify_EnvExample_ContainsAllVars` | — | repo | Diff `.env.example` keys vs `infra.md` table | All 10 variables present; no real secrets |

### 8.3 Manual vs Automated

These verifications are performed manually as part of the Definition of Done for Story 02. Automated CI will run `make lint test` from Story 11 onward.
