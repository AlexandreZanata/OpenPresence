# Official sources — WorkSchedule

## Repository documentation

| Document | Path |
|----------|------|
| Business rules BR-030–034 | `docs/BUSINESS-RULES.md` |
| Domain model | `docs/DOMAIN-MODEL.md` |
| Glossary | `docs/GLOSSARY.md` |

## Agent rules

```bash
./agent-harness/resolve-rules.sh domain schedule tdd business-rules
```

## Business rules

| ID | Summary |
|----|---------|
| BR-030 | Worked time = in/out pairs minus breaks |
| BR-031 | Lateness after scheduled start + tolerance |
| BR-032 | Overtime after scheduled end + tolerance |
| BR-033 | Split shifts / 12×36 independent windows |
| BR-034 | Time bank cumulative balance |

## Glossary terms

- `WorkSchedule`, `PunchRecord`, `AttendancePolicy`
