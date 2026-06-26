# Tasks — SubmitPunch use case

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh domain application layer punch integration`
- [ ] Tasks 05–11 domain packages and task 04 RLS in place

## Application layer

- [ ] `SubmitPunchCommand` — employeeId, type, GPS, deviceTime, frame bytes, device report
- [ ] `SubmitPunchHandler` — orchestration only; no business rules inline
- [ ] Ports: `EmployeeReader`, `PlacementReader`, `PolicyResolver`, `GeofenceResolver`, `BiometricClient`, `PunchRepository`

## Flow

- [ ] Load employee + active PRIMARY placement
- [ ] Resolve effective `AttendancePolicy` on placement node
- [ ] Resolve geofences for node subtree; run domain geofence checker
- [ ] Call biometric gRPC `VerifyPunch`
- [ ] Run `PunchValidator` + `FraudEvaluator`
- [ ] Persist via `postgres.WithTenant`

## Integration tests

- [ ] testcontainers: seed tenant, org tree, employee, geofence, placement
- [ ] Happy path → VALID row in `punch_records`
- [ ] Wrong tenant context → no row / error
- [ ] Invalid punch sequence → REJECTED, no VALID row

## Validation

- [ ] `go test -tags=integration ./internal/application/punch/...`
- [ ] `./scripts/verify-punch-usecase.sh`
- [ ] Update `docs/AGENT-IMPLEMENTATION-GUIDE.md` Task 14 status

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` active task
