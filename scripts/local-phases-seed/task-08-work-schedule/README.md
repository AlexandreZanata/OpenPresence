# Task 08 — WorkSchedule & time accounting (Go)

**Status:** pending  
**Phase ID:** task-08-work-schedule

## Goal

Implement `WorkSchedule` domain and **time accounting** logic per BR-030–BR-034: worked hours, lateness, overtime, split shifts (12×36), and time-bank accrual. Cover **public health shifts** and **private office** patterns.

## Scope

**In scope:**

- `internal/domain/workforce/schedule.go` — windows, tolerance, shift templates
- Pure functions: `CalculateWorkedMinutes`, `CalculateLateness`, `CalculateOvertime`, `EvaluateTimeBank`
- Templates: `Standard8h`, `Shift12x36`, `SplitShift`

**Out of scope:**

- REST reporting endpoints
- Legal payroll export formats (eSocial, etc.)

## Acceptance

- BR-030–BR-034 each covered by at least one table-driven test
- 12×36 and split-shift scenarios pass
- All [tasks.md](tasks.md) steps `[x]`

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
