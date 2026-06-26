# Tasks — Employee placement

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh domain workforce employee tdd`

## Domain model

- [ ] Create `internal/domain/workforce/` package
- [ ] `EmployeePlacement`: `employeeId`, `orgNodeId`, `type` (PRIMARY|SECONDARY), `validFrom`, `validUntil` (optional)
- [ ] `PlacementService` or aggregate methods: assign, transfer, end placement

## Rules

- [ ] At most one active PRIMARY per employee at any instant
- [ ] SECONDARY optional; multiple allowed if non-overlapping per policy
- [ ] Placement must reference node in same `tenantId`
- [ ] `ActivePlacementAt(employee, date)` query for punch validation

## Tests (TDD)

- [ ] Public: employee primary at Municipal Hospital / Nursing
- [ ] Public: transfer from Education to Health — old placement closed
- [ ] Private: primary HQ Sales + secondary Field Sales route team
- [ ] Reject second PRIMARY while first active

## Validation

- [ ] `go test ./internal/domain/workforce/...`
- [ ] `./scripts/verify-scaffold.sh`

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` active task
