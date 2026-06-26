# Attendance Service

Go service for the **Attendance** bounded context: punch validation, geofence rules, fraud orchestration, and offline sync.

## Layers

| Package | Responsibility |
|---------|----------------|
| `internal/domain/geofence` | Geofence validation (Haversine, circle, polygon) — BR-020–BR-024 |
| `internal/domain/organization` | Org tree (`OrgNode`, `OrgTree`), `AttendancePolicy`, ABAC subtree rules |
| `internal/application/authorization` | `PunchAuthorizationService`, `AuthorizePunchApprovalHandler` — manager/HR/auditor gates |
| `internal/application/enrollment` | `SaveFaceEmbeddingHandler` — persist embeddings after EnrollFace |
| `internal/application/punch` | `SubmitPunchHandler` — placement → policy → geofence → biometric → validate → fraud → lockout → persist |
| `internal/domain/punch` | `PunchRecord`, `PunchValidator` — BR-010–BR-015 |
| `internal/domain/fraud` | `FraudEvaluator`, `DeviceLockoutTracker` — BR-012–013 |
| `internal/domain/workforce` | Employee placement (*lotação*), `WorkSchedule`, time accounting BR-030–034 |
| `internal/domain` | PunchRecord, fraud flags (upcoming) |
| `internal/application` | Use cases, authorization orchestration |
| `internal/application/attendance` | `CalculateDayAttendanceHandler` — BR-030–034 from punches in DB |
| `internal/infrastructure/postgres` | sqlx, RLS migrations, `WithTenant`, `PunchRepository`, `FaceEmbeddingRepository` |
| `internal/interfaces` | HTTP handlers (Fiber), DTO mapping |

**Dependency rule:** domain does not import application, infrastructure, or interfaces.

## Commands

```bash
go build ./...
go test ./...
go vet ./...
go test -tags=integration ./internal/infrastructure/postgres/...
go test -tags=integration ./internal/application/punch/...
go test -tags=integration ./internal/application/punch/... -run E2E
go test -tags=integration ./internal/application/authorization/... -run E2E
go test -tags=integration ./internal/application/enrollment/... -run E2E_RLS
go test -cover ./internal/domain/geofence/...
go test -cover ./internal/domain/organization/...
go test -cover ./internal/domain/punch/...
go test -cover ./internal/domain/fraud/...
go test -cover ./internal/domain/workforce/...
```

From repo root:

```bash
./scripts/verify-scaffold.sh
./scripts/verify-geofence.sh
./scripts/verify-geofence-e2e.sh
./scripts/verify-organization.sh
./scripts/verify-attendance-policy.sh
./scripts/verify-workforce-placement.sh
./scripts/verify-work-schedule.sh
./scripts/verify-work-schedule-e2e.sh
./scripts/verify-punch.sh
./scripts/verify-punch-usecase.sh
./scripts/verify-fraud.sh
./scripts/verify-fraud-e2e.sh
./scripts/verify-authorization.sh
./scripts/verify-authorization-e2e.sh
./scripts/verify-rls.sh
./scripts/verify-rls-e2e.sh
```

## Migrations

Versioned SQL in `migrations/` (001–007). Apply with `postgres.ApplyMigrations` or your migration runner.

RLS policies use `current_setting('app.tenant_id')::uuid`. Application queries must run inside `postgres.WithTenant`.

## Related docs

- `docs/DATA-MODEL.md` — tables and RLS pattern
- `docs/SECURITY.md` — multi-tenancy
- `docs/ORGANIZATION.md` — hierarchy examples and node types
- `docs/BUSINESS-RULES.md` — BR-010–BR-034
- `docs/ARCHITECTURE.md` — service map
