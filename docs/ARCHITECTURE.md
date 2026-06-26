# System Architecture

## High-level diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                     MOBILE APP (KMP)                            │
│  Compose UI ─ ViewModel ─ UseCases ─ Domain ─ Repositories      │
│       ↕ HTTPS/mTLS              ↕ SQLDelight (offline)          │
└───────────────────────┬─────────────────────────────────────────┘
                        │
              ┌─────────▼──────────┐
              │  API GATEWAY (Go)  │
              │  Fiber v3 + JWT    │
              │  Rate limit, RBAC  │
              └──┬──────────┬──────┘
                 │          │
    ┌────────────▼──┐  ┌────▼──────────────────┐
    │  Auth Service  │  │  Attendance Service   │
    │  (Go/Fiber)    │  │  (Go/Fiber)           │
    │  JWT + Refresh │  │  Punch rules, policy  │
    └────────────────┘  └────────────┬──────────┘
                                     │ gRPC
                          ┌──────────▼──────────────┐
                          │ Biometric Service (Rust)│
                          │ Axum + ONNX Runtime     │
                          │ RetinaFace + AuraFace   │
                          │ MiniFASNet ensemble     │
                          └──────────┬──────────────┘
                                     │
              ┌──────────────────────┼─────────────────────┐
    ┌─────────▼──────┐   ┌───────────▼──────┐   ┌─────────▼───────┐
    │ PostgreSQL 16   │   │ Redis (Valkey)    │   │ NATS JetStream  │
    │ + pgvector      │   │ Sessions/Cache    │   │ Async events    │
    │ + TimescaleDB   │   └──────────────────┘   └─────────────────┘
    └─────────────────┘
```

## Layered design (per service)

All Go and Rust services follow the same layering:

| Layer | Responsibility |
|-------|----------------|
| **Interfaces** | HTTP/gRPC handlers, DTO mapping, middleware |
| **Application** | Use cases, orchestration, authorization checks |
| **Domain** | Aggregates, entities, value objects, domain services, events |
| **Infrastructure** | SQL (sqlx), Redis, NATS, gRPC clients, ONNX (Rust only) |

**Dependency rule:** Domain never depends on outer layers. Infrastructure implements domain ports.

## Services

| Service | Language | Responsibility |
|---------|----------|----------------|
| `api-gateway` | Go | Routing, JWT validation, rate limiting, tenant context |
| `auth` | Go | Login, refresh, device registration |
| `attendance` | Go | Punch validation, fraud orchestration, sync |
| `organization` | Go | Org tree, geofences, attendance policies |
| `workforce` | Go | Employees, biometric enrollment |
| `biometric` | Rust | Face detection, liveness, embedding extraction (gRPC) |

## Cross-cutting concerns

- **Multi-tenancy:** `tenant_id` on all rows; PostgreSQL RLS via `SET LOCAL app.tenant_id`
- **mTLS** between internal microservices
- **OpenTelemetry** spans on every punch (no raw biometric data in traces)
- **Circuit breaker** between gateway and biometric service

## Repository layout (monorepo)

```
openpresence/
├── services/
│   ├── api-gateway/
│   ├── attendance/
│   ├── organization/
│   ├── workforce/
│   └── biometric/
├── web/
│   └── admin/           # TanStack Start admin panel
├── mobile/
│   ├── shared/
│   ├── androidApp/
│   └── iosApp/
├── models/              # ONNX weights (not committed — see BIOMETRICS.md)
├── infra/
│   ├── k8s/
│   ├── terraform/
│   └── docker-compose.yml
└── docs/
```

See [ADR-003](adr/ADR-003-monorepo-structure.md).
