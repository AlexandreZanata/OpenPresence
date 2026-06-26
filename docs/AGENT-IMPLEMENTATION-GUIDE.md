# Agent Implementation Guide

Instructions for AI coding agents implementing OpenPresence. Read [AGENTS.md](../AGENTS.md) first.

## Prerequisites

```yaml
required_tools:
  go: ">=1.22"
  rust: ">=1.78"
  kotlin: ">=2.0"
  docker: ">=25"
  sqlx-cli: installed
  protoc: ">=3.21"
  buf: ">=1.30"

required_knowledge:
  - DDD aggregates, VO, domain events, repositories
  - TDD red-green-refactor
  - Clean architecture dependency rule
  - PostgreSQL RLS multi-tenancy
  - ONNX Runtime inference in Rust

forbidden_patterns:
  - Business logic in HTTP handlers
  - ORM hiding SQL (use sqlx)
  - Plaintext embeddings in APIs
  - device_time as official punch timestamp
  - Endpoints without tenant_id middleware
```

## Task 01 — Biometric Service (Rust)

**Goal:** gRPC server for biometric verification.

**Steps:**

1. `cargo new biometric-service` under `services/biometric/`
2. Dependencies: axum, tonic, ort, ndarray, image, tokio, serde_json
3. `FaceProcessor` struct with ONNX sessions loaded **once** at startup (`Arc<FaceProcessor>`)
4. Implement preprocess pipelines (liveness 80×80 BGR, recognition 112×112 RGB)
5. `ensemble_liveness` — average softmax, threshold 0.80
6. `cosine_similarity` for 512-dim vectors
7. gRPC `VerifyPunch` handler
8. Unit tests per pipeline function

**Note:** pgvector search stays in Go Attendance Service; Rust returns embedding to caller.

## Task 02 — Geofence Engine (Go)

**Goal:** `internal/domain/geofence/` with `GeofenceChecker` interface.

**Functions:** `HaversineDistance`, `IsInsideCircle`, `IsInsidePolygon` (ray casting), `NearestZone`

**TDD:** Write all tests in [TESTING.md](TESTING.md) geofence section before implementation.

## Task 03 — PunchViewModel (KMP)

**Goal:** `mobile/shared/.../PunchViewModel.kt` with sealed `PunchState` hierarchy.

**Flow:** device check → location → geofence → camera → liveness → submit → handle response

**Offline:** SQLDelight pending queue + background sync (stub repository in first pass).

**DI:** Koin module for ViewModel, Repository, BiometricProcessor, GeofenceValidator.

**Status:** implemented — see `mobile/shared/README.md` and `./scripts/verify-mobile.sh`.

## Task 04 — Row-Level Security (PostgreSQL)

Enable RLS on `employees`, `punch_records`, `face_embeddings`. Policy: `tenant_id = current_setting('app.tenant_id')::UUID`. Go: `SET LOCAL` inside transactions.

## Implementation constraints

### Security

- Embeddings: hash only in REST; never return raw vectors
- JWT 15min / refresh 7d
- mTLS internal services
- Argon2id passwords
- LGPD soft-delete + retention

### Performance

- IVFFlat pgvector `lists=100`
- Redis embedding cache TTL 5min
- Punch P99 < 500ms
- ONNX sessions at init

### Reliability

- Circuit breaker to biometric service
- NATS retry max 3
- `/health/live`, `/health/ready`
- Graceful shutdown

### Observability

OpenTelemetry span per punch with attributes listed in [INFRASTRUCTURE.md](INFRASTRUCTURE.md).

## Rule resolution by task

```bash
./agent-harness/resolve-rules.sh domain layer geofence    # Task 02
./agent-harness/resolve-rules.sh api endpoint auth          # REST handlers
./agent-harness/resolve-rules.sh owasp security biometric   # Security review
./agent-harness/generate-task-rules.sh domain geofence      # Cursor task scope
```

Delete `_task-active.mdc` when task completes.
