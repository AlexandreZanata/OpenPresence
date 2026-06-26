# Tasks — Geofence engine

> **Depends on:** [task-00-monorepo-scaffold](../task-00-monorepo-scaffold/tasks.md) complete.

## Preparation

- [x] Read [README.md](README.md) and [official_source.md](official_source.md)
- [x] Run `./agent-harness/resolve-rules.sh domain layer geofence tdd`
- [x] Confirm task-00 scaffold exists (`services/attendance/go.mod`)

## TDD — write failing tests first

- [x] Create `internal/domain/geofence/types.go` — `GpsCoordinate`, `GeofenceZone`, `GeofenceType`
- [x] Create `internal/domain/geofence/checker.go` — `GeofenceChecker` interface
- [x] Create `internal/domain/geofence/geofence_test.go` with tests:
  - [x] `TestHaversineDistance_KnownPair`
  - [x] `TestIsInsideCircle_CenterPoint`
  - [x] `TestIsInsideCircle_OnBoundary`
  - [x] `TestIsInsideCircle_JustOutside`
  - [x] `TestIsInsideCircle_WithDeviation` (BR-020)
  - [x] `TestIsInsidePolygon_InsideConvex`
  - [x] `TestIsInsidePolygon_Outside`
  - [x] `TestIsInsidePolygon_ConcaveShape` (BR-021)

## Implementation

- [x] Implement `HaversineDistance` (Earth radius 6371000m)
- [x] Implement `IsInsideCircle`
- [x] Implement `IsInsidePolygon` (ray casting + deviation buffer)
- [x] Implement `IsInsideZone` on `DefaultChecker`
- [x] Implement `NearestZone`
- [x] Implement `IsInsideAnyZone` (BR-023)
- [x] BR-024 zone validity (`ValidFrom` / `ValidUntil`)

## Validation

- [x] `go test ./internal/domain/geofence/... -v` — all green
- [x] `go test -cover ./internal/domain/geofence/...`
- [x] `./scripts/verify-geofence.sh` — manual verification
- [x] Size caps and no outer layer imports verified

## Completion

- [x] All steps above marked `[x]`
- [x] Update `.local/phases/README.md` active task
