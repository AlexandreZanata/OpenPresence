# New Project Checklist

> Complete **before writing the first line of application code**.
> Mirrors `agent-rules/AGENT-CORE-PRINCIPLES.md` checklist.
> If any item is blank, the agent **must ask** — never assume.

---

## Architecture and domain

- [x] **Layers defined** — [ARCHITECTURE.md](ARCHITECTURE.md)
- [x] **Entities and aggregates** — [DOMAIN-MODEL.md](DOMAIN-MODEL.md)
- [x] **Value Objects** — [DOMAIN-MODEL.md](DOMAIN-MODEL.md), [GLOSSARY.md](GLOSSARY.md)
- [x] **Business rules** — [BUSINESS-RULES.md](BUSINESS-RULES.md) (GIVEN/WHEN/THEN)
- [x] **State machines** — [DOMAIN-MODEL.md](DOMAIN-MODEL.md) (PunchStatus, PunchType sequence)
- [x] **Access roles** — [ORGANIZATION.md](ORGANIZATION.md) RBAC/ABAC matrix
- [x] **Domain events** — [DOMAIN-EVENTS.md](DOMAIN-EVENTS.md)
- [x] **Use cases** — [use-cases/](use-cases/) (UC-001–003; expand as needed)
- [x] **API contract** — [API-CONTRACT.md](API-CONTRACT.md)
- [x] **Glossary** — [GLOSSARY.md](GLOSSARY.md)

---

## Security (OWASP)

- [x] **OWASP Top 10:2025** — [SECURITY.md](SECURITY.md) summary mapping
- [x] **Agentic 2026 (ASI01–ASI10)** — referenced in [SECURITY.md](SECURITY.md); AI agents read-only + human confirm

---

## Agent harness

- [x] **Harness installed** — `agent-rules/`, `agent-harness/`, `.cursor/rules/` (restore via `./scripts/setup-local-ai.sh`)
- [x] **AGENTS.md** — project entry point for agent sessions
- [x] **Ponytail (static)** — `.cursor/rules/ponytail.mdc` via harness install
- [x] **Local AI workspace** — `.local/` tasks and `.cursor/` validation rules (gitignored)

---

## ADRs

- [x] ADR-001 Tech stack — [adr/ADR-001-tech-stack.md](adr/ADR-001-tech-stack.md)
- [x] ADR-002 Biometric stack — [adr/ADR-002-biometric-stack.md](adr/ADR-002-biometric-stack.md)
- [x] ADR-003 Monorepo — [adr/ADR-003-monorepo-structure.md](adr/ADR-003-monorepo-structure.md)

---

## Sign-off (pending)

| Role | Name | Date |
|------|------|------|
| Product / domain | | |
| Tech lead | | |

When sign-off is complete, implementation may begin per [AGENT-IMPLEMENTATION-GUIDE.md](AGENT-IMPLEMENTATION-GUIDE.md).
