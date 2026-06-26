# Official sources — SubmitPunch use case

## Repository documentation

| Document | Path |
|----------|------|
| UC-001 Clock in | `docs/use-cases/UC-001-punch-clock-in.md` |
| API contract — punch | `docs/API-CONTRACT.md` |
| Business rules BR-010–015 | `docs/BUSINESS-RULES.md` |
| Architecture | `docs/ARCHITECTURE.md` |
| Testing | `docs/TESTING.md` |

## Agent rules

```bash
./agent-harness/resolve-rules.sh domain application layer punch integration
```

| Rule file | Why |
|-----------|-----|
| `agent-rules/02-architecture/layering.md` | Application orchestrates domain |

## Business rules

| ID | Summary |
|----|---------|
| BR-010 | Valid punch |
| BR-014 | Sequence |
| BR-023 | Geofence from placement |

## Glossary terms

- `PunchRecord`, `Employee`, `AttendancePolicy`, `GeofenceZone`
