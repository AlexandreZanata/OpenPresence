# Attendance Service

Go service for the **Attendance** bounded context: punch validation, geofence rules, fraud orchestration, and offline sync.

## Layers

| Package | Responsibility |
|---------|----------------|
| `internal/domain/geofence` | Geofence validation (Haversine, circle, polygon) — BR-020–BR-024 |
| `internal/domain` | PunchRecord, fraud flags (upcoming) |
| `internal/application` | Use cases, authorization orchestration |
| `internal/infrastructure` | sqlx, Redis, NATS, gRPC clients |
| `internal/interfaces` | HTTP handlers (Fiber), DTO mapping |

**Dependency rule:** domain does not import application, infrastructure, or interfaces.

## Commands

```bash
go build ./...
go test ./...
go vet ./...
go test -cover ./internal/domain/geofence/...
```

From repo root:

```bash
./scripts/verify-scaffold.sh
./scripts/verify-geofence.sh
```

## Related docs

- `docs/DOMAIN-MODEL.md` — Attendance context
- `docs/BUSINESS-RULES.md` — BR-010–BR-024
- `docs/ARCHITECTURE.md` — service map
