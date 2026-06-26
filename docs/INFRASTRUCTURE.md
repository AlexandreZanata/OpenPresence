# Infrastructure

## Local development (Docker Compose)

Placeholder at `infra/docker-compose.yml` — services are commented until each microservice exists. Uncomment incrementally per [STACK.md](STACK.md).

**UC-001 E2E stack** (Postgres + attendance HTTP + biometric gRPC): `infra/docker-compose.e2e.yml`

```bash
docker compose -f infra/docker-compose.e2e.yml up -d --wait
./scripts/verify-uc001-e2e.sh --curl
docker compose -f infra/docker-compose.e2e.yml down
```

## Local development backend (host)

Faster iteration than full Docker compile (Rust/Go on host, Postgres in Docker):

```bash
./scripts/dev-backend.sh start    # Postgres :5433, attendance :8088, biometric :9090
./scripts/dev-backend.sh status
./scripts/dev-backend.sh stop
./scripts/verify-dev-backend.sh   # manual verification (starts stack if down)
```

Admin panel (planned `web/admin/`): `VITE_API_BASE_URL=http://127.0.0.1:8088`

```bash
cd web/admin && cp .env.example .env.local && npm install && npm run dev
```

See `web/admin/README.md`. Verify: `./scripts/verify-admin-scaffold.sh`.

Manual verification:

```bash
./scripts/verify-dev-backend.sh
```

**Note:** Auth service (`POST /v1/auth/login`) not implemented — admin UI will use dev mock until auth microservice exists.

Target stack (when enabled):

```yaml
# infra/docker-compose.yml
services:
  postgres:
    image: timescale/timescaledb-ha:pg16
  redis:
    image: valkey/valkey:7.2
  nats:
    image: nats:2.10-alpine
    command: -js
  api-gateway:
    build: ../services/api-gateway
  attendance:
    build: ../services/attendance
  biometric-service:
    build: ../services/biometric
```

Verify repo layout before enabling containers:

```bash
./scripts/verify-scaffold.sh
```

Ports (dev): Postgres 5432, Redis 6379, NATS 4222, API 8080, Biometric gRPC 9090.

## Production

```
GitHub Actions CI → container images → Helm charts (infra/k8s/) → Kubernetes
```

Optional IaC: `infra/terraform/` for cloud resources.

## Observability

| Signal | Stack |
|--------|-------|
| Traces | OpenTelemetry → Tempo |
| Logs | Loki |
| Metrics | Prometheus → Grafana |

### Required punch span attributes

`employee_id`, `tenant_id`, `punch_type`, `liveness_score`, `fraud_flags_count`, `geofence_id` — **never** raw biometric data.

### Prometheus metrics

- `openpresence_punch_total{status, tenant, fraud_type}`
- `openpresence_biometric_latency_seconds`
- `openpresence_liveness_score_histogram`

## Reliability

- Circuit breaker: API Gateway → Biometric Service
- NATS JetStream retries: max 3, exponential backoff
- Health: `/health/live`, `/health/ready` on every service
- Graceful shutdown: drain connections before exit

## Performance targets

| Metric | Target |
|--------|--------|
| Punch endpoint P99 | < 500ms (incl. biometric) |
| Embedding cache (Redis) | TTL 5min per tenant |
| pgvector IVFFlat | `lists=100` up to ~1M embeddings |
