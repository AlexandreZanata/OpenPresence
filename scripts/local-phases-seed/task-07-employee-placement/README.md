# Task 07 — Employee placement / lotação (Go)

**Status:** pending  
**Phase ID:** task-07-employee-placement

## Goal

Model **employee placement** (*lotação*): assignment of an employee to one or more `OrgNode`s with effective dates. Support **public sector** mobility between secretariats and **private sector** transfers between departments/teams. Primary placement drives geofence and schedule defaults.

## Scope

**In scope:**

- `internal/domain/workforce/placement.go` — `EmployeePlacement`, `PlacementType` (PRIMARY, SECONDARY)
- Rules: exactly one active PRIMARY per employee; overlapping placements validated
- Link placement → effective org subtree for geofence inheritance
- TDD fixtures: servidor lotado em Secretaria de Saúde; colaborador em Sales + Field team secondary

**Out of scope:**

- HR workflow APIs, approval chains for transfer
- SQL migrations (stub in-memory repo OK for tests)

## Acceptance

- Cannot have two active PRIMARY placements
- Transfer scenario: end old placement, start new — tests pass
- All [tasks.md](tasks.md) steps `[x]`

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
