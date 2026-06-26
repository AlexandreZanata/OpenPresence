# Official sources — PunchRecord domain

## Repository documentation

| Document | Path |
|----------|------|
| Business rules BR-010–015 | `docs/BUSINESS-RULES.md` |
| Domain model — PunchRecord | `docs/DOMAIN-MODEL.md` |
| Testing — PunchRecord table | `docs/TESTING.md` |
| Fraud | `docs/FRAUD-DETECTION.md` |

## Agent rules

```bash
./agent-harness/resolve-rules.sh domain punch state-machine tdd
```

| Rule file | Why |
|-----------|-----|
| `agent-rules/02-architecture/state-machines.md` | PunchStatus transitions |

## Business rules

| ID | Summary |
|----|---------|
| BR-010 | All valid punch criteria |
| BR-011 | Offline TTL (status PENDING) |
| BR-014 | Punch type sequence |
| BR-015 | Server time official |

## Glossary terms

- `PunchRecord`, `PunchType`, `PunchStatus`, `BiometricResult`, `GpsCoordinate`
