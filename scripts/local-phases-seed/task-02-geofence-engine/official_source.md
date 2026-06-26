# Official sources — Geofence engine

## Repository documentation

| Document | Path |
|----------|------|
| Business rules BR-020–024 | `docs/BUSINESS-RULES.md` |
| Domain model (GeofenceZone, GpsCoordinate) | `docs/DOMAIN-MODEL.md` |
| Testing (geofence test table) | `docs/TESTING.md` |
| Agent implementation Task 02 | `docs/AGENT-IMPLEMENTATION-GUIDE.md` |
| Glossary | `docs/GLOSSARY.md` |
| Architecture layers | `docs/ARCHITECTURE.md` |

## Agent rules

```bash
./agent-harness/resolve-rules.sh domain layer geofence tdd
./agent-harness/generate-task-rules.sh domain geofence tdd
```

| Rule file | Why |
|-----------|-----|
| `agent-rules/04-testing/tdd.md` | Red-green-refactor |
| `agent-rules/02-architecture/layering.md` | Domain-only package |
| `agent-rules/00-core/size-and-complexity-limits.md` | Function/file caps |
| `agent-rules/01-clean-code/functions.md` | Small pure functions |

## Business rules

| ID | Summary |
|----|---------|
| BR-020 | Circle: distance <= radius + allowedDeviation |
| BR-021 | Polygon: ray casting + deviation buffer |
| BR-022 | Low GPS accuracy → flag (Application layer later) |
| BR-023 | Match any assigned zone |
| BR-024 | Temporary zones by validFrom/validUntil |

## External references

| Topic | URL |
|-------|-----|
| Haversine formula | https://en.wikipedia.org/wiki/Haversine_formula |
| Ray casting (point in polygon) | https://en.wikipedia.org/wiki/Point_in_polygon |
| Earth radius (meters) | 6371000 — use constant in code |

## Glossary terms

- `GeofenceZone`, `GpsCoordinate`, `FraudType` (`OUT_OF_GEOFENCE`, `GPS_LOW_ACCURACY`)
