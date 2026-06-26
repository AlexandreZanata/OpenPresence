# Task 06 — AttendancePolicy inheritance (Go)

**Status:** pending  
**Phase ID:** task-06-attendance-policy

## Goal

Implement `AttendancePolicy` value object and **inheritance merge** along the org tree: child nodes override parent fields. Support presets for **public administration** (strict geofence, 12×36 tolerance) and **private enterprise** (flexible windows, overtime rules).

## Scope

**In scope:**

- `internal/domain/organization/policy.go` (or `internal/domain/policy/`)
- Merge algorithm: walk ancestors root → node, apply overrides
- Policy fields per `docs/DOMAIN-MODEL.md`: tolerance, offline TTL, biometric/geofence flags, overtime
- TDD with multi-level tree fixtures

**Out of scope:**

- Persisting policy JSONB (repository in later task)
- Payroll export

## Acceptance

- Effective policy at leaf node reflects parent defaults + local overrides
- Tests cover secretariat vs corporate preset scenarios
- All [tasks.md](tasks.md) steps `[x]`

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
