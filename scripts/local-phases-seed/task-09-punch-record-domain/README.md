# Task 09 — PunchRecord validation engine (Go)

**Status:** pending  
**Phase ID:** task-09-punch-record-domain

## Goal

Implement `PunchRecord` aggregate and **punch validation engine** for BR-010–BR-015: valid punch criteria, type sequence, anti-duplicate, clock manipulation, official server timestamp. Core enterprise punch logic.

## Scope

**In scope:**

- `internal/domain/punch/` — `PunchRecord`, `PunchValidator`, `PunchStatus` transitions
- Integrate inputs: biometric result VO, GPS, geofence match, device time, recent punches
- Reject / accept / suspicious decision hooks (fraud flags attached in task 10)

**Out of scope:**

- HTTP handler, database persistence
- Biometric gRPC (task 14)

## Acceptance

- BR-010–BR-015 each have dedicated tests
- Punch sequence CLOCK_IN → BREAK_START → BREAK_END → CLOCK_OUT enforced
- All [tasks.md](tasks.md) steps `[x]`

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
