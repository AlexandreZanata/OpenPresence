# Tasks — Hierarchy authorization

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh owasp authz authorization domain`
- [ ] Task 05 org tree + Task 07 placement available

## Domain

- [ ] `IsDescendant(ancestor, node, tree)` in organization package
- [ ] `ActorScope`: `role`, `assignedOrgNodeId`, `tenantId`
- [ ] `CanApprovePunch(actor, employeePlacement, tree)` — MANAGER subtree rule
- [ ] `CanExportPayroll(actor, tenant)` — HR_ANALYST tenant boundary

## Application

- [ ] `PunchAuthorizationService` in `internal/application/authorization/`
- [ ] Port `OrgTreeReader` for tests

## Tests (TDD)

- [ ] Public: Health manager approves hospital nurse; rejects education teacher
- [ ] Private: Sales manager approves inside sales; rejects IT employee
- [ ] AUDITOR: read allowed, write denied
- [ ] Cross-tenant always denied

## Validation

- [ ] `go test ./internal/domain/organization/... ./internal/application/authorization/...`
- [ ] `./scripts/verify-scaffold.sh`

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` active task
