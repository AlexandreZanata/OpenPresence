# Official sources — Hierarchy authorization

## Repository documentation

| Document | Path |
|----------|------|
| Organization — RBAC/ABAC | `docs/ORGANIZATION.md` |
| Security | `docs/SECURITY.md` |
| Use cases (manager approve) | `docs/use-cases/` |

## Agent rules

```bash
./agent-harness/resolve-rules.sh owasp authz authorization domain
```

| Rule file | Why |
|-----------|-----|
| `agent-rules/03-security/authorization.md` | Tenant + subtree isolation |

## Business rules

Supports BR-012 manager review workflow (authorization gate).

## Glossary terms

- `Role`, `OrgNode`, `Employee`, `Tenant`
