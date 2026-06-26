# Task 05 — Organization tree domain (Go)

**Status:** pending  
**Phase ID:** task-05-org-node-domain

## Goal

Implement the **Organization** bounded-context domain: tenant-scoped `OrgNode` tree with types for **public sector** (secretariats, locations, UBS) and **private sector** (divisions, departments, teams, work sites). Enforce tree invariants without infrastructure dependencies.

## Scope

**In scope:**

- `internal/domain/organization/` — `OrgNode`, `OrgNodeType`, tree operations
- Invariants: single root per tenant, no cycles, valid parent-child type rules
- Examples: municipality secretariat tree + corporate division tree (fixtures in tests)
- TDD unit tests (`*_test.go`)

**Out of scope:**

- HTTP APIs, SQL repositories, migrations for `org_nodes`
- RBAC / user assignment (task 11)

## Acceptance

- Tree invariants covered by tests (cycle, orphan, invalid type under parent)
- Public + private hierarchy fixtures documented in tests
- All [tasks.md](tasks.md) steps `[x]`

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
