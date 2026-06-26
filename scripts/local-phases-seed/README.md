# Implementation phases (local)

Each subfolder is **one task**. Work one folder at a time until every step in `tasks.md` is `[x]`.

## Active task

**Current:** [`task-05-org-node-domain`](task-05-org-node-domain/)

Completed: [`task-00-monorepo-scaffold`](task-00-monorepo-scaffold/) … [`task-04-row-level-security`](task-04-row-level-security/)

## Folder layout (every task)

| File | Purpose |
|------|---------|
| `README.md` | Goal, scope, acceptance summary |
| `official_source.md` | Canonical references — repo docs + external links for the agent to read |
| `tasks.md` | Step-by-step checklist; mark `[x]` only when validated |

## Workflow

1. Open the active task folder.
2. Read `README.md` → `official_source.md` → `tasks.md`.
3. Run `./agent-harness/resolve-rules.sh` with keywords from `official_source.md`.
4. Complete steps in order; do not skip validation steps.
5. When all steps are `[x]`, set next task as active in this file and in `.local/tasks/current.md`.

## Task index

### Foundation (done)

| Folder | Status | Description |
|--------|--------|-------------|
| `task-00-monorepo-scaffold` | done | Go module + `services/attendance` skeleton |
| `task-01-biometric-service` | done | Rust gRPC biometric service (stub mode) |
| `task-02-geofence-engine` | done | Domain geofence (TDD) — BR-020–BR-024 |
| `task-03-punch-viewmodel` | done | KMP PunchViewModel |
| `task-04-row-level-security` | done | PostgreSQL RLS multi-tenant |

### Enterprise domain & backend (planned)

Backend-only: business logic, validations, hierarchy, placements, image pipeline.

| Folder | Status | Description |
|--------|--------|-------------|
| `task-05-org-node-domain` | **active** | Org tree — public secretariats + private divisions |
| `task-06-attendance-policy` | pending | AttendancePolicy inheritance along tree |
| `task-07-employee-placement` | pending | Lotação / placement — primary & secondary |
| `task-08-work-schedule` | pending | WorkSchedule + BR-030–034 time accounting |
| `task-09-punch-record-domain` | pending | PunchRecord validation — BR-010–015 |
| `task-10-fraud-detection` | pending | Fraud flags + device lockout — BR-012–013 |
| `task-11-hierarchy-authorization` | pending | ABAC — manager subtree scope |
| `task-12-model-download` | pending | ONNX model download + checksum manifest |
| `task-13-onnx-inference` | pending | Rust real ONNX face/liveness pipeline |
| `task-14-punch-submission-usecase` | pending | SubmitPunch application orchestration |

## Dependency order

```
05 org tree → 06 policy → 07 placement ─┐
                         → 08 schedule   ├→ 09 punch domain → 10 fraud → 11 ABAC ─┐
02 geofence ─────────────────────────────┘                                      │
12 model download → 13 ONNX inference ──────────────────────────────────────────┤
04 RLS ─────────────────────────────────────────────────────────────────────────┴→ 14 SubmitPunch
```

## New task

Copy `_template/` to `task-NN-short-name/` and fill the three files. Also add to `scripts/local-phases-seed/` for version control.
