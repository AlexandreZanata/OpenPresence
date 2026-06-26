# Official sources — PunchViewModel

## Repository documentation

| Document | Path |
|----------|------|
| Mobile flows | `docs/MOBILE-FLOWS.md` |
| UC-001 Clock in | `docs/use-cases/UC-001-punch-clock-in.md` |
| Business rules BR-010–015 | `docs/BUSINESS-RULES.md` |
| API punch endpoint | `docs/API-CONTRACT.md` |
| Fraud (device layer) | `docs/FRAUD-DETECTION.md` |
| Agent guide Task 03 | `docs/AGENT-IMPLEMENTATION-GUIDE.md` |
| Stack (KMP) | `docs/STACK.md` |

## Agent rules

```bash
./agent-harness/resolve-rules.sh domain mobile state tdd
```

| Rule file | Why |
|-----------|-----|
| `agent-rules/02-architecture/state-machines.md` | PunchState transitions |
| `agent-rules/02-architecture/layering.md` | ViewModel in application/presentation |

## Business rules

| ID | Summary |
|----|---------|
| BR-010 | Valid punch criteria |
| BR-011 | Offline TTL |
| BR-014 | Punch type sequence |
| BR-015 | Server timestamp official |

## External references

| Topic | URL |
|-------|-----|
| Compose Multiplatform | https://www.jetbrains.com/compose-multiplatform/ |
| Koin KMP | https://insert-koin.io/docs/reference/koin-mp/kmp/ |
| Ktor Client | https://ktor.io/docs/client.html |

## Glossary terms

- `PunchRecord`, `PunchType`, `PunchStatus`, `DeviceIntegrityReport`, `GpsCoordinate`
