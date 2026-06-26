# Official sources — Employee placement

## Repository documentation

| Document | Path |
|----------|------|
| Domain model — Employee | `docs/DOMAIN-MODEL.md` |
| Organization | `docs/ORGANIZATION.md` |
| Glossary | `docs/GLOSSARY.md` |

## Agent rules

```bash
./agent-harness/resolve-rules.sh domain workforce employee tdd
```

## Business rules

Structural — enables BR-023 (geofence from assigned locations via placement node).

## Glossary terms

- `Employee`, `OrgNode`, `Tenant`
