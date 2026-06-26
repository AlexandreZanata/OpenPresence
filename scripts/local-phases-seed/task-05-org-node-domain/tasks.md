# Tasks — Organization tree domain

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh domain layer organization tree tdd`

## Domain model

- [ ] Create `internal/domain/organization/` package
- [ ] Define `OrgNodeType` enum: `DIVISION`, `DEPARTMENT`, `SECTION`, `TEAM`, `LOCATION`, `WORK_SITE`
- [ ] Define `OrgNode` entity: `id`, `tenantId`, `parentId`, `type`, `name`, `code` (optional slug)
- [ ] Document public-sector mapping: secretariat → `DIVISION`, hospital/UBS → `LOCATION`

## Tree invariants

- [ ] `OrgTree` (or equivalent) builder/validator: single root, no cycles
- [ ] Parent-child type rules (e.g. `TEAM` only under `DEPARTMENT` or `SECTION`)
- [ ] `WORK_SITE` may attach under `DIVISION` or `DEPARTMENT` (temporary sites)

## Tests (TDD)

- [ ] Fixture: municipality — Health Secretariat → Hospital → Nursing department
- [ ] Fixture: private company — HQ → Sales → Inside Sales team
- [ ] Reject cycle A→B→A
- [ ] Reject invalid child type under parent

## Validation

- [ ] `go test ./internal/domain/organization/...`
- [ ] `go vet ./internal/domain/organization/...`
- [ ] No imports from `infrastructure` or `interfaces`

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` active task
