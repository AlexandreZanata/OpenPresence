# Tasks — PunchRecord validation engine

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh domain punch state-machine tdd`

## Domain model

- [ ] Create `internal/domain/punch/` package
- [ ] `PunchRecord` aggregate with fields per `docs/DOMAIN-MODEL.md`
- [ ] `PunchValidationInput` — biometric, GPS, geofence, device time, server time, history
- [ ] `PunchValidator.Validate(input)` → `ValidationResult` (VALID | REJECTED + reasons)

## Rules (TDD — write tests first)

- [ ] BR-010: all criteria required for VALID
- [ ] BR-014: invalid sequence → REJECTED
- [ ] BR-015: `punchedAt` = server time; `deviceTime` audit only
- [ ] Anti-duplicate: second punch within 60s → REJECTED
- [ ] Clock manipulation: |device - server| > 300s → flag / reject per policy

## Validation

- [ ] `go test ./internal/domain/punch/... -v`
- [ ] `./scripts/verify-geofence.sh` still passes (no regression)
- [ ] Add punch domain section to `docs/TESTING.md`

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` active task
