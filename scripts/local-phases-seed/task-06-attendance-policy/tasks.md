# Tasks — AttendancePolicy inheritance

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh domain policy inheritance tdd`
- [ ] Task 05 org tree available (or minimal stub tree in tests)

## Domain model

- [ ] Define `AttendancePolicy` VO: `workdayDuration`, `toleranceMinutes`, `allowedPunchTypes`, `geofenceRequired`, `biometricRequired`, `offlineSyncMaxAge`, `overtimePolicy`
- [ ] `DefaultPolicy()` and `PublicSectorPreset()` / `PrivateSectorPreset()` factory helpers
- [ ] `MergePolicy(parent, override)` — field-level override semantics

## Inheritance

- [ ] `EffectivePolicy(node, ancestors[])` walks root → node
- [ ] Child `nil` fields inherit parent; explicit zero values documented
- [ ] Node may attach local policy partial override

## Tests (TDD)

- [ ] Tenant root policy → division override → department inherits merged result
- [ ] Public: secretariat sets `offlineSyncMaxAge` 8h; hospital node overrides tolerance only
- [ ] Private: HQ disables overtime; branch re-enables via override

## Validation

- [ ] `go test ./internal/domain/...` (policy package)
- [ ] `./scripts/verify-scaffold.sh` still passes

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` active task
