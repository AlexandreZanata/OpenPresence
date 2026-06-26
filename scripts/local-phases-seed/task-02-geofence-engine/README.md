# Task 02 — Geofence engine (Go domain)

**Status:** done  
**Phase ID:** task-02-geofence-engine

## Goal

Implement `internal/domain/geofence/` in the attendance service with TDD: Haversine distance, circle and polygon checks, and nearest zone lookup. Implements **BR-020** through **BR-024**.

## Scope

**In scope:**

- Pure domain package (no HTTP, no DB)
- `GeofenceChecker` interface and value types
- Unit tests written **before** implementation

**Out of scope:**

- REST endpoints
- Temporary zone date filtering (BR-024) — optional follow-up in same task if time permits
- Mobile geofence UI

## Acceptance

- All geofence tests in `docs/TESTING.md` pass
- All steps in [tasks.md](tasks.md) marked `[x]`
- Domain coverage contribution toward ≥90% target

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
