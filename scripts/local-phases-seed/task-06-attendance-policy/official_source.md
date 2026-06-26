# Official sources — AttendancePolicy inheritance

## Repository documentation

| Document | Path |
|----------|------|
| Domain model — AttendancePolicy | `docs/DOMAIN-MODEL.md` |
| Organization | `docs/ORGANIZATION.md` |
| Business rules — offline TTL | `docs/BUSINESS-RULES.md` (BR-011) |
| Glossary | `docs/GLOSSARY.md` |

## Agent rules

```bash
./agent-harness/resolve-rules.sh domain policy inheritance tdd
```

| Rule file | Why |
|-----------|-----|
| `agent-rules/02-architecture/layering.md` | Domain-only merge logic |
| `agent-rules/04-testing/tdd.md` | Policy merge tests |

## Business rules

| ID | Summary |
|----|---------|
| BR-011 | `offlineSyncMaxAge` default 8h |

## Glossary terms

- `AttendancePolicy`, `OrgNode`, `Tenant`
