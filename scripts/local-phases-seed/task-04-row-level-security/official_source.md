# Official sources — Row-level security

## Repository documentation

| Document | Path |
|----------|------|
| Data model | `docs/DATA-MODEL.md` |
| Security | `docs/SECURITY.md` |
| Agent guide Task 04 | `docs/AGENT-IMPLEMENTATION-GUIDE.md` |
| Infrastructure (Postgres) | `docs/INFRASTRUCTURE.md` |

## Agent rules

```bash
./agent-harness/resolve-rules.sh owasp authz tenant security migration
```

| Rule file | Why |
|-----------|-----|
| `agent-rules/03-security/authorization.md` | Tenant isolation |
| `agent-rules/07-data-management/migrations.md` | Versioned SQL |
| `agent-rules/03-security/ssrf-and-access-control.md` | No cross-tenant leak |

## Business rules

N/A — infrastructure concern; supports all tenant-scoped aggregates.

## External references

| Topic | URL |
|-------|-----|
| PostgreSQL RLS | https://www.postgresql.org/docs/current/ddl-rowsecurity.html |
| sqlx | https://github.com/launchbadge/sqlx |

## Glossary terms

- `Tenant`, `Employee`, `PunchRecord`
