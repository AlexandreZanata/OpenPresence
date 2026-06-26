# Task 14 — SubmitPunch use case (Go application)

**Status:** pending  
**Phase ID:** task-14-punch-submission-usecase

## Goal

Wire **application-layer** `SubmitPunch` use case: resolve employee placement → effective policy → geofence → work schedule context → biometric gRPC → punch validation + fraud → persist via tenant-scoped repository. End-to-end **backend punch processing** without mobile/UI.

## Scope

**In scope:**

- `internal/application/punch/submit_punch.go`
- Orchestrate domain packages from tasks 05–11 + geofence + biometric client port
- Integration test: mocked biometric gRPC + testcontainers Postgres + RLS
- `scripts/verify-punch-usecase.sh`

**Out of scope:**

- REST/Fiber HTTP layer
- Real ONNX (uses biometric stub or task 13 client)

## Acceptance

- Happy path punch → VALID in DB under correct tenant
- Cross-tenant submit rejected
- Out-of-geofence / invalid sequence rejected
- All [tasks.md](tasks.md) steps `[x]`

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
