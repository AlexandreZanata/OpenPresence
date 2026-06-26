# Attendance Service

Go service for the **Attendance** bounded context: punch validation, geofence rules, fraud orchestration, and offline sync.

## Layers

| Package | Responsibility |
|---------|----------------|
| `internal/domain/geofence` | Geofence validation (Haversine, circle, polygon) — BR-020–BR-024 |
| `internal/domain/organization` | Org tree (`OrgNode`, `OrgTree`) — type rules, cycles, orphans |
| `internal/domain` | PunchRecord, fraud flags (upcoming) |
| `internal/application` | Use cases, authorization orchestration |
| `internal/infrastructure/postgres` | sqlx, RLS migrations, `WithTenant` transactions |
| `internal/interfaces` | HTTP handlers (Fiber), DTO mapping |

**Dependency rule:** domain does not import application, infrastructure, or interfaces.

## Commands

```bash
go build ./...
go test ./...
go vet ./...
go test -tags=integration ./internal/infrastructure/postgres/...
go test -cover ./internal/domain/geofence/...
go test -cover ./internal/domain/organization/...
```

From repo root:

```bash
./scripts/verify-scaffold.sh
./scripts/verify-geofence.sh
./scripts/verify-organization.sh
./scripts/verify-rls.sh
```

## Migrations

Versioned SQL in `migrations/` (001–006). Apply with `postgres.ApplyMigrations` or your migration runner.

RLS policies use `current_setting('app.tenant_id')::uuid`. Application queries must run inside `postgres.WithTenant`.

## Related docs

- `docs/DATA-MODEL.md` — tables and RLS pattern
- `docs/SECURITY.md` — multi-tenancy
- `docs/ORGANIZATION.md` — hierarchy examples and node types
- `docs/BUSINESS-RULES.md` — BR-010–BR-024
- `docs/ARCHITECTURE.md` — service map
