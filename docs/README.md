# OpenPresence Documentation

> Biometric time and attendance platform with facial recognition, liveness detection, geofencing, and fraud detection.
> **Language:** English only (code, docs, commits, agent output).

## Index

| Document | Description |
|----------|-------------|
| [PRODUCT-OVERVIEW.md](PRODUCT-OVERVIEW.md) | Purpose, personas, product assumptions |
| [STACK.md](STACK.md) | Definitive technology stack |
| [BIOMETRICS.md](BIOMETRICS.md) | Open-source face recognition and liveness pipeline |
| [ARCHITECTURE.md](ARCHITECTURE.md) | System architecture and service map |
| [DOMAIN-MODEL.md](DOMAIN-MODEL.md) | Bounded contexts, aggregates, value objects |
| [BUSINESS-RULES.md](BUSINESS-RULES.md) | Business rules (BR-xxx) in GIVEN/WHEN/THEN |
| [FRAUD-DETECTION.md](FRAUD-DETECTION.md) | Fraud layers, matrix, responses |
| [ORGANIZATION.md](ORGANIZATION.md) | Org hierarchy, RBAC/ABAC |
| [MOBILE-FLOWS.md](MOBILE-FLOWS.md) | Mobile onboarding and punch flows |
| [API-CONTRACT.md](API-CONTRACT.md) | REST and gRPC contracts |
| [DATA-MODEL.md](DATA-MODEL.md) | PostgreSQL schema overview |
| [DOMAIN-EVENTS.md](DOMAIN-EVENTS.md) | Domain event catalog |
| [TESTING.md](TESTING.md) | TDD strategy and test cases by domain |
| [INFRASTRUCTURE.md](INFRASTRUCTURE.md) | Docker, K8s, observability |
| [SECURITY.md](SECURITY.md) | OWASP mapping, LGPD, implementation constraints |
| [COMMIT-CONVENTIONS.md](COMMIT-CONVENTIONS.md) | Commit and PR discipline |
| [CONTRIBUTING.md](CONTRIBUTING.md) | How to contribute |
| [AGENT-IMPLEMENTATION-GUIDE.md](AGENT-IMPLEMENTATION-GUIDE.md) | AI agent implementation tasks |
| [GLOSSARY.md](GLOSSARY.md) | Ubiquitous language |
| [NEW-PROJECT-CHECKLIST.md](NEW-PROJECT-CHECKLIST.md) | Pre-implementation checklist |
| [adr/](adr/) | Architecture Decision Records |
| [use-cases/](use-cases/) | Documented use cases |

## Quick links

- Agent entry point: [AGENTS.md](../AGENTS.md)
- Resolve rules: `./agent-harness/resolve-rules.sh <keywords>`
- Verify monorepo scaffold: `./scripts/verify-scaffold.sh`
