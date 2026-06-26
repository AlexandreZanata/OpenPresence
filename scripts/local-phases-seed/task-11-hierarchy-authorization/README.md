# Task 11 — Hierarchy authorization (ABAC) (Go)

**Status:** pending  
**Phase ID:** task-11-hierarchy-authorization

## Goal

Implement **org-subtree authorization** for enterprise punch operations: managers approve only punches for employees in their descendant nodes; HR scoped to tenant; auditors read-only. Supports **secretariat managers** (public) and **department heads** (private).

## Scope

**In scope:**

- `internal/domain/organization/authorization.go` — `CanAccessSubtree(actorNode, targetNode)`, `IsDescendant`
- `internal/application/authorization/` — `PunchAuthorizationService` using org tree port
- Pure domain tests + application tests with mocked tree

**Out of scope:**

- JWT middleware, Fiber handlers
- Full RBAC user store

## Acceptance

- Manager at Health Secretariat cannot approve punch for Education employee
- ORG_ADMIN at division approves subtree only
- All [tasks.md](tasks.md) steps `[x]`

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
