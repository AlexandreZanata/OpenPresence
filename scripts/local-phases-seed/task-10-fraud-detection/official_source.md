# Official sources — Fraud detection

## Repository documentation

| Document | Path |
|----------|------|
| Fraud detection | `docs/FRAUD-DETECTION.md` |
| Business rules BR-012–013 | `docs/BUSINESS-RULES.md` |
| Domain model — FraudFlag | `docs/DOMAIN-MODEL.md` |
| Glossary | `docs/GLOSSARY.md` |

## Agent rules

```bash
./agent-harness/resolve-rules.sh domain fraud security tdd
```

## Business rules

| ID | Summary |
|----|---------|
| BR-012 | FraudFlag → SUSPICIOUS when not auto-reject |
| BR-013 | 3 rejects / 10 min → 30 min device lockout |

## Glossary terms

- `FraudFlag`, `FraudType`, `PunchRecord`, `DeviceIntegrityReport`
