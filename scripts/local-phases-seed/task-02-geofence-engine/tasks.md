# Tasks — Geofence engine

> **Depends on:** [task-00-monorepo-scaffold](../task-00-monorepo-scaffold/tasks.md) complete (or equivalent `services/attendance` module).

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh domain layer geofence tdd`
- [ ] Confirm task-00 scaffold exists (`services/attendance/go.mod`)

## TDD — write failing tests first

- [ ] Create `internal/domain/geofence/types.go` — `GpsCoordinate`, `GeofenceZone`, `GeofenceType` (CIRCLE | POLYGON)
- [ ] Create `internal/domain/geofence/checker.go` — `GeofenceChecker` interface
- [ ] Create `internal/domain/geofence/geofence_test.go` with failing tests:
  - [ ] `TestHaversineDistance_KnownPair` — delta < 1m
  - [ ] `TestIsInsideCircle_CenterPoint` → true
  - [ ] `TestIsInsideCircle_OnBoundary` → true
  - [ ] `TestIsInsideCircle_JustOutside` → false
  - [ ] `TestIsInsideCircle_WithDeviation` → true (BR-020)
  - [ ] `TestIsInsidePolygon_InsideConvex` → true
  - [ ] `TestIsInsidePolygon_Outside` → false
  - [ ] `TestIsInsidePolygon_ConcaveShape` → true (BR-021)

## Implementation

- [ ] Implement `HaversineDistance(a, b GpsCoordinate) float64` (Earth radius 6371000m)
- [ ] Implement `IsInsideCircle(coord, center, radiusM, deviationM) bool`
- [ ] Implement `IsInsidePolygon(coord, polygon []GpsCoordinate) bool` (ray casting)
- [ ] Implement `IsInsideZone(coord, zone GeofenceZone) bool` on checker
- [ ] Implement `NearestZone(coord, zones []GeofenceZone) (*GeofenceZone, float64)`
- [ ] Implement `IsInsideAnyZone(coord, zones []GeofenceZone) bool` (BR-023)

## Validation

- [ ] `cd services/attendance && go test ./internal/domain/geofence/... -v` — all green
- [ ] `go test -cover ./internal/domain/geofence/...` — note coverage
- [ ] Verify each function ≤80 lines, each file ≤200 lines
- [ ] No imports from `internal/infrastructure` or `internal/interfaces` in domain package

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` — set `task-02` done, activate next task
- [ ] Run `./agent-harness/generate-task-rules.sh --clean` if task rules were generated
