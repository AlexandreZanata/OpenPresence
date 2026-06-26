# Task 10 — Fraud detection domain (Go)

**Status:** pending  
**Phase ID:** task-10-fraud-detection

## Goal

Implement **fraud flag** domain logic: classify anomalies, aggregate severity, decide SUSPICIOUS vs REJECTED (BR-012), and device lockout after consecutive failures (BR-013). Covers GPS, clock, device integrity, biometric failures.

## Scope

**In scope:**

- `internal/domain/fraud/` — `FraudFlag`, `FraudType`, `FraudEvaluator`
- Severity rules → `PunchStatus.SUSPICIOUS` vs auto-reject
- `DeviceLockoutTracker` — 3 rejects / 10 min → 30 min block

**Out of scope:**

- Manager review UI
- NATS alerts (infrastructure later)

## Acceptance

- Each `FraudType` in glossary has evaluation test
- BR-012 and BR-013 covered
- All [tasks.md](tasks.md) steps `[x]`

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
