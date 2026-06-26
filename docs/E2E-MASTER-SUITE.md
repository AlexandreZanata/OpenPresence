# E2E master suite

One command runs unit scripts, integration checks, and all executable E2E phases (`e2e-01` … `e2e-11`).

## Commands

```bash
./scripts/verify-all-business-rules.sh                         # unit + integration
./scripts/verify-all-business-rules.sh --quick                 # domain/unit only (no Docker)
./scripts/verify-all-business-rules.sh --e2e                   # full + E2E phases + BR report
./scripts/verify-all-business-rules.sh --e2e --fail-fast       # stop on first failure
./scripts/verify-all-business-rules.sh --e2e --continue-on-error
```

Local phase entry (gitignored): `./.local/phases/e2e-testing/e2e-99-master-suite/run.sh`

## Report

After `--e2e`, generates `.local/reports/e2e-last-run.md` with pass/fail per business rule. Source manifest: `scripts/business-rules-coverage.json`.

## E2E phases

| Phase | `run.sh` | Rules / scope |
|-------|----------|----------------|
| e2e-00-harness | quick smoke | harness only |
| e2e-01-enrollment | `verify-enrollment.sh` | BR-001–003 |
| e2e-02-punch-core | domain + SubmitPunch integration | BR-010–015 |
| e2e-03-geofence | `verify-geofence-e2e.sh` | BR-020–024 |
| e2e-04-fraud | `verify-fraud-e2e.sh` | BR-012–013 |
| e2e-05-work-schedule | `verify-work-schedule-e2e.sh` | BR-030–034 |
| e2e-06-authorization | `verify-authorization-e2e.sh` | SEC-ABAC |
| e2e-07-rls-security | `verify-rls-e2e.sh` | SEC-RLS, SEC-CROSS-TENANT |
| e2e-08-biometric-stack | `verify-biometric-e2e.sh` | BR-002, SEC-BIOMETRIC-STUB |
| e2e-09-submit-punch-stack | `verify-punch-stack-e2e.sh` | BR-010, BR-014, BR-023 |
| e2e-10-mobile-kmp | `verify-mobile-e2e.sh` | SEC-MOBILE |
| e2e-11-full-uc001 | `verify-uc001-e2e.sh` | SEC-UC001, BR-010–011 |

## Documented skips (N/A — 28/28 complete)

These rules have an explicit **N/A** path until REST/admin APIs exist:

| Rule | Reason |
|------|--------|
| BR-004 | Administrator-only enrollment — needs admin auth REST |
| BR-005 | Maximum 5 active embeddings — needs enrollment admin API |
| BR-006 | Soft-delete on deactivate — needs employee lifecycle API |

Optional (not counted in 28/28): **SEC-BIOMETRIC-ONNX** — run `./scripts/download-models.sh` then `ONNX_MODELS_PATH=./models ./scripts/verify-biometric-e2e.sh`.

## Prerequisites

Go, Rust/Cargo, Docker (integration + E2E), Gradle/JDK (mobile phase e2e-10).
