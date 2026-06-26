# Tasks — Fraud detection domain

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh domain fraud security tdd`

## Domain model

- [ ] Create `internal/domain/fraud/` package
- [ ] `FraudFlag` VO: `type`, `severity`, `detectedAt`, `metadata`
- [ ] `FraudEvaluator.Evaluate(input)` → flags + recommended status
- [ ] `DeviceLockoutTracker` state machine for BR-013

## Rules (TDD)

- [ ] `MOCK_GPS` → HIGH severity
- [ ] `CLOCK_MANIPULATION` when delta > 300s
- [ ] `IMPOSSIBLE_SPEED` between last punch and current GPS
- [ ] `DUPLICATE_PUNCH` within 60s
- [ ] `GPS_LOW_ACCURACY` → LOW (BR-022) — accept with flag
- [ ] BR-012: CRITICAL flags → REJECTED; else SUSPICIOUS
- [ ] BR-013: lockout after 3 REJECTED in 10 minutes

## Validation

- [ ] `go test ./internal/domain/fraud/... -v`
- [ ] Update `docs/TESTING.md` fraud section

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` active task
