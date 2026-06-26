# Implementation phases (local)

Each subfolder is **one task**. Work one folder at a time until every step in `tasks.md` is `[x]`.

## Active task

**Current:** [`task-03-punch-viewmodel`](task-03-punch-viewmodel/)

Completed: [`task-00-monorepo-scaffold`](task-00-monorepo-scaffold/), [`task-01-biometric-service`](task-01-biometric-service/), [`task-02-geofence-engine`](task-02-geofence-engine/)

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

| Folder | Status | Description |
|--------|--------|-------------|
| `task-00-monorepo-scaffold` | done | Go module + `services/attendance` skeleton |
| `task-01-biometric-service` | done | Rust gRPC biometric service |
| `task-02-geofence-engine` | done | Domain geofence (TDD) — BR-020–BR-024 |
| `task-03-punch-viewmodel` | **active** | KMP PunchViewModel |
| `task-04-row-level-security` | pending | PostgreSQL RLS multi-tenant |

## New task

Copy `_template/` to `task-NN-short-name/` and fill the three files.
