# Tasks — WorkSchedule & time accounting

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh domain schedule tdd business-rules`

## Domain model

- [ ] `WorkSchedule`: `scheduledStart`, `scheduledEnd`, `windows[]`, `toleranceMinutes`
- [ ] `WorkWindow` for split shifts and 12×36 cycles
- [ ] Attach schedule to employee (reference by ID in tests)

## Calculations (TDD first)

- [ ] `CalculateWorkedMinutes(punches)` — BR-030
- [ ] `CalculateLateness(clockIn, schedule, policy)` — BR-031
- [ ] `CalculateOvertime(clockOut, schedule, policy)` — BR-032
- [ ] `EvaluateWindows(day, schedule)` — BR-033 (12×36, split)
- [ ] `UpdateTimeBank(balance, overtimeMinutes, policy)` — BR-034

## Fixtures

- [ ] Public: nursing 12×36 — night window crosses midnight
- [ ] Private: 09:00–18:00 with 12:00–13:00 break
- [ ] Lateness 10 min with 5 min tolerance → 5 min late

## Validation

- [ ] `go test ./internal/domain/workforce/... -v`
- [ ] Document test matrix in `docs/TESTING.md` (work schedule section)

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` active task
