# Official sources — Monorepo scaffold

## Repository documentation

| Document | Path |
|----------|------|
| Architecture | `docs/ARCHITECTURE.md` |
| Monorepo ADR | `docs/adr/ADR-003-monorepo-structure.md` |
| Stack | `docs/STACK.md` |
| Commit conventions | `docs/COMMIT-CONVENTIONS.md` |
| Agent guide | `docs/AGENT-IMPLEMENTATION-GUIDE.md` |

## Agent rules

```bash
./agent-harness/resolve-rules.sh architecture layering domain
```

| Rule file | Why |
|-----------|-----|
| `agent-rules/02-architecture/layering.md` | Package layout |
| `agent-rules/02-architecture/bounded-contexts.md` | attendance service boundary |
| `agent-rules/00-core/change-discipline.md` | One logical change per commit |

## External references

| Topic | URL |
|-------|-----|
| Go modules | https://go.dev/doc/modules/managing-dependencies |
| Standard Go project layout | https://github.com/golang-standards/project-layout |

## Glossary terms

- `Tenant`, `PunchRecord`, `GeofenceZone` — `docs/GLOSSARY.md`
