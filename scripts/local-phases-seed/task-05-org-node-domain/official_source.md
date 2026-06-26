# Official sources — Organization tree domain

## Repository documentation

| Document | Path |
|----------|------|
| Organization hierarchy | `docs/ORGANIZATION.md` |
| Domain model | `docs/DOMAIN-MODEL.md` |
| Glossary | `docs/GLOSSARY.md` |
| Architecture | `docs/ARCHITECTURE.md` |

## Agent rules

```bash
./agent-harness/resolve-rules.sh domain layer organization tree tdd
```

| Rule file | Why |
|-----------|-----|
| `agent-rules/02-architecture/layering.md` | Pure domain, no infra imports |
| `agent-rules/04-testing/tdd.md` | Tests before implementation |
| `agent-rules/11-documentation-and-glossary/ubiquitous-language.md` | OrgNode terms |

## Business rules

N/A — structural domain; enables placement and ABAC in later tasks.

## External references

| Topic | URL |
|-------|-----|
| DDD aggregates | https://martinfowler.com/bliki/DDD_Aggregate.html |

## Glossary terms

- `Tenant`, `OrgNode`, `OrganizationNode`
